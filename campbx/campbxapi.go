// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

package campbx

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
	x.HighestBuy, err = strconv.Atof64(q.Buy)
	check(err)
	x.LowestSell, err = strconv.Atof64(q.Sell)
	check(err)
	x.Last, err = strconv.Atof64(q.Last)
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
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var b Balance
	err = restapi.PostJson(c, JsonBalance, map[string][]string{"user": {login.Username}, "pass": {login.Password}}, &b)
	check(err)
	x[xgen.BTC], err = strconv.Atof64(b.BtcLiquid)
	check(err)
	x[xgen.USD], err = strconv.Atof64(b.UsdLiquid)
	return
}

func (h OpenOrders) convert() (o xgen.OpenOrders) { // Convert campbx.OpenOrders to xgen.OpenOrders
	var err os.Error
	o.Sell = make(map[string]xgen.OpenOrder)
	o.Buy = make(map[string]xgen.OpenOrder)
	for _, order := range h.Buy {
		if order.Oid != "" {
			var t xgen.OpenOrder
			t.Date = time.Seconds() // Should use time.Parse to convert order.DateEntered to Unix time
			t.Price, err = strconv.Atof64(order.Price)
			check(err)
			t.Amount, err = strconv.Atof64(order.Quantity)
			check(err)
			o.Buy[order.Oid] = t
		}
	}
	for _, order := range h.Sell {
		if order.Oid != "" {
			var t xgen.OpenOrder
			t.Date = time.Seconds() // Should use time.Parse to convert order.DateEntered to Unix time
			t.Price, err = strconv.Atof64(order.Price)
			check(err)
			t.Amount, err = strconv.Atof64(order.Quantity)
			check(err)
			o.Sell[order.Oid] = t
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
	err = restapi.PostJson(c, JsonOrders, map[string][]string{"user": {login.Username}, "pass": {login.Password}}, &openOrders)
	check(err)
	x = openOrders.convert()
	return
}

// CancelOrder cancels an open order.
func CancelOrder(c appengine.Context, login xgen.Credentials, oid string, orderType string) (f interface{}, err os.Error) {
	err = restapi.PostJson(c, JsonCancel, map[string][]string{"user": {login.Username}, "pass": {login.Password}, "Type": {orderType}, "OrderID": {oid}}, &f)
	return
}

// Buy opens a new order to buy BTC.
func Buy(c appengine.Context, login xgen.Credentials, price float64, amount float64) (f interface{}, err os.Error) {
	err = restapi.PostJson(c, JsonBuySell, map[string][]string{"user": {login.Username}, "pass": {login.Password}, "TradeMode": {"QuickBuy"},
		"Price": {strconv.Ftoa64(price, 'f', -1)}, "Quantity": {strconv.Ftoa64(amount, 'f', -1)}}, &f)
	return
}

// Sell opens a new order to sell BTC.
func Sell(c appengine.Context, login xgen.Credentials, price float64, amount float64) (f interface{}, err os.Error) {
	err = restapi.PostJson(c, JsonBuySell, map[string][]string{"user": {login.Username}, "pass": {login.Password}, "TradeMode": {"QuickSell"},
		"Price": {strconv.Ftoa64(price, 'f', -1)}, "Quantity": {strconv.Ftoa64(amount, 'f', -1)}}, &f)
	return
}
