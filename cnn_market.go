// Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package TerminalStocks

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
	`regexp`
	`strings`
)

const marketURL = `http://money.cnn.com/data/markets/`

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
	regex     *regexp.Regexp // Regex to parse market data from HTML.
	errors    string         // Error(s), if any.
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

	market.errors = ``

	const any = `\s*(?:.+?)`
	const change = `>([\+\-]?[\d\.,]+)<\/span>`
	const price = `>([\d\.,]+)<\/span>`
	const percent = `>([\+\-]?[\d\.,]+%?)<`

	rules := []string{
		`>Dow<`, any, percent, any, price, any, change, any,
		`>Nasdaq<`, any, percent, any, price, any, change, any,
		`">S&P<`, any, percent, any, price, any, change, any,
		`>10\-year yield<`, any, price, any, percent, any,
		`>Oil<`, any, price, any, percent, any,
		`>Yen<`, any, price, any, percent, any,
		`>Euro<`, any, price, any, percent, any,
		`>Gold<`, any, price, any, percent, any,
		`>Nikkei 225<`, any, percent, any, price, any, change, any,
		`>Hang Seng<`, any, percent, any, price, any, change, any,
		`>FTSE 100<`, any, percent, any, price, any, change, any,
		`>DAX<`, any, percent, any, price, any, change, any,
	}

	market.regex = regexp.MustCompile(strings.Join(rules, ``))

	return market
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) Fetch() (self *Market) {
	self = market // <-- This ensures we return correct market after recover() from panic().
	defer func() {
		if err := recover(); err != nil {
			market.errors = fmt.Sprintf("Error fetching market data...\n%s", err)
		}
	}()

	response, err := http.Get(marketURL)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = market.isMarketOpen(body)
	return market.extract(market.trim(body))
}

// Ok returns two values: 1) boolean indicating whether the error has occured,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

//-----------------------------------------------------------------------------
func (market *Market) isMarketOpen(body []byte) []byte {
	// TBD -- CNN page doesn't seem to have market open/close indicator.
	return body
}

//-----------------------------------------------------------------------------
func (market *Market) trim(body []byte) []byte {
	start := bytes.Index(body, []byte(`Markets Overview`))
	finish := bytes.LastIndex(body, []byte(`Gainers`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (market *Market) extract(snippet []byte) *Market {
	matches := market.regex.FindStringSubmatch(string(snippet))

	if len(matches) < 31 {
		panic(`Unable to parse ` + marketURL)
	}

	market.Dow[`change`] = matches[1]
	market.Dow[`latest`] = matches[2]
	market.Dow[`percent`] = matches[3]

	market.Nasdaq[`change`] = matches[4]
	market.Nasdaq[`latest`] = matches[5]
	market.Nasdaq[`percent`] = matches[6]

	market.Sp500[`change`] = matches[7]
	market.Sp500[`latest`] = matches[8]
	market.Sp500[`percent`] = matches[9]


	market.Yield[`name`] = `10-year Yield`
	market.Yield[`latest`] = matches[10]
	market.Yield[`change`] = matches[11]

	market.Oil[`latest`] = matches[12]
	market.Oil[`change`] = matches[13]

	market.Yen[`latest`] = matches[14]
	market.Yen[`change`] = matches[15]

	market.Euro[`latest`] = matches[16]
	market.Euro[`change`] = matches[17]

	market.Gold[`latest`] = matches[18]
	market.Gold[`change`] = matches[19]

	market.Tokyo[`change`] = matches[20]
	market.Tokyo[`latest`] = matches[21]
	market.Tokyo[`percent`] = matches[22]

	market.HongKong[`change`] = matches[23]
	market.HongKong[`latest`] = matches[24]
	market.HongKong[`percent`] = matches[25]

	market.London[`change`] = matches[26]
	market.London[`latest`] = matches[27]
	market.London[`percent`] = matches[28]

	market.Frankfurt[`change`] = matches[29]
	market.Frankfurt[`latest`] = matches[30]
	market.Frankfurt[`percent`] = matches[31]

	return market
}

