// Copyright 2011 Teppo Salonen. All rights reserved.
// This program is distributed under the terms of the MIT/X11 license.

package arbit

import "xgen"

func init() {
	// Login credentials
	login[mtGox] = xgen.Credentials{Username: "<username>", Password: "<password>"}
	login[tradeHill] = xgen.Credentials{Username: "<username>", Password: "<password>"}
	//	login[campBx] = xgen.Credentials{Username: "<username>", Password: "<password>"}
	//	login[bitcoinica] = xgen.Credentials{Username: "<username>", Password: "<password>"}

	// Commissions per trade (volume discounts may be available, so adjust these to reflect your commission level)
	commission[mtGox] = 0.0060           // MtGox commission without volume discounts is currently 0.6% (=0.0060)
	commission[tradeHill] = 0.0060 * 0.9 // TradeHill commissions reduced by 10% if account is created via referral code/link
	//	commission[campBx] = 0.0055 * 0.9 // CampBX commissions reduced by 10% if account is created via referral code/link
	//	commission[bitcoinica] = 0 // Bitcoinica doesn't have a fixed commission - they adjust the spread between buy and sell instead

	/*
		If you don't have referral codes for TradeHill or CampBX but want to get the reduced commissions, feel free to use these:
		TradeHill: TH-R115773
		CampBX: https://CampBX.com/register.php?r=IPxoWpqzIq0
		Disclaimer! If you do use either of these codes, I will also receive 10% of the commissions you generate.
	*/

	// Minimum transaction sizes
	minTrade[mtGox][xgen.USD] = 0.01
	minTrade[mtGox][xgen.BTC] = 0.01
	minTrade[tradeHill][xgen.USD] = 1 // TradeHill's minimum transaction size is $1
	minTrade[tradeHill][xgen.BTC] = 0
	//	minTrade[campBx][xgen.USD] = 0
	//	minTrade[campBx][xgen.BTC] = 0.1 // CampBX's minimum transaction size is 0.1 BTC
	//	minTrade[bitcoinica][xgen.USD] = 0.02
	//	minTrade[bitcoinica][xgen.BTC] = 0.02 // Bitcoinica's minimum transaction size is 0.02 units

	onesidedArb = true
	paperTrade = false // Set true for testing/debugging only
}
