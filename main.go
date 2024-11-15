package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"

	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/initial"
	"github.com/avinashpandit/crypto-agg/pair"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	Init()

}

func Init() {
	e := InitExchange(exchange.PHEMEX)

	coins := e.GetCoins()
	pairs := e.GetPairs()
	log.Printf("%+v", pairs)
	for _, coin := range coins {
		pair := e.GetPairBySymbol("s" + coin.Code + "USDT")
		log.Printf("%+v", pair)
		maker, err := e.OrderBook(pair)
		if err != nil {
			log.Printf("OrderBook err: %v", err)
		}

		log.Printf("%+v", maker)
	}

	/*e.UpdateAllBalances()

	for _, coin := range coins {
		balances := e.GetBalance(coin)
		if balances > 0 {
			log.Printf("Coin Balance %s: %f ", coin.Code, balances)
			pair := e.GetPairBySymbol(coin.Code + "USD")
			log.Printf("%+v", pair)
			maker, err := e.OrderBook(pair)
			if err != nil {
				log.Printf("OrderBook err: %v", err)
			}

			log.Printf("%+v", maker)

			//Test_AODepositAddress(e, coin)
		}
	}*/
}

func InitExchange(exName exchange.ExchangeName) exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	//config.API_KEY = "U31c06G61OIMPi8lwMNHHhMKr6k+FhALjK8W4IEDTCbOjykuQUDGAlE6"
	//config.API_SECRET = "N2qzsRvOQ6d7Bmgl5riLT+PWnjR8jqqohg9TBqU+l+1JhFk5AGKU/ZytlVFk2k6bHibQ2SipdN8yAP5FzD6OPw=="
	config.ExName = exName

	inMan := initial.CreateInitManager()
	e := inMan.Init(config)
	log.Printf("Initial [ %v ] ", e.GetName())

	config = nil

	return e
}
