package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/avinashpandit/crypto-agg/cache"
	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/initial"
	"github.com/avinashpandit/crypto-agg/logger"
	"github.com/avinashpandit/crypto-agg/pair"
	"github.com/questdb/go-questdb-client/v3"
)

func main() {
	logger.InitLog()
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
	cache.InitCache()

	var ex []exchange.Exchange = make([]exchange.Exchange, 1)
	//ex[0] = InitExchange(exchange.BYBIT)
	ex[0] = InitExchange(exchange.KRAKEN)

	//ex[2] = InitExchange(exchange.KUCOIN)
	/*ex[2] = InitExchange(exchange.COINBASE)
	 */
	/*
		coins := e.GetCoins()
		pairs := e.GetPairs()
		logger.Info().Msgf("%+v", pairs)
		for _, coin := range coins {
			pair := e.GetPairBySymbol("s" + coin.Code + "USDT")
			maker, err := e.OrderBook(pair)
			if err != nil {
				logger.Info().Msgf("OrderBook err: %v", err)
			}

			if len(maker.Bids) > 0 && len(maker.Asks) > 0 {
				logger.Info().Msgf("%+v %+v %+v", pair, maker.Bids[0], maker.Asks[0])
			}

		}
	*/

	/*	coins := e.GetCoins()
		e.UpdateAllBalances()

		for _, coin := range coins {
			balances := e.GetBalance(coin)
			if balances > 0 {
				logger.Info().Msgf("Coin Balance %s: %f ", coin.Code, balances)
				Test_AODepositAddress(e, coin)
			}
		}
	*/

	baseCcys := []string{
		"SOL",
		"BTC",
		"ETH",
		"BNB",
		"USDT",
		"ADA",
		"XRP",
		"DOGE",
		"DOT",
	}

	pairs := make([]pair.Pair, 0)
	ctx1 := context.TODO()

	questdbClient, err := questdb.LineSenderFromConf(ctx1, "http::addr=localhost:9000;username=admin;password=quest;auto_flush_rows=100;auto_flush_interval=500;")
	if err != nil {
		panic("Failed to create client")
	}

	for _, e := range ex {

		for _, baseCcy := range baseCcys {
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

			if pair1 != nil {
				pairs = append(pairs, *pair1)
			}

		}
		var quoteHandler exchange.QuoteHandler = func(bid exchange.Quote, ask exchange.Quote, p string, e exchange.Exchange) error {
			if cache.SetQuote(string(e.GetName()), p, &bid, &ask) {
				logger.Info().Msgf("Received: %s %v  %v from exchange %s", p, bid, ask, e.GetName())
				timestamp := time.Now()
				err = questdbClient.Table("test_crypto_prices").
					Symbol("symbol", p).
					Float64Column("bid", bid.Rate).
					Float64Column("ask", ask.Rate).
					At(ctx1, timestamp)

				if err != nil {
					panic("Failed to insert data")
				}
			}

			return nil
		}

		e.SubscribeAndProcessQuoteMessage(pairs, quoteHandler)

		/*
			maker, err := e.OrderBook(pair1)
			if err != nil {
				logger.Info().Msgf("OrderBook err: %v", err)
			}

			if maker != nil && len(maker.Bids) > 0 && len(maker.Asks) > 0 {
				logger.Info().Msgf("%+v %+v %+v", pair1, maker.Bids[0], maker.Asks[0])
				}*/
	}

	/*e.UpdateAllBalances()

	for _, coin := range coins {
		balances := e.GetBalance(coin)
		if balances > 0 {
			logger.Info().Msgf("Coin Balance %s: %f ", coin.Code, balances)
			pair := e.GetPairBySymbol(coin.Code + "USD")
			logger.Info().Msgf("%+v", pair)
			maker, err := e.OrderBook(pair)
			if err != nil {
				logger.Info().Msgf("OrderBook err: %v", err)
			}

			logger.Info().Msgf("%+v", maker)

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
	logger.Info().Msgf("Initial [ %v ] ", e.GetName())

	config = nil

	return e
}
