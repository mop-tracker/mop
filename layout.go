// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var currencies = map[string]string{
	"RUB": "₽",
	"GBP": "£",
	"GBp": "p",
	"SEK": "kr",
	"EUR": "€",
	"JPY": "¥",
}

// Column describes formatting rules for individual column within the list
// of stock quotes.
type Column struct {
	width     int                    // Column width.
	name      string                 // The name of the field in the Stock struct.
	title     string                 // Column title to display in the header.
	formatter func(...string) string // Optional function to format the contents of the column.
}

// Layout is used to format and display all the collected data, i.e. market
// updates and the list of stock quotes.
type Layout struct {
	columns        []Column           // List of stock quotes columns.
	sorter         *Sorter            // Pointer to sorting receiver.
	filter         *Filter            // Pointer to filtering receiver.
	regex          *regexp.Regexp     // Pointer to regular expression to align decimal points.
	marketTemplate *template.Template // Pointer to template to format market data.
	quotesTemplate *template.Template // Pointer to template to format the list of stock quotes.
}

// Creates the layout and assigns the default values that stay unchanged.
func NewLayout() *Layout {
	layout := &Layout{}
	layout.columns = []Column{
		{-10, `Ticker`, `Ticker`, nil},
		{10, `LastTrade`, `Last`, currency},
		{10, `Change`, `Change`, currency},
		{10, `ChangePct`, `Change%`, last},
		{10, `Open`, `Open`, currency},
		{10, `Low`, `Low`, currency},
		{10, `High`, `High`, currency},
		{10, `Low52`, `52w Low`, currency},
		{10, `High52`, `52w High`, currency},
		{11, `Volume`, `Volume`, integer},
		{11, `AvgVolume`, `AvgVolume`, integer},
		{9, `PeRatio`, `P/E`, blank},
		{9, `Dividend`, `Dividend`, zero},
		{9, `Yield`, `Yield`, percent},
		{11, `MarketCap`, `MktCap`, currency},
		{13, `PreOpen`, `PreMktChg%`, percent},
		{13, `AfterHours`, `AfterMktChg%`, percent},
	}
	layout.regex = regexp.MustCompile(`(\.\d+)[TBMK]?$`)
	layout.marketTemplate = buildMarketTemplate()
	layout.quotesTemplate = buildQuotesTemplate()

	return layout
}

// Market merges given market data structure with the market template and
// returns formatted string that includes highlighting markup.
func (layout *Layout) Market(market *Market) string {
	if ok, err := market.Ok(); !ok { // If there was an error fetching market data...
		return err // then simply return the error string.
	}

	highlight(market.Dow, market.Sp500, market.Nasdaq,
		market.Tokyo, market.HongKong, market.London, market.Frankfurt,
		market.Yield, market.Oil, market.Euro, market.Yen, market.Gold)
	buffer := new(bytes.Buffer)
	layout.marketTemplate.Execute(buffer, market)

	return buffer.String()
}

// Quotes uses quotes template to format timestamp, stock quotes header,
// and the list of given stock quotes. It returns formatted string with
// all the necessary markup.
func (layout *Layout) Quotes(quotes *Quotes) string {
	zonename, _ := time.Now().In(time.Local).Zone()
	if ok, err := quotes.Ok(); !ok { // If there was an error fetching stock quotes...
		return err // then simply return the error string.
	}

	vars := struct {
		Now    string  // Current timestamp.
		Header string  // Formatted header line.
		Stocks []Stock // List of formatted stock quotes.
	}{
		time.Now().Format(`3:04:05pm ` + zonename),
		layout.Header(quotes.profile),
		layout.prettify(quotes),
	}

	buffer := new(bytes.Buffer)
	layout.quotesTemplate.Execute(buffer, vars)

	return buffer.String()
}

// Header iterates over column titles and formats the header line. The
// formatting includes placing an arrow next to the sorted column title.
// When the column editor is active it knows how to highlight currently
// selected column title.
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

// TotalColumns is the utility method for the column editor that returns
// total number of columns.
func (layout *Layout) TotalColumns() int {
	return len(layout.columns)
}

// -----------------------------------------------------------------------------
func (layout *Layout) prettify(quotes *Quotes) []Stock {
	pretty := make([]Stock, len(quotes.stocks))

	//
	// Iterate over the list of stocks to get the longest ticker name (some tickers will exceed the allotted 10 char length for the Ticker column)
	// Save the longest ticker length and use max(longestlength, column.width) later in the second loop to keep the ticker indentations consistent
	//
	tickerWidth := 0
	for _, stock := range quotes.stocks {
		value := reflect.ValueOf(&stock).Elem().FieldByName(`Ticker`).String()
		currentLength := len(value)
		if currentLength > tickerWidth {
			tickerWidth = currentLength
		}
	}

	//
	// Iterate over the list of stocks and properly format all its columns.
	//
	for i, stock := range quotes.stocks {
		pretty[i].Direction = stock.Direction
		//
		// Iterate over the list of stock columns. For each column name:
		// - Get current column value.
		// - If the column has the formatter method then call it.
		// - Set the column value padding it to the given width.
		//
		for _, column := range layout.columns {
			// ex. value = stock.Change
			value := reflect.ValueOf(&stock).Elem().FieldByName(column.name).String()
			if column.formatter != nil {
				// ex. value = currency(value)
				value = column.formatter(value, stock.Currency)
			}
			// ex. pretty[i].Change = layout.pad(value, 10)
			if column.name == `Ticker` && (0-tickerWidth) < column.width {
				column.width = (0 - tickerWidth)
			}
			reflect.ValueOf(&pretty[i]).Elem().FieldByName(column.name).SetString(layout.pad(value, column.width))
		}
	}

	profile := quotes.profile

	if profile.Filter != "" { // Fix for blank display if invalid filter expression was cleared.
		if profile.filterExpression != nil {
			if layout.filter == nil { // Initialize filter on first invocation.
				layout.filter = NewFilter(profile)
			}
			pretty = layout.filter.Apply(pretty)
		}
	}

	if layout.sorter == nil { // Initialize sorter on first invocation.
		layout.sorter = NewSorter(profile)
	}
	layout.sorter.SortByCurrentColumn(pretty)
	//
	// Group stocks by advancing/declining unless sorted by Chanage or Change%
	// in which case the grouping has been done already.
	//
	if profile.Grouped && (profile.SortColumn < 2 || profile.SortColumn > 3) {
		pretty = group(pretty)
	}

	return pretty
}

// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
func buildMarketTemplate() *template.Template {
	markup := `<tag>Dow</> {{.Dow.change}} ({{.Dow.percent}}) at {{.Dow.latest}} <tag>S&P 500</> {{.Sp500.change}} ({{.Sp500.percent}}) at {{.Sp500.latest}} <tag>NASDAQ</> {{.Nasdaq.change}} ({{.Nasdaq.percent}}) at {{.Nasdaq.latest}}
<tag>Tokyo</> {{.Tokyo.change}} ({{.Tokyo.percent}}) at {{.Tokyo.latest}} <tag>HK</> {{.HongKong.change}} ({{.HongKong.percent}}) at {{.HongKong.latest}} <tag>London</> {{.London.change}} ({{.London.percent}}) at {{.London.latest}} <tag>Frankfurt</> {{.Frankfurt.change}} ({{.Frankfurt.percent}}) at {{.Frankfurt.latest}} {{if .IsClosed}}<right>U.S. markets closed</right>{{end}}
<tag>10-Year Yield</> {{.Yield.latest}} ({{.Yield.change}}) <tag>Euro</> ${{.Euro.latest}} ({{.Euro.change}}) <tag>Yen</> ¥{{.Yen.latest}} ({{.Yen.change}}) <tag>Oil</> ${{.Oil.latest}} ({{.Oil.change}}) <tag>Gold</> ${{.Gold.latest}} ({{.Gold.change}})`

	return template.Must(template.New(`market`).Parse(markup))
}

// -----------------------------------------------------------------------------
func buildQuotesTemplate() *template.Template {
	markup := `<right><time>{{.Now}}</></right>



<header>{{.Header}}</>
{{range.Stocks}}{{if eq .Direction 1}}<gain>{{else if eq .Direction -1}}<loss>{{end}}{{.Ticker}}{{.LastTrade}}{{.Change}}{{.ChangePct}}{{.Open}}{{.Low}}{{.High}}{{.Low52}}{{.High52}}{{.Volume}}{{.AvgVolume}}{{.PeRatio}}{{.Dividend}}{{.Yield}}{{.MarketCap}}{{.PreOpen}}{{.AfterHours}}</>
{{end}}`

	return template.Must(template.New(`quotes`).Parse(markup))
}

// -----------------------------------------------------------------------------
func highlight(collections ...map[string]string) {
	for _, collection := range collections {
		change := collection[`change`]
		if change[len(change)-1:] == `%` {
			change = change[0 : len(change)-1]
		}
		adv, err := strconv.ParseFloat(change, 64)
		if err == nil {
			if adv < 0.0 {
				collection[`change`] = `<loss>` + collection[`change`] + `</>`
			} else if adv > 0.0 {
				collection[`change`] = `<gain>` + collection[`change`] + `</>`
			}
		}
	}
}

// -----------------------------------------------------------------------------
func group(stocks []Stock) []Stock {
	grouped := make([]Stock, len(stocks))
	current := 0

	for _, stock := range stocks {
		if stock.Direction >= 0 {
			grouped[current] = stock
			current++
		}
	}
	for _, stock := range stocks {
		if stock.Direction < 0 {
			grouped[current] = stock
			current++
		}
	}

	return grouped
}

// -----------------------------------------------------------------------------
func arrowFor(column int, profile *Profile) string {
	if column == profile.SortColumn {
		if profile.Ascending {
			return string('▲')
		}
		return string('▼')
	}
	return ``
}

// -----------------------------------------------------------------------------
func blank(str ...string) string {
	if len(str) < 1 {
		return "ERR"
	}
	if (len(str[0]) == 3 && str[0][0:3] == `N/A`) || len(str[0]) == 0 {
		return `-`
	}

	return str[0]
}

// -----------------------------------------------------------------------------
func zero(str ...string) string {
	if len(str) < 2 {
		return "ERR"
	}
	if str[0] == `0.00` {
		return `-`
	}

	return currency(str[0], str[1])
}

// -----------------------------------------------------------------------------
func last(str ...string) string {
	if len(str) < 1 {
		return "ERR"
	}
	if len(str[0]) >= 6 && str[0][0:6] == `N/A - ` {
		return str[0][6:]
	}

	return percent(str[0])
}

// -----------------------------------------------------------------------------
func currency(str ...string) string {
	if len(str) < 2 {
		return "ERR"
	}
	//default to $
	symbol := "$"
	c, ok := currencies[str[1]]
	if ok {
		symbol = c
	}
	if str[0] == `N/A` || len(str[0]) == 0 {
		return `-`
	}
	if sign := str[0][0:1]; sign == `+` || sign == `-` {
		return sign + symbol + str[0][1:]
	}

	return symbol + str[0]
}

// Returns percent value truncated at 2 decimal points.
// -----------------------------------------------------------------------------
func percent(str ...string) string {
	if len(str) < 1 {
		return "ERR"
	}
	if str[0] == `N/A` || len(str[0]) == 0 {
		return `-`
	}

	split := strings.Split(str[0], ".")
	if len(split) == 2 {
		digits := len(split[1])
		if digits > 2 {
			digits = 2
		}
		str[0] = split[0] + "." + split[1][0:digits]
	}
	if str[0][len(str)-1] != '%' {
		str[0] += `%`
	}
	return str[0]
}

// Returns value as integer (no trailing digits after a '.').
// -----------------------------------------------------------------------------
func integer(str ...string) string {
	if len(str) < 1 {
		return "ERR"
	}
	if str[0] == `N/A` || len(str[0]) == 0 {
		return `-`
	}

	// Don't strip after the '.' if we have a value such as 123.45M
	if unicode.IsDigit(rune(str[0][len(str[0])-1])) {
		split := strings.Split(str[0], ".")
		if len(split) == 2 {
			return split[0]
		}
	}
	return str[0]
}
