// Copyright (c) 2013-2023 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const quotesURL = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`

// const quotesURLv7QueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`
const quotesURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

const noDataIndicator = `N/A`

// Stock stores quote information for the particular stock ticker. The data
// for all the fields except 'Direction' is fetched using Yahoo market API.
type Stock struct {
	Ticker     string `json:"symbol"`                      // Stock ticker.
	LastTrade  string `json:"regularMarketPrice"`          // l1: last trade.
	Change     string `json:"regularMarketChange"`         // c6: change real time.
	ChangePct  string `json:"regularMarketChangePercent"`  // k2: percent change real time.
	Open       string `json:"regularMarketOpen"`           // o: market open price.
	Low        string `json:"regularMarketDayLow"`         // g: day's low.
	High       string `json:"regularMarketDayHigh"`        // h: day's high.
	Low52      string `json:"fiftyTwoWeekLow"`             // j: 52-weeks low.
	High52     string `json:"fiftyTwoWeekHigh"`            // k: 52-weeks high.
	Volume     string `json:"regularMarketVolume"`         // v: volume.
	AvgVolume  string `json:"averageDailyVolume10Day"`     // a2: average volume.
	PeRatio    string `json:"trailingPE"`                  // r2: P/E ration real time.
	PeRatioX   string `json:"trailingPE"`                  // r: P/E ration (fallback when real time is N/A).
	Dividend   string `json:"trailingAnnualDividendRate"`  // d: dividend.
	Yield      string `json:"trailingAnnualDividendYield"` // y: dividend yield.
	MarketCap  string `json:"marketCap"`                   // j3: market cap real time.
	MarketCapX string `json:"marketCap"`                   // j1: market cap (fallback when real time is N/A).
	Currency   string `json:"currency"`                    // String code for currency of stock.
	Direction  int    // -1 when change is < $0, 0 when change is = $0, 1 when change is > $0.
	PreOpen    string `json:"preMarketChangePercent,omitempty"`
	AfterHours string `json:"postMarketChangePercent,omitempty"`
}

// Quotes stores relevant pointers as well as the array of stock quotes for
// the tickers we are tracking.
type Quotes struct {
	market  *Market  // Pointer to Market.
	profile *Profile // Pointer to Profile.
	stocks  []Stock  // Array of stock quote data.
	errors  string   // Error string if any.
}

// Sets the initial values and returns new Quotes struct.
func NewQuotes(market *Market, profile *Profile) *Quotes {
	return &Quotes{
		market:  market,
		profile: profile,
		errors:  ``,
	}
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
func (quotes *Quotes) Fetch() (self *Quotes) {
	self = quotes // <-- This ensures we return correct quotes after recover() from panic().
	if quotes.isReady() {
		defer func() {
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			} else {
				quotes.errors = ""
			}
		}()

		url := fmt.Sprintf(quotesURL, quotes.market.crumb, strings.Join(quotes.profile.Tickers, `,`))

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
		quotes.stocks = nil // Force fetch.
	}
	return
}

// RemoveTickers saves the list of tickers and refreshes the stock data if some
// tickers have been removed. The function gets called from the line editor
// when user removes existing stock tickers.
func (quotes *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.profile.RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.stocks = nil // Force fetch.
	}
	return
}

// isReady returns true if we haven't fetched the quotes yet *or* the stock
// market is still open and we might want to grab the latest quotes. In both
// cases we make sure the list of requested tickers is not empty.
func (quotes *Quotes) isReady() bool {
	return (quotes.stocks == nil || !quotes.market.IsClosed) && len(quotes.profile.Tickers) > 0
}

// this will parse the json objects
func (quotes *Quotes) parse2(body []byte) (*Quotes, error) {
	var d struct {
		QuoteResponse struct {
			Result []struct {
				Symbol                      string `json:"symbol"`
				RegularMarketPrice          string `json:"regularMarketPrice"`
				RegularMarketChange         string `json:"regularMarketChange"`
				RegularMarketChangePercent  string `json:"regularMarketChangePercent"`
				RegularMarketOpen           string `json:"regularMarketOpen"`
				RegularMarketDayLow         string `json:"regularMarketDayLow"`
				RegularMarketDayHigh        string `json:"regularMarketDayHigh"`
				FiftyTwoWeekLow             string `json:"fiftyTwoWeekLow"`
				FiftyTwoWeekHigh            string `json:"fiftyTwoWeekHigh"`
				RegularMarketVolume         string `json:"regularMarketVolume"`
				AverageDailyVolume10Day     string `json:"averageDailyVolume10Day"`
				TrailingPE                  string `json:"trailingPE"`
				TrailingAnnualDividendRate  string `json:"trailingAnnualDividendRate"`
				TrailingAnnualDividendYield string `json:"trailingAnnualDividendYield"`
				MarketCap                   string `json:"marketCap"`
				Currency                    string `json:"currency"`
				PreMarketChangePercent      string `json:"preMarketChangePercent,omitempty"`
				PostMarketChangePercent     string `json:"postMarketChangePercent,omitempty"`
			} `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := json.Unmarshal(body, &d); err != nil {
		return nil, err
	}

	quotes.stocks = make([]Stock, len(d.QuoteResponse.Result))
	for i, stock := range d.QuoteResponse.Result {
		quotes.stocks[i].Ticker = stock.Symbol
		quotes.stocks[i].LastTrade = stock.RegularMarketPrice
		quotes.stocks[i].Change = stock.RegularMarketChange
		quotes.stocks[i].ChangePct = stock.RegularMarketChangePercent
		quotes.stocks[i].Open = stock.RegularMarketOpen
		quotes.stocks[i].Low = stock.RegularMarketDayLow
		quotes.stocks[i].High = stock.RegularMarketDayHigh
		quotes.stocks[i].Low52 = stock.FiftyTwoWeekLow
		quotes.stocks[i].High52 = stock.FiftyTwoWeekHigh
		quotes.stocks[i].Volume = stock.RegularMarketVolume
		quotes.stocks[i].AvgVolume = stock.AverageDailyVolume10Day
		quotes.stocks[i].PeRatio = stock.TrailingPE
		quotes.stocks[i].PeRatioX = stock.TrailingPE
		quotes.stocks[i].Dividend = stock.TrailingAnnualDividendRate
		// The value here is returned in decimal representation but we want to display it as a percentage.
		val, err := strconv.ParseFloat(stock.TrailingAnnualDividendRate, 64)
		if err != nil {
			// I think this might break if the case actually triggers no idea how to do it more robustly.
			quotes.stocks[i].Yield = "N/A"
		} else {
			quotes.stocks[i].Yield = strconv.FormatFloat(val*100, 'f', 2, 64)
		}
		quotes.stocks[i].Yield = stock.TrailingAnnualDividendYield
		quotes.stocks[i].MarketCap = stock.MarketCap
		quotes.stocks[i].MarketCapX = stock.MarketCap
		quotes.stocks[i].Currency = stock.Currency
		quotes.stocks[i].PreOpen = stock.PreMarketChangePercent
		quotes.stocks[i].AfterHours = stock.PostMarketChangePercent
	}

	return quotes, nil
}

// Use reflection to parse and assign the quotes data fetched using the Yahoo
// market API.
func (quotes *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes.stocks = make([]Stock, len(lines))
	//
	// Get the total number of fields in the Stock struct. Skip the last
	// Advancing field which is not fetched.
	//
	fieldsCount := reflect.ValueOf(quotes.stocks[0]).NumField() - 1
	//
	// Split each line into columns, then iterate over the Stock struct
	// fields to assign column values.
	//
	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		for j := 0; j < fieldsCount; j++ {
			// ex. quotes.stocks[i].Ticker = string(columns[0])
			reflect.ValueOf(&quotes.stocks[i]).Elem().Field(j).SetString(string(columns[j]))
		}
		//
		// Try realtime value and revert to the last known if the
		// realtime is not available.
		//
		if quotes.stocks[i].PeRatio == `N/A` && quotes.stocks[i].PeRatioX != `N/A` {
			quotes.stocks[i].PeRatio = quotes.stocks[i].PeRatioX
		}
		if quotes.stocks[i].MarketCap == `N/A` && quotes.stocks[i].MarketCapX != `N/A` {
			quotes.stocks[i].MarketCap = quotes.stocks[i].MarketCapX
		}
		//
		// Get the direction of the stock
		//
		adv, err := strconv.ParseFloat(quotes.stocks[i].Change, 64)
		quotes.stocks[i].Direction = 0
		if err == nil {
			if adv < 0 {
				quotes.stocks[i].Direction = -1
			} else if adv > 0 {
				quotes.stocks[i].Direction = 1
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
