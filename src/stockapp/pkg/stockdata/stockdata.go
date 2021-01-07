package stockdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"stockapp/pkg/errorhandlers"
	"time"
)

//Stock struct for stock data
type Stock struct {
	Symbol                 string
	ClosePrice             float32
	OpenPrice              float32
	TotalVolume            float32
	HighPrice              float32
	LowPrice               float32
	RegularMarketLastPrice float32
	Mark                   float32
	Delayed                bool
	Shortable              bool
}

//Candles struct for tracking past data
type Candles struct {
	Candle []Data `json:"candles"`
}

//Data struct for tracking stock data
type Data struct {
	Close           float32 `json:"close"`
	Datetime        int64   `json:"datetime"`
	High            float32 `json:"high"`
	Low             float32 `json:"low"`
	Open            float32 `json:"open"`
	Volume          float32 `json:"volume"`
	MarketLastPrice float32
}

//GetStockData returns a data for a single stock
func GetStockData(symbl string, consumerKey string) (Stock, error) {
	//Delay to not hit rate limit during larger individual stock pulls
	time.Sleep(500 * time.Millisecond)

	//Create HTTP Client need for GET call below
	client := &http.Client{}

	//Create map for Stock struct
	var stockMap map[string]Stock

	//GET creation request for Quote data, variables used for stock symbl and key
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/quotes?apikey=%s", symbl, consumerKey), nil)

	//Send HTTP request to server
	resp, err := client.Do(request)
	if err != nil {
		return Stock{}, err
	}
	//fmt.Println(resp) //#Debugging print line

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Stock{}, err
	}
	//fmt.Println(string(body)) //Debugging print line

	//Confirm Body data
	bodyCheck := errorhandlers.ConfirmBody(string(body))

	if bodyCheck != "" {
		return Stock{}, errors.New(bodyCheck)
	}

	//Unmarshal response into struct to be used later
	err = json.Unmarshal(body, &stockMap)
	if err != nil {
		return Stock{}, err
	}

	//Convert Map to struct
	stockStruct := makeStockStructFromMap(symbl, stockMap)

	//Return new Stock map
	return stockStruct, nil
}

//GetMultipleStocks returns data for more then 1 stock
func GetMultipleStocks(tickers string, consumerKey string) ([]Stock, error) {
	//Delay to not hit rate limit during larger individual stock pulls
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}
	var groupStocks []Stock
	var parsed map[string]Stock
	var newStock Stock

	v := url.Values{}

	v.Add("symbol", tickers)
	//fmt.Println(v.Encode())

	//fmt.Println(tickers)
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/quotes?apikey=%s&%s", consumerKey, v.Encode()), nil)

	//Add this header to enable non delayed quotes
	//request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", configuration.TDKEY))

	resp, err := client.Do(request)
	//fmt.Println(resp) //#Debugging print line

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body)) //Debugging print line

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	for _, item := range parsed {
		newStock.HighPrice = parsed[item.Symbol].HighPrice
		newStock.ClosePrice = parsed[item.Symbol].ClosePrice
		newStock.LowPrice = parsed[item.Symbol].LowPrice
		newStock.Mark = parsed[item.Symbol].Mark
		newStock.OpenPrice = parsed[item.Symbol].OpenPrice
		newStock.TotalVolume = parsed[item.Symbol].TotalVolume
		newStock.Mark = parsed[item.Symbol].Mark
		newStock.Symbol = item.Symbol
		newStock.Delayed = parsed[item.Symbol].Delayed

		//stockStruct := makeStockStructFromMap(item.Symbol, item)
		groupStocks = append(groupStocks, newStock)
	}

	return groupStocks, nil
}

//GetPriceHistory get past data
func GetPriceHistory(symbl string, periodType string, period string, frequencyType string, frequency string, consumerKey string) ([]float32, error) {
	time.Sleep(500 * time.Millisecond)
	var candles Candles
	var numbers []float32
	resp, err := http.Get(fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/pricehistory?apikey=%s&periodType=%s&period=%s&frequencyType=%s&frequency=%s", symbl, consumerKey, periodType, period, frequencyType, frequency))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//Confirm Body data
	bodyCheck := errorhandlers.ConfirmBody(string(body))

	if bodyCheck != "" {
		return nil, errors.New(bodyCheck)
	}
	err = json.Unmarshal(body, &candles)
	if err != nil {
		return nil, err
	}

	//fmt.Println(candles.Candle)
	for _, item := range candles.Candle {
		numbers = append(numbers, item.Close)
	}
	return numbers, nil
}

func makeStockStructFromMap(symbl string, dataMap map[string]Stock) Stock {
	//Stock struct variable for recieving map
	var stockData Stock

	//Mapping all relevant fields
	stockData.ClosePrice = dataMap[symbl].ClosePrice
	stockData.OpenPrice = dataMap[symbl].OpenPrice
	stockData.TotalVolume = dataMap[symbl].TotalVolume
	stockData.HighPrice = dataMap[symbl].HighPrice
	stockData.LowPrice = dataMap[symbl].LowPrice
	stockData.RegularMarketLastPrice = dataMap[symbl].RegularMarketLastPrice
	stockData.Mark = dataMap[symbl].Mark
	stockData.Delayed = dataMap[symbl].Delayed
	stockData.Shortable = dataMap[symbl].Shortable
	stockData.Symbol = dataMap[symbl].Symbol

	//Return stock struct
	return stockData
}

//ChunkStockSymbls Breaks up symbols so that you are not calling more then x(chunksize) at a time
func ChunkStockSymbls(symbols []string, chunkSize int) [][]string {
	var chunk [][]string
	numChunks := len(symbols) / chunkSize
	for i := 0; i < numChunks; i++ {
		chunk = append(chunk, symbols[i*chunkSize:(i+1)*chunkSize])
	}
	chunk = append(chunk, symbols[numChunks*chunkSize:])
	return chunk
}

//MakeTickersToString turns a slice into a string
func MakeTickersToString(tickers []string) string {
	var fullTickers string
	//fmt.Println(tickers)
	for i, item := range tickers {
		if i < len(tickers) {
			fullTickers += fmt.Sprintf(`%s,`, item)
		} else {
			fullTickers += fmt.Sprintf("%s", item)
		}
	}
	//fmt.Println(fullTickers)
	return fullTickers
}
