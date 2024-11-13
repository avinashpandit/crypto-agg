package phemex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JSONResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
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
	NeedAddrTag     string `json:"needAddrTag"`
	Status          string `json:"status"`
	DisplayCurrency string `json:"displayCurrency"`
	InAssetsDisplay int    `json:"inAssetsDisplay"`
	Perpetual       string `json:"perpetual"`
	StableCoin      string `json:"stableCoin"`
	AssetsPrecision int    `json:"assetsPrecision"`
}

type Product struct {
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
	ContractSize             int    `json:"contractSize"`
	LotSize                  int    `json:"lotSize"`
	TickSize                 int    `json:"tickSize"`
	PriceScale               int    `json:"priceScale"`
	RatioScale               int    `json:"ratioScale"`
	PricePrecision           int    `json:"pricePrecision"`
	MinPriceEp               int    `json:"minPriceEp"`
	MaxPriceEp               int    `json:"maxPriceEp"`
	MaxOrderQty              int    `json:"maxOrderQty"`
	Description              string `json:"description"`
	Status                   string `json:"status"`
	TipOrderQty              int    `json:"tipOrderQty"`
	ListTime                 int    `json:"listTime"`
	MajorSymbol              bool   `json:"majorSymbol"`
	DefaultLeverage          string `json:"defaultLeverage"`
	FundingInterval          int    `json:"fundingInterval"`
	MaxLeverage              int    `json:"maxLeverage"`
	LeverageMargin           int    `json:"leverageMargin"`
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
	PriceScale               string `json:"priceScale"`
	RatioScale               string `json:"ratioScale"`
	PricePrecision           int    `json:"pricePrecision"`
	BaseCurrency             string `json:"baseCurrency"`
	Description              string `json:"description"`
	Status                   string `json:"status"`
	TipOrderQty              string `json:"tipOrderQty"`
	ListTime                 int    `json:"listTime"`
	MajorSymbol              bool   `json:"majorSymbol"`
	DefaultLeverage          string `json:"defaultLeverage"`
	FundingInterval          int    `json:"fundingInterval"`
	MaxLeverage              int    `json:"maxLeverage"`
	LeverageMargin           int    `json:"leverageMargin"`
	MaxOrderQtyRq            string `json:"maxOrderQtyRq"`
	MaxPriceRp               string `json:"maxPriceRp"`
	MinOrderValueRv          string `json:"minOrderValueRv"`
	MinPriceRp               string `json:"minPriceRp"`
	QtyPrecision             int    `json:"qtyPrecision"`
	QtyStepSize              string `json:"qtyStepSize"`
	TipOrderQtyRq            string `json:"tipOrderQtyRq"`
	MaxOpenPosLeverage       int    `json:"maxOpenPosLeverage"`
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
	Options         []int  `json:"options"`
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
	Options         []int  `json:"options"`
	InitialMarginRr string `json:"initialMarginRr"`
}

type LeverageMargin struct {
	IndexID int                  `json:"index_id"`
	Items   []LeverageMarginItem `json:"items"`
}

type LeverageMarginItem struct {
	NotionalValueRv         int    `json:"notionalValueRv"`
	MaxLeverage             int    `json:"maxLeverage"`
	MaintenanceMarginRateRr string `json:"maintenanceMarginRateRr"`
	MaintenanceAmountRv     string `json:"maintenanceAmountRv"`
}

type OrderBook struct {
	Asks []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"asks"`
	Bids []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"bids"`
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
