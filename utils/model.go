package utils

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/pair"
	cmap "github.com/orcaman/concurrent-map"
)

type CommonData struct {
	Coins []*coin.Coin `json: "coins"`
	Pairs []*pair.Pair `json: "pairs"`
}

type ExchangeData struct {
	CoinConstraint cmap.ConcurrentMap
	PairConstraint cmap.ConcurrentMap
}

type JsonData struct {
	CoinConstraint []*exchange.CoinConstraint `json: "coinconstraint"`
	PairConstraint []*exchange.PairConstraint `json: "pairconstraint"`
}
