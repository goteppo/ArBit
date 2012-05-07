// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the MIT/X11 license.

// Package tradehill implements functions for sending and receiving data via TradeHill API.
package tradehill

// TradeHill Trading API: https://www.tradehill.com/Support/TradingAPI/
// Description of the data: http://bitcoincharts.com/about/exchanges/
const (
	// Public Market Data
	JsonTicker = "https://api.tradehill.com/APIv1/USD/Ticker"
	JsonDepth  = "https://api.tradehill.com/APIv1/USD/Orderbook"
	JsonRecent = "https://api.tradehill.com/APIv1/USD/Trades"

	// Authenticated Trading Functions
	JsonBalance = "https://api.tradehill.com/APIv1/USD/GetBalance"
	JsonOrders  = "https://api.tradehill.com/APIv1/USD/GetOrders"
	JsonBuy     = "https://api.tradehill.com/APIv1/USD/BuyBTC"
	JsonSell    = "https://api.tradehill.com/APIv1/USD/SellBTC"
	JsonCancel  = "https://api.tradehill.com/APIv1/USD/CancelOrder"
)

type quote struct {
	Buy       string
	Sell      string
	Last      string
	Vol       string
	High      string
	Low       string
	Last_when string // E.g. "5 minutes ago"
}

// Quote is a struct representing the best available buy and sell prices at the time.
type Quote struct {
	Ticker quote
}

// OrderBook is a struct representing a limit order book.
type OrderBook struct {
	Asks [][2]string // Sell orders (price and amount)
	Bids [][2]string // Buy orders (price and amount)
}

const (
	Price = iota
	Amount
)

// Trade is a struct representing a historical trade.
type Trade struct {
	Date   int64  // Unix timestamp of the trade
	Tid    int64  // Unique trade id (monotonically increasing integer) for each trade	
	Price  string // Price in your markets currency (e.g. USD)
	Amount string // Amount of bitcoins exchanged in that trade
}

// RecentTrades is a struct representing a slice of historical trades.
type RecentTrades struct {
	Trades []Trade
}

// Balance contains the amount of each currency in the account.
type Balance struct {
	USD           string // Can be "0E-10" ?
	USD_Available string
	USD_Reserved  string
	BTC           string
	BTC_Available string
	BTC_Reserved  string
}

type openOrder struct {
	OrderType         int8 `json:"type"` // 1 = Sell order, 2 = Buy order
	Status            int8   // 1 = Active (which is the value for all returned order, any other value here indicates an error)
	Oid               int64  // Unique order id
	Symbol            string // Currency symbol (currently always "BTC")
	Price             string // Limit price of the order (in USD)
	Amount_orig       string // The original amount (size) of your order
	Amount            string // The amount (size) of your order that remains to be filled (amount = amount_orig - amount filled)
	Reserved_amount   string // The amount of money reserved in order to fill this order (includes estimated commissions)
	Reserved_currency string // The currency of the reserved money (USD for buy orders, BTC for sell orders)
	Date              int64  // The datetime this order was created
}

// OpenOrders is a struct representing all our open buy and sell orders in the account.
type OpenOrders struct {
	Orders []openOrder
}
