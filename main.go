package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/initial"
	"github.com/avinashpandit/crypto-agg/logger"
	"github.com/avinashpandit/crypto-agg/pair"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// create wait group
	wg := sync.WaitGroup{}
	wg.Add(1)

	Init()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		buf := make([]byte, 1<<20)
		for {
			<-sigs
			stacklen := runtime.Stack(buf, true)
			logger.Info().Msgf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		}
	}()

	wg.Wait()
}

func Init() {
	var ex []exchange.Exchange = make([]exchange.Exchange, 2)
	ex[0] = InitExchange(exchange.BYBIT)
	ex[1] = InitExchange(exchange.KRAKEN)
	//ex[2] = InitExchange(exchange.KUCOIN)
	/*ex[2] = InitExchange(exchange.COINBASE)
	 */
	/*
		coins := e.GetCoins()
		pairs := e.GetPairs()
		log.Printf("%+v", pairs)
		for _, coin := range coins {
			pair := e.GetPairBySymbol("s" + coin.Code + "USDT")
			maker, err := e.OrderBook(pair)
			if err != nil {
				log.Printf("OrderBook err: %v", err)
			}

			if len(maker.Bids) > 0 && len(maker.Asks) > 0 {
				log.Printf("%+v %+v %+v", pair, maker.Bids[0], maker.Asks[0])
			}

		}
	*/

	/*	coins := e.GetCoins()
		e.UpdateAllBalances()

		for _, coin := range coins {
			balances := e.GetBalance(coin)
			if balances > 0 {
				log.Printf("Coin Balance %s: %f ", coin.Code, balances)
				Test_AODepositAddress(e, coin)
			}
		}
	*/
	for _, e := range ex {

		baseCcy := "SOL"
		quoteCcy := "USDT"

		var pair1 *pair.Pair
		pair1 = pair.GetPairByKey(baseCcy + "|" + quoteCcy)
		if pair1 == nil {
			pair1 = pair.GetPairByKey(quoteCcy + "|" + baseCcy)
			if pair1 == nil {
				quoteCcy = "USD"
				pair1 = pair.GetPairByKey(quoteCcy + "|" + baseCcy)
			}
		}

		e.SubscribeAndProcessWebsocketMessage([]pair.Pair{*pair1}, func(message string) error {
			fmt.Println("Received:", message)
			return nil
		})

		maker, err := e.OrderBook(pair1)
		if err != nil {
			log.Printf("OrderBook err: %v", err)
		}

		if maker != nil && len(maker.Bids) > 0 && len(maker.Asks) > 0 {
			log.Printf("%+v %+v %+v", pair1, maker.Bids[0], maker.Asks[0])
		}
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
