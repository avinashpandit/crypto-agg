package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"

	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/conf"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/initial"
	"github.com/avinashpandit/crypto-agg/logger"
	"github.com/avinashpandit/crypto-agg/pair"
	"github.com/avinashpandit/crypto-agg/utils"
	"github.com/davecgh/go-spew/spew"
)

func InitEx(exName exchange.ExchangeName) exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)
	config.API_KEY = "U31c06G61OIMPi8lwMNHHhMKr6k+FhALjK8W4IEDTCbOjykuQUDGAlE6"
	config.API_SECRET = "N2qzsRvOQ6d7Bmgl5riLT+PWnjR8jqqohg9TBqU+l+1JhFk5AGKU/ZytlVFk2k6bHibQ2SipdN8yAP5FzD6OPw=="
	config.ExName = exName

	inMan := initial.CreateInitManager()
	e := inMan.Init(config)
	logger.Info().Msgf("Initial [ %v ] ", e.GetName())

	config = nil

	return e
}

func InitExFromJson(exName exchange.ExchangeName) exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	// config.Source = exchange.EXCHANGE_API
	config.Source = exchange.JSON_FILE
	config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	utils.GetCommonDataFromJSON(config.SourceURI)
	conf.Exchange(exName, config)

	inMan := initial.CreateInitManager()
	e := inMan.Init(config)
	logger.Info().Msgf("Initial [ %v ] from JSON", e.GetName())

	config = nil

	return e
}

/********************Public API********************/
func Test_Coins(e exchange.Exchange) {
	coins := e.GetCoins()
	if len(coins) > 0 {
		for _, coin := range coins {
			logger.Info().Msgf("%s Coin %+v", e.GetName(), coin)
		}
	} else {
		log.Panicf("%s didn't get coins' data.", e.GetName())
	}
}

func Test_Pairs(e exchange.Exchange) {
	pairs := e.GetPairs()
	if len(pairs) > 0 {
		for _, pair := range pairs {
			logger.Info().Msgf("%s Pair %+v", e.GetName(), pair)
		}
	} else {
		log.Panicf("%s didn't get pairs' data.", e.GetName())
	}
}

func Test_Pair(e exchange.Exchange, pair *pair.Pair) {
	logger.Info().Msgf("%s Pair: %+v", e.GetName(), pair)
	logger.Info().Msgf("%s Pair Code: %s", e.GetName(), e.GetSymbolByPair(pair))
	logger.Info().Msgf("%s Coin Codes: %s, %s", e.GetName(), e.GetSymbolByCoin(pair.Base), e.GetSymbolByCoin(pair.Target))
}

func Test_Orderbook(e exchange.Exchange, p *pair.Pair) {
	maker, err := e.OrderBook(p)
	logger.Info().Msgf("%s OrderBook %+v   error:%v", e.GetName(), maker, err)
}

/********************Private API********************/
func Test_Balance(e exchange.Exchange, p *pair.Pair) {
	e.UpdateAllBalances()

	base := e.GetBalance(p.Base)
	target := e.GetBalance(p.Target)
	logger.Info().Msgf("Pair: %12s  Base %s: %f | Target %s: %f", p.Name, p.Base.Code, base, p.Target.Code, target)
}

func Test_Trading(e exchange.Exchange, p *pair.Pair, rate, quantity float64) {
	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		logger.Info().Msgf("%s Limit Buy: %v", e.GetName(), order)

		err = e.OrderStatus(order)
		if err == nil {
			logger.Info().Msgf("%s Order Status: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Order Status Err: %s", e.GetName(), err)
		}

		err = e.CancelOrder(order)
		if err == nil {
			logger.Info().Msgf("%s Cancel Order: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Cancel Err: %s", e.GetName(), err)
		}

		err = e.OrderStatus(order)
		if err == nil {
			logger.Info().Msgf("%s Order Status: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Order Status Err: %s", e.GetName(), err)
		}
	} else {
		logger.Info().Msgf("%s Limit Buy Err: %s", e.GetName(), err)
	}
}

func Test_Trading_Sell(e exchange.Exchange, p *pair.Pair, rate, quantity float64) {
	order, err := e.LimitSell(p, quantity, rate)
	if err == nil {
		logger.Info().Msgf("%s Limit Sell: %+v", e.GetName(), order)

		err = e.OrderStatus(order)
		if err == nil {
			logger.Info().Msgf("%s Order Status: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Order Status Err: %s", e.GetName(), err)
		}

		err = e.CancelOrder(order)
		if err == nil {
			logger.Info().Msgf("%s Cancel Order: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Cancel Err: %s", e.GetName(), err)
		}

		err = e.OrderStatus(order)
		if err == nil {
			logger.Info().Msgf("%s Order Status: %+v", e.GetName(), order)
		} else {
			logger.Info().Msgf("%s Order Status Err: %s", e.GetName(), err)
		}
	} else {
		logger.Info().Msgf("%s Limit Sell Err: %s", e.GetName(), err)
	}
}

// check auth only
func Test_OrderStatus(e exchange.Exchange, p *pair.Pair, orderID string) {
	order := &exchange.Order{
		Pair:      p,
		OrderID:   orderID,
		Rate:      0.001,
		Quantity:  100,
		Direction: exchange.Buy,
		Status:    exchange.New,
	}

	err := e.OrderStatus(order)
	if err == nil {
		logger.Info().Msgf("%s Order Status: %v", e.GetName(), order)
	} else {
		logger.Info().Msgf("%s Order Status Err: %s", e.GetName(), err)
	}
}

func Test_CancelOrder(e exchange.Exchange, p *pair.Pair, orderID string) {
	order := &exchange.Order{
		Pair:      p,
		OrderID:   orderID,
		Rate:      0.001,
		Quantity:  10,
		Direction: exchange.Buy,
		Status:    exchange.New,
	}

	err := e.CancelOrder(order)
	if err == nil {
		logger.Info().Msgf("%s Cancel Order: %v", e.GetName(), order)
	} else {
		logger.Info().Msgf("%s Cancel Order Err: %s", e.GetName(), err)
	}
}

func Test_Withdraw(e exchange.Exchange, c *coin.Coin, amount float64, addr string) {
	if e.Withdraw(c, amount, addr, "") {
		logger.Info().Msgf("%s %s Withdraw Successful!", e.GetName(), c.Code)
	} else {
		logger.Info().Msgf("%s %s Withdraw Failed!", e.GetName(), c.Code)
	}
}

func Test_DoWithdraw(e exchange.Exchange, c *coin.Coin, amount string, addr string, tag string) {
	opWithdraw := &exchange.AccountOperation{
		Type:            exchange.Withdraw,
		Coin:            c,
		WithdrawAmount:  amount,
		WithdrawAddress: addr,
		WithdrawTag:     tag,
		DebugMode:       true,
	}
	err := e.DoAccountOperation(opWithdraw)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	logger.Info().Msgf("WithdrawID: %v, err: %v", opWithdraw.WithdrawID, opWithdraw.Error)
}

func Test_DoTransfer(e exchange.Exchange, c *coin.Coin, amount string, from, to exchange.WalletType) {
	opTransfer := &exchange.AccountOperation{
		Type:                exchange.Transfer,
		Coin:                c,
		TransferAmount:      amount,
		TransferFrom:        from,
		TransferDestination: to,
		DebugMode:           true,
	}
	err := e.DoAccountOperation(opTransfer)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
}

func Test_CheckBalance(e exchange.Exchange, c *coin.Coin, balanceType exchange.WalletType) {
	opBalance := &exchange.AccountOperation{
		Type:      exchange.Balance,
		Coin:      c,
		Wallet:    balanceType,
		DebugMode: true,
	}
	err := e.DoAccountOperation(opBalance)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	logger.Info().Msgf("%v Account available: %v, frozen: %v", opBalance.Coin.Code, opBalance.BalanceAvailable, opBalance.BalanceFrozen)
}

func Test_CheckAllBalance(e exchange.Exchange, balanceType exchange.WalletType) {
	opAllBalance := &exchange.AccountOperation{
		Type:      exchange.BalanceList,
		Wallet:    balanceType,
		DebugMode: true,
	}
	err := e.DoAccountOperation(opAllBalance)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	for i, balance := range opAllBalance.BalanceList {
		if balance.BalanceAvailable+balance.BalanceFrozen == 0 {
			continue
		}
		logger.Info().Msgf("AllAccount balance: %v Coin: %v, avaliable: %v, frozen: %v", i+1, balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opAllBalance.BalanceList) == 0 {
		log.Println("AllAccount 0 balance")
	}
}

func Test_TradeHistory(e exchange.Exchange, pair *pair.Pair) {
	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.TradeHistory,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	for _, trade := range opTradeHistory.TradeHistory {
		logger.Info().Msgf("TradeHistory: %+v", trade)
	}
}

func Test_NewOrderBook(e exchange.Exchange, pair *pair.Pair) {
	opOrderBook := &exchange.PublicOperation{
		Type:      exchange.Orderbook,
		Wallet:    exchange.SpotWallet,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opOrderBook)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}

	logger.Info().Msgf("%s OrderBook %+v   error:%v", e.GetName(), opOrderBook.Maker, opOrderBook.Error)
}

func Test_CoinChainType(e exchange.Exchange, coin *coin.Coin) {
	opCoinChainType := &exchange.PublicOperation{
		Type:      exchange.CoinChainType,
		EX:        e.GetName(),
		Coin:      coin,
		DebugMode: true,
	}

	err := e.LoadPublicData(opCoinChainType)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}

	logger.Info().Msgf("%s %s Chain Type: %s", opCoinChainType.EX, opCoinChainType.Coin.Code, opCoinChainType.CoinChainType)
}

func Test_DoOrderbook(e exchange.Exchange, pair *pair.Pair) {
	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.Orderbook,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	logger.Info().Msgf("%s OrderBook %+v", e.GetName(), opTradeHistory.Maker)
}

func Test_AOOpenOrder(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetOpenOrder,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		for _, o := range op.OpenOrders {
			logger.Info().Msgf("%s OpenOrders: %v %+v", e.GetName(), o.Pair.Name, o)
		}
		if len(op.OpenOrders) == 0 {
			logger.Info().Msgf("%s OpenOrder Response: %v", e.GetName(), op.CallResponce)
		}
	}
}

func Test_AOOrderHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetOrderHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		for _, o := range op.OrderHistory {
			logger.Info().Msgf("%s OrderHistory %+v", e.GetName(), o)
		}
		if len(op.OrderHistory) == 0 {
			logger.Info().Msgf("%s OrderHistory Response: %v", e.GetName(), op.CallResponce)
		}
	}
}

func Test_AODepositAddress(e exchange.Exchange, coin *coin.Coin) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetDepositAddress,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      coin,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		for chain, addr := range op.DepositAddresses {
			logger.Info().Msgf("%s DepositAddresses: %v - %v %+v", e.GetName(), chain, addr.Coin.Code, addr)
		}
	}
}

func Test_AODepositHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetDepositHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		if len(op.DepositHistory) == 0 {
			logger.Info().Msgf("%s DepositHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, his := range op.DepositHistory {
			logger.Info().Msgf("%s DepositHistory: %v %+v", e.GetName(), i, his)
		}
	}
}

func Test_AOWithdrawalHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetWithdrawalHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		if len(op.WithdrawalHistory) == 0 {
			logger.Info().Msgf("%s WithdrawalHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, his := range op.WithdrawalHistory {
			logger.Info().Msgf("%s WithdrawalHistory: %v %+v", e.GetName(), i, his)
		}
	}
}

func Test_AOTransferHistory(e exchange.Exchange) {
	op := &exchange.AccountOperation{
		Type:        exchange.GetTransferHistory,
		Wallet:      exchange.SpotWallet,
		SubUserName: "sub1", // coinex only
		Ex:          e.GetName(),
		DebugMode:   true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		logger.Info().Msgf("%+v", err)
	} else {
		if len(op.TransferInHistory)+len(op.TransferOutHistory) == 0 {
			logger.Info().Msgf("%s TransferInHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, tIn := range op.TransferInHistory {
			logger.Info().Msgf("%s TransferInHistory: %v %+v", e.GetName(), i, tIn)
		}
		for i, tOut := range op.TransferOutHistory {
			logger.Info().Msgf("%s TransferOutHistory: %v %+v", e.GetName(), i, tOut)
		}
	}
}

func Test_TickerPrice(e exchange.Exchange) {
	opTickerPrice := &exchange.PublicOperation{
		Type:      exchange.GetTickerPrice,
		EX:        e.GetName(),
		Wallet:    exchange.SpotWallet,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTickerPrice)
	if err != nil {
		logger.Info().Msgf("%v", err)
		return
	}
	for _, ticker := range opTickerPrice.TickerPrice {
		logger.Info().Msgf("TickerPrice: %v, %v", ticker.Pair.Name, ticker.Price)
	}
}

func SubBalances(e exchange.Exchange, subID string) {
	// Sub Spot AllBalance
	opSubBalance := &exchange.AccountOperation{
		Wallet:       exchange.SpotWallet,
		Type:         exchange.SubBalanceList,
		SubAccountID: subID,
		Ex:           e.GetName(),
		DebugMode:    true,
	}
	err := e.DoAccountOperation(opSubBalance)
	if err != nil {
		logger.Info().Msgf("==%v", err)
		return
	}
	for _, balance := range opSubBalance.BalanceList {
		logger.Info().Msgf("SubBalances balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opSubBalance.BalanceList) == 0 {
		log.Println("SubBalances 0 balance")
	}
	logger.Info().Msgf("SubBalances JSON RESPONSE: %v", opSubBalance.CallResponce)
	logger.Info().Msgf("SubBalances done")
}

func SubAllBalances(e exchange.Exchange) {
	// Sub All Spot AllBalance
	opSubAllBalance := &exchange.AccountOperation{
		Wallet:    exchange.SpotWallet,
		Type:      exchange.SubAllBalanceList,
		Ex:        e.GetName(),
		DebugMode: true,
	}
	err := e.DoAccountOperation(opSubAllBalance)
	if err != nil {
		logger.Info().Msgf("==%v", err)
		return
	}
	for _, balance := range opSubAllBalance.BalanceList {
		logger.Info().Msgf("SubAllBalances balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opSubAllBalance.BalanceList) == 0 {
		log.Println("SubAllBalances 0 balance")
	}
	logger.Info().Msgf("SubAllBalances JSON RESPONSE: %v", opSubAllBalance.CallResponce)
	logger.Info().Msgf("SubAllBalances done")
}

func SubAccountList(e exchange.Exchange) {
	// Sub account list
	opSubAccountList := &exchange.AccountOperation{
		Wallet:    exchange.SpotWallet,
		Type:      exchange.GetSubAccountList,
		Ex:        e.GetName(),
		DebugMode: true,
	}
	err := e.DoAccountOperation(opSubAccountList)
	if err != nil {
		logger.Info().Msgf("==%v", err)
		return
	}
	for _, account := range opSubAccountList.SubAccountList {
		logger.Info().Msgf("AllSubAccount account: %+v", account)
	}
	if len(opSubAccountList.SubAccountList) == 0 {
		log.Println("No Sub Account Info")
	}
	logger.Info().Msgf("SubAccountList JSON RESPONSE: %v", opSubAccountList.CallResponce)
	logger.Info().Msgf("AllSubAccount done")
}

/********************General********************/
func Test_ConstraintFetch(e exchange.Exchange, p *pair.Pair) {
	status := e.GetConstraintFetchMethod(p)
	spew.Dump(status)
}

func Test_Constraint(e exchange.Exchange, p *pair.Pair) {
	baseConstraint := e.GetCoinConstraint(p.Base)
	targerConstraint := e.GetCoinConstraint(p.Target)
	pairConstrinat := e.GetPairConstraint(p)

	logger.Info().Msgf("%s %s Coin Constraint: %+v, %v", e.GetName(), p.Base.Code, baseConstraint, baseConstraint.Coin)
	logger.Info().Msgf("%s %s Coin Constraint: %+v, %v", e.GetName(), p.Target.Code, targerConstraint, targerConstraint.Coin)
	logger.Info().Msgf("%s %s Pair Constraint: %+v", e.GetName(), p.Name, pairConstrinat)
}
