// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

package tradehill

import (
	"appengine"
	"os"
	"restapi"
	"xgen"
	"strconv"
	"sort"
	"appengine/datastore"
	"time"
)

func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

// GetQuote retrieves the "ticker" data.
func GetQuote(c appengine.Context) (x xgen.Quote, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var q Quote
	err = restapi.GetJson(c, JsonTicker, &q)
	check(err)
	x.Date = datastore.SecondsToTime(time.Seconds())
	x.HighestBuy, err = strconv.Atof64(q.Ticker.Buy)
	check(err)
	x.LowestSell, err = strconv.Atof64(q.Ticker.Sell)
	check(err)
	x.Last, err = strconv.Atof64(q.Ticker.Last)
	check(err)
	if !x.Validate() {
		panic("Invalid Ticker")
	}
	return
}

// GetOrderBook retrieves the limit order book.
func GetOrderBook(c appengine.Context) (x xgen.OrderBook, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var b OrderBook
	err = restapi.GetJson(c, JsonDepth, &b)
	check(err)
	for _, ask := range b.Asks {
		var o xgen.Order
		o.Price, err = strconv.Atof64(ask[Price])
		check(err)
		o.Amount, err = strconv.Atof64(ask[Amount])
		check(err)
		x.SellTree = append(x.SellTree, o)
	}
	for _, bid := range b.Bids {
		var o xgen.Order
		o.Price, err = strconv.Atof64(bid[Price])
		check(err)
		o.Amount, err = strconv.Atof64(bid[Amount])
		check(err)
		x.BuyTree = append(x.BuyTree, o)
	}
	if !x.Validate() {
		panic("Invalid Depth")
	}
	sort.Sort(x.BuyTree)
	sort.Sort(x.SellTree)
	x.BuyTree.Reverse()
	return
}

// GetBalance retrieves the account balance.
func GetBalance(c appengine.Context, login xgen.Credentials) (x xgen.Balance, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var b Balance
	err = restapi.PostJson(c, JsonBalance, map[string][]string{"name": {login.Username}, "pass": {login.Password}}, &b)
	check(err)
	x[xgen.BTC], err = strconv.Atof64(b.BTC)
	check(err)
	x[xgen.USD], err = strconv.Atof64(b.USD)
	return
}

func (h OpenOrders) convert() (o xgen.OpenOrders) { // Convert tradehill.OpenOrders to xgen.OpenOrders
	var err os.Error
	o.Sell = make(map[string]xgen.OpenOrder)
	o.Buy = make(map[string]xgen.OpenOrder)
	for _, order := range h.Orders {
		var t xgen.OpenOrder
		t.Date = order.Date
		t.Price, err = strconv.Atof64(order.Price)
		check(err)
		t.Amount, err = strconv.Atof64(order.Amount)
		check(err)
		if order.OrderType == 1 { // Sell order
			o.Sell[strconv.Itoa64(order.Oid)] = t
		} else if order.OrderType == 2 { // Buy order
			o.Buy[strconv.Itoa64(order.Oid)] = t
		} else {
			panic("Invalid order type")
		}
	}
	return
}

// GetOpenOrders retrieves a list of all open orders.
func GetOpenOrders(c appengine.Context, login xgen.Credentials) (x xgen.OpenOrders, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var openOrders OpenOrders
	err = restapi.PostJson(c, JsonOrders, map[string][]string{"name": {login.Username}, "pass": {login.Password}}, &openOrders)
	check(err)
	x = openOrders.convert()
	return
}

// CancelOrder cancels an open order.
func CancelOrder(c appengine.Context, login xgen.Credentials, oid string) (x xgen.OpenOrders, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var openOrders OpenOrders
	err = restapi.PostJson(c, JsonCancel, map[string][]string{"name": {login.Username}, "pass": {login.Password}, "oid": {oid}}, &openOrders)
	check(err)
	x = openOrders.convert() // Note: Canceling an order may take up to 1 second. During that time the order is considered active and will be returned as active by GetOrders.
	return
}

// Buy opens a new order to buy BTC.
func Buy(c appengine.Context, login xgen.Credentials, price float64, amount float64) (x xgen.OpenOrders, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var openOrders OpenOrders
	err = restapi.PostJson(c, JsonBuy, map[string][]string{"name": {login.Username}, "pass": {login.Password},
		"price": {strconv.Ftoa64(price, 'f', -1)}, "amount": {strconv.Ftoa64(amount, 'f', -1)}}, &openOrders)
	check(err)
	x = openOrders.convert()
	return
}

// Sell opens a new order to sell BTC.
func Sell(c appengine.Context, login xgen.Credentials, price float64, amount float64) (x xgen.OpenOrders, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var openOrders OpenOrders
	err = restapi.PostJson(c, JsonSell, map[string][]string{"name": {login.Username}, "pass": {login.Password},
		"price": {strconv.Ftoa64(price, 'f', -1)}, "amount": {strconv.Ftoa64(amount, 'f', -1)}}, &openOrders)
	check(err)
	x = openOrders.convert()
	return
}
