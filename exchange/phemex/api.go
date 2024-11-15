package phemex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/avinashpandit/crypto-agg/coin"
	"github.com/avinashpandit/crypto-agg/exchange"
	"github.com/avinashpandit/crypto-agg/pair"
)

const (
	API_URL string = "https://api.phemex.com" //"https://www.Phemex.com"
)

/*API Base Knowledge
Path: API function. Usually after the base endpoint URL
Method:
	Get - Call a URL, API return a response
	Post - Call a URL & send a request, API return a response
Public API:
	It doesn't need authorization/signature , can be called by browser to get response.
	using exchange.HttpGetRequest/exchange.HttpPostRequest
Private API:
	Authorization/Signature is requried. The signature request should look at Exchange API Document.
	using ApiKeyGet/ApiKeyPost
Response:
	Response is a json structure.
	Copy the json to https://transform.now.sh/json-to-go/ convert to go Struct.
	Add the go Struct to model.go

ex. Get /api/v1/depth
Get - Method
/api/v1/depth - Path*/

/*************** Public API ***************/
/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Phemex) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := &Data{}

	strRequestUrl := "/public/products"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, currency := range coinsData.Currencies {
		base := &coin.Coin{}
		base = coin.GetCoin(currency.Currency)
		if base == nil {
			base = &coin.Coin{}
			base.Code = currency.Currency
			coin.AddCoin(base)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     currency.Currency,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = currency.Currency
			}
			e.SetCoinConstraint(coinConstraint)
		}

	}

	for _, product := range coinsData.Products {
		if product.Type == "Perpetual" {
			continue
		}

		p := &pair.Pair{}
		currencies := strings.Split(product.DisplaySymbol, "/")
		base := coin.GetCoin(strings.TrimSpace(currencies[0]))
		target := coin.GetCoin(strings.TrimSpace(currencies[1]))
		if base != nil && target != nil {
			p = pair.GetPair(base, target)
		}
		pairConstraint := e.GetPairConstraint(p)
		if pairConstraint == nil {
			pairConstraint = &exchange.PairConstraint{
				PairID:   p.ID,
				Pair:     p,
				ExSymbol: product.Symbol,
				LotSize:  float64(product.PriceScale),
				Listed:   DEFAULT_LISTED,
			}
		}
		e.SetPairConstraint(pairConstraint)
	}

	return nil

}

func (e *Phemex) GetPairsData() error {

	return nil
}

/*
Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct
*/
func (e *Phemex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	pairConstraint := e.GetPairConstraint(pair)

	strRequestUrl := "/md/orderbook"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	//mapParams["depth"] = "10"
	mapParams["symbol"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Result.Book.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate = float64(bid[0]) / math.Pow(10, float64(pairConstraint.LotSize))
		buydata.Quantity = float64(bid[1]) / math.Pow(10, float64(pairConstraint.LotSize))

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Result.Book.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate = float64(ask[0]) / math.Pow(10, float64(pairConstraint.LotSize))
		selldata.Quantity = float64(ask[1]) / math.Pow(10, float64(pairConstraint.LotSize))

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Public API ***************/

/*************** Private API ***************/
func (e *Phemex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Phemex) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	symbol := e.GetSymbolByCoin(operation.Coin)
	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Json Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}

	for key, v := range accountBalance {
		if key != symbol {
			continue
		}

		freeAmount, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
		}
		frozen, err := strconv.ParseFloat(v.Frozen, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
		}

		operation.BalanceAvailable = freeAmount
		operation.BalanceFrozen = frozen
	}

	return nil
}

func (e *Phemex) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
	}
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}

	for key, v := range accountBalance {
		c := e.GetCoinBySymbol(key)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
			}
			frozen, err := strconv.ParseFloat(v.Frozen, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeAmount,
				BalanceFrozen:    frozen,
			}
			operation.BalanceList = append(operation.BalanceList, b)
		}
	}

	return nil
}

func (e *Phemex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} /* else if jsonResponse.Code != 200 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse)
		return
	} */
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for key, v := range accountBalance {
		c := e.GetCoinBySymbol(key)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Available, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
				return
			}
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Phemex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Phemex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/v1/private/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["trade_type"] = "2"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	// log.Printf("mapParams: %+v", mapParams)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprint(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Phemex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/v1/private/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["trade_type"] = "1"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	// log.Printf("mapParams: %+v", mapParams)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprint(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Phemex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/open/api/v1/private/order"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(order.Pair)
	mapParams["trade_no"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.ID == order.OrderID {
		switch orderStatus.Status {
		case "1":
			order.Status = exchange.New
		case "3":
			order.Status = exchange.Partial
		case "5":
			order.Status = exchange.Canceling
		case "4":
			order.Status = exchange.Cancelled
		case "2":
			order.Status = exchange.Filled
		default:
			order.Status = exchange.Other
		}
	}

	return nil
}

func (e *Phemex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Phemex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/spot/orders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderID"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Phemex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Phemex) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]string) string {
	strUrl := API_URL + strRequestPath // + "?" + exchange.Map2UrlQuery(mapParams)

	mapParams["api_key"] = e.API_KEY
	mapParams["req_time"] = fmt.Sprintf("%d", time.Now().Unix()) //"1234567890"
	authParams := exchange.Map2UrlQuery(mapParams)
	authParams += "&api_secret=" + e.API_SECRET

	signature := exchange.ComputeMD5(authParams)
	mapParams["sign"] = signature

	strUrl = strUrl + "?" + exchange.Map2UrlQuery(mapParams)

	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Accept", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
