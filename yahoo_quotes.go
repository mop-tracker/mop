// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

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
// Also http://query.yahooapis.com/v1/public/yql
// ?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in(%22ALU%22,%22AAPL%22)
// &env=http%3A%2F%2Fstockstables.org%2Falltables.env
// &format=json'
//
// Current, Change, Open, High, Low, 52-W High, 52-W Low, Volume, AvgVolume, P/E, Yield, Market Cap.
// l1: last trade
// c6: change rt
// k2: change % rt
// o: open
// g: day's low
// h: day's high
// j: 52w low
// k: 52w high
// v: volume
// a2: avg volume
// r2: p/e rt
// r: p/e
// d: dividend/share
// y: wield
// j3: market cap rt
// j1: market cap

const quotes_url = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,l1c6k2oghjkva2r2rdyj3j1`

type Stock struct {
	Ticker      string
	LastTrade   string
	Change      string
	ChangePct   string
	Open        string
	Low         string
	High        string
	Low52       string
	High52      string
	Volume      string
	AvgVolume   string
	PeRatio     string
	PeRatioX    string
	Dividend    string
	Yield       string
	MarketCap   string
	MarketCapX  string
	Advancing   bool
}

type Quotes struct {
	market	    *Market
	profile	    *Profile
	stocks	    []Stock
	errors      string
}

//-----------------------------------------------------------------------------
func (self *Quotes) Initialize(market *Market, profile *Profile) *Quotes {
	self.market = market
	self.profile = profile
	self.errors = ``

	return self
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
//-----------------------------------------------------------------------------
func (self *Quotes) Fetch() (this *Quotes) {
	this = self // <-- This ensures we return correct self after recover() from panic() attack.
	if self.is_ready() {
		defer func() {
			if err := recover(); err != nil {
				self.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			}
		}()

		url := fmt.Sprintf(quotes_url, strings.Join(self.profile.Tickers, `+`))
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		self.parse(sanitize(body))
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *Quotes) Ok() (bool, string) {
	return self.errors == ``, self.errors
}

//-----------------------------------------------------------------------------
func (self *Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = self.profile.AddTickers(tickers); err == nil && added > 0 {
		self.stocks = nil	// Force fetch.
	}
	return
}

//-----------------------------------------------------------------------------
func (self *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = self.profile.RemoveTickers(tickers); err == nil && removed > 0 {
		self.stocks = nil	// Force fetch.
	}
	return
}

// "Private" methods.

// Return true if we haven't fetched the quotes yet *or* the stock market is
// still open and we might want to grab the latest quotes. In both cases we
// make sure the list of requested tickers is not empty.
//-----------------------------------------------------------------------------
func (self *Quotes) is_ready() bool {
	return (self.stocks == nil || !self.market.IsClosed) && len(self.profile.Tickers) > 0
}


//-----------------------------------------------------------------------------
func (self *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	self.stocks = make([]Stock, len(lines))
	//
	// Get the total number of fields in the Stock struct. Skip the last
	// Advanicing field which is not fetched.
	//
	number_of_fields := reflect.ValueOf(self.stocks[0]).NumField() - 1
	//
	// Split each line into columns, then iterate over the Stock struct
	// fields to assign column values.
	//
	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		for j := 0; j < number_of_fields; j++ {
			// ex. self.stocks[i].Ticker = string(columns[0])
			reflect.ValueOf(&self.stocks[i]).Elem().Field(j).SetString(string(columns[j]))
		}
		//
		// Try realtime value and revert to the last known if the
		// realtime is not available.
		//
		if self.stocks[i].PeRatio == `N/A` && self.stocks[i].PeRatioX != `N/A` {
			self.stocks[i].PeRatio = self.stocks[i].PeRatioX
		}
		if self.stocks[i].MarketCap == `N/A` && self.stocks[i].MarketCapX != `N/A` {
			self.stocks[i].MarketCap = self.stocks[i].MarketCapX
		}
		//
		// Stock is advancing if the change is not negative (i.e. $0.00
		// is also "advancing").
		//
		self.stocks[i].Advancing = (self.stocks[i].Change[0:1] != `-`)
	}

	return self
}

// Utility methods.

//-----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
