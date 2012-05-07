// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the MIT/X11 license.

// Package campbx implements functions for sending and receiving data via CampBX API.
package campbx

// CampBX API: https://campbx.com/api.php
const (
	// Public Market Data
	JsonTicker = "http://campbx.com/api/xticker.php"
	JsonDepth  = "http://CampBX.com/api/xdepth.php"

	// Authenticated Trading Functions
	JsonBalance = "https://CampBX.com/api/myfunds.php"
	JsonOrders  = "https://CampBX.com/api/myorders.php"
	JsonBuySell = "https://CampBX.com/api/tradeenter.php"
	JsonCancel  = "https://CampBX.com/api/tradecancel.php"
)

// Quote is a struct representing the best available buy and sell prices at the time.
type Quote struct {
	Buy  string `json:"Best Bid"`
	Sell string `json:"Best Ask"`
	Last string `json:"Last Trade"`
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

// Balance contains the amount of each currency in the account.
type Balance struct {
	UsdTotal  string `json:"Total USD"`
	UsdLiquid string `json:"Liquid USD"`
	UsdMargin string `json:"Margin Account USD"`
	BtcTotal  string `json:"Total BTC"`
	BtcLiquid string `json:"Liquid BTC"`
	BtcMargin string `json:"Margin Account BTC"`
}

type openOrder struct {
	Info        string // Optional, e.g. "No open Buy Orders"
	OrderType   string `json:"Order Type"` // e.g. "Quick Sell"
	Oid         string `json:"Order ID"`   // Unique order id	
	Price       string // Price of BTC in USD
	Quantity    string // Amount of BTC
	MarginPct   string `json:"Margin Percent"`
	StopLoss    string `json:"Stop-loss"`
	FillType    string `json:"Fill Type"`
	DarkPool    string `json:"Dark Pool"`
	DateEntered string `json:"Order Entered"`
	DateExpires string `json:"Order Expiry"`
}

// OpenOrders is a struct representing all our open buy and sell orders in the account.
type OpenOrders struct {
	Buy  []openOrder
	Sell []openOrder
}
