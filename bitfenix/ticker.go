package bitfenix

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// APIURL points to Bitfinex API URL, found at https://www.bitfinex.com/pages/API
	APIURL = "https://api.bitfinex.com"
	// LEND ...
	LEND = "lend"
	// BORROW ...
	BORROW = "borrow"
)

// API structure stores Bitfinex API credentials
type API struct {
	APIKey    string
	APISecret string
}

// ErrorMessage ...
type ErrorMessage struct {
	Message string `json:"message"` // Returned only on error
}

// Ticker ...
type Ticker struct {
	Mid       float64 `json:"mid,string"`        // mid (price): (bid + ask) / 2
	Bid       float64 `json:"bid,string"`        // bid (price): Innermost bid.
	Ask       float64 `json:"ask,string"`        // ask (price): Innermost ask.
	LastPrice float64 `json:"last_price,string"` // last_price (price) The price at which the last order executed.
	Low       float64 `json:"low,string"`        // low (price): Lowest trade price of the last 24 hours
	High      float64 `json:"high,string"`       // high (price): Highest trade price of the last 24 hours
	Volume    float64 `json:"volume,string"`     // volume (price): Trading volume of the last 24 hours
	Timestamp float64 `json:"timestamp,string"`  // timestamp (time) The timestamp at which this information was valid.
}

func New(key, secret string) (api *API) {
	api = &API{
		APIKey:    key,
		APISecret: secret,
	}
	return api
}

///////////////////////////////////////
// Main API methods
///////////////////////////////////////

// Ticker returns innermost bid and asks and information on the most recent trade,
//
//	as well as high, low and volume of the last 24 hours.
func (api *API) Ticker(symbol string) (ticker Ticker, err error) {
	symbol = strings.ToLower(symbol)

	body, err := api.get("/v2/pubticker/" + symbol)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ticker)
	if err != nil || ticker.LastPrice == 0 { // Failed to unmarshal expected message
		// Attempt to unmarshal the error message
		errorMessage := ErrorMessage{}
		err = json.Unmarshal(body, &errorMessage)
		if err != nil { // Not expected message and not expected error, bailing...
			return
		}

		return ticker, errors.New("API: " + errorMessage.Message)
	}

	return
}

func (api *API) get(url string) (body []byte, err error) {
	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	resp, err := client.Get(APIURL + url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
