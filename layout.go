// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`bytes`
	`fmt`
	`regexp`
	`strings`
	`text/template`
	`time`
)

const TotalColumns = 15

type Column struct {
	width	int
	title	string
}

type Layout struct {
	columns          []Column
	regex            *regexp.Regexp
	market_template  *template.Template
	quotes_template  *template.Template
}

//-----------------------------------------------------------------------------
func (self *Layout) Initialize() *Layout {
	self.columns = []Column{
		{ -7, `Ticker` },
		{ 10, `Last` },
		{ 10, `Change` },
		{ 10, `Change%` },
		{ 10, `Open` },
		{ 10, `Low` },
		{ 10, `High` },
		{ 10, `52w Low` },
		{ 10, `52w High` },
		{ 11, `Volume` },
		{ 11, `AvgVolume` },
		{  9, `P/E` },
		{  9, `Dividend` },
		{  9, `Yield` },
		{ 11, `MktCap` },
	}
	self.regex = regexp.MustCompile(`(\.\d+)[MB]?$`)
	self.market_template = build_market_template()
	self.quotes_template = build_quotes_template()

	return self
}

//-----------------------------------------------------------------------------
func (self *Layout) Market(market *Market) string {
	if ok, err := market.Ok(); !ok {
		return err
	}

	highlight(market.Dow, market.Sp500, market.Nasdaq)
	buffer := new(bytes.Buffer)
	self.market_template.Execute(buffer, market)

	return buffer.String()
}

//-----------------------------------------------------------------------------
func (self *Layout) Quotes(quotes *Quotes) string {
	if ok, err := quotes.Ok(); !ok {
		return err
	}

	vars := struct {
		Now    string
		Header string
		Stocks []Stock
	}{
		time.Now().Format(`3:04:05pm PST`),
		self.Header(quotes.profile),
		self.prettify(quotes),
	}

	buffer := new(bytes.Buffer)
	self.quotes_template.Execute(buffer, vars)

	return buffer.String()
}

//-----------------------------------------------------------------------------
func (self *Layout) Header(profile *Profile) string {
	str, selected_column := ``, profile.selected_column

	for i,col := range self.columns {
		arrow := arrow_for(i, profile)
		if i != selected_column {
			str += fmt.Sprintf(`%*s`, col.width, arrow + col.title)
		} else {
			str += fmt.Sprintf(`<r>%*s</r>`, col.width, arrow + col.title)
		}
	}

	return `<u>` + str + `</u>`
}

//-----------------------------------------------------------------------------
func (self *Layout) prettify(quotes *Quotes) []Stock {
	pretty := make([]Stock, len(quotes.stocks))

	for i, q := range quotes.stocks {
		pretty[i].Ticker    = self.pad(q.Ticker,                    0)
		pretty[i].LastTrade = self.pad(currency(q.LastTrade),       1)
		pretty[i].Change    = self.pad(currency(q.Change),          2)
		pretty[i].ChangePct = self.pad(last(q.ChangePct),           3)
		pretty[i].Open      = self.pad(currency(q.Open),            4)
		pretty[i].Low       = self.pad(currency(q.Low),             5)
		pretty[i].High      = self.pad(currency(q.High),            6)
		pretty[i].Low52     = self.pad(currency(q.Low52),           7)
		pretty[i].High52    = self.pad(currency(q.High52),          8)
		pretty[i].Volume    = self.pad(q.Volume,                    9)
		pretty[i].AvgVolume = self.pad(q.AvgVolume,                10)
		pretty[i].PeRatio   = self.pad(blank(q.PeRatioX),          11)
		pretty[i].Dividend  = self.pad(blank_currency(q.Dividend), 12)
		pretty[i].Yield     = self.pad(percent(q.Yield),           13)
		pretty[i].MarketCap = self.pad(currency(q.MarketCapX),     14)
		pretty[i].Advancing = q.Advancing
	}

	profile := quotes.profile
	new(Sorter).Initialize(profile).SortByCurrentColumn(pretty)
	//
	// Group stocks by advancing/declining unless sorted by Chanage or Change%
	// in which case the grouping is done already.
	//
	if profile.Grouped && (profile.SortColumn < 2 || profile.SortColumn > 3) {
		pretty = group(pretty)
	}

	return pretty
}

//-----------------------------------------------------------------------------
func (self *Layout) pad(str string, col int) string {
	match := self.regex.FindStringSubmatch(str)
	if len(match) > 0 {
		switch len(match[1]) {
		case 2:
			str = strings.Replace(str, match[1], match[1] + `0`, 1)
		case 4, 5:
			str = strings.Replace(str, match[1], match[1][0:3], 1)
		}
	}

	return fmt.Sprintf(`%*s`, self.columns[col].width, str)
}

//-----------------------------------------------------------------------------
func build_market_template() *template.Template {
	markup := `{{.Dow.name}}: {{.Dow.change}} ({{.Dow.percent}}) at {{.Dow.latest}}, {{.Sp500.name}}: {{.Sp500.change}} ({{.Sp500.percent}}) at {{.Sp500.latest}}, {{.Nasdaq.name}}: {{.Nasdaq.change}} ({{.Nasdaq.percent}}) at {{.Nasdaq.latest}}
{{.Advances.name}}: {{.Advances.nyse}} ({{.Advances.nysep}}) on NYSE and {{.Advances.nasdaq}} ({{.Advances.nasdaqp}}) on Nasdaq. {{.Declines.name}}: {{.Declines.nyse}} ({{.Declines.nysep}}) on NYSE and {{.Declines.nasdaq}} ({{.Declines.nasdaqp}}) on Nasdaq {{if .IsClosed}}<right>U.S. markets closed</right>{{end}}
New highs: {{.Highs.nyse}} on NYSE and {{.Highs.nasdaq}} on Nasdaq. New lows: {{.Lows.nyse}} on NYSE and {{.Lows.nasdaq}} on Nasdaq.`

	return template.Must(template.New(`market`).Parse(markup))
}

//-----------------------------------------------------------------------------
func build_quotes_template() *template.Template {
	markup := `<right><white>{{.Now}}</></right>



{{.Header}}
{{range.Stocks}}{{if .Advancing}}<green>{{end}}{{.Ticker}}{{.LastTrade}}{{.Change}}{{.ChangePct}}{{.Open}}{{.Low}}{{.High}}{{.Low52}}{{.High52}}{{.Volume}}{{.AvgVolume}}{{.PeRatio}}{{.Dividend}}{{.Yield}}{{.MarketCap}}</>
{{end}}`

	return template.Must(template.New(`quotes`).Parse(markup))
}

//-----------------------------------------------------------------------------
func highlight(collections ...map[string]string) {
	for _, collection := range collections {
		if collection[`change`][0:1] != `-` {
			collection[`change`] = `<green>` + collection[`change`] + `</>`
		}
	}
}

//-----------------------------------------------------------------------------
func group(stocks []Stock) []Stock {
	grouped := make([]Stock, len(stocks))
	current := 0

	for _,stock := range stocks {
		if stock.Advancing {
			grouped[current] = stock
			current++
		}
	}
	for _,stock := range stocks {
		if !stock.Advancing {
			grouped[current] = stock
			current++
		}
	}

	return grouped
}

//-----------------------------------------------------------------------------
func arrow_for(column int, profile *Profile) string {
	if column == profile.SortColumn {
		if profile.Ascending {
			return string('\U00002191')
		}
		return string('\U00002193')
	}
	return ``
}

//-----------------------------------------------------------------------------
func blank(str string) string {
	if len(str) == 3 && str[0:3] == `N/A` {
		return `-`
	}

	return str
}

//-----------------------------------------------------------------------------
func blank_currency(str string) string {
	if str == `0.00` {
		return `-`
	}

	return currency(str)
}

//-----------------------------------------------------------------------------
func last(str string) string {
	if len(str) >= 6 && str[0:6] != `N/A - ` {
		return str
	}

	return str[6:]
}

//-----------------------------------------------------------------------------
func currency(str string) string {
	if str == `N/A` {
		return `-`
	}
	if sign := str[0:1]; sign == `+` || sign == `-` {
		return sign + `$` + str[1:]
	}

	return `$` + str
}

//-----------------------------------------------------------------------------
func percent(str string) string {
	if str == `N/A` {
		return `-`
	}

	return str + `%`
}
