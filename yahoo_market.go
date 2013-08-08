// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`fmt`
	`bytes`
	`io/ioutil`
	`net/http`
	`regexp`
	`strings`
)

const url = `http://finance.yahoo.com/marketupdate/overview`

type Market struct {
	IsClosed   bool
	Dow        map[string]string
	Nasdaq     map[string]string
	Sp500      map[string]string
	Advances   map[string]string
	Declines   map[string]string
	Unchanged  map[string]string
	Highs      map[string]string
	Lows       map[string]string
	regex      *regexp.Regexp
	errors     string
}

//-----------------------------------------------------------------------------
func (self *Market) Initialize() *Market {
	self.IsClosed   = false
	self.Dow        = make(map[string]string)
	self.Nasdaq     = make(map[string]string)
	self.Sp500      = make(map[string]string)
	self.Advances   = make(map[string]string)
	self.Declines   = make(map[string]string)
	self.Unchanged  = make(map[string]string)
	self.Highs      = make(map[string]string)
	self.Lows       = make(map[string]string)
	self.errors     = ``

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

	self.regex = regexp.MustCompile(strings.Join(rules, ``))

	return self
}

//-----------------------------------------------------------------------------
func (self *Market) Fetch() (this *Market) {
	this = self	// <-- This ensures we return correct self in case of panic attack.
	defer func() {
		if err := recover(); err != nil {
			self.errors = fmt.Sprintf("Error fetching market data...\n%s", err)
		}
	}()

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = self.check_if_market_open(body)
	return self.extract(self.trim(body))
}

//-----------------------------------------------------------------------------
func (self *Market) Format() string {
	return new(Layout).Initialize().Market(self)
}

//-----------------------------------------------------------------------------
func (self *Market) Ok() (bool, string) {
	return self.errors == ``, self.errors
}

// private
//-----------------------------------------------------------------------------
func (self *Market) check_if_market_open(body []byte) []byte {
	start := bytes.Index(body, []byte(`id="yfs_market_time"`))
	finish := start + bytes.Index(body[start:], []byte(`</span>`))
	snippet := body[start:finish]
	self.IsClosed = bytes.Contains(snippet, []byte(`closed`)) || bytes.Contains(snippet, []byte(`open in`))

	return body[finish:]
}

//-----------------------------------------------------------------------------
func (self *Market) trim(body []byte) []byte {
	start := bytes.Index(body, []byte(`<table id="yfimktsumm"`))
	finish := bytes.LastIndex(body, []byte(`<table id="yfimktsumm"`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (self *Market) extract(snippet []byte) *Market {
	matches := self.regex.FindAllStringSubmatch(string(snippet), -1)
	if len(matches) < 1 || len(matches[0]) < 37 {
		panic(`Unable to parse ` + url)
	}

	self.Dow[`name`] = matches[0][1]
	self.Dow[`latest`] = matches[0][2]
	self.Dow[`change`] = matches[0][4]
	switch matches[0][3] {
	case `008800`:
		self.Dow[`change`] = `+` + matches[0][4]
		self.Dow[`percent`] = `+` + matches[0][5]
	case `cc0000`:
		self.Dow[`change`] = `-` + matches[0][4]
		self.Dow[`percent`] = `-` + matches[0][5]
	default:
		self.Dow[`change`] = matches[0][4]
		self.Dow[`percent`] = matches[0][5]
	}

	self.Nasdaq[`name`] = matches[0][6]
	self.Nasdaq[`latest`] = matches[0][7]
	switch matches[0][8] {
	case `008800`:
		self.Nasdaq[`change`] = `+` + matches[0][9]
		self.Nasdaq[`percent`] = `+` + matches[0][10]
	case `cc0000`:
		self.Nasdaq[`change`] = `-` + matches[0][9]
		self.Nasdaq[`percent`] = `-` + matches[0][10]
	default:
		self.Nasdaq[`change`] = matches[0][9]
		self.Nasdaq[`percent`] = matches[0][10]
	}

	self.Sp500[`name`] = matches[0][11]
	self.Sp500[`latest`] = matches[0][12]
	switch matches[0][13] {
	case `008800`:
		self.Sp500[`change`] = `+` + matches[0][14]
		self.Sp500[`percent`] = `+` + matches[0][15]
	case `cc0000`:
		self.Sp500[`change`] = `-` + matches[0][14]
		self.Sp500[`percent`] = `-` + matches[0][15]
	default:
		self.Sp500[`change`] = matches[0][14]
		self.Sp500[`percent`] = matches[0][15]
	}

	self.Advances[`name`] = matches[0][16]
	self.Advances[`nyse`] = matches[0][17]
	self.Advances[`nysep`] = matches[0][18]
	self.Advances[`nasdaq`] = matches[0][19]
	self.Advances[`nasdaqp`] = matches[0][20]

	self.Declines[`name`] = matches[0][21]
	self.Declines[`nyse`] = matches[0][22]
	self.Declines[`nysep`] = matches[0][23]
	self.Declines[`nasdaq`] = matches[0][24]
	self.Declines[`nasdaqp`] = matches[0][25]

	self.Unchanged[`name`] = matches[0][26]
	self.Unchanged[`nyse`] = matches[0][27]
	self.Unchanged[`nysep`] = matches[0][28]
	self.Unchanged[`nasdaq`] = matches[0][29]
	self.Unchanged[`nasdaqp`] = matches[0][30]

	self.Highs[`name`] = matches[0][31]
	self.Highs[`nyse`] = matches[0][32]
	self.Highs[`nasdaq`] = matches[0][33]
	self.Lows[`name`] = matches[0][34]
	self.Lows[`nyse`] = matches[0][35]
	self.Lows[`nasdaq`] = matches[0][36]

	return self
}
