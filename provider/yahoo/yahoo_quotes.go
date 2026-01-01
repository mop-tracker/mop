// Copyright (c) 2013-2023 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package yahoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/mop-tracker/mop/provider"
)

const quotesURL = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`

// const quotesURLv7QueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`
const quotesURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

const noDataIndicator = `N/A`

// Profile interface defines the methods required from the profile configuration
type Profile = provider.Profile

// Stock stores quote information for the particular stock ticker. The data
// for all the fields except 'Direction' is fetched using Yahoo market API.
type Stock = provider.Stock

// Quotes stores relevant pointers as well as the array of stock quotes for
// the tickers we are tracking.
type Quotes struct {
	market  *Market // Pointer to Market.
	profile Profile // Pointer to Profile (Interface).
	Stocks  []Stock // Array of stock quote data.
	errors  string  // Error string if any.
}

// Sets the initial values and returns new Quotes struct.
func NewQuotes(market *Market, profile Profile) *Quotes {
	return &Quotes{
		market:  market,
		profile: profile,
		errors:  ``,
	}
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
func (quotes *Quotes) Fetch() provider.Quotes {
	if quotes.isReady() {
		defer func() {
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			} else {
				quotes.errors = ""
			}
		}()

		url := fmt.Sprintf(quotesURL, quotes.market.crumb, strings.Join(quotes.profile.GetTickers(), `,`))

		client := http.Client{}
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}

		request.Header = http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"en-US,en;q=0.5"},
			"Connection":      {"keep-alive"},
			"Content-Type":    {"application/json"},
			"Cookie":          {quotes.market.cookies},
			"Host":            {"query1.finance.yahoo.com"},
			"Origin":          {"https://finance.yahoo.com"},
			"Referer":         {"https://finance.yahoo.com"},
			"Sec-Fetch-Dest":  {"empty"},
			"Sec-Fetch-Mode":  {"cors"},
			"Sec-Fetch-Site":  {"same-site"},
			"TE":              {"trailers"},
			"User-Agent":      {userAgent},
		}

		response, err := client.Do(request)
		// response, err := http.Get(url + quotesURLQueryParts)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		quotes.parse2(body)
	}

	return quotes
}

// Ok returns two values: 1) boolean indicating whether the error has occurred,
// and 2) the error text itself.
func (quotes *Quotes) Ok() (bool, string) {
	return quotes.errors == ``, quotes.errors
}

// AddTickers saves the list of tickers and refreshes the stock data if new
// tickers have been added. The function gets called from the line editor
// when user adds new stock tickers.
func (quotes *Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = quotes.profile.AddTickers(tickers); err == nil && added > 0 {
		quotes.Stocks = nil // Force fetch.
	}
	return
}

// RemoveTickers saves the list of tickers and refreshes the stock data if some
// tickers have been removed. The function gets called from the line editor
// when user removes existing stock tickers.
func (quotes *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.profile.RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.Stocks = nil // Force fetch.
	}
	return
}

// isReady returns true if we haven't fetched the quotes yet *or* the stock
// market is still open and we might want to grab the latest quotes. In both
// cases we make sure the list of requested tickers is not empty.
func (quotes *Quotes) isReady() bool {
	return (quotes.Stocks == nil || !quotes.market.IsClosed()) && len(quotes.profile.GetTickers()) > 0
}

func (q *Quotes) GetStocks() []provider.Stock {
	return q.Stocks
}

func (q *Quotes) RefreshAdvice() int {
	return 1
}

func (q *Quotes) BindOnUpdate(f func()) {
	// Yahoo has no background updates
}

// this will parse the json objects
func (quotes *Quotes) parse2(body []byte) (*Quotes, error) {
	var d struct {
		QuoteResponse struct {
			Result []struct {
				Symbol                      string  `json:"symbol"`
				RegularMarketPrice          float64 `json:"regularMarketPrice"`
				RegularMarketChange         float64 `json:"regularMarketChange"`
				RegularMarketChangePercent  float64 `json:"regularMarketChangePercent"`
				RegularMarketOpen           float64 `json:"regularMarketOpen"`
				RegularMarketDayLow         float64 `json:"regularMarketDayLow"`
				RegularMarketDayHigh        float64 `json:"regularMarketDayHigh"`
				FiftyTwoWeekLow             float64 `json:"fiftyTwoWeekLow"`
				FiftyTwoWeekHigh            float64 `json:"fiftyTwoWeekHigh"`
				RegularMarketVolume         float64 `json:"regularMarketVolume"`
				AverageDailyVolume10Day     float64 `json:"averageDailyVolume10Day"`
				TrailingPE                  float64 `json:"trailingPE"`
				TrailingAnnualDividendRate  float64 `json:"trailingAnnualDividendRate"`
				TrailingAnnualDividendYield float64 `json:"trailingAnnualDividendYield"`
				MarketCap                   float64 `json:"marketCap"`
				Currency                    string  `json:"currency"`
				PreMarketChangePercent      float64 `json:"preMarketChangePercent,omitempty"`
				PostMarketChangePercent     float64 `json:"postMarketChangePercent,omitempty"`
			} `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := json.Unmarshal(body, &d); err != nil {
		// Can't unmarshal the data.
		// Let's try to figure out what went wrong.
		var data interface{}
		json.Unmarshal(body, &data)
		return nil, fmt.Errorf("JSON unmarshal failed: %w\n%+v", err, data)
	}

	quotes.Stocks = make([]Stock, len(d.QuoteResponse.Result))
	for i, stock := range d.QuoteResponse.Result {
		quotes.Stocks[i].Ticker = stock.Symbol
		quotes.Stocks[i].LastTrade = float2Str(stock.RegularMarketPrice)
		quotes.Stocks[i].Change = float2Str(stock.RegularMarketChange)
		quotes.Stocks[i].ChangePct = float2Str(stock.RegularMarketChangePercent)
		quotes.Stocks[i].Open = float2Str(stock.RegularMarketOpen)
		quotes.Stocks[i].Low = float2Str(stock.RegularMarketDayLow)
		quotes.Stocks[i].High = float2Str(stock.RegularMarketDayHigh)
		quotes.Stocks[i].Low52 = float2Str(stock.FiftyTwoWeekLow)
		quotes.Stocks[i].High52 = float2Str(stock.FiftyTwoWeekHigh)
		quotes.Stocks[i].Volume = float2Str(stock.RegularMarketVolume)
		quotes.Stocks[i].AvgVolume = float2Str(stock.AverageDailyVolume10Day)
		quotes.Stocks[i].PeRatio = float2Str(stock.TrailingPE)
		quotes.Stocks[i].PeRatioX = float2Str(stock.TrailingPE)
		quotes.Stocks[i].Dividend = float2Str(stock.TrailingAnnualDividendRate)
		if stock.TrailingAnnualDividendYield != 0 {
			quotes.Stocks[i].Yield = float2Str(stock.TrailingAnnualDividendYield * 100)
		} else {
			quotes.Stocks[i].Yield = noDataIndicator
		}
		quotes.Stocks[i].MarketCap = float2Str(stock.MarketCap)
		quotes.Stocks[i].MarketCapX = float2Str(stock.MarketCap)
		quotes.Stocks[i].Currency = stock.Currency
		quotes.Stocks[i].PreOpen = float2Str(stock.PreMarketChangePercent)
		quotes.Stocks[i].AfterHours = float2Str(stock.PostMarketChangePercent)

		adv := stock.RegularMarketChange
		quotes.Stocks[i].Direction = 0
		if adv < 0.0 {
			quotes.Stocks[i].Direction = -1
		} else if adv > 0.0 {
			quotes.Stocks[i].Direction = 1
		}
	}

	return quotes, nil
}

// Use reflection to parse and assign the quotes data fetched using the Yahoo
// market API.
func (quotes *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes.Stocks = make([]Stock, len(lines))
	//
	// Get the total number of fields in the Stock struct. Skip the last
	// Advancing field which is not fetched.
	//
	fieldsCount := reflect.ValueOf(quotes.Stocks[0]).NumField() - 1
	//
	// Split each line into columns, then iterate over the Stock struct
	// fields to assign column values.
	//
	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		for j := 0; j < fieldsCount; j++ {
			// ex. quotes.Stocks[i].Ticker = string(columns[0])
			reflect.ValueOf(&quotes.Stocks[i]).Elem().Field(j).SetString(string(columns[j]))
		}
		//
		// Try realtime value and revert to the last known if the
		// realtime is not available.
		//
		if quotes.Stocks[i].PeRatio == `N/A` && quotes.Stocks[i].PeRatioX != `N/A` {
			quotes.Stocks[i].PeRatio = quotes.Stocks[i].PeRatioX
		}
		if quotes.Stocks[i].MarketCap == `N/A` && quotes.Stocks[i].MarketCapX != `N/A` {
			quotes.Stocks[i].MarketCap = quotes.Stocks[i].MarketCapX
		}
		//
		// Get the direction of the stock
		//
		adv, err := strconv.ParseFloat(quotes.Stocks[i].Change, 64)
		quotes.Stocks[i].Direction = 0
		if err == nil {
			if adv < 0 {
				quotes.Stocks[i].Direction = -1
			} else if adv > 0 {
				quotes.Stocks[i].Direction = 1
			}
		}
	}

	return quotes
}

// -----------------------------------------------------------------------------
func float2Str(v float64) string {
	unit := ""
	switch {
	case v > 1.0e12:
		v /= 1.0e12
		unit = "T"
	case v > 1.0e9:
		v /= 1.0e9
		unit = "B"
	case v > 1.0e6:
		v /= 1.0e6
		unit = "M"
	case v > 1.0e5:
		v /= 1.0e3
		unit = "K"
	default:
		unit = ""
	}
	// parse
	return fmt.Sprintf("%0.3f%s", v, unit)
}
