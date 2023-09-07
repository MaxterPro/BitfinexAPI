package bitfenix

import (
	"os"
	//	"strconv"
	"testing"
)

var APIKey = os.Getenv("BITFINEX_API_KEY")
var APISecret = os.Getenv("BITFINEX_API_SECRET")

var apiPublic = New("", "")
var apiPrivate = New(APIKey, APISecret)

func checkEnv(t *testing.T) {
	if APIKey == "" || APISecret == "" {
		t.Skip("Skipping test because because APIKey and/or APISecret env variables are not set")
	}
}

func TestTicker(t *testing.T) {
	// Test normal request
	ticker, err := apiPublic.Ticker("btcusd")
	if err != nil || ticker.LastPrice == 0 {
		t.Error("Failed: " + err.Error())
		return
	}

	// Test bad request,
	// which must return an error
	ticker, err = apiPublic.Ticker("random")
	if err == nil {
		t.Error("Failed")
		return
	}
}
