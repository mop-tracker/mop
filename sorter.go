// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`sort`
	`strconv`
	`strings`
)

// Sorter gets called to sort stock quotes by one of the columns. The
// setup is rather lengthy; there should probably be more concise way
// that uses reflection and avoids hardcoding the column names.
type Sorter struct {
	profile *Profile // Pointer to where we store sort column and order.
}

type sortable []Stock

func (list sortable) Len() int      { return len(list) }
func (list sortable) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

type byTickerAsc struct{ sortable }
type byLastTradeAsc struct{ sortable }
type byChangeAsc struct{ sortable }
type byChangePctAsc struct{ sortable }
type byOpenAsc struct{ sortable }
type byLowAsc struct{ sortable }
type byHighAsc struct{ sortable }
type byLow52Asc struct{ sortable }
type byHigh52Asc struct{ sortable }
type byVolumeAsc struct{ sortable }
type byAvgVolumeAsc struct{ sortable }
type byPeRatioAsc struct{ sortable }
type byDividendAsc struct{ sortable }
type byYieldAsc struct{ sortable }
type byMarketCapAsc struct{ sortable }
type byPreOpenAsc struct{ sortable }
type byAfterHoursAsc struct{ sortable }

type byTickerDesc struct{ sortable }
type byLastTradeDesc struct{ sortable }
type byChangeDesc struct{ sortable }
type byChangePctDesc struct{ sortable }
type byOpenDesc struct{ sortable }
type byLowDesc struct{ sortable }
type byHighDesc struct{ sortable }
type byLow52Desc struct{ sortable }
type byHigh52Desc struct{ sortable }
type byVolumeDesc struct{ sortable }
type byAvgVolumeDesc struct{ sortable }
type byPeRatioDesc struct{ sortable }
type byDividendDesc struct{ sortable }
type byYieldDesc struct{ sortable }
type byMarketCapDesc struct{ sortable }
type byPreOpenDesc struct{ sortable }
type byAfterHoursDesc struct{ sortable }

func (list byTickerAsc) Less(i, j int) bool {
	return list.sortable[i].Ticker < list.sortable[j].Ticker
}
func (list byLastTradeAsc) Less(i, j int) bool {
	return list.sortable[i].LastTrade < list.sortable[j].LastTrade
}
func (list byChangeAsc) Less(i, j int) bool {
	return c(list.sortable[i].Change) < c(list.sortable[j].Change)
}
func (list byChangePctAsc) Less(i, j int) bool {
	return c(list.sortable[i].ChangePct) < c(list.sortable[j].ChangePct)
}
func (list byOpenAsc) Less(i, j int) bool {
	return list.sortable[i].Open < list.sortable[j].Open
}
func (list byLowAsc) Less(i, j int) bool { 
	return list.sortable[i].Low < list.sortable[j].Low
}
func (list byHighAsc) Less(i, j int) bool {
	return list.sortable[i].High < list.sortable[j].High
}
func (list byLow52Asc) Less(i, j int) bool {
	return list.sortable[i].Low52 < list.sortable[j].Low52
}
func (list byHigh52Asc) Less(i, j int) bool {
	return list.sortable[i].High52 < list.sortable[j].High52
}
func (list byVolumeAsc) Less(i, j int) bool {
	return m(list.sortable[i].Volume) < m(list.sortable[j].Volume)
}
func (list byAvgVolumeAsc) Less(i, j int) bool {
	return m(list.sortable[i].AvgVolume) < m(list.sortable[j].AvgVolume)
}
func (list byPeRatioAsc) Less(i, j int) bool {
	return list.sortable[i].PeRatio < list.sortable[j].PeRatio
}
func (list byDividendAsc) Less(i, j int) bool {
	return list.sortable[i].Dividend < list.sortable[j].Dividend
}
func (list byYieldAsc) Less(i, j int) bool {
	return list.sortable[i].Yield < list.sortable[j].Yield
}
func (list byMarketCapAsc) Less(i, j int) bool {
	return m(list.sortable[i].MarketCap) < m(list.sortable[j].MarketCap)
}
func (list byPreOpenAsc) Less(i, j int) bool {
	return c(list.sortable[i].PreOpen) < c(list.sortable[j].PreOpen)
}
func (list byAfterHoursAsc) Less(i, j int) bool {
	return c(list.sortable[i].AfterHours) < c(list.sortable[j].AfterHours)
}


func (list byTickerDesc) Less(i, j int) bool {
	return list.sortable[j].Ticker < list.sortable[i].Ticker
}
func (list byLastTradeDesc) Less(i, j int) bool {
	return list.sortable[j].LastTrade < list.sortable[i].LastTrade
}
func (list byChangeDesc) Less(i, j int) bool {
	return c(list.sortable[j].Change) < c(list.sortable[i].Change)
}
func (list byChangePctDesc) Less(i, j int) bool {
	return c(list.sortable[j].ChangePct) < c(list.sortable[i].ChangePct)
}
func (list byOpenDesc) Less(i, j int) bool {
	return list.sortable[j].Open < list.sortable[i].Open
}
func (list byLowDesc) Less(i, j int) bool { 
	return list.sortable[j].Low < list.sortable[i].Low
}
func (list byHighDesc) Less(i, j int) bool {
	return list.sortable[j].High < list.sortable[i].High
}
func (list byLow52Desc) Less(i, j int) bool {
	return list.sortable[j].Low52 < list.sortable[i].Low52
}
func (list byHigh52Desc) Less(i, j int) bool {
	return list.sortable[j].High52 < list.sortable[i].High52
}
func (list byVolumeDesc) Less(i, j int) bool {
	return m(list.sortable[j].Volume) < m(list.sortable[i].Volume)
}
func (list byAvgVolumeDesc) Less(i, j int) bool {
	return m(list.sortable[j].AvgVolume) < m(list.sortable[i].AvgVolume)
}
func (list byPeRatioDesc) Less(i, j int) bool {
	return list.sortable[j].PeRatio < list.sortable[i].PeRatio
}
func (list byDividendDesc) Less(i, j int) bool {
	return list.sortable[j].Dividend < list.sortable[i].Dividend
}
func (list byYieldDesc) Less(i, j int) bool {
	return list.sortable[j].Yield < list.sortable[i].Yield
}
func (list byMarketCapDesc) Less(i, j int) bool {
	return m(list.sortable[j].MarketCap) < m(list.sortable[i].MarketCap)
}
func (list byPreOpenDesc) Less(i, j int) bool {
	return c(list.sortable[j].PreOpen) < c(list.sortable[i].PreOpen)
}
func (list byAfterHoursDesc) Less(i, j int) bool {
	return c(list.sortable[j].AfterHours) < c(list.sortable[i].AfterHours)
}

// Returns new Sorter struct.
func NewSorter(profile *Profile) *Sorter {
	return &Sorter{
		profile: profile,
	}
}

// SortByCurrentColumn builds a list of sort interface based on current sort
// order, then calls sort.Sort to do the actual job.
func (sorter *Sorter) SortByCurrentColumn(stocks []Stock) *Sorter {
	var interfaces []sort.Interface

	if sorter.profile.Ascending {
		interfaces = []sort.Interface{
			byTickerAsc{stocks},
			byLastTradeAsc{stocks},
			byChangeAsc{stocks},
			byChangePctAsc{stocks},
			byOpenAsc{stocks},
			byLowAsc{stocks},
			byHighAsc{stocks},
			byLow52Asc{stocks},
			byHigh52Asc{stocks},
			byVolumeAsc{stocks},
			byAvgVolumeAsc{stocks},
			byPeRatioAsc{stocks},
			byDividendAsc{stocks},
			byYieldAsc{stocks},
			byMarketCapAsc{stocks},
			byPreOpenAsc{stocks},
			byAfterHoursAsc{stocks},
		}
	} else {
		interfaces = []sort.Interface{
			byTickerDesc{stocks},
			byLastTradeDesc{stocks},
			byChangeDesc{stocks},
			byChangePctDesc{stocks},
			byOpenDesc{stocks},
			byLowDesc{stocks},
			byHighDesc{stocks},
			byLow52Desc{stocks},
			byHigh52Desc{stocks},
			byVolumeDesc{stocks},
			byAvgVolumeDesc{stocks},
			byPeRatioDesc{stocks},
			byDividendDesc{stocks},
			byYieldDesc{stocks},
			byMarketCapDesc{stocks},
			byPreOpenDesc{stocks},
			byAfterHoursDesc{stocks},
		}
	}

	sort.Sort(interfaces[sorter.profile.SortColumn])

	return sorter
}

// The same exact method is used to sort by $Change and Change%. In both cases
// we sort by the value of Change% so that multiple $0.00s get sorted properly.
func c(str string) float32 {
	c := "$"
	for _, v := range currencies {
		if strings.Contains(str,v) {
			c = v
		}
	}
	trimmed := strings.Replace(strings.Trim(str, ` %`), c, ``, 1)
	value, _ := strconv.ParseFloat(trimmed, 32)
	return float32(value)
}

// When sorting by the market value we must first convert 42B etc. notations
// to proper numeric values.
func m(str string) float32 {
	if len(str) == 0 {
		return 0
	}

	multiplier := 1.0

	switch str[len(str)-1:] { // Check the last character.
	case `T`:
		multiplier = 1000000000000.0
	case `B`:
		multiplier = 1000000000.0
	case `M`:
		multiplier = 1000000.0
	case `K`:
		multiplier = 1000.0
	}

	trimmed := strings.Trim(str, ` $TBMK`) // Get rid of non-numeric characters.
	value, _ := strconv.ParseFloat(trimmed, 32)

	return float32(value * multiplier)
}
