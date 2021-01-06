package main

import (
	"flag"
	"fmt"
	"stockapp/pkg/calculations"
	"stockapp/pkg/keyfile"
	"stockapp/pkg/managefile"
	"stockapp/pkg/stockdata"
	"strings"
)

func main() {
	//Create variable for storing passed in command line data
	var stockTicker string
	var sector string

	//Create flag
	flag.StringVar(&stockTicker, "stock", "", "Stock symbol")
	flag.StringVar(&sector, "sector", "", "Sector/ filename without the .txt extension. Example: top100")

	//Must parse flags or flags will not work
	flag.Parse()

	//Upercase symbl as this will help later
	stockTicker = strings.ToUpper(stockTicker)

	//If individual stock symbol provided execute this block
	if stockTicker != "" {
		//Execute get Stock Data function and store data in response and print if there are no errors
		fullStock, err := stockdata.GetStockData(stockTicker, keyfile.ConsumerKey)
		if err != nil {
			fmt.Printf("Please see the following error: %s\n", err)
			return
		} else {
			fmt.Printf("%s regular market price is %.2f\n", fullStock.Symbol, fullStock.RegularMarketLastPrice)
		}

		//Print stock symbol and regular market price based on the returned struct

		//Determine treding from median and print if there are no errors
		trendingMedian, err := calculations.CalcTrendingMedian(fullStock, keyfile.ConsumerKey)
		if err != nil {
			fmt.Printf("Please see the following error: %s\n", err)
			return
		} else {
			fmt.Println(trendingMedian)
		}
	}

	//If a sector is provided, execute the multi ticker code only if a single stock is not provided
	if sector != "" {
		stockList := []stockdata.Stock{}
		multiStockList := []stockdata.Stock{}
		var chunkStocks [][]string

		stocks, err := managefile.GetSectorData(sector)
		if err != nil {
			fmt.Println(err)
		}
		if len(stocks) > 1 {
			if len(stocks) > 500 {
				//Error when running
				/* &{406  406 HTTP/1.1 1 1 map[Access-Control-Allow-Headers:[origin, x-requested-with, accept, authorization, content-type] Access-Control-Allow-Methods:[GET, PUT, POST, DELETE, OPTIONS, HEAD, PATCH] Access-Control-Allow-Origin:[] Access-Control-Max-Age:[3628800] Cache-Control:[no-cache,no-store,must-revalidate] Connection:[keep-alive] Content-Security-Policy:[frame-ancestors 'self'] Content-Type:[application/json;charset=UTF-8] Date:[Wed, 06 Jan 2021 02:17:17 GMT] Strict-Transport-Security:[max-age=31536000; includeSubDomains max-age=31536000] Vary:[Accept-Encoding] X-Content-Type-Options:[nosniff] X-Frame-Options:[SAMEORIGIN] X-Xss-Protection:[1; mode=block]] 0xc000462060 -1 [chunked] false true map[] 0xc0000d0100 0xc0003f0000}
				{"error":"Symbol size is over the limit"}
				json: cannot unmarshal string into Go value of type stockdata.Stock */

				/* "error":"Individual App's transactions per seconds restriction reached. Please contact us with further questions"*/

				chunkStocks = stockdata.ChunkStockSymbls(stocks, 500)
				for _, item := range chunkStocks {
					fmt.Println(len(chunkStocks))
					stockList, err = stockdata.GetMultipleStocks(stockdata.MakeTickersToString(item), keyfile.ConsumerKey)
					if err != nil {
						fmt.Println(err)
					}
					for _, item2 := range stockList {
						trendingMedian, err := calculations.CalcTrendingMedian(item2, keyfile.ConsumerKey)
						if err != nil {
							fmt.Printf("Please see the following error: %s\n", err)
							return
						} else {
							fmt.Println(trendingMedian)
						}
						multiStockList = append(multiStockList, item2)
					}
				}
			} else {
				stockList, err = stockdata.GetMultipleStocks(stockdata.MakeTickersToString(stocks), keyfile.ConsumerKey)
				if err != nil {
					fmt.Println(err)
				}
				for _, item2 := range stockList {
					trendingMedian, err := calculations.CalcTrendingMedian(item2, keyfile.ConsumerKey)
					if err != nil {
						fmt.Printf("Please see the following error: %s\n", err)
						return
					} else {
						fmt.Println(trendingMedian)
					}
					multiStockList = append(multiStockList, item2)
				}
			}

		}

	}
}
