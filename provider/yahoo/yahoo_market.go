// Copyright (c) 2013-2023 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package yahoo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mop-tracker/mop/provider"
)

const (
	marketURL           = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`
	marketURLQueryParts = `range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`
	symbols             = `^DJI,^IXIC,^GSPC,^N225,^HSI,^FTSE,^GDAXI,^TNX,CL=F,JPY=X,EUR=X,GC=F`
)

// MarketIndex stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type MarketIndex = provider.MarketIndex

type Market struct {
	provider.MarketData
	errors  string // Error(s), if any.
	url     string // URL with symbols to fetch data
	cookies string // cookies for auth
	crumb   string // crumb for the cookies, to be applied as a query param
}

// Returns new initialized Market struct.
func NewMarket() *Market {
	market := &Market{}
	market.Closed = false

	market.cookies = fetchCookies()
	market.crumb = fetchCrumb(market.cookies)

	// Construct URL with query parameters using url.Values
	params := url.Values{}
	params.Add("range", "1d")
	params.Add("interval", "5m")
	params.Add("indicators", "close")
	params.Add("includeTimestamps", "false")
	params.Add("includePrePost", "false")
	params.Add("corsDomain", "finance.yahoo.com")
	params.Add(".tsrc", "finance")

	market.url = fmt.Sprintf(marketURL, market.crumb, symbols) + "&" + params.Encode()

	market.errors = ""

	return market
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) Fetch() provider.Market {
	defer func() {
		if err := recover(); err != nil {
			market.errors = fmt.Sprintf("Error fetching market data...\n%s", err)
		} else {
			market.errors = ""
		}
	}()

	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, market.url, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"application/json"},
		"Cookie":          {market.cookies},
		"Host":            {"query1.finance.yahoo.com"},
		"Origin":          {"https://finance.yahoo.com"},
		"Referer":         {"https://finance.yahoo.com"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
		"TE":              {"trailers"},
		"User-Agent":      {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = market.isMarketOpen(body)
	return market.extract(body)
}

// Ok returns two values: 1) boolean indicating whether the error has occurred,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

func (market *Market) IsClosed() bool {
	return market.MarketData.Closed
}

func (market *Market) GetData() *provider.MarketData {
	return &market.MarketData
}

func (market *Market) RefreshAdvice() int {
	return 1
}

// -----------------------------------------------------------------------------
func (market *Market) isMarketOpen(body []byte) []byte {
	// TBD -- CNN page doesn't seem to have market open/close indicator.
	return body
}

// -----------------------------------------------------------------------------
func assign(result struct {
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
}, changeAsPercent bool) MarketIndex {
	change := float2Str(result.RegularMarketChange)
	latest := float2Str(result.RegularMarketPrice)
	percent := ""

	if changeAsPercent {
		change = float2Str(result.RegularMarketChangePercent) + `%`
	} else {
		percent = float2Str(result.RegularMarketChangePercent) + `%`
	}

	return MarketIndex{
		Change:  change,
		Latest:  latest,
		Percent: percent,
	}
}

// -----------------------------------------------------------------------------
func (market *Market) extract(body []byte) *Market {
	var d struct {
		MarketResponse struct {
			Result []struct {
				RegularMarketChange        float64 `json:"regularMarketChange"`
				RegularMarketPrice         float64 `json:"regularMarketPrice"`
				RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
			} `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := json.Unmarshal(body, &d); err != nil {
		panic(fmt.Sprintf("JSON unmarshal failed: %v", err))
	}

	results := d.MarketResponse.Result
	if len(results) < 12 {
		panic(fmt.Sprintf("unexpected number of results: got %d, expected at least 12", len(results)))
	}

	market.Dow = assign(results[0], false)
	market.Nasdaq = assign(results[1], false)
	market.Sp500 = assign(results[2], false)
	market.Tokyo = assign(results[3], false)
	market.HongKong = assign(results[4], false)
	market.London = assign(results[5], false)
	market.Frankfurt = assign(results[6], false)
	market.Yield.Name = "10-year Yield"
	market.Yield = assign(results[7], false)

	market.Oil.Name = "Crude Oil"
	market.Oil = assign(results[8], true)
	market.Yen.Name = "Yen"
	market.Yen = assign(results[9], true)
	market.Euro.Name = "Euro"
	market.Euro = assign(results[10], true)
	market.Gold.Name = "Gold"
	market.Gold = assign(results[11], true)

	return market
}
