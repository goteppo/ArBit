// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

// Package exchb implements functions for sending and receiving data via ExchB API.
package exchb

// ExchB Trade API: https://www.exchangebitcoins.com/api
const (
	// Public Market Data
	JsonTicker = "https://www.exchangebitcoins.com/data/ticker"
	JsonDepth  = "https://www.exchangebitcoins.com/data/depth"
	JsonRecent = "https://www.exchangebitcoins.com/data/recent"

	// Authenticated Trading Functions
	JsonBalance = "https://www.exchangebitcoins.com/data/getFunds"
	JsonOrders  = "https://www.exchangebitcoins.com/data/getOrders"
	JsonBuy     = "https://www.exchangebitcoins.com/data/buyBTC"
	JsonSell    = "https://www.exchangebitcoins.com/data/sellBTC"
	JsonCancel  = "https://www.exchangebitcoins.com/data/cancelOrder"
)

type quote struct {
	Buy  float64
	Sell float64
	Last float64
	Vol  float64
	High float64
	Low  float64
}

// Quote is a struct representing the best available buy and sell prices at the time.
type Quote struct {
	Ticker quote
}

// OrderBook is a struct representing a limit order book.
type OrderBook struct {
	Asks [][2]float64 // Sell orders (price and amount)
	Bids [][2]float64 // Buy orders (price and amount)
}

const (
	Price = iota
	Amount
)

// Trade is a struct representing a historical trade.
type Trade struct {
	Date   int64   // Unix timestamp of the trade
	Tid    int64   // Unique trade id (monotonically increasing integer) for each trade	
	Price  float64 // Price in USD
	Amount float64 // Amount of bitcoins exchanged in that trade
}

// RecentTrades is a struct representing a slice of historical trades.
type RecentTrades struct {
	Trades []Trade
}

// Balance contains the amount of each currency in the account.
type Balance struct {
	Usds float64
	Btcs float64
}

type openOrder struct {
	OrderType string "type" // "Buy", "Sell"
	//	OrderType	string	`json:"type"`// r59 release
	Oid    string // Unique order id
	Price  string // Price of BTC in USD (the string contains a dollar sign!)
	Amount string // Amount of BTC
	Date   string
}

// OpenOrders is a struct representing all our open buy and sell orders in the account.
type OpenOrders struct {
	//Status		string
	Usds   float64
	Btcs   float64
	Ticker quote
	Orders []openOrder
}
