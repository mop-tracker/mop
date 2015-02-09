// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package util

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
	`regexp`
	`strings`

)

const marketURL = `http://finance.yahoo.com`

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
	market.Sp500     = make(map[string]string)
	market.Dow       = make(map[string]string)
	market.Nasdaq    = make(map[string]string)
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
		`Nasdaq</span>`,     any, price, color, price, any, percent, any,
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
	finish := bytes.LastIndex(body, []byte(`<div id='ms-market-strip-0' class='ms-panel-ft Pt-6 Pb-6 Pos-r W-100'>`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (market *Market) extract(snippet []byte) *Market {
	matches := market.regex.FindStringSubmatch(string(snippet))
	fmt.Sprintf("cao ni ma %d", len(matches))
	fmt.Sprintf("cao ni ma %d", len(matches))

	fmt.Sprintf("cao ni ma %d", len(matches))
	fmt.Sprintf("cao ni ma %d", len(matches))
	fmt.Sprintf("cao ni ma %d", len(matches))
	fmt.Sprintf("cao ni ma %d", len(matches))

        // fmt.Printf("\n\n\n%q\n\n\n", matches)
    // if len(matches) < 37 {
    //     panic(`Unable to parse ` + marketURL)
    // }

	market.Sp500[`name`] = `S&P 500`
	market.Sp500[`latest`] = matches[1]

	market.Sp500[`change`] = matches[3]
	market.Sp500[`percent`] = matches[4]
	if matches[2] == `green` {
		market.Sp500[`change`] = `+` + market.Sp500[`change`]
		market.Sp500[`percent`] = `+` + market.Sp500[`percent`]
	} else if matches[2] == `red` {
		market.Sp500[`change`] = `-` + market.Sp500[`change`]
		market.Sp500[`percent`] = `-` + market.Sp500[`percent`]
	}

	market.Dow[`name`] = `Dow`
	market.Dow[`latest`] = matches[5]
	market.Dow[`change`] = matches[7]
	market.Dow[`percent`] = matches[8]
	if matches[6] == `green` {
		market.Dow[`change`] = `+` + market.Dow[`change`]
		market.Dow[`percent`] = `+` + market.Dow[`percent`]
	} else if matches[6] == `red` {
		market.Dow[`change`] = `-` + market.Dow[`change`]
		market.Dow[`percent`] = `-` + market.Dow[`percent`]
	}

	market.Nasdaq[`name`] = `NASDAQ`
	market.Nasdaq[`latest`] = matches[9]
	market.Nasdaq[`change`] = matches[11]
	market.Nasdaq[`percent`] = matches[12]
	if matches[10] == `green` {
		market.Nasdaq[`change`] = `+` + market.Nasdaq[`change`]
		market.Nasdaq[`percent`] = `+` + market.Nasdaq[`percent`]
	} else if matches[10] == `red` {
		market.Nasdaq[`change`] = `-` + market.Nasdaq[`change`]
		market.Nasdaq[`percent`] = `-` + market.Nasdaq[`percent`]
	}

	market.London[`name`] = `London`
	market.London[`latest`] = matches[13]
	market.London[`change`] = matches[15]
	market.London[`percent`] = matches[16]
	if matches[14] == `green` {
		market.London[`change`] = `+` + market.London[`change`]
		market.London[`percent`] = `+` + market.London[`percent`]
	} else if matches[14] == `red` {
		market.London[`change`] = `-` + market.London[`change`]
		market.London[`percent`] = `-` + market.London[`percent`]
	}

	market.Frankfurt[`name`] = `Frankfurt`
	market.Frankfurt[`latest`] = matches[17]
	market.Frankfurt[`change`] = matches[19]
	market.Frankfurt[`percent`] = matches[20]
	if matches[18] == `green` {
		market.Frankfurt[`change`] = `+` + market.Frankfurt[`change`]
		market.Frankfurt[`percent`] = `+` + market.Frankfurt[`percent`]
	} else if matches[18] == `red` {
		market.Frankfurt[`change`] = `-` + market.Frankfurt[`change`]
		market.Frankfurt[`percent`] = `-` + market.Frankfurt[`percent`]
	}

	market.Paris[`name`] = `Paris`
	market.Paris[`latest`] = matches[21]
	market.Paris[`change`] = matches[23]
	market.Paris[`percent`] = matches[24]
	if matches[22] == `green` {
		market.Paris[`change`] = `+` + market.Paris[`change`]
		market.Paris[`percent`] = `+` + market.Paris[`percent`]
	} else if matches[22] == `red` {
		market.Paris[`change`] = `-` + market.Paris[`change`]
		market.Paris[`percent`] = `-` + market.Paris[`percent`]
	}

	market.Tokyo[`name`] = `Tokyo`
	market.Tokyo[`latest`] = matches[25]
	market.Tokyo[`change`] = matches[27]
	market.Tokyo[`percent`] = matches[28]
	if matches[26] == `green` {
		market.Tokyo[`change`] = `+` + market.Tokyo[`change`]
		market.Tokyo[`percent`] = `+` + market.Tokyo[`percent`]
	} else if matches[26] == `red` {
		market.Tokyo[`change`] = `-` + market.Tokyo[`change`]
		market.Tokyo[`percent`] = `-` + market.Tokyo[`percent`]
	}

	market.HongKong[`name`] = `Hong Kong`
	market.HongKong[`latest`] = matches[29]
	market.HongKong[`change`] = matches[31]
	market.HongKong[`percent`] = matches[32]
	if matches[30] == `green` {
		market.HongKong[`change`] = `+` + market.HongKong[`change`]
		market.HongKong[`percent`] = `+` + market.HongKong[`percent`]
	} else if matches[30] == `red` {
		market.HongKong[`change`] = `-` + market.HongKong[`change`]
		market.HongKong[`percent`] = `-` + market.HongKong[`percent`]
	}

	market.Shanghai[`name`] = `Shanghai`
	market.Shanghai[`latest`] = matches[33]
	market.Shanghai[`change`] = matches[35]
	market.Shanghai[`percent`] = matches[36]
	if matches[34] == `green` {
		market.Shanghai[`change`] = `+` + market.Shanghai[`change`]
		market.Shanghai[`percent`] = `+` + market.Shanghai[`percent`]
	} else if matches[34] == `red` {
		market.Shanghai[`change`] = `-` + market.Shanghai[`change`]
		market.Shanghai[`percent`] = `-` + market.Shanghai[`percent`]
	}

	return market
}
