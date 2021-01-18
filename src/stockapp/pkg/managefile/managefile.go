package managefile

import (
	"io/ioutil"
	"os"
	"strings"
)

func writeDataToText(data string) error {
	//Append the file and if it doesn't exist create it
	file, err := os.OpenFile("stock.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	//fmt.Printf("Wrote %d bytes\n", fileWrite)
	file.Sync()
	return nil
}

//GetSectorData Returns tickers from txt file
func GetSectorData(sector string) ([]string, error) {
	var sectorFileLocation string
	switch sector {
	//only one case statement here but you could have many case statements and many txt files to gather an unlimited number of stock sectors and symbols.
	case "toptickers":
		sectorFileLocation = "pkg/sectors/toptickers.txt"
	case "construction":
		sectorFileLocation = "pkg/sectors/construction.txt"
	case "consumer":
		sectorFileLocation = "pkg/sectors/consumer.txt"
	case "finance":
		sectorFileLocation = "pkg/sectors/finance.txt"
	case "medical":
		sectorFileLocation = "pkg/sectors/medical.txt"
	case "oilsenergy":
		sectorFileLocation = "pkg/sectors/oilsenergy.txt"
	case "all":
		sectorFileLocation = "pkg/sectors/all.txt"
	}

	stocklist, err := getStockList(sectorFileLocation)
	if err != nil {
		return nil, err
	}
	return stocklist, nil
}

func getStockList(file string) ([]string, error) {
	var stockData []string
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	//fmt.Println("\"" + string(data) + "\"")
	stockSplit := strings.Split(string(data), " ")
	for _, item := range stockSplit {
		if item != "" {
			//fmt.Printf("Number: %d\n", i)
			//fmt.Println(item)
			stockData = append(stockData, item)
		}
	}
	//fmt.Println(stockData)

	return stockSplit, nil
}
