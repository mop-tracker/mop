// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package provider

// Profile interface defines the methods required from the profile configuration
// used by the data providers.
type Profile interface {
	GetTickers() []string
	AddTickers(tickers []string) (int, error)
	RemoveTickers(tickers []string) (int, error)
}

// MarketIndex stores current market information for a single index.
type MarketIndex struct {
	Change  string
	Latest  string
	Percent string
	Name    string // optional
}

// Stock stores quote information for a particular stock ticker.
type Stock struct {
	Ticker     string // Stock ticker.
	LastTrade  string // last trade.
	Change     string // change real time.
	ChangePct  string // percent change real time.
	Open       string // market open price.
	Low        string // day's low.
	High       string // day's high.
	Low52      string // 52-weeks low.
	High52     string // 52-weeks high.
	Volume     string // volume.
	AvgVolume  string // average volume.
	PeRatio    string // P/E ration real time.
	PeRatioX   string // P/E ration (fallback when real time is N/A).
	Dividend   string // dividend.
	Yield      string // dividend yield.
	MarketCap  string // market cap real time.
	MarketCapX string // market cap (fallback when real time is N/A).
	Currency   string // String code for currency of stock.
	Direction  int    // -1 when change is < $0, 0 when change is = $0, 1 when change is > $0.
	PreOpen    string // pre-market change percent.
	AfterHours string // after-hours change percent.
}

// MarketData stores all the market-wide information.
type MarketData struct {
	Closed    bool
	Dow       MarketIndex
	Nasdaq    MarketIndex
	Sp500     MarketIndex
	Tokyo     MarketIndex
	HongKong  MarketIndex
	London    MarketIndex
	Frankfurt MarketIndex
	Yield     MarketIndex
	Oil       MarketIndex
	Yen       MarketIndex
	Euro      MarketIndex
	Gold      MarketIndex
}

// Market defines the interface for fetching and accessing market-wide data.
type Market interface {
	Fetch() Market
	Ok() (bool, string)
	IsClosed() bool
	GetData() *MarketData
	RefreshAdvice() int // Minimum refresh interval in seconds
}

// Quotes defines the interface for fetching and accessing stock-specific data.
type Quotes interface {
	Fetch() Quotes
	Ok() (bool, string)
	AddTickers(tickers []string) (int, error)
	RemoveTickers(tickers []string) (int, error)
	GetStocks() []Stock
	RefreshAdvice() int // Minimum refresh interval in seconds
	BindOnUpdate(func())
}
