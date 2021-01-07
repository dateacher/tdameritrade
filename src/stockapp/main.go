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
			//Run multi stock gathering if you would like multiple
			if len(stocks) > 10 {
				//Set chunk size to 500
				chunkStocks = stockdata.ChunkStockSymbls(stocks, 500)
				for i, item := range chunkStocks {
					fmt.Printf("Chunk %d out of %d\n", i+1, len(chunkStocks))
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
							fmt.Printf("%s:\nMark price: %.2f\n%s\n", item2.Symbol, item2.Mark, trendingMedian)
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
						fmt.Printf("%s:\nMark price: %.2f\n%s\n", item2.Symbol, item2.Mark, trendingMedian)
					}
					multiStockList = append(multiStockList, item2)
				}
			}

		}

	}
}
