package mop

import (
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"
	"fmt"
)

// See http://dev.markitondemand.com/
// Ex. http://dev.markitondemand.com/Api/v2/Quote/json?symbol=AAPL

type M map[string]interface{}

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

type MarkitStock struct {
	Ticker                  string  `json:"Symbol"`
	LastTrade               float32 `json:"LastPrice"`
	Change                  float32 `json:"Change"`
	ChangePct               float32 `json:"ChangePercent"`
	Open, Low, High, Volume float32
}

// Quotes stores relevant pointers as well as the array of stock
// quotes for the tickers we are tracking.
type MarkitQuotes struct {
	Quotes //Anonymous type (embedding)
}

func NewMarkitQuotes(market *Market, profile *Profile) *MarkitQuotes {
	return &MarkitQuotes{*NewQuotes(market, profile)}
}

// Fetch the latest stock quotes and parse raw fetched data into
// array of []Stock structs.
func (quotes *MarkitQuotes) Fetch() {
	endpoint := "http://dev.markitondemand.com/Api/v2/Quote/json?symbol="
	if quotes.isReady() {
		defer func() {
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			}
		}()

		//url := fmt.Sprintf(endpoint, strings.Join(quotes.GetProfile().Tickers, "+"))
		tickers := quotes.GetProfile().Tickers
		urls := make([]string, len(tickers))
		for i, symbol := range tickers {
			urls[i] = endpoint + symbol
		}
		results := asyncHttpGets(urls)
		quotes.parse(results)
	}
	return
}

func (quotes *MarkitQuotes) parse(results []*HttpResponse) {
	data := make([]MarkitStock, len(results))
	quotes.stocks = make([]Stock, len(results))

	for i, result := range results {
		byte, err := ioutil.ReadAll(result.response.Body)
		if err != nil {
			panic(err)
		} else {
			err := json.Unmarshal(byte, &data[i])

			quotes.stocks[i].Ticker = result.url[strings.IndexRune(result.url, '=')+1:]
			//quotes.stocks[i].Ticker = fmt.Sprintf("%v", data[i].Ticker)
			quotes.stocks[i].LastTrade = fmt.Sprintf("%v", data[i].LastTrade)
			quotes.stocks[i].Change = fmt.Sprintf("%v", data[i].Change)
			quotes.stocks[i].ChangePct = fmt.Sprintf("%v", data[i].ChangePct)
			quotes.stocks[i].Open = fmt.Sprintf("%v", data[i].Open)
			quotes.stocks[i].Low = fmt.Sprintf("%v", data[i].Low)
			quotes.stocks[i].High = fmt.Sprintf("%v", data[i].High)
			quotes.stocks[i].Volume = fmt.Sprintf("%v", data[i].Volume)

			if err != nil {
				panic(err)
			}

		}
	}
}

func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			//fmt.Printf("%s was fetched\n", r.url)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
			//case <-time.After(50 * time.Millisecond):
			//	fmt.Printf(".")
		}
	}
	return responses
}
