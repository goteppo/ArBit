// Copyright 2011 Teppo Salonen. All rights reserved.
// This program is distributed under the terms of the GNU LGPLv3 license.

// Package arbit is the "main" package of ArBit.
package arbit

// TODO: Send out an email with the arbitrage strategy executed?
// TODO: Show the list of executed trades and calculated strategies in the dashboard (current dahsboard is pretty useless)
// TODO: Replace hardcoded 500ms delays between API calls with exchange specific sleep times (100ms for TradeHill, 500ms for CampBX, etc)
// TODO: Improve debugging, logging, and add more data to the dashboard

import (
	"os"
	"fmt"
	"http"
	"xgen"
	"mtgox"
	"tradehill"
	//	"campbx"
	"appdb"
	"arbitrage"
	"time"
	"appengine"
	"appengine/datastore"
)

// Bitcoin exchanges to be used for arbitrage
const (
	mtGox = iota
	tradeHill
	//	campBx
	//	bitcoinica

	numExchanges
)

var exchangeName = [numExchanges]string{"MtGox", "TradeHill" /*, "CampBX", "Bitcoinica"*/ }

var login [numExchanges]xgen.Credentials

var commission [numExchanges]float64                   // Commission per trade (by exchange)
var minTrade [numExchanges][xgen.NumCurrencies]float64 // Minimum transaction size (by exchange by currency)

var onesidedArb bool // Take one side of an arbitrage even if not enough funds on the other account (used for balancing USD and BTC within the account)
var paperTrade bool  // In Paper Trade mode trades will not be executed

func init() {
	http.HandleFunc("/cron/", errorHandlerLog(cronjob))
	http.HandleFunc("/dashboard/", errorHandlerWeb(dashboard))
	http.HandleFunc("/testing/", errorHandlerWeb(unittests))
}

func errorHandlerLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(os.Error); ok {
				c := appengine.NewContext(r)
				c.Criticalf(e.String())
			}
		}()
		fn(w, r)
	}
}

func errorHandlerWeb(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(os.Error); ok {
				w.WriteHeader(500)
				fmt.Fprintln(w, e)
			}
		}()
		fn(w, r)
	}
}

func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func dashboard(w http.ResponseWriter, r *http.Request) { // Simple monitoring dashboard
	//var err os.Error
	c := appengine.NewContext(r)

	fmt.Fprintln(w, "<table>")
	fmt.Fprintln(w, "<tr><th>Exchange</th><th>Updated</th><th>Highest Buy</th><th>Lowest Sell</th><th>Last Trade</th></tr>")

	// Read ticker data from datastore
	var ticker [numExchanges][]datastore.Map
	for i := int8(0); i < numExchanges; i++ {
		ticker[i], _ = appdb.Query(c, "Ticker_"+exchangeName[i], "", nil, "-Date", 0, 1) // Get the last ticker only
		fmt.Fprintln(w, "<tr><td>", exchangeName[i], "</td><td>", time.SecondsToLocalTime(int64(ticker[i][0]["Date"].(datastore.Time))/1e6),
			"</td><td>", ticker[i][0]["HighestBuy"], "</td><td>", ticker[i][0]["LowestSell"], "</td><td>", ticker[i][0]["Last"], "</td></tr>")
	}
	fmt.Fprintln(w, "</table>")
}

func cronjob(w http.ResponseWriter, r *http.Request) { // Main program (to be run as a cron job)
	var err os.Error
	c := appengine.NewContext(r)

	// Quotes/Tickers by exchange
	var quote [numExchanges]xgen.Quote
	quote[mtGox], err = mtgox.GetQuote(c)
	check(err)
	quote[tradeHill], err = tradehill.GetQuote(c)
	check(err)
	//	quote[campBx], err = campbx.GetQuote(c) // For some reason the Unmarshal in GetJson function in api.go causes: "runtime error: invalid memory address or nil pointer dereference"
	//	check(err)

	// Store ticker data in datastore
	for i := int8(0); i < numExchanges; i++ {
		err = appdb.KeyPut(c, "Ticker_"+exchangeName[i], &quote[i], "", time.Seconds())
		check(err)
	}

	// Check if arbitrage exists
	var maxBid, minAsk float64
	for ex, q := range quote {
		bid := q.HighestBuy * (1 - commission[ex])
		if maxBid == 0 || bid > maxBid {
			maxBid = bid
		}
		ask := q.LowestSell / (1 - commission[ex])
		if minAsk == 0 || ask < minAsk {
			minAsk = ask
		}
	}
	if maxBid < minAsk {
		fmt.Fprintln(w, "No Arbitrage Exists: Highest Buy ", maxBid, " < Lowest Sell ", minAsk, "<br>")
		return
	}

	time.Sleep(0.5 * 1e9) // Wait for half a second before the next API calls

	// Account balances by exchange
	var funds [numExchanges]xgen.Balance
	funds[mtGox], err = mtgox.GetBalance(c, login[mtGox])
	check(err)
	funds[tradeHill], err = tradehill.GetBalance(c, login[tradeHill])
	check(err)
	//	funds[campBx], err = campbx.GetBalance(c, login[campBx])
	//	check(err)

	time.Sleep(0.5 * 1e9) // Wait for half a second before the next API calls

	// Open orders by exchange
	var pending [numExchanges]xgen.OpenOrders
	pending[mtGox], err = mtgox.GetOpenOrders(c, login[mtGox])
	check(err)
	pending[tradeHill], err = tradehill.GetOpenOrders(c, login[tradeHill])
	check(err)
	//	pending[campBx], err = campbx.GetOpenOrders(c, login[campBx])
	//	check(err)

	time.Sleep(0.5 * 1e9) // Wait for half a second before the next API calls

	// Cancel any open orders
	for i := int8(0); i < numExchanges; i++ {
		for oid, _ := range pending[i].Buy {
			if !paperTrade {
				switch i {
				case mtGox:
					_, err = mtgox.CancelOrder(c, login[i], oid, 2)
				case tradeHill:
					_, err = tradehill.CancelOrder(c, login[i], oid)
					//					case campBx:  _, err = campbx.CancelOrder(c, login[i], oid, "Buy")
				}
				check(err)
				// Store the canceled order in datastore
				err = appdb.KeyPut(c, "Cancel_"+exchangeName[i], &pending[i].Buy, "", time.Seconds())
				check(err)
				time.Sleep(0.5 * 1e9) // Wait for half a second before the next API calls
			}
		}
		for oid, _ := range pending[i].Sell {
			if !paperTrade {
				switch i {
				case mtGox:
					_, err = mtgox.CancelOrder(c, login[i], oid, 1)
				case tradeHill:
					_, err = tradehill.CancelOrder(c, login[i], oid)
					//					case campBx:  _, err = campbx.CancelOrder(c, login[i], oid, "Sell")
				}
				check(err)
				// Store the canceled order in datastore
				err = appdb.KeyPut(c, "Cancel_"+exchangeName[i], &pending[i].Sell, "", time.Seconds())
				check(err)
				time.Sleep(0.5 * 1e9) // Wait for half a second before the next API calls
			}
		}
	}

	// Limit order books by exchange
	var book [numExchanges]xgen.OrderBook
	book[mtGox], err = mtgox.GetOrderBook(c)
	check(err)
	book[tradeHill], err = tradehill.GetOrderBook(c)
	check(err)
	//	book[campBx], err = campbx.GetOrderBook(c)
	//	check(err)

	// Find the arbitrage strategy
	strategy := arbitrage.Calculate(book[:], funds[:], commission[:], minTrade[:])

	// Check for internal arbitrage within the same exchange, since those should not happen if the data is correct and the exchange is working correctly
	for i := int8(0); i < numExchanges; i++ {
		if strategy.Buy[i].Amount > 0 && strategy.Sell[i].Amount > 0 {
			panic(fmt.Sprintln("Arbitrage within", exchangeName[i], "order books"))
		}
	}

	// If one-sided trades allowed, use them for balancing the total USD and BTC within accounts
	if onesidedArb {
		strategy = arbitrage.Onesided(strategy, funds[:], commission[:])
	}

	time.Sleep(0.5 * 1e9) // One second is 1e9 nanoseconds

	// Execute trades
	for i := int8(0); i < numExchanges; i++ {
		if strategy.Buy[i].Amount > 0 {
			fmt.Fprintln(w, exchangeName[i], ": Buy", strategy.Buy[i].Amount, "bitcoins for", strategy.Buy[i].Price, "USD per BTC <br>")

			// Execute trade:
			if !paperTrade {
				switch i {
				case mtGox:
					_, err = mtgox.Buy(c, login[i], strategy.Buy[i].Price, strategy.Buy[i].Amount)
				case tradeHill:
					_, err = tradehill.Buy(c, login[i], strategy.Buy[i].Price, strategy.Buy[i].Amount)
					//					case campBx: _, err = campbx.Buy(c, login[i], strategy.Buy[i].Price, strategy.Buy[i].Amount)
				}
				check(err)
				// Store the order in datastore
				err = appdb.KeyPut(c, "Buy_"+exchangeName[i], &strategy.Buy[i], "", time.Seconds())
				check(err)
			}
		} else if strategy.Sell[i].Amount > 0 {
			fmt.Fprintln(w, exchangeName[i], ": Sell", strategy.Sell[i].Amount, "bitcoins for", strategy.Sell[i].Price, "USD per BTC <br>")

			// Execute trade:
			if !paperTrade {
				switch i {
				case mtGox:
					_, err = mtgox.Sell(c, login[i], strategy.Sell[i].Price, strategy.Sell[i].Amount)
				case tradeHill:
					_, err = tradehill.Sell(c, login[i], strategy.Sell[i].Price, strategy.Sell[i].Amount)
					//					case campBx: _, err = campbx.Sell(c, login[i], strategy.Sell[i].Price, strategy.Sell[i].Amount)
				}
				check(err)
				// Store the order in datastore
				err = appdb.KeyPut(c, "Sell_"+exchangeName[i], &strategy.Sell[i], "", time.Seconds())
				check(err)
			}
		} else {
			fmt.Fprintln(w, "No arbirage opportunities at:", exchangeName[i], "<br>")
		}
		fmt.Fprintln(w, exchangeName[i], ": Bid", book[i].BuyTree[0].Price*(1-commission[i]), "Ask",
			book[i].SellTree[0].Price/(1-commission[i]), "<br><br>")
	}
}

func unittests(w http.ResponseWriter, r *http.Request) { // Delete and use 'gotest' instead!
	err := arbitrage.TestCalculate()
	if err != nil {
		fmt.Fprintln(w, err.String())
	} else {
		fmt.Fprintln(w, "arbitrage.Calculate: OK<br>")
	}

	err = arbitrage.TestOnesided()
	if err != nil {
		fmt.Fprintln(w, err.String())
	} else {
		fmt.Fprintln(w, "arbitrage.TestOnesided: OK<br>")
	}
}
