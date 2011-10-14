// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

// Package mtgox implements functions for sending and receiving data via Mt Gox API.
package mtgox

// Mt. Gox Trade API: https://mtgox.com/support/tradeAPI
// API documentation: https://en.bitcoin.it/wiki/MtGox/API
const (
	// Public Market Data
	JsonTicker = "https://mtgox.com/api/0/data/ticker.php"
	JsonDepth  = "https://mtgox.com/api/0/data/getDepth.php"
	JsonRecent = "https://mtgox.com/api/0/data/getTrades.php"

	// Authenticated Trading Functions
	JsonBalance = "https://mtgox.com/api/0/getFunds.php"
	JsonOrders  = "https://mtgox.com/api/0/getOrders.php"
	JsonBuy     = "https://mtgox.com/api/0/buyBTC.php"
	JsonSell    = "https://mtgox.com/api/0/sellBTC.php"
	JsonCancel  = "https://mtgox.com/api/0/cancelOrder.php"
)

type quote struct {
	Buy  float64 // Highest Bid
	Sell float64 // Lowest Ask
	Vol  int64   // Traded volume since last close
	High float64 // Maximum rate (?) since last close
	Low  float64 // Minimum rate (?) since last close
	Last float64
	Avg  float64
	Vwap float64
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
	Date      int64  // Unix timestamp of the trade	
	Tid       string // Trade id (big integer, which is in fact the trade timestamp in microseconds)	
	Price_int string // Price per unit * 1E5 or 1E3 (1E5 for USD, 1E3 for JPY)
	// USD: multiply by 0.00001 to get the actual price of BTC in USD, JPY: multiply by 0.001
	Amount_int     string // Traded amount * 1E8 (multiply by 0.00000001 to get the actual number of BTC traded)
	Price_currency string // Currency in which trade was completed (USD? JPY?)
	Item           string // What was this trade about (BTC?)
	Trade_type     string // Did this trade result from the execution of a 'bid' or 'ask'?

	//	Old API values with float types - deprecated
	//	Price	string 
	//	Amount	string
}

// RecentTrades is a struct representing a slice of historical trades.
type RecentTrades struct {
	Trades []Trade
}

// Balance contains the amount of each currency in the account.
type Balance struct {
	Usds string
	Btcs string
}

type newBalance struct {
	Currency  string // USD / BTC
	Value     string // Float
	Value_int string // Int *1E5 or *1E8 ?
	Display   string // String format including the currency symbol (e.g. "$")
}

type openOrder struct {
	OrderType   int8 `json:"type"` // 1 = Sell order, 2 = Buy order
	Status      int8   // 1 = Active
	Real_status string // "open" or ?
	Oid         string // Unique order id (but not int, e.g. "2d8dfdd1-5342-42bc-9dde-a81fdfa63920")
	Currency    string // ? Was "USD" when trying to buy BTC with USD
	Item        string // ? Was "BTC" when trying to buy BTC
	Price       string // Price of BTC in USD ?
	Price_int   string // Price of BTC in USD ? times 1E5
	Amount      string // Amount of BTC ?
	Amount_int  string // Amount of BTC ? times 1E8
	Dark        int8   // Dark Pool ?
	Priority    string // ? int
	Date        int64
}

// OpenOrders is a struct representing all our open buy and sell orders in the account.
type OpenOrders struct {
	Usds     string
	Btcs     string
	New_Usds newBalance
	New_btcs newBalance
	Orders   []openOrder
}
