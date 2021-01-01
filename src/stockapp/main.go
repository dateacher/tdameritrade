package main

import (
	"flag"
	"fmt"
	"stockapp/pkg/calculations"
	"stockapp/pkg/keyfile"
	"stockapp/pkg/stockdata"
	"strings"
)

func main() {
	//Create variable for storing passed in command line data
	var stockTicker string

	//Create flag
	flag.StringVar(&stockTicker, "stock", "", "Stock symbol")

	//Must parse flags or flags will not work
	flag.Parse()

	//Upercase symbl as this will help later
	stockTicker = strings.ToUpper(stockTicker)

	//Execute get Stock Data function and store data in response
	fullStock, err := stockdata.GetStockData(stockTicker, keyfile.ConsumerKey)
	if err != nil {
		fmt.Println(err)
	}

	//Print stock symbol and regular market price based on the returned struct
	fmt.Printf("%s regular market price is %.2f\n", fullStock.Symbol, fullStock.RegularMarketLastPrice)

	//Determine treding from median
	fmt.Println(calculations.CalcTrendingMedian(fullStock, keyfile.ConsumerKey))
	//Print response to command/ Terminal
}
