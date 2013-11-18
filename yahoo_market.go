// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
	`regexp`
	`strings`
)

const marketURL = `http://finance.yahoo.com/` // `http://finance.yahoo.com/marketupdate/overview`

// Market stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type Market struct {
	IsClosed   bool		      // True when U.S. markets are closed.
	Dow        map[string]string  // Hash of Dow Jones indicators.
	Nasdaq     map[string]string  // Hash of NASDAQ indicators.
	Sp500      map[string]string  // Hash of S&P 500 indicators.
	London     map[string]string
	Frankfurt  map[string]string
	Paris      map[string]string
	Tokyo      map[string]string
	HongKong   map[string]string
	Shanghai   map[string]string
	regex      *regexp.Regexp     // Regex to parse market data from HTML.
	errors     string	      // Error(s), if any.
}

// Initialize creates empty hashes and builds regular expression used to parse
// market data from HTML page.
func (market *Market) Initialize() *Market {
	market.IsClosed  = false
	market.Dow       = make(map[string]string)
	market.Nasdaq    = make(map[string]string)
	market.Sp500     = make(map[string]string)
	market.London    = make(map[string]string)
	market.Frankfurt = make(map[string]string)
	market.Paris     = make(map[string]string)
	market.Tokyo     = make(map[string]string)
	market.HongKong  = make(map[string]string)
	market.Shanghai  = make(map[string]string)
	market.errors    = ``

	const any = `\s*<.+?>`
	const color = `<.+?price-change-([a-z]+)'>[\+\-]?`
	const price = `([\d\.,]+)`
	const percent = `\(([\d\.,%]+)\)`

	rules := []string{
		`S&P 500</span>`,    any, price, color, price, any, percent, any,
		`Dow</span>`,        any, price, color, price, any, percent, any,
		`NASDAQ</span>`,     any, price, color, price, any, percent, any,
		`FTSE</span>`,       any, price, color, price, any, percent, any,
		`DAX</span>`,        any, price, color, price, any, percent, any,
		`CAC 40</span>`,     any, price, color, price, any, percent, any,
		`NIKKEI 225</span>`, any, price, color, price, any, percent, any,
		`Hang Seng</span>`,  any, price, color, price, any, percent, any,
		`SSE Comp</span>`,   any, price, color, price, any, percent, any,
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

	body = market.checkIfMarketIsOpen(body)
	return market.extract(market.trim(body))
}

// Ok returns two values: 1) boolean indicating whether the error has occured,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

//-----------------------------------------------------------------------------
func (market *Market) checkIfMarketIsOpen(body []byte) []byte {
	start := bytes.Index(body, []byte(`id='yfs_market_time'`))
	finish := start + bytes.Index(body[start:], []byte(`</span>`))
	snippet := body[start:finish]
	market.IsClosed = bytes.Contains(snippet, []byte(`closed`)) || bytes.Contains(snippet, []byte(`open in`))

	return body[finish:]
}

//-----------------------------------------------------------------------------
func (market *Market) trim(body []byte) []byte {
	start := bytes.Index(body, []byte(`>S&P 500<`))
	finish := bytes.LastIndex(body, []byte(`id="mediafinancesuperherogs"`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (market *Market) extract(snippet []byte) *Market {
	matches := market.regex.FindAllStringSubmatch(string(snippet), -1)
	//fmt.Printf("\n\n\n%q\n\n\n", matches[0])
	if len(matches) < 1 || len(matches[0]) < 37 {
		panic(`Unable to parse ` + marketURL)
	}

	market.Dow[`name`] = `Dow`
	market.Dow[`latest`] = matches[0][1]
	market.Dow[`change`] = matches[0][3]
	market.Dow[`percent`] = matches[0][4]
	if matches[0][2] == `green` {
		market.Dow[`change`] = `+` + market.Dow[`change`]
		market.Dow[`percent`] = `+` + market.Dow[`percent`]
	} else if matches[0][2] == `?` {
		market.Dow[`change`] = `-` + market.Dow[`change`]
		market.Dow[`percent`] = `-` + market.Dow[`percent`]
	}

	market.Nasdaq[`name`] = `NASDAQ`
	market.Nasdaq[`latest`] = matches[0][5]
	market.Nasdaq[`change`] = matches[0][7]
	market.Nasdaq[`percent`] = matches[0][8]
	if matches[0][6] == `green` {
		market.Nasdaq[`change`] = `+` + market.Nasdaq[`change`]
		market.Nasdaq[`percent`] = `+` + market.Nasdaq[`percent`]
	} else if matches[0][2] == `?` {
		market.Nasdaq[`change`] = `-` + market.Nasdaq[`change`]
		market.Nasdaq[`percent`] = `-` + market.Nasdaq[`percent`]
	}

	market.Sp500[`name`] = `S&P 500`
	market.Sp500[`latest`] = matches[0][9]
	market.Sp500[`change`] = matches[0][11]
	market.Sp500[`percent`] = matches[0][12]
	if matches[0][10] == `green` {
		market.Sp500[`change`] = `+` + market.Sp500[`change`]
		market.Sp500[`percent`] = `+` + market.Sp500[`percent`]
	} else if matches[0][2] == `?` {
		market.Sp500[`change`] = `-` + market.Sp500[`change`]
		market.Sp500[`percent`] = `-` + market.Sp500[`percent`]
	}

	market.London[`name`] = `London`
	market.London[`latest`] = matches[0][13]
	market.London[`change`] = matches[0][15]
	market.London[`percent`] = matches[0][16]
	if matches[0][14] == `green` {
		market.London[`change`] = `+` + market.London[`change`]
		market.London[`percent`] = `+` + market.London[`percent`]
	} else if matches[0][2] == `?` {
		market.London[`change`] = `-` + market.London[`change`]
		market.London[`percent`] = `-` + market.London[`percent`]
	}

	market.Frankfurt[`name`] = `Frankfurt`
	market.Frankfurt[`latest`] = matches[0][17]
	market.Frankfurt[`change`] = matches[0][19]
	market.Frankfurt[`percent`] = matches[0][20]
	if matches[0][18] == `green` {
		market.Frankfurt[`change`] = `+` + market.Frankfurt[`change`]
		market.Frankfurt[`percent`] = `+` + market.Frankfurt[`percent`]
	} else if matches[0][2] == `?` {
		market.Frankfurt[`change`] = `-` + market.Frankfurt[`change`]
		market.Frankfurt[`percent`] = `-` + market.Frankfurt[`percent`]
	}

	market.Paris[`name`] = `Paris`
	market.Paris[`latest`] = matches[0][21]
	market.Paris[`change`] = matches[0][23]
	market.Paris[`percent`] = matches[0][24]
	if matches[0][22] == `green` {
		market.Paris[`change`] = `+` + market.Paris[`change`]
		market.Paris[`percent`] = `+` + market.Paris[`percent`]
	} else if matches[0][2] == `?` {
		market.Paris[`change`] = `-` + market.Paris[`change`]
		market.Paris[`percent`] = `-` + market.Paris[`percent`]
	}

	market.Tokyo[`name`] = `Tokyo`
	market.Tokyo[`latest`] = matches[0][25]
	market.Tokyo[`change`] = matches[0][27]
	market.Tokyo[`percent`] = matches[0][28]
	if matches[0][26] == `green` {
		market.Tokyo[`change`] = `+` + market.Tokyo[`change`]
		market.Tokyo[`percent`] = `+` + market.Tokyo[`percent`]
	} else if matches[0][2] == `?` {
		market.Tokyo[`change`] = `-` + market.Tokyo[`change`]
		market.Tokyo[`percent`] = `-` + market.Tokyo[`percent`]
	}

	market.HongKong[`name`] = `Hong Kong`
	market.HongKong[`latest`] = matches[0][29]
	market.HongKong[`change`] = matches[0][31]
	market.HongKong[`percent`] = matches[0][32]
	if matches[0][30] == `green` {
		market.HongKong[`change`] = `+` + market.HongKong[`change`]
		market.HongKong[`percent`] = `+` + market.HongKong[`percent`]
	} else if matches[0][2] == `?` {
		market.HongKong[`change`] = `-` + market.HongKong[`change`]
		market.HongKong[`percent`] = `-` + market.HongKong[`percent`]
	}

	market.Shanghai[`name`] = `Shanghai`
	market.Shanghai[`latest`] = matches[0][33]
	market.Shanghai[`change`] = matches[0][36]
	market.Shanghai[`percent`] = matches[0][36]
	if matches[0][34] == `green` {
		market.Shanghai[`change`] = `+` + market.Shanghai[`change`]
		market.Shanghai[`percent`] = `+` + market.Shanghai[`percent`]
	} else if matches[0][2] == `?` {
		market.Shanghai[`change`] = `-` + market.Shanghai[`change`]
		market.Shanghai[`percent`] = `-` + market.Shanghai[`percent`]
	}

	return market
}
