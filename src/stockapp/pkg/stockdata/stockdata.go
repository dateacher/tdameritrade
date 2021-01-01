package stockdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func GetStockData(symbl string, consumerKey string) (Stock, error) {
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
	//fmt.Println(body) //Debugging print line

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

//GetPriceHistory get past data
func GetPriceHistory(symbl string, periodType string, period string, frequencyType string, frequency string, consumerKey string) []float32 {
	var candles Candles
	var numbers []float32
	resp, err := http.Get(fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/pricehistory?apikey=%s&periodType=%s&period=%s&frequencyType=%s&frequency=%s", symbl, consumerKey, periodType, period, frequencyType, frequency))
	if err != nil {
		//fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println(err)
	}
	err = json.Unmarshal(body, &candles)
	if err != nil {
		//fmt.Println(err)
	}

	//fmt.Println(candles.Candle)
	for _, item := range candles.Candle {
		numbers = append(numbers, item.Close)
	}
	return numbers
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
