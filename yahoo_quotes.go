// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

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
//
const quotesURL = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,l1c6k2oghjkva2r2rdyj3j1`

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

type Quotes struct {
	market	    *Market
	profile	    *Profile
	stocks	    []Stock
	errors      string
}

//-----------------------------------------------------------------------------
func (quotes *Quotes) Initialize(market *Market, profile *Profile) *Quotes {
	quotes.market = market
	quotes.profile = profile
	quotes.errors = ``

	return quotes
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
func (quotes *Quotes) Fetch() (this *Quotes) {
	this = quotes // <-- This ensures we return correct quotes after recover() from panic() attack.
	if quotes.isReady() {
		defer func() {
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

		quotes.parse(sanitize(body))
	}

	return quotes
}

//-----------------------------------------------------------------------------
func (quotes *Quotes) Ok() (bool, string) {
	return quotes.errors == ``, quotes.errors
}

//-----------------------------------------------------------------------------
func (quotes *Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = quotes.profile.AddTickers(tickers); err == nil && added > 0 {
		quotes.stocks = nil	// Force fetch.
	}
	return
}

//-----------------------------------------------------------------------------
func (quotes *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.profile.RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.stocks = nil	// Force fetch.
	}
	return
}

// isReady returns true if we haven't fetched the quotes yet *or* the stock
// market is still open and we might want to grab the latest quotes. In both
// cases we make sure the list of requested tickers is not empty.
func (quotes *Quotes) isReady() bool {
	return (quotes.stocks == nil || !quotes.market.IsClosed) && len(quotes.profile.Tickers) > 0
}


//-----------------------------------------------------------------------------
func (quotes *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes.stocks = make([]Stock, len(lines))
	//
	// Get the total number of fields in the Stock struct. Skip the last
	// Advanicing field which is not fetched.
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
		// Stock is advancing if the change is not negative (i.e. $0.00
		// is also "advancing").
		//
		quotes.stocks[i].Advancing = (quotes.stocks[i].Change[0:1] != `-`)
	}

	return quotes
}

//-----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
