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

const marketURL = `http://finance.yahoo.com/marketupdate/overview`

// Market stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type Market struct {
	IsClosed  bool              // True when U.S. markets are closed.
	Dow       map[string]string // Hash of Dow Jones indicators.
	Nasdaq    map[string]string // Hash of NASDAQ indicators.
	Sp500     map[string]string // Hash of S&P 500 indicators.
	Advances  map[string]string // Number of advanced stocks on NYSE and NASDAQ.
	Declines  map[string]string // Ditto for declines.
	Unchanged map[string]string // Ditto for unchanged.
	Highs     map[string]string // Number of new highs on NYSE and NASDAQ.
	Lows      map[string]string // Ditto for new lows.
	regex     *regexp.Regexp    // Regex to parse market data from HTML.
	errors    string            // Error(s), if any.
}

// Initialize creates empty hashes and builds regular expression used to parse
// market data from HTML page.
func (market *Market) Initialize() *Market {
	market.IsClosed = false
	market.Dow = make(map[string]string)
	market.Nasdaq = make(map[string]string)
	market.Sp500 = make(map[string]string)
	market.Advances = make(map[string]string)
	market.Declines = make(map[string]string)
	market.Unchanged = make(map[string]string)
	market.Highs = make(map[string]string)
	market.Lows = make(map[string]string)
	market.errors = ``

	const any = `\s*<.+?>`
	const some = `<.+?`
	const space = `\s*`
	const color = `#([08c]{6});">\s*`
	const price = `([\d\.,]+)`
	const percent = `\(([\d\.,%]+)\)`

	rules := []string{
		`(Dow)`, any, price, some, color, price, some, percent, any,
		`(Nasdaq)`, any, price, some, color, price, some, percent, any,
		`(S&P 500)`, any, price, some, color, price, some, percent, any,
		`(Advances)`, any, price, space, percent, any, price, space, percent, any,
		`(Declines)`, any, price, space, percent, any, price, space, percent, any,
		`(Unchanged)`, any, price, space, percent, any, price, space, percent, any,
		`(New Hi's)`, any, price, any, price, any,
		`(New Lo's)`, any, price, any, price, any,
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
	start := bytes.Index(body, []byte(`id="yfs_market_time"`))
	finish := start + bytes.Index(body[start:], []byte(`</span>`))
	snippet := body[start:finish]
	market.IsClosed = bytes.Contains(snippet, []byte(`closed`)) || bytes.Contains(snippet, []byte(`open in`))

	return body[finish:]
}

//-----------------------------------------------------------------------------
func (market *Market) trim(body []byte) []byte {
	start := bytes.Index(body, []byte(`<table id="yfimktsumm"`))
	finish := bytes.LastIndex(body, []byte(`<table id="yfimktsumm"`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (market *Market) extract(snippet []byte) *Market {
	matches := market.regex.FindAllStringSubmatch(string(snippet), -1)
	if len(matches) < 1 || len(matches[0]) < 37 {
		panic(`Unable to parse ` + marketURL)
	}

	market.Dow[`name`] = matches[0][1]
	market.Dow[`latest`] = matches[0][2]
	market.Dow[`change`] = matches[0][4]
	switch matches[0][3] {
	case `008800`:
		market.Dow[`change`] = `+` + matches[0][4]
		market.Dow[`percent`] = `+` + matches[0][5]
	case `cc0000`:
		market.Dow[`change`] = `-` + matches[0][4]
		market.Dow[`percent`] = `-` + matches[0][5]
	default:
		market.Dow[`change`] = matches[0][4]
		market.Dow[`percent`] = matches[0][5]
	}

	market.Nasdaq[`name`] = matches[0][6]
	market.Nasdaq[`latest`] = matches[0][7]
	switch matches[0][8] {
	case `008800`:
		market.Nasdaq[`change`] = `+` + matches[0][9]
		market.Nasdaq[`percent`] = `+` + matches[0][10]
	case `cc0000`:
		market.Nasdaq[`change`] = `-` + matches[0][9]
		market.Nasdaq[`percent`] = `-` + matches[0][10]
	default:
		market.Nasdaq[`change`] = matches[0][9]
		market.Nasdaq[`percent`] = matches[0][10]
	}

	market.Sp500[`name`] = matches[0][11]
	market.Sp500[`latest`] = matches[0][12]
	switch matches[0][13] {
	case `008800`:
		market.Sp500[`change`] = `+` + matches[0][14]
		market.Sp500[`percent`] = `+` + matches[0][15]
	case `cc0000`:
		market.Sp500[`change`] = `-` + matches[0][14]
		market.Sp500[`percent`] = `-` + matches[0][15]
	default:
		market.Sp500[`change`] = matches[0][14]
		market.Sp500[`percent`] = matches[0][15]
	}

	market.Advances[`name`] = matches[0][16]
	market.Advances[`nyse`] = matches[0][17]
	market.Advances[`nysep`] = matches[0][18]
	market.Advances[`nasdaq`] = matches[0][19]
	market.Advances[`nasdaqp`] = matches[0][20]

	market.Declines[`name`] = matches[0][21]
	market.Declines[`nyse`] = matches[0][22]
	market.Declines[`nysep`] = matches[0][23]
	market.Declines[`nasdaq`] = matches[0][24]
	market.Declines[`nasdaqp`] = matches[0][25]

	market.Unchanged[`name`] = matches[0][26]
	market.Unchanged[`nyse`] = matches[0][27]
	market.Unchanged[`nysep`] = matches[0][28]
	market.Unchanged[`nasdaq`] = matches[0][29]
	market.Unchanged[`nasdaqp`] = matches[0][30]

	market.Highs[`name`] = matches[0][31]
	market.Highs[`nyse`] = matches[0][32]
	market.Highs[`nasdaq`] = matches[0][33]
	market.Lows[`name`] = matches[0][34]
	market.Lows[`nyse`] = matches[0][35]
	market.Lows[`nasdaq`] = matches[0][36]

	return market
}
