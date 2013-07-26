// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
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
type ByPeRatioXAsc     struct { Sortable }
type ByDividendAsc     struct { Sortable }
type ByYieldAsc        struct { Sortable }
type ByMarketCapAsc    struct { Sortable }
type ByMarketCapXAsc   struct { Sortable }
                       
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
type ByPeRatioXDesc    struct { Sortable }
type ByDividendDesc    struct { Sortable }
type ByYieldDesc       struct { Sortable }
type ByMarketCapDesc   struct { Sortable }
type ByMarketCapXDesc  struct { Sortable }


func (list ByTickerAsc)       Less(i, j int) bool { return list.Sortable[i].Ticker        < list.Sortable[j].Ticker }
func (list ByLastTradeAsc)    Less(i, j int) bool { return list.Sortable[i].LastTrade     < list.Sortable[j].LastTrade }
func (list ByChangeAsc)       Less(i, j int) bool { return z(list.Sortable[i].Change)     < z(list.Sortable[j].Change) }
func (list ByChangePctAsc)    Less(i, j int) bool { return z(list.Sortable[i].ChangePct)  < z(list.Sortable[j].ChangePct) }
func (list ByOpenAsc)         Less(i, j int) bool { return list.Sortable[i].Open          < list.Sortable[j].Open }
func (list ByLowAsc)          Less(i, j int) bool { return list.Sortable[i].Low           < list.Sortable[j].Low }
func (list ByHighAsc)         Less(i, j int) bool { return list.Sortable[i].High          < list.Sortable[j].High }
func (list ByLow52Asc)        Less(i, j int) bool { return list.Sortable[i].Low52         < list.Sortable[j].Low52 }
func (list ByHigh52Asc)       Less(i, j int) bool { return list.Sortable[i].High52        < list.Sortable[j].High52 }
func (list ByVolumeAsc)       Less(i, j int) bool { return list.Sortable[i].Volume        < list.Sortable[j].Volume }
func (list ByAvgVolumeAsc)    Less(i, j int) bool { return list.Sortable[i].AvgVolume     < list.Sortable[j].AvgVolume }
func (list ByPeRatioAsc)      Less(i, j int) bool { return list.Sortable[i].PeRatio       < list.Sortable[j].PeRatio }
func (list ByPeRatioXAsc)     Less(i, j int) bool { return list.Sortable[i].PeRatioX      < list.Sortable[j].PeRatioX }
func (list ByDividendAsc)     Less(i, j int) bool { return list.Sortable[i].Dividend      < list.Sortable[j].Dividend }
func (list ByYieldAsc)        Less(i, j int) bool { return list.Sortable[i].Yield         < list.Sortable[j].Yield }
func (list ByMarketCapAsc)    Less(i, j int) bool { return list.Sortable[i].MarketCap     < list.Sortable[j].MarketCap }
func (list ByMarketCapXAsc)   Less(i, j int) bool { return list.Sortable[i].MarketCapX    < list.Sortable[j].MarketCapX }
                                  
func (list ByTickerDesc)      Less(i, j int) bool { return list.Sortable[j].Ticker        < list.Sortable[i].Ticker }
func (list ByLastTradeDesc)   Less(i, j int) bool { return list.Sortable[j].LastTrade     < list.Sortable[i].LastTrade }
func (list ByChangeDesc)      Less(i, j int) bool { return z(list.Sortable[j].Change)     < z(list.Sortable[i].Change) }
func (list ByChangePctDesc)   Less(i, j int) bool { return z(list.Sortable[j].ChangePct)  < z(list.Sortable[i].ChangePct) }
func (list ByOpenDesc)        Less(i, j int) bool { return list.Sortable[j].Open          < list.Sortable[i].Open }
func (list ByLowDesc)         Less(i, j int) bool { return list.Sortable[j].Low           < list.Sortable[i].Low }
func (list ByHighDesc)        Less(i, j int) bool { return list.Sortable[j].High          < list.Sortable[i].High }
func (list ByLow52Desc)       Less(i, j int) bool { return list.Sortable[j].Low52         < list.Sortable[i].Low52 }
func (list ByHigh52Desc)      Less(i, j int) bool { return list.Sortable[j].High52        < list.Sortable[i].High52 }
func (list ByVolumeDesc)      Less(i, j int) bool { return list.Sortable[j].Volume        < list.Sortable[i].Volume }
func (list ByAvgVolumeDesc)   Less(i, j int) bool { return list.Sortable[j].AvgVolume     < list.Sortable[i].AvgVolume }
func (list ByPeRatioDesc)     Less(i, j int) bool { return list.Sortable[j].PeRatio       < list.Sortable[i].PeRatio }
func (list ByPeRatioXDesc)    Less(i, j int) bool { return list.Sortable[j].PeRatioX      < list.Sortable[i].PeRatioX }
func (list ByDividendDesc)    Less(i, j int) bool { return list.Sortable[j].Dividend      < list.Sortable[i].Dividend }
func (list ByYieldDesc)       Less(i, j int) bool { return list.Sortable[j].Yield         < list.Sortable[i].Yield }
func (list ByMarketCapDesc)   Less(i, j int) bool { return list.Sortable[j].MarketCap     < list.Sortable[i].MarketCap }
func (list ByMarketCapXDesc)  Less(i, j int) bool { return list.Sortable[j].MarketCapX    < list.Sortable[i].MarketCapX }

func z(str string) float32 {
	float := strings.Replace(strings.Trim(str, ` %`), `$`, ``, 1)
	value,_ := strconv.ParseFloat(float, 32)
	return float32(value)
}