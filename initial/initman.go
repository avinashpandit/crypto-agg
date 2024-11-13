package initial

import (
	"sync"

	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/exchange/abcc"
	"github.com/avinashpandit/crypto-agg/exchange/bcex"
	"github.com/avinashpandit/crypto-agg/exchange/bgogo"
	"github.com/avinashpandit/crypto-agg/exchange/bibox"
	"github.com/avinashpandit/crypto-agg/exchange/bigone"
	"github.com/avinashpandit/crypto-agg/exchange/biki"
	"github.com/avinashpandit/crypto-agg/exchange/binance"
	"github.com/avinashpandit/crypto-agg/exchange/binancedex"
	"github.com/avinashpandit/crypto-agg/exchange/bitbay"
	"github.com/avinashpandit/crypto-agg/exchange/bitbns"
	"github.com/avinashpandit/crypto-agg/exchange/bitfinex"
	"github.com/avinashpandit/crypto-agg/exchange/bitforex"
	"github.com/avinashpandit/crypto-agg/exchange/bithumb"
	"github.com/avinashpandit/crypto-agg/exchange/bitmart"
	"github.com/avinashpandit/crypto-agg/exchange/bitmax"
	"github.com/avinashpandit/crypto-agg/exchange/bitmex"
	"github.com/avinashpandit/crypto-agg/exchange/bitpie"
	"github.com/avinashpandit/crypto-agg/exchange/bitrue"
	"github.com/avinashpandit/crypto-agg/exchange/bitstamp"
	"github.com/avinashpandit/crypto-agg/exchange/bittrex"
	"github.com/avinashpandit/crypto-agg/exchange/bitz"
	"github.com/avinashpandit/crypto-agg/exchange/bkex"
	"github.com/avinashpandit/crypto-agg/exchange/blocktrade"
	"github.com/avinashpandit/crypto-agg/exchange/bw"
	"github.com/avinashpandit/crypto-agg/exchange/bybit"
	"github.com/avinashpandit/crypto-agg/exchange/coinbase"
	"github.com/avinashpandit/crypto-agg/exchange/coinbene"
	"github.com/avinashpandit/crypto-agg/exchange/coindeal"
	"github.com/avinashpandit/crypto-agg/exchange/coineal"
	"github.com/avinashpandit/crypto-agg/exchange/coinex"
	"github.com/avinashpandit/crypto-agg/exchange/cointiger"
	"github.com/avinashpandit/crypto-agg/exchange/dcoin"
	"github.com/avinashpandit/crypto-agg/exchange/deribit"
	"github.com/avinashpandit/crypto-agg/exchange/digifinex"
	"github.com/avinashpandit/crypto-agg/exchange/dragonex"
	"github.com/avinashpandit/crypto-agg/exchange/ftx"
	"github.com/avinashpandit/crypto-agg/exchange/gateio"
	"github.com/avinashpandit/crypto-agg/exchange/goko"
	"github.com/avinashpandit/crypto-agg/exchange/hibitex"
	"github.com/avinashpandit/crypto-agg/exchange/hitbtc"
	"github.com/avinashpandit/crypto-agg/exchange/homiex"
	"github.com/avinashpandit/crypto-agg/exchange/hoo"
	"github.com/avinashpandit/crypto-agg/exchange/huobi"
	"github.com/avinashpandit/crypto-agg/exchange/huobidm"
	"github.com/avinashpandit/crypto-agg/exchange/huobiotc"
	"github.com/avinashpandit/crypto-agg/exchange/ibankdigital"
	"github.com/avinashpandit/crypto-agg/exchange/idcm"
	"github.com/avinashpandit/crypto-agg/exchange/idex"
	"github.com/avinashpandit/crypto-agg/exchange/kraken"
	"github.com/avinashpandit/crypto-agg/exchange/kucoin"
	"github.com/avinashpandit/crypto-agg/exchange/latoken"
	"github.com/avinashpandit/crypto-agg/exchange/lbank"
	"github.com/avinashpandit/crypto-agg/exchange/liquid"
	"github.com/avinashpandit/crypto-agg/exchange/mxc"
	"github.com/avinashpandit/crypto-agg/exchange/newcapital"
	"github.com/avinashpandit/crypto-agg/exchange/okex"
	"github.com/avinashpandit/crypto-agg/exchange/okexdm"
	"github.com/avinashpandit/crypto-agg/exchange/oksim"
	"github.com/avinashpandit/crypto-agg/exchange/otcbtc"
	"github.com/avinashpandit/crypto-agg/exchange/phemex"
	"github.com/avinashpandit/crypto-agg/exchange/poloniex"
	"github.com/avinashpandit/crypto-agg/exchange/probit"
	"github.com/avinashpandit/crypto-agg/exchange/stex"
	"github.com/avinashpandit/crypto-agg/exchange/switcheo"
	"github.com/avinashpandit/crypto-agg/exchange/tagz"
	"github.com/avinashpandit/crypto-agg/exchange/tokok"
	"github.com/avinashpandit/crypto-agg/exchange/tradeogre"
	"github.com/avinashpandit/crypto-agg/exchange/txbit"
	"github.com/avinashpandit/crypto-agg/exchange/virgocx"
	"github.com/avinashpandit/crypto-agg/exchange/zebitex"
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

	case exchange.BITTREX:
		ex := bittrex.CreateBittrex(config)
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

	case exchange.STEX:
		ex := stex.CreateStex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITMEX:
		ex := bitmex.CreateBitmex(config)
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

	case exchange.HUOBIOTC:
		ex := huobiotc.CreateHuobiOTC(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITMAX:
		ex := bitmax.CreateBitmax(config)
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

	case exchange.OTCBTC:
		ex := otcbtc.CreateOtcbtc(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.HUOBI:
		ex := huobi.CreateHuobi(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BIBOX:
		ex := bibox.CreateBibox(config)
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

	case exchange.HITBTC:
		ex := hitbtc.CreateHitbtc(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.DRAGONEX:
		ex := dragonex.CreateDragonex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BIGONE:
		ex := bigone.CreateBigone(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITFINEX:
		ex := bitfinex.CreateBitfinex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.GATEIO:
		ex := gateio.CreateGateio(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.IDEX:
		ex := idex.CreateIdex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.LIQUID:
		ex := liquid.CreateLiquid(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITFOREX:
		ex := bitforex.CreateBitforex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.TOKOK:
		ex := tokok.CreateTokok(config)
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

	case exchange.BITRUE:
		ex := bitrue.CreateBitrue(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	// case exchange.TRADESATOSHI:
	// 	ex := tradesatoshi.CreateTradeSatoshi(config)
	// 	if ex != nil {
	// 		e.exMan.Add(ex)
	// 	}
	// 	return ex

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

	case exchange.POLONIEX:
		ex := poloniex.CreatePoloniex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINEAL:
		ex := coineal.CreateCoineal(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.TRADEOGRE:
		ex := tradeogre.CreateTradeogre(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINBENE:
		ex := coinbene.CreateCoinbene(config)
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

	case exchange.BITMART:
		ex := bitmart.CreateBitmart(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BIKI:
		ex := biki.CreateBiki(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.DCOIN:
		ex := dcoin.CreateDcoin(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINTIGER:
		ex := cointiger.CreateCointiger(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.HUOBIDM:
		ex := huobidm.CreateHuobidm(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BW:
		ex := bw.CreateBw(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITBAY:
		ex := bitbay.CreateBitbay(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.DERIBIT:
		ex := deribit.CreateDeribit(config)
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

	case exchange.GOKO:
		ex := goko.CreateGoko(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BCEX:
		ex := bcex.CreateBcex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.DIGIFINEX:
		ex := digifinex.CreateDigifinex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.LATOKEN:
		ex := latoken.CreateLatoken(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.VIRGOCX:
		ex := virgocx.CreateVirgocx(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.ABCC:
		ex := abcc.CreateAbcc(config)
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

	case exchange.ZEBITEX:
		ex := zebitex.CreateZebitex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BITHUMB:
		ex := bithumb.CreateBithumb(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.SWITCHEO:
		ex := switcheo.CreateSwitcheo(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BLOCKTRADE:
		ex := blocktrade.CreateBlocktrade(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BKEX:
		ex := bkex.CreateBkex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.NEWCAPITAL:
		ex := newcapital.CreateNewcapital(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.COINDEAL:
		ex := coindeal.CreateCoindeal(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.HIBITEX:
		ex := hibitex.CreateHibitex(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.BGOGO:
		ex := bgogo.CreateBgogo(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.FTX:
		ex := ftx.CreateFtx(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.TXBIT:
		ex := txbit.CreateTxbit(config)
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

	case exchange.BITPIE:
		ex := bitpie.CreateBitpie(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.TAGZ:
		ex := tagz.CreateTagz(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.IDCM:
		ex := idcm.CreateIdcm(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.HOO:
		ex := hoo.CreateHoo(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	case exchange.HOMIEX:
		ex := homiex.CreateHomiex(config)
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

	case exchange.OKSIM:
		ex := oksim.CreateOksim(config)
		if ex != nil {
			e.exMan.Add(ex)
		}
		return ex

	}
	return nil
}
