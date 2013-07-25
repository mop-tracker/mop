// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
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
	columns []Column
}

//-----------------------------------------------------------------------------
func (self *Layout) Initialize() *Layout {
	self.columns = make([]Column, TotalColumns)

	self.columns[0]  = Column{ -7, `Ticker`}
	self.columns[1]  = Column{ 10, `Last`}
	self.columns[2]  = Column{ 10, `Change`}
	self.columns[3]  = Column{ 10, `%Change`}
	self.columns[4]  = Column{ 10, `Open`}
	self.columns[5]  = Column{ 10, `Low`}
	self.columns[6]  = Column{ 10, `High`}
	self.columns[7]  = Column{ 10, `52w Low`}
	self.columns[8]  = Column{ 10, `52w High`}
	self.columns[9]  = Column{ 11, `Volume`}
	self.columns[10] = Column{ 11, `AvgVolume`}
	self.columns[11] = Column{ 10, `P/E`}
	self.columns[12] = Column{ 10, `Dividend`}
	self.columns[13] = Column{ 10, `Yield`}
	self.columns[14] = Column{ 11, `MktCap`}

	return self
}

//-----------------------------------------------------------------------------
func (self *Layout) Market(m *Market) string {
	markup := `{{.Dow.name}}: `
	if m.Dow[`change`][0:1] != `-` {
		markup += `<green>{{.Dow.change}} ({{.Dow.percent}})</green> at {{.Dow.latest}}, `
	} else {
		markup += `{{.Dow.change}} ({{.Dow.percent}}) at {{.Dow.latest}}, `
	}
	markup += `{{.Sp500.name}}: `
	if m.Sp500[`change`][0:1] != `-` {
		markup += `<green>{{.Sp500.change}} ({{.Sp500.percent}})</green> at {{.Sp500.latest}}, `
	} else {
		markup += `{{.Sp500.change}} ({{.Sp500.percent}}) at {{.Sp500.latest}}, `
	}
	markup += `{{.Nasdaq.name}}: `
	if m.Nasdaq[`change`][0:1] != `-` {
		markup += `<green>{{.Nasdaq.change}} ({{.Nasdaq.percent}})</green> at {{.Nasdaq.latest}}`
	} else {
		markup += `{{.Nasdaq.change}} ({{.Nasdaq.percent}}) at {{.Nasdaq.latest}}`
	}
	markup += "\n"
	markup += `{{.Advances.name}}: {{.Advances.nyse}} ({{.Advances.nysep}}) on NYSE and {{.Advances.nasdaq}} ({{.Advances.nasdaqp}}) on Nasdaq. `
	markup += `{{.Declines.name}}: {{.Declines.nyse}} ({{.Declines.nysep}}) on NYSE and {{.Declines.nasdaq}} ({{.Declines.nasdaqp}}) on Nasdaq`
	if !m.Open {
		markup += `<right>U.S. markets closed</right>`
	}
	markup += "\n"
	markup += `New highs: {{.Highs.nyse}} on NYSE and {{.Highs.nasdaq}} on Nasdaq. `
	markup += `New lows: {{.Lows.nyse}} on NYSE and {{.Lows.nasdaq}} on Nasdaq.`
	template, err := template.New(`market`).Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, m)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

//-----------------------------------------------------------------------------
func (self *Layout) Quotes(quotes *Quotes) string {
	vars := struct {
		Now    string
		Header string
		Stocks []Stock
	}{
		time.Now().Format(`3:04:05pm PST`),
		self.Header(quotes.profile.selected_column),
		self.prettify(quotes),
	}

	markup := `<right><white>{{.Now}}</white></right>



{{.Header}}
{{range.Stocks}}{{.Color}}{{.Ticker}}{{.LastTrade}}{{.Change}}{{.ChangePercent}}{{.Open}}{{.Low}}{{.High}}{{.Low52}}{{.High52}}{{.Volume}}{{.AvgVolume}}{{.PeRatio}}{{.Dividend}}{{.Yield}}{{.MarketCap}}
{{end}}`
	//markup += fmt.Sprintf("[%v]", quotes.profile.Grouped)
	template, err := template.New(`quotes`).Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, vars)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

//-----------------------------------------------------------------------------
func (self *Layout) Header(selected_column int) string {
	str := `<u>`
	for i,col := range self.columns {
		if i != selected_column {
			str += fmt.Sprintf(`%*s`, col.width, col.title)
		} else {
			str += fmt.Sprintf(`<r>%*s</r>`, col.width, col.title)
		}
	}
	str += `</u>`

	return str
}

//-----------------------------------------------------------------------------
func (self *Layout) prettify(quotes *Quotes) []Stock {
	pretty := make([]Stock, len(quotes.stocks))
	for i, q := range group(quotes) {
		pretty[i].Ticker        = pad(q.Ticker,                      self.columns[0].width)
		pretty[i].LastTrade     = pad(with_currency(q.LastTrade),    self.columns[1].width)
		pretty[i].Change        = pad(with_currency(q.Change),       self.columns[2].width)
		pretty[i].ChangePercent = pad(last_of_pair(q.ChangePercent), self.columns[3].width)
		pretty[i].Open          = pad(with_currency(q.Open),         self.columns[4].width)
		pretty[i].Low           = pad(with_currency(q.Low),          self.columns[5].width)
		pretty[i].High          = pad(with_currency(q.High),         self.columns[6].width)
		pretty[i].Low52         = pad(with_currency(q.Low52),        self.columns[7].width)
		pretty[i].High52        = pad(with_currency(q.High52),       self.columns[8].width)
		pretty[i].Volume        = pad(q.Volume,                      self.columns[9].width)
		pretty[i].AvgVolume     = pad(q.AvgVolume,                   self.columns[10].width)
		pretty[i].PeRatio       = pad(nullify(q.PeRatioX),           self.columns[11].width)
		pretty[i].Dividend      = pad(with_currency(q.Dividend),     self.columns[12].width)
		pretty[i].Yield         = pad(with_percent(q.Yield),         self.columns[13].width)
		pretty[i].MarketCap     = pad(with_currency(q.MarketCapX),   self.columns[14].width)
	}
	return pretty
}

//-----------------------------------------------------------------------------
func group(quotes *Quotes) []Stock {
	if !quotes.profile.Grouped {
		return quotes.stocks
	} else {
		grouped := make([]Stock, len(quotes.stocks))
		current := 0
		for _,stock := range quotes.stocks {
			if strings.Index(stock.Change, "-") == -1 {
				grouped[current] = stock
				current++
			}
		}
		for _,stock := range quotes.stocks {
			if strings.Index(stock.Change, "-") != -1 {
				grouped[current] = stock
				current++
			}
		}
		return grouped
	}
}

//-----------------------------------------------------------------------------
func nullify(str string) string {
	if len(str) == 3 && str[0:3] == `N/A` {
		return `-`
	} else {
		return str
	}
}

//-----------------------------------------------------------------------------
func last_of_pair(str string) string {
	if len(str) >= 6 && str[0:6] != `N/A - ` {
		return str
	} else {
		return str[6:]
	}
}

//-----------------------------------------------------------------------------
func with_currency(str string) string {
	if str == `N/A` {
		return `-`
	} else {
		switch str[0:1] {
		case `+`, `-`:
			return str[0:1] + `$` + str[1:]
		default:
			return `$` + str
		}
	}
}

//-----------------------------------------------------------------------------
func with_percent(str string) string {
	if str == `N/A` {
		return `-`
	} else {
		return str + `%`
	}
}

//-----------------------------------------------------------------------------
func colorize(str string) string {
	if str == `N/A` {
		return `-`
	} else if str[0:1] == `-` {
		return `<red>` + str + `</red>`
	} else {
		return `<green>` + str + `</green>`
	}
}

//-----------------------------------------------------------------------------
func ticker(str string, change string) string {
	if change[0:1] == `-` {
		return `<red>` + str + `</red>`
	} else {
		return `<green>` + str + `</green>`
	}
}

//-----------------------------------------------------------------------------
func pad(str string, width int) string {
	re := regexp.MustCompile(`(\.\d+)[MB]?$`)
	match := re.FindStringSubmatch(str)
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
