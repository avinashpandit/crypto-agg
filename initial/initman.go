package initial

import (
	"sync"

	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/exchange/binance"
	"github.com/avinashpandit/crypto-agg/exchange/binancedex"
	"github.com/avinashpandit/crypto-agg/exchange/bitbns"
	"github.com/avinashpandit/crypto-agg/exchange/bitstamp"
	"github.com/avinashpandit/crypto-agg/exchange/bitz"
	"github.com/avinashpandit/crypto-agg/exchange/bybit"
	"github.com/avinashpandit/crypto-agg/exchange/coinbase"
	"github.com/avinashpandit/crypto-agg/exchange/coinex"
	"github.com/avinashpandit/crypto-agg/exchange/ibankdigital"
	"github.com/avinashpandit/crypto-agg/exchange/kraken"
	"github.com/avinashpandit/crypto-agg/exchange/kucoin"
	"github.com/avinashpandit/crypto-agg/exchange/lbank"
	"github.com/avinashpandit/crypto-agg/exchange/mxc"
	"github.com/avinashpandit/crypto-agg/exchange/okex"
	"github.com/avinashpandit/crypto-agg/exchange/okexdm"
	"github.com/avinashpandit/crypto-agg/exchange/phemex"
	"github.com/avinashpandit/crypto-agg/exchange/probit"
)

var instance *InitManager
var once sync.Once

type InitManager struct {
	exMan *exchange.ExchangeManager
}

func CreateInitManager() *InitManager {
	once.Do(func() {
		instance = &InitManager{
			exMan: exchange.CreateExchangeManager(),
		}
	})
	return instance
}

func (e *InitManager) Init(config *exchange.Config) exchange.Exchange {
	switch config.ExName {
	case exchange.BINANCE:
		ex := binance.CreateBinance(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINEX:
		ex := coinex.CreateCoinex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.KUCOIN:
		ex := kucoin.CreateKucoin(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITSTAMP:
		ex := bitstamp.CreateBitstamp(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.OKEX:
		ex := okex.CreateOkex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITZ:
		ex := bitz.CreateBitz(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.MXC:
		ex := mxc.CreateMxc(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.KRAKEN:
		ex := kraken.CreateKraken(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.PHEMEX:
		ex := phemex.CreatePhemex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.IBANKDIGITAL:
		ex := ibankdigital.CreateIbankdigital(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.LBANK:
		ex := lbank.CreateLbank(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BINANCEDEX:
		ex := binancedex.CreateBinanceDex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.OKEXDM:
		ex := okexdm.CreateOkexdm(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BYBIT:
		ex := bybit.CreateBybit(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.PROBIT:
		ex := probit.CreateProbit(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINBASE:
		ex := coinbase.CreateCoinbase(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITBNS:
		ex := bitbns.CreateBitbns(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex
	}
	return nil
}
