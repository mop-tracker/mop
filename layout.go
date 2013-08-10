// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`bytes`
	`fmt`
	`reflect`
	`regexp`
	`strings`
	`text/template`
	`time`
)

const TotalColumns = 15

type Column struct {
	width	   int
	name       string
	title	   string
	formatter  func(string)string
}

type Layout struct {
	columns          []Column
	sorter           *Sorter
	regex            *regexp.Regexp
	market_template  *template.Template
	quotes_template  *template.Template
}

//-----------------------------------------------------------------------------
func (self *Layout) Initialize() *Layout {
	self.columns = []Column{
		{ -7, `Ticker`,    `Ticker`,    nil            },
		{ 10, `LastTrade`, `Last`,      currency       },
		{ 10, `Change`,    `Change`,    currency       },
		{ 10, `ChangePct`, `Change%`,   last           },
		{ 10, `Open`,      `Open`,      currency       },
		{ 10, `Low`,       `Low`,       currency       },
		{ 10, `High`,      `High`,      currency       },
		{ 10, `Low52`,     `52w Low`,   currency       },
		{ 10, `High52`,    `52w High`,  currency       },
		{ 11, `Volume`,    `Volume`,    nil            },
		{ 11, `AvgVolume`, `AvgVolume`, nil            },
		{  9, `PeRatio`,   `P/E`,       blank          },
		{  9, `Dividend`,  `Dividend`,  blank_currency },
		{  9, `Yield`,     `Yield`,     percent        },
		{ 11, `MarketCap`, `MktCap`,    currency       },
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
	//
	// Iterate over the list of stocks and properly format all its columns.
	//
	for i, stock := range quotes.stocks {
		pretty[i].Advancing = stock.Advancing
		//
		// Iterate over the list of stock columns. For each column name:
		// - Get current column value.
		// - If the column has the formatter method then call it.
		// - Set the column value padding it to the given width.
		//
		for _,column := range self.columns {
			// ex. value = stock.Change
			value := reflect.ValueOf(&stock).Elem().FieldByName(column.name).String()
			if column.formatter != nil {
				// ex. value = currency(value)
				value = column.formatter(value)
			}
			// ex. pretty[i].Change = self.pad(value, 10)
			reflect.ValueOf(&pretty[i]).Elem().FieldByName(column.name).SetString(self.pad(value, column.width))
		}
	}

	profile := quotes.profile
	if self.sorter == nil { // Initialize sorter on first invocation.
		self.sorter = new(Sorter).Initialize(profile)
	}
	self.sorter.SortByCurrentColumn(pretty)
	//
	// Group stocks by advancing/declining unless sorted by Chanage or Change%
	// in which case the grouping has been done already.
	//
	if profile.Grouped && (profile.SortColumn < 2 || profile.SortColumn > 3) {
		pretty = group(pretty)
	}

	return pretty
}

//-----------------------------------------------------------------------------
func (self *Layout) pad(str string, width int) string {
	match := self.regex.FindStringSubmatch(str)
	if len(match) > 0 {
		switch len(match[1]) {
		case 2:
			str = strings.Replace(str, match[1], match[1] + `0`, 1)
		case 4, 5:
			str = strings.Replace(str, match[1], match[1][0:3], 1)
		}
	}

	return fmt.Sprintf(`%*s`, width, str)
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
