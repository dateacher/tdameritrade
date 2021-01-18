package authentication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"stockapp/pkg/keyfile"
)

type token struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	TokenType          string `json:"token_type"`
	ExpireTime         int    `json:"expires_in"`
	Scope              string `json:"scope"`
	RefreshTokenExpire int    `json:"refresh_token_expires_in"`
}

func RefreshToken() error {
	//Assign a token variable to unmarshal too
	var authkey token
	//Create http client object
	client := &http.Client{}
	//Assign a url.Values object to start building on
	data := url.Values{}
	//Assign struct data to be sent in HTTP POST
	data.Set("grant_type", "refresh_token")
	data.Add("refresh_token", fmt.Sprintf("%s", keyfile.TDREFRESH))
	data.Add("client_id", fmt.Sprintf("%s", keyfile.CONSUMERKEY))
	//Create HTTP Request
	request, err := http.NewRequest("POST", "https://api.tdameritrade.com/v1/oauth2/token?", bytes.NewBufferString(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//Run request
	resp, err := client.Do(request)
	//Handle request
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//Turn body json response into created variable of struct type token
	err = json.Unmarshal(body, &authkey)
	if err != nil {
		return err
	}
	//Notify user Token updated
	fmt.Println("*Token needed a refresh*")
	//Print new access token (if necessary can toggle so you can manually update the keyfile for other purposes etc. Commenting out as we are not concerend about using the key manually later.)
	//Update configuration key with new key
	keyfile.TDKEY = string(authkey.AccessToken)
	return nil
}
