// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`sort`
	`strings`
	`strconv`
)

type Sortable []Stock
func (list Sortable) Len() int { return len(list) }
func (list Sortable) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

type ByTickerAsc       struct { Sortable }
type ByLastTradeAsc    struct { Sortable }
type ByChangeAsc       struct { Sortable }
type ByChangePctAsc    struct { Sortable }
type ByOpenAsc         struct { Sortable }
type ByLowAsc          struct { Sortable }
type ByHighAsc         struct { Sortable }
type ByLow52Asc        struct { Sortable }
type ByHigh52Asc       struct { Sortable }
type ByVolumeAsc       struct { Sortable }
type ByAvgVolumeAsc    struct { Sortable }
type ByPeRatioAsc      struct { Sortable }
type ByDividendAsc     struct { Sortable }
type ByYieldAsc        struct { Sortable }
type ByMarketCapAsc    struct { Sortable }
                       
type ByTickerDesc      struct { Sortable }
type ByLastTradeDesc   struct { Sortable }
type ByChangeDesc      struct { Sortable }
type ByChangePctDesc   struct { Sortable }
type ByOpenDesc        struct { Sortable }
type ByLowDesc         struct { Sortable }
type ByHighDesc        struct { Sortable }
type ByLow52Desc       struct { Sortable }
type ByHigh52Desc      struct { Sortable }
type ByVolumeDesc      struct { Sortable }
type ByAvgVolumeDesc   struct { Sortable }
type ByPeRatioDesc     struct { Sortable }
type ByDividendDesc    struct { Sortable }
type ByYieldDesc       struct { Sortable }
type ByMarketCapDesc   struct { Sortable }

func (list ByTickerAsc)       Less(i, j int) bool { return list.Sortable[i].Ticker        < list.Sortable[j].Ticker }
func (list ByLastTradeAsc)    Less(i, j int) bool { return list.Sortable[i].LastTrade     < list.Sortable[j].LastTrade }
func (list ByChangeAsc)       Less(i, j int) bool { return c(list.Sortable[i].Change)     < c(list.Sortable[j].Change) }
func (list ByChangePctAsc)    Less(i, j int) bool { return c(list.Sortable[i].ChangePct)  < c(list.Sortable[j].ChangePct) }
func (list ByOpenAsc)         Less(i, j int) bool { return list.Sortable[i].Open          < list.Sortable[j].Open }
func (list ByLowAsc)          Less(i, j int) bool { return list.Sortable[i].Low           < list.Sortable[j].Low }
func (list ByHighAsc)         Less(i, j int) bool { return list.Sortable[i].High          < list.Sortable[j].High }
func (list ByLow52Asc)        Less(i, j int) bool { return list.Sortable[i].Low52         < list.Sortable[j].Low52 }
func (list ByHigh52Asc)       Less(i, j int) bool { return list.Sortable[i].High52        < list.Sortable[j].High52 }
func (list ByVolumeAsc)       Less(i, j int) bool { return list.Sortable[i].Volume        < list.Sortable[j].Volume }
func (list ByAvgVolumeAsc)    Less(i, j int) bool { return list.Sortable[i].AvgVolume     < list.Sortable[j].AvgVolume }
func (list ByPeRatioAsc)      Less(i, j int) bool { return list.Sortable[i].PeRatio       < list.Sortable[j].PeRatio }
func (list ByDividendAsc)     Less(i, j int) bool { return list.Sortable[i].Dividend      < list.Sortable[j].Dividend }
func (list ByYieldAsc)        Less(i, j int) bool { return list.Sortable[i].Yield         < list.Sortable[j].Yield }
func (list ByMarketCapAsc)    Less(i, j int) bool { return m(list.Sortable[i].MarketCap)  < m(list.Sortable[j].MarketCap) }
                                  
func (list ByTickerDesc)      Less(i, j int) bool { return list.Sortable[j].Ticker        < list.Sortable[i].Ticker }
func (list ByLastTradeDesc)   Less(i, j int) bool { return list.Sortable[j].LastTrade     < list.Sortable[i].LastTrade }
func (list ByChangeDesc)      Less(i, j int) bool { return c(list.Sortable[j].ChangePct)  < c(list.Sortable[i].ChangePct) }
func (list ByChangePctDesc)   Less(i, j int) bool { return c(list.Sortable[j].ChangePct)  < c(list.Sortable[i].ChangePct) }
func (list ByOpenDesc)        Less(i, j int) bool { return list.Sortable[j].Open          < list.Sortable[i].Open }
func (list ByLowDesc)         Less(i, j int) bool { return list.Sortable[j].Low           < list.Sortable[i].Low }
func (list ByHighDesc)        Less(i, j int) bool { return list.Sortable[j].High          < list.Sortable[i].High }
func (list ByLow52Desc)       Less(i, j int) bool { return list.Sortable[j].Low52         < list.Sortable[i].Low52 }
func (list ByHigh52Desc)      Less(i, j int) bool { return list.Sortable[j].High52        < list.Sortable[i].High52 }
func (list ByVolumeDesc)      Less(i, j int) bool { return list.Sortable[j].Volume        < list.Sortable[i].Volume }
func (list ByAvgVolumeDesc)   Less(i, j int) bool { return list.Sortable[j].AvgVolume     < list.Sortable[i].AvgVolume }
func (list ByPeRatioDesc)     Less(i, j int) bool { return list.Sortable[j].PeRatio       < list.Sortable[i].PeRatio }
func (list ByDividendDesc)    Less(i, j int) bool { return list.Sortable[j].Dividend      < list.Sortable[i].Dividend }
func (list ByYieldDesc)       Less(i, j int) bool { return list.Sortable[j].Yield         < list.Sortable[i].Yield }
func (list ByMarketCapDesc)   Less(i, j int) bool { return m(list.Sortable[j].MarketCap)  < m(list.Sortable[i].MarketCap) }

type Sorter struct {
	profile  *Profile
}

func (self *Sorter) Initialize(profile *Profile) *Sorter {
	self.profile = profile

	return self
}

func (self *Sorter) SortByCurrentColumn(stocks []Stock) *Sorter {
	var interfaces []sort.Interface

	if self.profile.Ascending {
		interfaces = []sort.Interface{
			ByTickerAsc       { stocks },
			ByLastTradeAsc    { stocks },
			ByChangeAsc       { stocks },
			ByChangePctAsc    { stocks },
			ByOpenAsc         { stocks },
			ByLowAsc          { stocks },
			ByHighAsc         { stocks },
			ByLow52Asc        { stocks },
			ByHigh52Asc       { stocks },
			ByVolumeAsc       { stocks },
			ByAvgVolumeAsc    { stocks },
			ByPeRatioAsc      { stocks },
			ByDividendAsc     { stocks },
			ByYieldAsc        { stocks },
			ByMarketCapAsc    { stocks },
		}
	} else {
		interfaces = []sort.Interface{
			ByTickerDesc      { stocks },
			ByLastTradeDesc   { stocks },
			ByChangeDesc      { stocks },
			ByChangePctDesc   { stocks },
			ByOpenDesc        { stocks },
			ByLowDesc         { stocks },
			ByHighDesc        { stocks },
			ByLow52Desc       { stocks },
			ByHigh52Desc      { stocks },
			ByVolumeDesc      { stocks },
			ByAvgVolumeDesc   { stocks },
			ByPeRatioDesc     { stocks },
			ByDividendDesc    { stocks },
			ByYieldDesc       { stocks },
			ByMarketCapDesc   { stocks },
		}
	}

	sort.Sort(interfaces[self.profile.SortColumn])

	return self
}

// The same exact method is used to sort by Change and Change%. In both cases
// we sort by the value of Change% so that $0.00 change gets sorted proferly.
func c(str string) float32 {
	trimmed := strings.Replace(strings.Trim(str, ` %`), `$`, ``, 1)
	value, _ := strconv.ParseFloat(trimmed, 32)
	return float32(value)
}

func m(str string) float32 {
	multiplier := 1.0
	switch str[len(str)-1:len(str)] {
	case `B`:
		multiplier = 1000000000.0
	case `M`:
		multiplier = 1000000.0
	case `K`:
		multiplier = 1000.0
	}
	trimmed := strings.Trim(str, ` $BMK`)
	value, _ := strconv.ParseFloat(trimmed, 32)
	return float32(value * multiplier)
}
