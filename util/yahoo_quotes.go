// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package util

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
	`reflect`
	`strings`
)

// See http://www.gummy-stuff.org/Yahoo-stocks.htm
//
// Also http://query.yahooapis.com/v1/public/yql
// ?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in(%22ALU%22,%22AAPL%22)
// &env=http%3A%2F%2Fstockstables.org%2Falltables.env&format=json'

const quotesURL = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,l1c6k2oghjkva2r2rdyj3j1`

// Stock stores quote information for the particular stock ticker. The data
// for all the fields except 'Advancing' is fetched using Yahoo market API.
type Stock struct {
	Ticker      string  // Stock ticker.
	LastTrade   string  // l1: last trade.
	Change      string  // c6: change real time.
	ChangePct   string  // k2: percent change real time.
	Open        string  // o: market open price.
	Low         string  // g: day's low.
	High        string  // h: day's high.
	Low52       string  // j: 52-weeks low.
	High52      string  // k: 52-weeks high.
	Volume      string  // v: volume.
	AvgVolume   string  // a2: average volume.
	PeRatio     string  // r2: P/E ration real time.
	PeRatioX    string  // r: P/E ration (fallback when real time is N/A).
	Dividend    string  // d: dividend.
	Yield       string  // y: dividend yield.
	MarketCap   string  // j3: market cap real time.
	MarketCapX  string  // j1: market cap (fallback when real time is N/A).
	Advancing   bool    // True when change is >= $0.
}

// Quotes stores relevant pointers as well as the array of stock quotes for
// the tickers we are tracking.
type Quotes struct {
	market	  *Market   // Pointer to Market.
	Quotes_profile	  *Profile  // Pointer to Profile.
	Stocks	   []Stock  // Array of stock quote data.
	errors     string   // Error string if any.
}

// Initialize sets the initial values for the Quotes struct. It returns "self"
// so that the next function call could be chained.
func (quotes *Quotes) Initialize(market *Market, profile *Profile) *Quotes {
	quotes.market = market
	quotes.Quotes_profile = profile
	quotes.errors = ``

	return quotes
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
func (quotes *Quotes) Fetch() (self *Quotes) {
	self = quotes // <-- This ensures we return correct quotes after recover() from panic().
	if quotes.isReady() {
		defer func() {
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			}
		}()

		url := fmt.Sprintf(quotesURL, strings.Join(quotes.Quotes_profile.Tickers, `+`))
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		quotes.parse(sanitize(body))
	}

	return quotes
}

// Ok returns two values: 1) boolean indicating whether the error has occured,
// and 2) the error text itself.
func (quotes *Quotes) Ok() (bool, string) {
	return quotes.errors == ``, quotes.errors
}

// AddTickers saves the list of tickers and refreshes the stock data if new
// tickers have been added. The function gets called from the line editor
// when user adds new stock tickers.
func (quotes *Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = quotes.Quotes_profile.AddTickers(tickers); err == nil && added > 0 {
		quotes.Stocks = nil	// Force fetch.
	}
	return
}

// RemoveTickers saves the list of tickers and refreshes the stock data if some
// tickers have been removed. The function gets called from the line editor
// when user removes existing stock tickers.
func (quotes *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.Quotes_profile.RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.Stocks = nil	// Force fetch.
	}
	return
}

// isReady returns true if we haven't fetched the quotes yet *or* the stock
// market is still open and we might want to grab the latest quotes. In both
// cases we make sure the list of requested tickers is not empty.
func (quotes *Quotes) isReady() bool {
	return (quotes.Stocks == nil || !quotes.market.IsClosed) && len(quotes.Quotes_profile.Tickers) > 0
}

func (quotes *Quotes) GetProfile() *Profile {
	return quotes.Quotes_profile
}

// Use reflection to parse and assign the quotes data fetched using the Yahoo
// market API.
func (quotes *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes.Stocks = make([]Stock, len(lines))
	//
	// Get the total number of fields in the Stock struct. Skip the last
	// Advanicing field which is not fetched.
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
		// Stock is advancing if the change is not negative (i.e. $0.00
		// is also "advancing").
		//
		quotes.Stocks[i].Advancing = (quotes.Stocks[i].Change[0:1] != `-`)
	}

	return quotes
}

//-----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
