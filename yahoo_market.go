// Copyright (c) 2013-2019 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

const marketURL = `https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s`
const marketURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

// Market stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type Market struct {
	IsClosed  bool              // True when U.S. markets are closed.
	Dow       map[string]string // Hash of Dow Jones indicators.
	Nasdaq    map[string]string // Hash of NASDAQ indicators.
	Sp500     map[string]string // Hash of S&P 500 indicators.
	Tokyo     map[string]string
	HongKong  map[string]string
	London    map[string]string
	Frankfurt map[string]string
	Yield     map[string]string
	Oil       map[string]string
	Yen       map[string]string
	Euro      map[string]string
	Gold      map[string]string
	errors    string // Error(s), if any.
	url       string // URL with symbols to fetch data
}

// Returns new initialized Market struct.
func NewMarket() *Market {
	market := &Market{}
	market.IsClosed = false
	market.Dow = make(map[string]string)
	market.Nasdaq = make(map[string]string)
	market.Sp500 = make(map[string]string)

	market.Tokyo = make(map[string]string)
	market.HongKong = make(map[string]string)
	market.London = make(map[string]string)
	market.Frankfurt = make(map[string]string)

	market.Yield = make(map[string]string)
	market.Oil = make(map[string]string)
	market.Yen = make(map[string]string)
	market.Euro = make(map[string]string)
	market.Gold = make(map[string]string)

	market.url = fmt.Sprintf(marketURL, `^DJI,^IXIC,^GSPC,^N225,^HSI,^FTSE,^GDAXI,^TNX,CL=F,JPY=X,EUR=X,GC=F`) + marketURLQueryParts

	market.errors = ``

	return market
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) Fetch() (self *Market) {
	self = market // <-- This ensures we return correct market after recover() from panic().
	defer func() {
		if err := recover(); err != nil {
			market.errors = fmt.Sprintf("Error fetching market data...\n%s", err)
		} else {
			market.errors = ""
		}
	}()

	response, err := http.Get(market.url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = market.isMarketOpen(body)
	return market.extract(body)
}

// Ok returns two values: 1) boolean indicating whether the error has occurred,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

// -----------------------------------------------------------------------------
func (market *Market) isMarketOpen(body []byte) []byte {
	// TBD -- CNN page doesn't seem to have market open/close indicator.
	return body
}

// -----------------------------------------------------------------------------
func assign(results []map[string]interface{}, position int, changeAsPercent bool) map[string]string {
    out := make(map[string]string)
	out[`change`] = float2Str(results[position]["regularMarketChange"].(float64))
        out[`latest`] = float2Str(results[position]["regularMarketPrice"].(float64))
    if changeAsPercent{
        out[`change`] = float2Str(results[position]["regularMarketChangePercent"].(float64)) + `%`
    } else { 
        out[`percent`] = float2Str(results[position]["regularMarketChangePercent"].(float64))
    }
    return out
}

// -----------------------------------------------------------------------------
func (market *Market) extract(body []byte) *Market {
    d := map[string]map[string][]map[string]interface{}{}
    err := json.Unmarshal(body, &d)
    if err != nil {
        panic(err)
    }
    results := d["quoteResponse"]["result"]
    market.Dow = assign(results, 0, false)
    market.Nasdaq = assign(results, 1, false)
    market.Sp500 = assign(results, 2, false)
    market.Tokyo = assign(results, 3, false)
    market.HongKong = assign(results, 4, false)
    market.London = assign(results, 5, false)
    market.Frankfurt = assign(results, 6, false)
    market.Yield[`name`] = `10-year Yield`
    market.Yield = assign(results, 7, false)

    market.Oil = assign(results, 8, true)
    market.Yen = assign(results, 9, true)
    market.Euro = assign(results, 10, true)
    market.Gold = assign(results, 11, true)

    return market
}
