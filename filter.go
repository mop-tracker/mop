// Copyright (c) 2013-2019 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"strings"
	"log"
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

// Apply builds a list of sort interface based on current sort
// order, then calls sort.Sort to do the actual job.
func (filter *Filter) Apply(stocks []Stock) []Stock {
	var filteredStocks []Stock

	for _, stock := range stocks {
		//var values = map[string]interface{}{
		//        "ticker":        strings.TrimSpace(stock.Ticker),
		//        "last":          m(stock.LastTrade),
		//        "change":        c(stock.Change),
		//        "changePercent": c(stock.ChangePct),
		//        "open":          m(stock.Open),
		//        "low":           m(stock.Low),
		//        "high":          m(stock.High),
		//        "low52":         m(stock.Low52),
		//        "high52":        m(stock.High52),
		//        "volume":        m(stock.Volume),
		//        "avgVolume":     m(stock.AvgVolume),
		//        "pe":            m(stock.PeRatio),
		//        "peX":           m(stock.PeRatioX),
		//        "dividend":      m(stock.Dividend),
		//        "yield":         m(stock.Yield),
		//        "mktCap":        m(stock.MarketCap),
		//        "mktCapX":       m(stock.MarketCapX),
		//        "advancing":     stock.Advancing,
		//}

		var values = make(map[string]interface{})
		//values["pe"] = 8;
		var err error
		values["ticker"]            = strings.TrimSpace(stock.Ticker)
		values["last"],err          = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.LastTrade),"$",""),64)
		values["change"],err        = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Change),"$",""),64)
		values["changePercent"],err = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.ChangePct),"$",""),64)
		values["open"],err          = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Open),"$",""),64)
		values["low"],err           = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Low),"$",""),64)
		values["high"],err          = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.High),"$",""),64)
		values["low52"],err         = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Low52),"$",""),64)
		values["high52"],err        = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.High52),"$",""),64)
		values["volume"],err        = strconv.ParseFloat(stock.Volume,64)
		values["avgVolume"],err     = strconv.ParseFloat(stock.AvgVolume,64)
		values["pe"],err            = strconv.ParseFloat(stock.PeRatio,64)
		values["peX"],err           = strconv.ParseFloat(stock.PeRatioX,64)
		values["dividend"],err      = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Dividend),"$",""),64)
		values["yield"],err         = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.Yield),"$",""),64)
		values["mktCap"],err        = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.MarketCap),"$",""),64)
		values["mktCapX"],err       = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(stock.MarketCapX),"$",""),64)
		values["advancing"]         = stock.Advancing

		log.Println(stock.Dividend)
		log.Println("m9: ",m("$9"))
		log.Printf("type m($9) is %T\n",m("$9"))
		log.Println("dividend:",stock.Dividend,m(stock.Dividend))
		log.Printf("dividend is %T\n",stock.Dividend)
		log.Printf("m(dividend) is %T\n",m(stock.Dividend))
		log.Printf("values[pe] is type %T\n",values["pe"])
		result, err := filter.profile.filterExpression.Evaluate(values)

		if err != nil {
			log.Println("In filter.go, err:",err)
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
