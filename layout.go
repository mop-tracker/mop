// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`bytes`
	`fmt`
	`log`
	"math"
	`reflect`
	`regexp`
	`strings`
	`text/template`
	`time`
)

// Column describes formatting rules for individual column within
// the list of stock quotes.
type Column struct {
	width     int                 // Column width.
	name      string              // The name of the field in the Stock struct.
	title     string              // Column title to display in the header.
	formatter func(string) string // Optional function to format the contents of the column.
}

// Layout is used to format and display all the collected data,
// i.e. market updates and the list of stock quotes.
type Layout struct {
	columns        []Column           // List of stock quotes columns.
	sorter         *Sorter            // Pointer to sorting receiver.
	regex          *regexp.Regexp     // Pointer to regular expression to align decimal points.
	marketTemplate *template.Template // Pointer to template to format market data.
	quotesTemplate *template.Template // Pointer to template to format the list of stock quotes.
}

// Initialize assigns the default values that stay unchanged for
// the life of allocated Layout struct.
func (layout *Layout) Initialize() *Layout {
	layout.columns = []Column{
		{-7, `Ticker`, `Ticker`, nil},
		{10, `LastTrade`, `Last`, currency},
		{10, `Change`, `Change`, currency},
		{10, `ChangePct`, `Change%`, last},
		{10, `Open`, `Open`, currency},
		{10, `Low`, `Low`, currency},
		{10, `High`, `High`, currency},
		{10, `Low52`, `52w Low`, currency},
		{10, `High52`, `52w High`, currency},
		{11, `Volume`, `Volume`, nil},
		{11, `AvgVolume`, `AvgVolume`, nil},
		{9, `PeRatio`, `P/E`, blank},
		{9, `Dividend`, `Dividend`, zero},
		{9, `Yield`, `Yield`, percent},
		{11, `MarketCap`, `MktCap`, currency},
	}
	layout.regex = regexp.MustCompile(`(\.\d+)[BMK]?$`)
	layout.marketTemplate = buildMarketTemplate()
	layout.quotesTemplate = buildQuotesTemplate()

	return layout
}

// Market merges given market data structure with the market
// template and returns formatted string that includes
// highlighting markup.
func (layout *Layout) Market(market *Market) string {
	if ok, err := market.Ok(); !ok { // If there was an error fetching market data...
		return err // then simply return the error string.
	}

	highlight(market.Dow, market.Sp500, market.Nasdaq, market.London, market.Frankfurt,
		market.Paris, market.Tokyo, market.HongKong, market.Shanghai)
	buffer := new(bytes.Buffer)
	layout.marketTemplate.Execute(buffer, market)

	return buffer.String()
}

// Quotes uses quotes template to format timestamp, stock quotes
// header, and the list of given stock quotes. It returns
// formatted string with all the necessary markup.
func (layout *Layout) Quotes(quotes Quoter) string {
	if ok, err := quotes.Ok(); !ok { // If there was an error
		//fetching stock quotes...
		return err // then simply return the error string.
	}

	vars := struct {
		Now    string  // Current timestamp.
		Header string  // Formatted header line.
		Stocks []Stock // List of formatted stock quotes.
	}{
		time.Now().Format(`3:04:05pm PST`),
		layout.Header(quotes.GetProfile()),
		layout.prettify(quotes),
	}

	buffer := new(bytes.Buffer)
	layout.quotesTemplate.Execute(buffer, vars)

	return buffer.String()
}

// Header iterates over column titles and formats the header line.
// The formatting includes placing an arrow next to the sorted
// column title. When the column editor is active it knows how to
// highlight currently selected column title.
func (layout *Layout) Header(profile *Profile) string {
	str, selectedColumn := ``, profile.selectedColumn

	for i, col := range layout.columns {
		arrow := arrowFor(i, profile)
		if i != selectedColumn {
			str += fmt.Sprintf(`%*s`, col.width, arrow+col.title)
		} else {
			str += fmt.Sprintf(`<r>%*s</r>`, col.width, arrow+col.title)
		}
	}

	return `<u>` + str + `</u>`
}

// TotalColumns is the utility method for the column editor that
// returns total number of columns.
func (layout *Layout) TotalColumns() int {
	return len(layout.columns)
}

//-----------------------------------------------------------------------------
func (layout *Layout) prettify(quotes Quoter) []Stock {
	stocks := quotes.GetStocks()

	pretty := make([]Stock, len(stocks))
	//
	// Iterate over the list of stocks and properly format all its columns.
	//
	for i, stock := range stocks {
		pretty[i].SetAdvance(stock.GetAdvance())
		//
		// Iterate over the list of stock columns. For each
		// column name:
		// - Get current column value.
		// - If the column has the formatter method then call it.
		// - Set the column value padding it to the given width.
		//
		for _, column := range layout.columns {
			// ex. value = stock.Change
			value := reflect.ValueOf(&stock).Elem().FieldByName(column.name).String()
			if column.formatter != nil {
				// ex. value = currency(value)
				value = column.formatter(value)
			}
			// ex. pretty[i].Change = layout.pad(value, 10)
			reflect.ValueOf(&pretty[i]).Elem().FieldByName(column.name).SetString(layout.pad(value, column.width))
		}
	}

	profile := quotes.GetProfile()
	if layout.sorter == nil { // Initialize sorter on first invocation.
		layout.sorter = new(Sorter).Initialize(profile)
	}
	layout.sorter.SortByCurrentColumn(pretty)
	//
	// Group stocks by advancing/declining unless sorted by
	// Chanage or Change% in which case the grouping has been
	// done already.
	//
	if profile.Grouped && (profile.SortColumn < 2 || profile.SortColumn > 3) {
		pretty = group(pretty)
	}

	return pretty
}

//-----------------------------------------------------------------------------
func (layout *Layout) pad(str string, width int) string {
	match := layout.regex.FindStringSubmatch(str)
	if len(match) > 0 {
		switch len(match[1]) {
		case 2:
			str = strings.Replace(str, match[1], match[1]+`0`, 1)
		case 4, 5:
			str = strings.Replace(str, match[1], match[1][0:3], 1)
		}
	}

	return fmt.Sprintf(`%*s`, width, str)
}

//-----------------------------------------------------------------------------
func buildMarketTemplate() *template.Template {
	markup := `<yellow>{{.Dow.name}}</> {{.Dow.change}} ({{.Dow.percent}}) at {{.Dow.latest}} <yellow>{{.Sp500.name}}</> {{.Sp500.change}} ({{.Sp500.percent}}) at {{.Sp500.latest}} <yellow>{{.Nasdaq.name}}</> {{.Nasdaq.change}} ({{.Nasdaq.percent}}) at {{.Nasdaq.latest}}
<yellow>{{.London.name}}</> {{.London.change}} ({{.London.percent}}) at {{.London.latest}} <yellow>{{.Frankfurt.name}}</> {{.Frankfurt.change}} ({{.Frankfurt.percent}}) at {{.Frankfurt.latest}} <yellow>{{.Paris.name}}</> {{.Paris.change}} ({{.Paris.percent}}) at {{.Paris.latest}} {{if .IsClosed}}<right>U.S. markets closed</right>{{end}}
<yellow>{{.Tokyo.name}}</> {{.Tokyo.change}} ({{.Tokyo.percent}}) at {{.Tokyo.latest}} <yellow>{{.HongKong.name}}</> {{.HongKong.change}} ({{.HongKong.percent}}) at {{.HongKong.latest}} <yellow>{{.Shanghai.name}}</> {{.Shanghai.change}} ({{.Shanghai.percent}}) at {{.Shanghai.latest}}`

	return template.Must(template.New(`market`).Parse(markup))
}

//-----------------------------------------------------------------------------
func buildQuotesTemplate() *template.Template {
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

	for _, stock := range stocks {
		if stock.GetAdvance() {
			grouped[current] = stock
			current++
		}
	}
	for _, stock := range stocks {
		if !stock.GetAdvance() {
			grouped[current] = stock
			current++
		}
	}

	return grouped
}

//-----------------------------------------------------------------------------
func arrowFor(column int, profile *Profile) string {
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
func zero(str string) string {
	if str == `0.00` {
		return `-`
	}

	return currency(str)
}

//-----------------------------------------------------------------------------
func last(str string) string {
	if len(str) >= 6 && str[0:6] != `N/A - ` {
		return str[0:int(math.Min(float64(len(str)), 9))]
	}

	if len(str) < 6 {
		return ""
	}

	return str[6:]
}

//-----------------------------------------------------------------------------
func currency(str string) string {
	if str == `N/A` {
		return `-`
	}
	if str == "" {
		return ""
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
