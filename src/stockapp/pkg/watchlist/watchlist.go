package watchlist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"stockapp/pkg/authentication"
	"stockapp/pkg/keyfile"
	"strings"
	"time"
)

type WatchList struct {
	Name           string           `json:"name"`
	WatchlistID    string           `json:"watchlistId,omitempty"`
	AccountID      string           `json:"accountId,omitempty"`
	Status         string           `json:"status,omitempty"`
	WatchlistItems []WatchlistItems `json:"watchlistItems"`
}

type WatchlistItems struct {
	SequenceID    int        `json:"sequenceId,omitempty"`
	Quantity      float32    `json:"quantity"`
	AveragePrice  float32    `json:"averagePrice"`
	Commission    float32    `json:"commission"`
	PurchasedDate string     `json:"purchasedDate,omitempty"`
	Instrument    Instrument `json:"instrument"`
}

type Instrument struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description,omitempty"`
	AssetType   string `json:"assetType"`
}

//CreateWatchList initializes a Watchlist, if the watchlist named is already created, wachlist ID will be returned
func CreateWatchList(listName string, ticker string) (string, error) {
	client := &http.Client{}
	watchListID := ""
	watchlistData := WatchList{
		Name: listName,
		WatchlistItems: []WatchlistItems{
			{Quantity: 1,
				Commission: 0,
				Instrument: Instrument{
					Symbol:    ticker,
					AssetType: "EQUITY",
				},
			},
		},
	}

	//fmt.Println(watchlistCreate.WatchlistItems[0].Instrument.Symbol)

	watchlistMarshal, err := json.Marshal(watchlistData)
	//fmt.Printf("Watchlist:\n%s\n", watchlistMarshal)

	request, err := http.NewRequest("POST", fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/watchlists", keyfile.ACCTNUM), bytes.NewBuffer([]byte(watchlistMarshal)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", keyfile.TDKEY))

	resp, err := client.Do(request)
	if resp.StatusCode == 401 {
		err := authentication.RefreshToken()
		if err != nil {
			return "", errors.New(fmt.Sprintf("Authentication Token failed for watchlist creation\n%v", err))
		}
		request, err := http.NewRequest("POST", fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/watchlists", keyfile.ACCTNUM), bytes.NewBuffer([]byte(watchlistMarshal)))
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", keyfile.TDKEY))
		resp, err = client.Do(request)
	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n", body)
	if strings.Contains(string(body), "Watchlist Name already in use.") {
		fmt.Printf("Watchlist already exists\n")
		watchList, _, err := GetWatchLists() //List watch lists
		if err != nil {
			return "", err
		}
		for _, item := range watchList {
			if item.Name == listName {
				watchListID = item.WatchlistID
			}
		}
		return watchListID, nil
	}
	watchList, _, err := GetWatchLists() //List watch lists
	if err != nil {
		return "", err
	}
	for _, item := range watchList {
		if item.Name == listName {
			watchListID = item.WatchlistID
		}
	}
	fmt.Printf("Watchlist created with ID: %s\nName: %s\n", watchListID, listName)
	return watchListID, nil
}

func UpdateWatchList(ticker []string, watchlistName string, watchlistID string) error {
	client := &http.Client{}

	for _, item := range ticker {
		time.Sleep(1000 * time.Millisecond)
		watchlistData := WatchList{
			Name:        watchlistName,
			WatchlistID: watchlistID,
			WatchlistItems: []WatchlistItems{
				{SequenceID: 100000},
			},
		}
		watchlistData.WatchlistItems[0].Instrument.Symbol = item
		watchlistData.WatchlistItems[0].Instrument.AssetType = "EQUITY"
		watchlistMarshal, err := json.Marshal(watchlistData)
		if err != nil {
			return err
		}
		request, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/watchlists/%s", keyfile.ACCTNUM, watchlistID), bytes.NewBuffer([]byte(watchlistMarshal)))
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", keyfile.TDKEY))
		resp, err := client.Do(request)
		if resp.StatusCode != 204 {
			fmt.Printf("Watchlist update failed")
			if resp.StatusCode == 429 {
				fmt.Println("Slow down the API calls")
			}
		}
		if err != nil {
			return err
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s added to watchlist\n", item)
	}
	return nil
}

func GetWatchLists() ([]WatchList, string, error) {
	client := &http.Client{}
	var whatchList []WatchList

	time.Sleep(1000 * time.Millisecond)

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/watchlists", keyfile.ACCTNUM), nil)
	if err != nil {
		return nil, "", err
	}
	//Add this header to enable non delayed quotes
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", keyfile.TDKEY))

	resp, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	//If Token expired, refresh token and re-run function
	if resp.StatusCode == 401 {
		err := authentication.RefreshToken()
		if err != nil {
			return nil, "", nil
		}
		return GetWatchLists()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(body, &whatchList)
	if err != nil {
		return nil, "", err
	}
	return whatchList, string(body), nil
}
