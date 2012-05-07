// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the MIT/X11 license.

package xgen

// TODO: Store monetary values as fixed points (instead of floating points), at least if native support gets added to Go

// List of currencies.
const (
	BTC = iota
	USD

	NumCurrencies
)

// Credentials stores the username and password for the account.
type Credentials struct {
	Username string
	Password string
}

// Balance contains the amount of each currency in the account.
type Balance [NumCurrencies]float64

// OpenOrder contains basic information of an order.
type OpenOrder struct {
	Date   int64   // Timestamp for when the order was created
	Price  float64 // Price per BTC (usually in USD)
	Amount float64 // Amount of BTC
}

// OpenOrders is a struct representing all our open buy and sell orders in the account.
type OpenOrders struct {
	Buy  map[string]OpenOrder // Each order has a unique ID, used as the map key
	Sell map[string]OpenOrder // Some exchanges don't use integers for the Order ID's, therefore using string instead
}
