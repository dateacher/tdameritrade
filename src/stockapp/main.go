package main

import (
	"flag"
	"fmt"
	"stockapp/pkg/calculations"
	"stockapp/pkg/keyfile"
	"stockapp/pkg/managefile"
	"stockapp/pkg/prompts"
	"stockapp/pkg/stockdata"
	"stockapp/pkg/watchlist"
	"strings"
	"time"
)

func main() {
	//Create variable for storing passed in command line data
	var stockTicker string
	var sector string
	var getWatchlist bool
	var setWatchlist bool

	//Create flag
	flag.StringVar(&stockTicker, "stock", "", "Stock symbol")
	flag.StringVar(&sector, "sector", "", "Sector/ filename without the .txt extension. Example: top100")
	flag.BoolVar(&getWatchlist, "getwatchlist", false, "Boolean default false  - Triggers return of watchlists associated with account")
	flag.BoolVar(&setWatchlist, "setwatchlist", false, "Boolean default false - Triggers creation of watchlist based on search patterns")

	//Must parse flags or flags will not work
	flag.Parse()

	//Upercase symbl as this will help later
	stockTicker = strings.ToUpper(stockTicker)

	//If not input show menu
	if stockTicker == "" && sector == "" && getWatchlist == false {
		prompts.ShowMenu()
	}
	if getWatchlist == true {
		_, listString, err := watchlist.GetWatchLists()
		if err != nil {
			fmt.Println(err)
			return
		}
		//Print watchlists as an example, though you could do anything with the watchlist struct.
		fmt.Println(listString)
	}
	//If individual stock symbol provided execute this block
	if stockTicker != "" {
		//Execute get Stock Data function and store data in response and print if there are no errors
		fullStock, err := stockdata.GetStockData(stockTicker, keyfile.CONSUMERKEY)
		if err != nil {
			fmt.Printf("Please see the following error: %s\n", err)
			return
		} else {
			fmt.Printf("%s regular market price is %.2f\n", fullStock.Symbol, fullStock.RegularMarketLastPrice)
		}

		//Print stock symbol and regular market price based on the returned struct

		//Determine treding from median and print if there are no errors
		trendingMedian, err := calculations.CalcTrendingMedian(fullStock, keyfile.CONSUMERKEY)
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
					stockList, err = stockdata.GetMultipleStocks(stockdata.MakeTickersToString(item), keyfile.CONSUMERKEY)
					if err != nil {
						fmt.Println(err)
					}
					for _, item2 := range stockList {
						trendingMedian, err := calculations.CalcTrendingMedian(item2, keyfile.CONSUMERKEY)
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
				stockList, err = stockdata.GetMultipleStocks(stockdata.MakeTickersToString(stocks), keyfile.CONSUMERKEY)
				if err != nil {
					fmt.Println(err)
				}
				for _, item2 := range stockList {
					trendingMedian, err := calculations.CalcTrendingMedian(item2, keyfile.CONSUMERKEY)
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

		//If setwatchlist boolean true, create a watchlist based on the first 10 tickers in stockList created above.
		if setWatchlist == true {
			//Create string array for list of ticker symbols to add to watchlist
			var tickerSymbolOnly []string
			//Specify a watchlist name, in this example I am naming it Autolist with the current date/ time
			watchListName := fmt.Sprintf("AutoList-%d-%d-%d-%d", time.Now().Day(), time.Now().Year(), time.Now().Hour(), time.Now().Minute())
			//Creating watchlist, if it already exists, the watchlistid will be returned.
			watchListId, err := watchlist.CreateWatchList(watchListName, stockList[0].Symbol)
			if err != nil {
				fmt.Printf("Please see the following error: %s\n", err)
				return
			}
			fmt.Printf("%s added to watchlist\n", stockList[0].Symbol)
			for i, item := range stockList[1:] {
				//Store the first 10 stock ticker symbols
				if i < 9 {
					tickerSymbolOnly = append(tickerSymbolOnly, item.Symbol)
				}
			}
			err = watchlist.UpdateWatchList(tickerSymbolOnly, watchListName, watchListId)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Watchlist update complete")
		}
		//Update the newly created watchlist
	}
}
