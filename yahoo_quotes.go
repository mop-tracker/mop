// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
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

const yahoo_quotes_url = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,l1c6k2oghjkva2r2rdyj3j1`

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
	market	      *Market
	profile	      *Profile
	stocks	      []Stock
}

//-----------------------------------------------------------------------------
func (self *Quotes) Initialize(market *Market, profile *Profile) *Quotes {
	self.market = market
	self.profile = profile

	return self
}

// Fetch the latest stock quotes and parse raw fetched data into array of
// []Stock structs.
func (self *Quotes) Fetch() *Quotes {
	if self.Ready() {
		// Format the URL and send the request.
		url := fmt.Sprintf(yahoo_quotes_url, self.profile.ListOfTickers())
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		// Fetch response and get its body.
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		self.parse(self.sanitize(body))
	}

	return self
}

// Return true if we haven't fetched the quotes yet *or* the stock market is
// still open and we might want to grab the latest quotes. In both cases we
// make sure the list of requested tickers is not empty.
func (self *Quotes) Ready() bool {
	return (self.stocks == nil || !self.market.IsClosed) && len(self.profile.Tickers) > 0
}


//-----------------------------------------------------------------------------
func (self *Quotes) Format() string {
	return new(Layout).Initialize().Quotes(self)
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

//-----------------------------------------------------------------------------
func (self *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	self.stocks = make([]Stock, len(lines))

	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		self.stocks[i].Ticker      = string(columns[0])
		self.stocks[i].LastTrade   = string(columns[1])
		self.stocks[i].Change      = string(columns[2])
		self.stocks[i].ChangePct   = string(columns[3])
		self.stocks[i].Open        = string(columns[4])
		self.stocks[i].Low         = string(columns[5])
		self.stocks[i].High        = string(columns[6])
		self.stocks[i].Low52       = string(columns[7])
		self.stocks[i].High52      = string(columns[8])
		self.stocks[i].Volume      = string(columns[9])
		self.stocks[i].AvgVolume   = string(columns[10])
		self.stocks[i].PeRatio     = string(columns[11])
		self.stocks[i].PeRatioX    = string(columns[12])
		self.stocks[i].Dividend    = string(columns[13])
		self.stocks[i].Yield       = string(columns[14])
		self.stocks[i].MarketCap   = string(columns[15])
		self.stocks[i].MarketCapX  = string(columns[16])
		self.stocks[i].Advancing   = self.stocks[i].Change[0:1] != `-`
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *Quotes) sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
