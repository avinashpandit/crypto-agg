package phemex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    uint8           `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

type Data struct {
	Currencies      []Currency       `json:"currencies"`
	Products        []Product        `json:"products"`
	PerpProductsV2  []PerpProductV2  `json:"perpProductsV2"`
	RiskLimits      []RiskLimit      `json:"riskLimits"`
	Leverages       []Leverage       `json:"leverages"`
	RiskLimitsV2    []RiskLimitV2    `json:"riskLimitsV2"`
	LeveragesV2     []LeverageV2     `json:"leveragesV2"`
	LeverageMargins []LeverageMargin `json:"leverageMargins"`
	RatioScale      int              `json:"ratioScale"`
	Md5Checksum     string           `json:"md5Checksum"`
}

type Currency struct {
	Currency        string `json:"currency"`
	Name            string `json:"name"`
	Code            int    `json:"code"`
	ValueScale      int    `json:"valueScale"`
	MinValueEv      int    `json:"minValueEv"`
	MaxValueEv      int    `json:"maxValueEv"`
	Status          string `json:"status"`
	DisplayCurrency string `json:"displayCurrency"`
	InAssetsDisplay int    `json:"inAssetsDisplay"`
	Perpetual       int    `json:"perpetual"`
	StableCoin      int    `json:"stableCoin"`
	AssetsPrecision int    `json:"assetsPrecision"`
}

type Product struct {
	Symbol                   string  `json:"symbol"`
	Code                     int     `json:"code"`
	Type                     string  `json:"type"`
	DisplaySymbol            string  `json:"displaySymbol"`
	IndexSymbol              string  `json:"indexSymbol"`
	MarkSymbol               string  `json:"markSymbol"`
	FundingRateSymbol        string  `json:"fundingRateSymbol"`
	FundingRate8hSymbol      string  `json:"fundingRate8hSymbol"`
	ContractUnderlyingAssets string  `json:"contractUnderlyingAssets"`
	SettleCurrency           string  `json:"settleCurrency"`
	QuoteCurrency            string  `json:"quoteCurrency"`
	LotSize                  int64   `json:"lotSize"`
	TickSize                 float64 `json:"tickSize"`
	PriceScale               int64   `json:"priceScale"`
	RatioScale               int64   `json:"ratioScale"`
	PricePrecision           int64   `json:"pricePrecision"`
	MinPriceEp               int64   `json:"minPriceEp"`
	MaxPriceEp               int64   `json:"maxPriceEp"`
	MaxOrderQty              int64   `json:"maxOrderQty"`
	Status                   string  `json:"status"`
	TipOrderQty              int64   `json:"tipOrderQty"`
	Description              string  `json:"description"`
}

type PerpProductV2 struct {
	Symbol                   string `json:"symbol"`
	Code                     int    `json:"code"`
	Type                     string `json:"type"`
	DisplaySymbol            string `json:"displaySymbol"`
	IndexSymbol              string `json:"indexSymbol"`
	MarkSymbol               string `json:"markSymbol"`
	FundingRateSymbol        string `json:"fundingRateSymbol"`
	FundingRate8hSymbol      string `json:"fundingRate8hSymbol"`
	ContractUnderlyingAssets string `json:"contractUnderlyingAssets"`
	SettleCurrency           string `json:"settleCurrency"`
	QuoteCurrency            string `json:"quoteCurrency"`
	TickSize                 string `json:"tickSize"`
}

type RiskLimit struct {
	Symbol     string          `json:"symbol"`
	Steps      string          `json:"steps"`
	RiskLimits []RiskLimitItem `json:"riskLimits"`
}

type RiskLimitItem struct {
	Limit               int    `json:"limit"`
	InitialMargin       string `json:"initialMargin"`
	InitialMarginEr     int    `json:"initialMarginEr"`
	MaintenanceMargin   string `json:"maintenanceMargin"`
	MaintenanceMarginEr int    `json:"maintenanceMarginEr"`
}

type Leverage struct {
	InitialMargin   string `json:"initialMargin"`
	InitialMarginEr int    `json:"initialMarginEr"`
}

type RiskLimitV2 struct {
	Symbol     string            `json:"symbol"`
	Steps      string            `json:"steps"`
	RiskLimits []RiskLimitV2Item `json:"riskLimits"`
}

type RiskLimitV2Item struct {
	Limit               int    `json:"limit"`
	InitialMarginRr     string `json:"initialMarginRr"`
	MaintenanceMarginRr string `json:"maintenanceMarginRr"`
}

type LeverageV2 struct {
	InitialMarginRr string `json:"initialMarginRr"`
}

type LeverageMargin struct {
	IndexID int                  `json:"index_id"`
	Items   []LeverageMarginItem `json:"items"`
}

type LeverageMarginItem struct {
	NotionalValueRv int `json:"notionalValueRv"`
}

type OrderBook struct {
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
	Result struct {
		Book struct {
			Asks [][2]int64 `json:"asks"`
			Bids [][2]int64 `json:"bids"`
		} `json:"book"`
		Depth     int    `json:"depth"`
		Sequence  int64  `json:"sequence"`
		Symbol    string `json:"symbol"`
		Timestamp int64  `json:"timestamp"`
		Type      string `json:"type"`
	} `json:"result"`
}

type AccountBalances struct {
	Frozen    string `json:"frozen"`
	Available string `json:"available"`
}

// type AutoGenerated struct {
// 	BTC struct {
// 		Frozen    string `json:"frozen"`
// 		Available string `json:"available"`
// 	} `json:"BTC"`
// 	ETH struct {
// 		Frozen    string `json:"frozen"`
// 		Available string `json:"available"`
// 	} `json:"ETH"`
// }

type PlaceOrder struct {
	OrderID string `json:"data"`
}

type OrderStatus struct {
	ID             string `json:"id"`
	Market         string `json:"market"`
	Price          string `json:"price"`
	Status         string `json:"status"`
	TotalQuantity  string `json:"totalQuantity"`
	TradedQuantity string `json:"tradedQuantity"`
	TradedAmount   string `json:"tradedAmount"`
	CreateTime     string `json:"createTime"`
	Type           int    `json:"type"`
}

type TradeHistory []struct {
	TradeTime     string `json:"tradeTime"`
	TradePrice    string `json:"tradePrice"`
	TradeQuantity string `json:"tradeQuantity"`
	TradeType     string `json:"tradeType"`
}
