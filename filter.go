// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"strings"
	"strconv"
)

// Filter gets called to sort stock quotes by one of the columns. The
// setup is rather lengthy; there should probably be more concise way
// that uses reflection and avoids hardcoding the column names.
type Filter struct {
	profile *Profile // Pointer to where we store sort column and order.
}

// Returns new Filter struct.
func NewFilter(profile *Profile) *Filter {
	return &Filter{
		profile: profile,
	}
}

// Changes money and % notation to a plain float for math, comparisons.
func stringToNumber (numberString string) float64 {
	// If the string "$3.6B" is passed in, the returned float will be 3.6E+09.
	// If 0.03% is passed in, the returned float will be 0.03 (NOT 0.0003!).
	newString := strings.TrimSpace(numberString) // Take off whitespace.
	newString = strings.Replace(newString,"$","",1)      // Remove the $ symbol.
	newString = strings.Replace(newString,"%","",1)      // Remove the $ symbol.
	newString = strings.Replace(newString,"K","E+3",1)   // Thousand (kilo)
	newString = strings.Replace(newString,"M","E+6",1)   // Million
	newString = strings.Replace(newString,"B","E+9",1)   // Billion
	newString = strings.Replace(newString,"T","E+12",1)  // Trillion
	finalValue, _ := strconv.ParseFloat(newString, 64)
	return finalValue
}

// Apply builds a list of sort interface based on current sort
// order, then calls sort.Sort to do the actual job.
func (filter *Filter) Apply(stocks []Stock) []Stock {
	var filteredStocks []Stock

	for _, stock := range stocks {
		var values = make(map[string]interface{})
		// Make conversions from the strings to floats where necessary.
		values["ticker"]        = strings.TrimSpace(stock.Ticker) // Remains string
		values["last"]          = stringToNumber(stock.LastTrade)
		values["change"]        = stringToNumber(stock.Change)
		values["changePercent"] = stringToNumber(stock.ChangePct)
		values["open"]          = stringToNumber(stock.Open)
		values["low"]           = stringToNumber(stock.Low)
		values["high"]          = stringToNumber(stock.High)
		values["low52"]         = stringToNumber(stock.Low52)
		values["high52"]        = stringToNumber(stock.High52)
		values["dividend"]      = stringToNumber(stock.Dividend)
		values["yield"]         = stringToNumber(stock.Yield)
		values["mktCap"]        = stringToNumber(stock.MarketCap)
		values["mktCapX"]       = stringToNumber(stock.MarketCapX)
		values["volume"]        = stringToNumber(stock.Volume)
		values["avgVolume"]     = stringToNumber(stock.AvgVolume)
		values["pe"]            = stringToNumber(stock.PeRatio)
		values["peX"]           = stringToNumber(stock.PeRatioX)
		values["direction"]     = stock.Direction                 // Remains int.

		result, err := filter.profile.filterExpression.Evaluate(values)

		if err != nil {
                        // The filter isn't working, so reset to no filter.
                        filter.profile.Filter = ""
                        // Return an empty list.  The next main loop cycle will
                        // show unfiltered.
                        return filteredStocks
		}

		truthy, ok := result.(bool)

		if !ok {
                        // The filter isn't working, so reset to no filter.
                        filter.profile.Filter = ""
                        // Return an empty list.  The next main loop cycle will
                        // show unfiltered.
                        return filteredStocks
		}

		if truthy {
			filteredStocks = append(filteredStocks, stock)
		}
	}

	return filteredStocks
}
