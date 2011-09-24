// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

package exchb

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
	x.HighestBuy = q.Ticker.Buy
	x.LowestSell = q.Ticker.Sell
	x.Last = q.Ticker.Last
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
		o := xgen.Order{Price: ask[Price], Amount: ask[Amount]}
		x.SellTree = append(x.SellTree, o)
	}
	for _, bid := range b.Bids {
		o := xgen.Order{Price: bid[Price], Amount: bid[Amount]}
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
	var b Balance
	err = restapi.PostJson(c, JsonBalance, map[string][]string{"name": {login.Username}, "pass": {login.Password}}, &b)
	x[xgen.BTC] = b.Btcs
	x[xgen.USD] = b.Usds
	return
}

func (h OpenOrders) convert() (o xgen.OpenOrders) { // Convert exchb.OpenOrders to xgen.OpenOrders
	var err os.Error
	o.Sell = make(map[string]xgen.OpenOrder)
	o.Buy = make(map[string]xgen.OpenOrder)
	for _, order := range h.Orders {
		var t xgen.OpenOrder
		t.Date, err = strconv.Atoi64(order.Date)
		check(err)
		//t.Price, err = strconv.Atof64(order.Price) // fails because of the dollar sign in the string
		//check(err)
		t.Amount, err = strconv.Atof64(order.Amount)
		check(err)
		if order.OrderType == "Sell" {
			o.Sell[order.Oid] = t
		} else if order.OrderType == "Buy" {
			o.Buy[order.Oid] = t
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
	x = openOrders.convert()
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
