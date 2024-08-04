package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"os"
	"time"

	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/conf"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/exchange/coinbase"
	"github.com/avinashpandit/crypto-agg/exchange/kraken"
	"github.com/avinashpandit/crypto-agg/pair"
	"github.com/avinashpandit/crypto-agg/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	exMan := exchange.CreateExchangeManager()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		/* case "export":
		Init(exchange.EXCHANGE_API, "")
		utils.ConvertBaseDataToJson("./data")
		for _, ex := range exMan.GetExchanges() {
			utils.ConvertExchangeDataToJson("./data", ex)
		}
		break */
		case "json":
			Init(exchange.JSON_FILE, "./data")
			for _, ex := range exMan.GetExchanges() {
				for _, coin := range ex.GetCoins() {
					log.Printf("%s Coin %+v", ex.GetName(), coin)
				}
				for _, pair := range ex.GetPairs() {
					log.Printf("%s Pair %+v", ex.GetName(), pair)
				}
			}
			break
		case "renew":
			Init(exchange.JSON_FILE, "./data")
			updateConfig := &exchange.Update{
				ExNames: exMan.GetSupportExchanges(),
				Method:  exchange.TIME_TIGGER,
				Time:    10 * time.Second,
			}
			exMan.UpdateExData(updateConfig)
			break
		case "test":
			base := coin.Coin{
				Code: "BTC",
			}
			target := coin.Coin{
				Code: "ETH",
			}
			pair := pair.Pair{
				Base:   &base,
				Target: &target,
			}
			log.Println(pair)

			// okex.Socket(&pair)
			// stex.Socket()
			// bitfinex.Socket()
		}
	}
}

func Init(source exchange.DataSource, sourceURI string) {
	coin.Init()
	pair.Init()
	if source == exchange.JSON_FILE {
		utils.GetCommonDataFromJSON(sourceURI)
	}
	config := &exchange.Config{}
	config.Source = source
	config.SourceURI = sourceURI

	InitKraken(config)
	InitCoinbase(config)
	// InitBitbns(config)
}

func InitKraken(config *exchange.Config) {
	conf.Exchange(exchange.KRAKEN, config)
	ex := kraken.CreateKraken(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitCoinbase(config *exchange.Config) {
	conf.Exchange(exchange.COINBASE, config)
	ex := coinbase.CreateCoinbase(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}
