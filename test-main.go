package main

import (
	"log"

	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/pair"
	// "../../exchange/kraken"
	// "../../utils"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func test() {

	e := InitEx(exchange.KRAKEN)
	pair := pair.GetPairByKey("USDT|ETH")

	coins := e.GetCoins()
	e.UpdateAllBalances()

	for _, coin := range coins {
		balances := e.GetBalance(coin)
		if balances > 0 {
			log.Printf("Coin Balance %s: %f ", coin.Code, balances)
			Test_AODepositAddress(e, coin)
		}
	}
	//Test_Coins(e)
	//Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.0001, 100)
	// Test_Trading_Sell(e, pair, 0.05, 0.04)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
	//Test_TradeHistory(e, pair)
	//Test_NewOrderBook(e, pair)
}
