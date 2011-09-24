// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

// Package xgen provides generic structures for storing data from various Bitcoin exchanges
package xgen

// TODO: Store monetary values as fixed points (instead of floating points), at least if native support gets added to Go

import "appengine/datastore"

// Quote is a struct representing the best available buy and sell prices at the time.
type Quote struct {
	Date       datastore.Time
	HighestBuy float64
	LowestSell float64
	Last       float64
}

// Validate method checks that |Quote| values are valid (non-zero).
func (t Quote) Validate() bool {
	if t.HighestBuy == 0 || t.LowestSell == 0 || t.Last == 0 {
		return false
	}
	return true
}

// Order is a struct representing a single limit order from the order book.
type Order struct {
	Price  float64 // Price per BTC
	Amount float64 // Number of BTC
}
type orders []Order

// OrderBook is a struct representing a limit order book.
type OrderBook struct {
	BuyTree  orders // Bids
	SellTree orders // Asks
}

// Methods needed for sorting the orders in the order book (by price).
func (m orders) Len() int           { return len(m) }
func (m orders) Less(i, j int) bool { return m[i].Price < m[j].Price }
func (m orders) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// Reverse will reverse the order of orders in the order book.
func (m orders) Reverse() {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}
}

// Validate method checks that |OrderBook| is valid (not empty and non-zero values).
func (m OrderBook) Validate() bool {
	if len(m.SellTree) == 0 || len(m.BuyTree) == 0 {
		return false
	}
	for _, ask := range m.SellTree {
		if ask.Price == 0 || ask.Amount == 0 {
			return false
		}
	}
	for _, bid := range m.BuyTree {
		if bid.Price == 0 || bid.Amount == 0 {
			return false
		}
	}
	return true
}

// Trade is a struct representing a historical trade.
type Trade struct {
	Date   int64   // Unix timestamp
	Tid    int64   // Unique trade id
	Price  float64 // Price in USD (or other currency)
	Amount float64 // Amount of BTC
}

// UniqueKey is a method identifying the trade id, to be used as a key by the datastore.
func (t Trade) UniqueKey() (string, int64) {
	return "", t.Tid
}

// RecentTrades is a struct representing a slice of historical trades.
type RecentTrades struct {
	Trades []Trade
}

// Validate method checks that |RecentTrades| is valid (not empty and non-zero values).
func (t RecentTrades) Validate() bool {
	if len(t.Trades) == 0 {
		return false
	}
	for _, trade := range t.Trades {
		// In theory zero is a valid value for trade.Date, but in practice it means something went wrong (same for trade.Tid)
		if trade.Date == 0 || trade.Tid == 0 || trade.Price == 0 || trade.Amount == 0 {
			return false
		}
	}
	return true
}
