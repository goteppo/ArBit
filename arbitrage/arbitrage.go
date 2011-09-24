// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

// Package arbitrage implements functions for calculating Bitcoin arbitrage strategies.
package arbitrage

import (
	"xgen"
	"sort"
	//	"math"
)

// Strategy is a struct for storing the calculated trading strategy.
type Strategy struct {
	Buy  []xgen.Order
	Sell []xgen.Order
}

type arbOrder struct {
	order    xgen.Order
	exchange int8
}
type arbOrders []arbOrder
type arbOrderBook struct { // Master limit order book (after combining limit orders from all exchanges)
	buyTree  arbOrders
	sellTree arbOrders
}

func (m arbOrders) Len() int           { return len(m) }
func (m arbOrders) Less(i, j int) bool { return m[i].order.Price < m[j].order.Price }
func (m arbOrders) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m arbOrders) Reverse() {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}
}

// Calculate calculates an optimal arbitrage strategy.
func Calculate(book []xgen.OrderBook, funds []xgen.Balance, commission []float64, minTrade [][xgen.NumCurrencies]float64) (arb Strategy) {
	arb.Buy = make([]xgen.Order, len(book))
	arb.Sell = make([]xgen.Order, len(book))

	// Combine all order books into one
	var arbBook arbOrderBook
	for i, b := range book {
		if !b.Validate() {
			continue
		}
		for _, o := range b.BuyTree {
			var a arbOrder
			a.order = o
			a.exchange = int8(i)
			arbBook.buyTree = append(arbBook.buyTree, a)
		}
		for _, o := range b.SellTree {
			var a arbOrder
			a.order = o
			a.exchange = int8(i)
			arbBook.sellTree = append(arbBook.sellTree, a)
		}
	}
	sort.Sort(arbBook.buyTree)
	sort.Sort(arbBook.sellTree)
	arbBook.buyTree.Reverse()

	fundsLeft := make([]xgen.Balance, len(funds))
	copy(fundsLeft, funds)

	// Find the arbitrage trades
	buyer, seller := 0, 0
	buyerExchange := arbBook.buyTree[buyer].exchange
	sellerExchange := arbBook.sellTree[seller].exchange
	buyerAmount := arbBook.buyTree[buyer].order.Amount
	sellerAmount := arbBook.sellTree[seller].order.Amount
	for arbBook.buyTree[buyer].order.Price*(1-commission[buyerExchange]) > arbBook.sellTree[seller].order.Price/(1-commission[sellerExchange]) {
		// If not enough BTC in the account for the minimum allowed trade size (sell), it's the same as if the account was empty of BTC
		if fundsLeft[buyerExchange][xgen.BTC] < minTrade[buyerExchange][xgen.BTC] ||
			fundsLeft[buyerExchange][xgen.BTC] < minTrade[buyerExchange][xgen.USD]/arbBook.buyTree[buyer].order.Price {
			fundsLeft[buyerExchange][xgen.BTC] = 0
		}
		// If not enough USD in the account for the minimum allowed trade size (buy), it's the same as if the account was empty of USD
		if fundsLeft[sellerExchange][xgen.USD] < minTrade[sellerExchange][xgen.USD] ||
			fundsLeft[sellerExchange][xgen.USD] < minTrade[sellerExchange][xgen.BTC]*arbBook.sellTree[seller].order.Price {
			fundsLeft[sellerExchange][xgen.USD] = 0
		}

		// Arbitrage found -> update the execution plan as long as funds still left
		if fundsLeft[buyerExchange][xgen.BTC] > 0 {
			arb.Sell[buyerExchange].Price = arbBook.buyTree[buyer].order.Price
		}
		if fundsLeft[sellerExchange][xgen.USD] > 0 {
			arb.Buy[sellerExchange].Price = arbBook.sellTree[seller].order.Price
		}

		// Can't make a bigger trades than we have funds for
		if buyerAmount > fundsLeft[buyerExchange][xgen.BTC] {
			buyerAmount = fundsLeft[buyerExchange][xgen.BTC]
		}
		if sellerAmount*arbBook.sellTree[seller].order.Price > fundsLeft[sellerExchange][xgen.USD] {
			sellerAmount = fundsLeft[sellerExchange][xgen.USD] / arbBook.sellTree[seller].order.Price
		}

		// Available arbitrage is limited to the volume of the smaller side (buyer/seller).
		// For the sake of simplicity, we keep the (absolute) amount of BTC same after the trade as it was before the trade (i.e. all the profit will be in USD),
		// (even though theoretically a more correct way might be to keep the relative balances the same for all currencies).
		switch {
		case buyerAmount > sellerAmount*(1-commission[buyerExchange]):
			buyerCapped := sellerAmount * (1 - commission[buyerExchange])
			arb.Sell[buyerExchange].Amount += buyerCapped
			arb.Buy[sellerExchange].Amount += sellerAmount
			fundsLeft[buyerExchange][xgen.BTC] -= buyerCapped
			fundsLeft[sellerExchange][xgen.USD] -= sellerAmount * arb.Buy[sellerExchange].Price
			buyerAmount -= buyerCapped
			seller++
			if seller == len(arbBook.sellTree) {
				return
			}
			sellerAmount = arbBook.sellTree[seller].order.Amount
		case buyerAmount < sellerAmount*(1-commission[buyerExchange]):
			sellerCapped := buyerAmount / (1 - commission[buyerExchange])
			arb.Sell[buyerExchange].Amount += buyerAmount
			arb.Buy[sellerExchange].Amount += sellerCapped
			fundsLeft[buyerExchange][xgen.BTC] -= buyerAmount
			fundsLeft[sellerExchange][xgen.USD] -= sellerCapped * arb.Buy[sellerExchange].Price
			sellerAmount -= sellerCapped
			buyer++
			if buyer == len(arbBook.buyTree) {
				return
			}
			buyerAmount = arbBook.buyTree[buyer].order.Amount
		default:
			arb.Sell[buyerExchange].Amount += buyerAmount
			arb.Buy[sellerExchange].Amount += sellerAmount
			fundsLeft[buyerExchange][xgen.BTC] -= buyerAmount
			fundsLeft[sellerExchange][xgen.USD] -= sellerAmount * arb.Buy[sellerExchange].Price
			buyer++
			seller++
			if buyer == len(arbBook.buyTree) || seller == len(arbBook.sellTree) {
				return
			}
			buyerAmount = arbBook.buyTree[buyer].order.Amount
			sellerAmount = arbBook.sellTree[seller].order.Amount
		}
		buyerExchange = arbBook.buyTree[buyer].exchange
		sellerExchange = arbBook.sellTree[seller].exchange
	}
	return
}

// Onesided adjusts the amounts in an existing strategy to balance the USD and BTC amounts within each exchange.
func Onesided(strategy Strategy, funds []xgen.Balance, commission []float64) (newStgy Strategy) {
	newStgy.Buy = make([]xgen.Order, len(strategy.Buy))
	newStgy.Sell = make([]xgen.Order, len(strategy.Sell))
	copy(newStgy.Buy, strategy.Buy)
	copy(newStgy.Sell, strategy.Sell)

	for i := 0; i < len(funds); i++ {
		usdLeft := funds[i][xgen.USD] + strategy.Sell[i].Amount*strategy.Sell[i].Price*(1-commission[i]) - strategy.Buy[i].Amount*strategy.Buy[i].Price
		btcLeft := funds[i][xgen.BTC] + strategy.Buy[i].Amount*(1-commission[i]) - strategy.Sell[i].Amount
		if strategy.Buy[i].Price > 0 && usdLeft > btcLeft*strategy.Buy[i].Price {
			newStgy.Buy[i].Amount = strategy.Buy[i].Amount + (usdLeft/strategy.Buy[i].Price-btcLeft)/(2-commission[i])
		}
		if strategy.Sell[i].Price > 0 && btcLeft*strategy.Sell[i].Price > usdLeft {
			newStgy.Sell[i].Amount = strategy.Sell[i].Amount + (btcLeft-usdLeft/strategy.Sell[i].Price)/(2-commission[i])
		}
	}
	return
}
