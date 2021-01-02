[Video Setup](https://www.youtube.com/watch?v=bb2kAlzkt20) | [TDAmeritrade Getting Strated](https://developer.tdameritrade.com/content/getting-started)

# Setup Development Environment
Download and install your development environment of choice. I used visual studio code, however, a simple text editor is all that is needed as we are working with GO and it is simple.
https://golang.org/doc/tutorial/getting-started

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


#  Obtain Consumer Key from TDAmeritrade

## Step 1 - Create a TDAmeritrade Developer account 

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
    Response


        PLTR regular market price is 23.55
        30 day average price 26.36, shows stock heading *down*

        GE regular market price is 10.80
        30 day average price 10.79, shows stock heading *up*

        TSLA regular market price is 705.67
        30 day average price 635.72, shows stock heading *up*

<br><br>

# Additional Notes

Helpful Links:

    https://developer.tdameritrade.com/content/getting-started

