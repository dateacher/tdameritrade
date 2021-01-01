package calculations

import (
	"fmt"
	"stockapp/pkg/stockdata"
)

//CalculateMedian allows the median calculation of any slice of float32
func CalculateMedian(nums []float32) float32 {
	var total float32

	//sum all numbers up to return the total devided by len of slice
	for _, item := range nums {
		total += item
	}

	return total / float32(len(nums))
}

//CalcTrendingMedian calculates median price for months worth of data, compare this price to the current price and
func CalcTrendingMedian(stock stockdata.Stock, consumerKey string) string {
	//String variable to return statement
	var fullStatement string

	//Calculate median price
	medianPriceData := CalculateMedian(stockdata.GetPriceHistory(stock.Symbol, "month", "1", "daily", "1", consumerKey))

	//Determine if median is lower the last stock price. This can help us determine if we are higher or lower then the trending median
	if medianPriceData < stock.RegularMarketLastPrice {
		fullStatement += fmt.Sprintf("30 day average price %.2f, shows stock heading *up*\n", medianPriceData)
	} else {
		fullStatement += fmt.Sprintf("30 day average price %.2f, shows stock heading *down*\n", medianPriceData)
	}

	//Returns statement but could return anything
	return fullStatement
}