[Video Setup Part 1](https://www.youtube.com/watch?v=bb2kAlzkt20) | [Video watchlist add Part 2](https://youtu.be/sqwO9HTQmQ0) | [TD Ameritrade Getting Strated](https://developer.tdameritrade.com/content/getting-started)

# Setup Development Environment
Download and install your development environment of choice. I used visual studio code, however, a simple text editor is all that is needed as we are working with GO and it is simple.

https://code.visualstudio.com/

https://golang.org/doc/install

Create a simple test to confirm Go is working

    package main

    import (
        "fmt"
    )

    func main(){
        fmt.Println("Hello World")
    }

#  Obtain Consumer Key from TD Ameritrade

## Step 1 - Create a TD Ameritrade Developer account 

* Visit https://developer.tdameritrade.com/
* Register if this is your first time here - free and easy **Be sure to ACCEPT the API terms and conditions** and then follow the link in your email.
* Set your password and follow the link at the bottom of the page for the initial login

## Step 2 - Register App (Necessary to generate Consumer key)

* Select `My Apps`
* Select `+ Add a new App`
* Input required information and select `Create App`
* Copy `Consumer Key` so that we can test. You will be able to get to it anytime so don't worry about saving it *offline*.

Example: 

    App Name - ezmoneyapp
    Callback URL - http://localhost
    What is the purpose of your application? - Swim in a sea of green.

## Step 3 - Test new app via Web URL
Use the 'Consumer Key' on a random stock via the web UI to confirm the app has been created properly.

* Replace `<CONSUMER_KEY>` with your copied or visit the Web Interface to test your newly created key - Example used is for the stock symbol GE - ".../marketdata/`GE`/quotes..."

    Example:

         https://api.tdameritrade.com/v1/marketdata/GE/quotes?apikey=<CONSUMER_KEY>

* If you recieved `Stock Data` proceed, however if you received an `Error` check your key/ stock symbol and don't proceed until you recieve `Stock Data`.

    Stock Data:

        {"GE":{"assetType":"EQUITY","assetMainType":"EQUITY","cusip":"3123103","symbol":"GE","description":"General Electric Company Common 

    |Error|Recommended Fix|
    |-----|-----|
    |{}|Invalid Stock Symbol|
    |"error":"Invalid ApiKey"| Consumer Key needs attention|

# Test TDAmeritrade API using GO

## Step 1 - Clone the repository
    git clone https://github.com/dateacher/tdameritrade.git
## Step 2 - Update Consumer Key

* Update the variable ConsumerKey inside of keyfile.go with the key from the above (Register App step).
        
       var ConsumerKey string = "<Consumer Key>" 
## Step 3 - Run application
* From within your terminal or if you are using it, the visual studio terminal execute the application with any stock symbol.
        
        $ go run .\main.go -stock pltr
        $ go run .\main.go -stock ge
        $ go run .\main.go -stock tsla
        $ go run .\main.go -sector toptickers
        $ go run .\main.go -sector all
    Response

        PLTR regular market price is 23.55
        30 day average price 26.36, shows stock heading *down*

        GE regular market price is 10.80
        30 day average price 10.79, shows stock heading *up*

        TSLA regular market price is 705.67
        30 day average price 635.72, shows stock heading *up*

<br><br>

# Authenticated Account Configuration
## Step 1 - Obtain Account Number
    Account number can be found at in TD AMERITRADE website under `Client Service` > `My Profile` > `Personal Information` > `Account number` - Show

    In thinkorswim desktop software, account number should show in window bar `Account:XXXXXXXX123TDA`, You want just the numbers before TDA.
## Step 2 - Creating Authentication string
This authentication will allow the ability for you to get non delayed quote data, the ability to create watchlists within your account and the ability to configure automatic trading.

Authenitcation guide: https://developer.tdameritrade.com/content/authentication-faq

Quick rundown if Part 1 was followed and you are using localhost as the endpoint.
* Obtain `Auth Code`

    *  Generate Auth URL following the example https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=http%3A%2F%2Flocalhost&client_id=*\<Replace with Consumer Key>*%40AMER.OAUTHAP 
        * Important Note - Consumer key menetioned above is the consumer key of the app you created in part 1. Replace all of `<Replace with Consumer Key>` including `<>` with your key (example `AUGFNDO124`). Then executed the command.

    * Login to `TD Ameritrade Secure Log-in` using the account you would like to manage via the API.
        * As you may have noticed, creating the `App` via the TD Ameritrade Development site used a different set of credentials. These 2 users are independent of one another. You create an app to use the API, you create this `Auth Code` and Tokens to make changes within your TD Ameritrade account directly.
    * Once logged in, select `Allow`, since we used local host we should get a page that looks like an offline, or error page. This is expected. You want to copy the URL because this is were the `code` variable is returned. We specifacly need the information after `code=` from the URL, however copy the full path for review to ensure you get the full `code` variable data.
    * The returned `code` variable is URL encoded, we need to decode it first, then we can use the information in the following command to generate the necessary authentication token. This process is needed for your initially authentication and then every 90 days, this is necessary to link your Developer application with your TD Ameritrade account.
        * I used https://www.urldecoder.org/ decode option for decoding my `code` variable

## Step 3 - Obtaining initial access token pairs
*  Visit https://developer.tdameritrade.com/authentication/apis/post/token-0 and fill out the fields:
    
    grant_type: `authorization_code`    
    refresh_token: *Leave blank*\
    access_type: `offline`
    code: Decoded code from above step\
    client_id: `Consumer Key`\
    redirect_uri: `http://localhost`

    Click `Send`

    * The resonse should be `HTTP/1.1 200 OK`. We now need to copy both the access_token and refresh_token. The current access_token will be good for the next 30 mins, however we need the refresh_token to get a new token whenever needed (ie. every 30 mins).

* Add the `refresh_token` to the variable TDREFRESH and the `access_token` to the variable TDKEY, unless you made a new keyfile, both variables can be found in `../keyfile/keyfile.go` in the repository.

## Step 4 - Execute program to create watchlist
As seen in part 2 of the youtube configuration guide, this simple example will identify the median for a given sector, create an array of the first 10 stock symbols and then creating a watchlist using the first symbol while appending the remaining 9 symbols to the watchlist. Once this watchlist is created you can find the watchlist in your TD Ameritrade accounts thinkorswim personal watchlist locations.
* Example command:
    ```
    $ go run .\main.go -sector oilsenergy -setwatchlist
    ```   

## Step 5 - Get Watchlists
On occasion you may want to list all watchlists associated with your account
* Example command:
    ```
    $ go run .\main.go -getwatchlist
    ```

# Additional Notes

Helpful Links:

    https://golang.org/doc/tutorial/getting-started
    https://developer.tdameritrade.com/content/getting-started

# Disclaimer
I am not a tax, stock, or financial advisor. You use the knowledge laid out in this application at your own discretion. The executer of this application acknowledges that any money, made or lost is subject to all legal liabilities and conditions set in place by the executioners’ state and government laws. 

Please customize this application to make you money, do not expect this application to make money on your behalf.


