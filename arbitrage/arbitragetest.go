// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

package arbitrage

// TODO: Unit tests need to be written using "gotest" (and file renamed to arbitrage_test.go).

import (
	"xgen"
	//	"testing"
	"os"
	"fmt"
)

type arbTest struct {
	book       []xgen.OrderBook
	funds      []xgen.Balance
	commission []float64
	minTrade   [][xgen.NumCurrencies]float64
	out        Strategy
}

type onesidedTest struct {
	funds      []xgen.Balance
	commission []float64
	in         Strategy
	out        Strategy
}

var arbTests = []arbTest{
	// Test #1: Arbitrage between 2 exchanges with 20% commissions: After commissions we would pay $7.50 to get 1.5 BTC, and sell 1.5 BTC for $8.40.
	arbTest{
		[]xgen.OrderBook{
			{BuyTree: []xgen.Order{{8.0, 1.0}, {7.0, 2.0}}, SellTree: []xgen.Order{{9.0, 10.0}}}, // Price, Amount
			{BuyTree: []xgen.Order{{1.0, 10.0}}, SellTree: []xgen.Order{{4.0, 2.0}}},
		},
		[]xgen.Balance{
			[xgen.NumCurrencies]float64{1.5, 10.0}, // BTC, USD
			[xgen.NumCurrencies]float64{10.0, 10.0},
		},
		[]float64{0.2, 0.2},                           // commissions
		[][xgen.NumCurrencies]float64{{0, 0}, {0, 0}}, // minimum allowed trade amounts
		Strategy{Buy: []xgen.Order{{0, 0}, {4.0, 1.875}}, Sell: []xgen.Order{{7.0, 1.5}, {0, 0}}},
	},
	// Test #2: 3-way arbitrage with no commissions (for simplicity).
	arbTest{
		[]xgen.OrderBook{
			{BuyTree: []xgen.Order{{6.5, 1.0}, {5.5, 2.0}, {4.5, 4.0}}, SellTree: []xgen.Order{{10.0, 100.0}}}, // Price, Amount
			{BuyTree: []xgen.Order{{1.0, 100.0}}, SellTree: []xgen.Order{{4.5, 1.0}, {5.5, 2.0}, {6.5, 4.0}}},
			{BuyTree: []xgen.Order{{4.9, 2.0}, {4.8, 4.0}, {4.7, 8.0}}, SellTree: []xgen.Order{{5.0, 2.0}, {5.1, 4.0}, {5.2, 8.0}}},
		},
		[]xgen.Balance{
			[xgen.NumCurrencies]float64{2.0, 0.0}, // BTC, USD
			[xgen.NumCurrencies]float64{0.0, 20.0},
			[xgen.NumCurrencies]float64{0.5, 10.0},
		},
		[]float64{0.0, 0.0, 0.0},                              // commissions
		[][xgen.NumCurrencies]float64{{0, 0}, {0, 0}, {0, 0}}, // minimum allowed trade amounts
		Strategy{Buy: []xgen.Order{{0, 0}, {4.5, 1.0}, {5.0, 1.0}}, Sell: []xgen.Order{{5.5, 2.0}, {0, 0}, {0, 0}}},
	},
	// Test #3: Same as test #2 but adjusted account balances enough to change to output
	arbTest{
		[]xgen.OrderBook{
			{BuyTree: []xgen.Order{{6.5, 1.0}, {5.5, 2.0}, {4.5, 4.0}}, SellTree: []xgen.Order{{10.0, 100.0}}}, // Price, Amount
			{BuyTree: []xgen.Order{{1.0, 100.0}}, SellTree: []xgen.Order{{4.5, 1.0}, {5.5, 2.0}, {6.5, 4.0}}},
			{BuyTree: []xgen.Order{{4.9, 2.0}, {4.8, 4.0}, {4.7, 8.0}}, SellTree: []xgen.Order{{5.0, 2.0}, {5.1, 4.0}, {5.2, 8.0}}},
		},
		[]xgen.Balance{
			[xgen.NumCurrencies]float64{2.0, 0.0}, // BTC, USD
			[xgen.NumCurrencies]float64{0.0, 20.0},
			[xgen.NumCurrencies]float64{0.5, 2.5},
		},
		[]float64{0.0, 0.0, 0.0},                              // commissions
		[][xgen.NumCurrencies]float64{{0, 0}, {0, 0}, {0, 0}}, // minimum allowed trade amounts
		Strategy{Buy: []xgen.Order{{0, 0}, {4.5, 1.0}, {5.0, 0.5}}, Sell: []xgen.Order{{5.5, 1.5}, {0, 0}, {0, 0}}},
	},
}

var onesidedTests = []onesidedTest{
	onesidedTest{
		[]xgen.Balance{
			[xgen.NumCurrencies]float64{2.0, 0.0}, // BTC, USD
			[xgen.NumCurrencies]float64{1.0, 20.0},
			[xgen.NumCurrencies]float64{1.0, 2.0},
		},
		[]float64{0.4, 0.4, 0.4}, // commissions
		Strategy{Buy: []xgen.Order{{0, 0}, {4.0, 1.0}, {4.0, 0.5}}, Sell: []xgen.Order{{6.0, 0.50}, {0, 0}, {0, 0}}},
		Strategy{Buy: []xgen.Order{{0, 0}, {4.0, 2.5}, {4.0, 0.5}}, Sell: []xgen.Order{{6.0, 1.25}, {0, 0}, {0, 0}}},
	},
}

func TestCalculate( /*t *testing.T*/ ) os.Error {
	for i, at := range arbTests {
		v := Calculate(at.book, at.funds, at.commission, at.minTrade)
		if fmt.Sprint(v) != fmt.Sprint(at.out) {
			//t.Errorf("arbitrageStrategy = %d, want %d.", v, at.out)
			return os.NewError(fmt.Sprint("ArbitrageStrategy (#", (i + 1), ")<br>", at.book, "<br>", at.funds, "<br>=<br>", v, "<br>want<br>", at.out))
		}
	}
	return nil
}

func TestOnesided( /*t *testing.T*/ ) os.Error {
	for i, ot := range onesidedTests {
		v := Onesided(ot.in, ot.funds, ot.commission)
		if fmt.Sprint(v) != fmt.Sprint(ot.out) {
			//t.Errorf("onesidedArbitrage = %d, want %d.", v, ot.out)
			return os.NewError(fmt.Sprint("OnesidedArbitrage (#", (i + 1), ")<br>", ot.in, "<br>=<br>", v, "<br>want<br>", ot.out))
		}
	}
	return nil
}
