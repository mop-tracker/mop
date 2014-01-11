// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package mop

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

// Stock stores quote information for the particular stock ticker.
// The data for all the fields except 'Advancing' is fetched using
// Yahoo market API.
type YahooStock struct {
	Stock
}

// Quotes stores relevant pointers as well as the array of stock quotes for
// the tickers we are tracking.
type YahooQuotes struct {
	Quotes
}

func NewYahooQuotes(market *Market, profile *Profile) *YahooQuotes {
	return &YahooQuotes{*NewQuotes(market, profile)}
}

// Fetch the latest stock quotes and parse raw fetched data into
// array of []Stock structs.
func (quotes *YahooQuotes) Fetch() {
	if quotes.isReady() {
		defer func() {
			if len(quotes.stocks) == 0 {
				quotes.errors = fmt.Sprint("No stocks returned from server")
			}
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			}
		}()

		url := fmt.Sprintf(quotesURL, strings.Join(quotes.profile.Tickers, `+`))
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		quotes.stocks = quotes.parse(sanitize(body))
	}
	return
}

// Use reflection to parse and assign the quotes data fetched
// using the Yahoo market API.
func (quotes *YahooQuotes) parse(body []byte) []Stock {

	lines := bytes.Split(body, []byte{'\n'})
	stocks := make([]Stock, len(lines))
	// Get the total number of fields in the Stock struct.
	// Skip the last
	// Advancing field which is not fetched.
	fieldsCount := reflect.ValueOf(stocks[0]).NumField() - 1
	// Split each line into columns, then iterate over the Stock
	// struct fields to assign column values.
	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		for j := 0; j < fieldsCount; j++ {
			// ex. quotes.stocks[i].Ticker = string(columns[0])
			reflect.ValueOf(&stocks[i]).Elem().Field(j).SetString(string(columns[j]))

		}
		// Try realtime value and revert to the last known if the
		// realtime is not available.
		if stocks[i].PeRatio == `N/A` && stocks[i].PeRatioX != `N/A` {
			stocks[i].PeRatio = stocks[i].PeRatioX
		}
		if stocks[i].MarketCap == `N/A` && stocks[i].MarketCapX != `N/A` {
			stocks[i].MarketCap = stocks[i].MarketCapX
		}
		// Stock is advancing if the change is not negative
		// (i.e. $0.00 is also "advancing").
		stocks[i].Advancing = (stocks[i].Change[0:1] != `-`)
	}
	return stocks
}
