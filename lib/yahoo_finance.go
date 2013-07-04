// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
        "fmt"
        "time"
        "bytes"
        "net/http"
        "io/ioutil"
        // "strings"
)

// See http://www.gummy-stuff.org/Yahoo-data.htm
// Current, Change, Open, High, Low, 52-W High, 52-W Low, Volume, AvgVolume, P/E, Yield, Market Cap.
// b2: ask rt
// c6: change rt
// k2: change % rt
// o: open
// g: day's low
// h: day's high
// j: 52w low
// k: 52w high
// v: volume
// a2: avg volume
// r2: p/e rt
// r: p/e
// y: wield
// j3: market cap rt
// j1: market cap

const yahoo_finance_url = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,b2c6k2oghjkva2r2ryj3j1`

// "AAPL", 602.93, "+2.31", "N/A - +0.55%", 420.95, 417.45,  422.98,  385.10, 705.07,  8604594, 15205700, N/A,  9.99, 2.63, N/A, 395.0B
// "GOOG",   0.00, "+4.12", "N/A - +0.47%", 879.90, 878.50,  889.17,  562.09, 920.60,  1048628,  2353530, N/A, 26.40,  N/A, N/A, 294.1B
// "PG",    94.58, "+0.13", "N/A - +0.17%",  78.28,  77.4301, 78.75,   60.86,  82.54,  5347846,  9929320, N/A, 17.58, 2.92, N/A, 215.3B

type Quote struct {
	Ticker          []byte
	Ask             []byte
	Change          []byte
	ChangePercent   []byte
	Open            []byte
	Low             []byte
	High            []byte
	Low52           []byte
	High52          []byte
	Volume          []byte
	AvgVolume       []byte
	PeRatio         []byte
	PeRatioX        []byte
	Yield           []byte
        MarketCap       []byte
        MarketCapX      []byte
}
type Quotes []Quote

var quotes Quotes

// func Get(tickers []string) Quotes {
func Get(tickers string) Quotes {
	if len(quotes) > 0 && time.Now().Second() % 5 != 0 { // Fetch quotes every 5 seconds.
		return quotes
	}

        // Format the URL and send the request.
        // url := fmt.Sprintf(yahoo_finance_url, strings.Join(tickers, "+"))
        url := fmt.Sprintf(yahoo_finance_url, tickers)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
        fmt.Println("\n\n\n\n\n\rFetched quotes: " + time.Now().Format("3:04:05pm PST"))
	if err != nil {
		panic(err)
	}
        quotes = parse(sanitize(body))

	return quotes
}

func sanitize(body []byte) []byte {
        return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}

// "AAPL", 602.93, "+2.31", "N/A - +0.55%", 420.95, 417.45,  422.98,  385.10, 705.07,  8604594, 15205700, N/A,  9.99, 2.63, N/A, 395.0B
func parse(body []byte) Quotes {
        // fmt.Printf("[%s]\n", body)
        lines := bytes.Split(body, []byte{'\n'})
        quotes := make(Quotes, len(lines))

        for i,line := range lines {
                // fmt.Printf("\n\n{%d} -> [%s]\n\n", i, string(line))
                parse_line(line, &quotes[i])
        }

        return quotes
}

func parse_line(line []byte, quote *Quote) {
        // var quote Quote
        columns := bytes.Split(line, []byte{','})
        // fmt.Printf("{%s} -> [%d]", string(line), len(columns))

        quote.Ticker          = columns[0]
        quote.Ask             = columns[1]
        quote.Change          = columns[2]
        quote.ChangePercent   = columns[3]
        quote.Open            = columns[4]
        quote.Low             = columns[5]
        quote.High            = columns[6]
        quote.Low52           = columns[7]
        quote.High52          = columns[8]
        quote.Volume          = columns[9]
        quote.AvgVolume       = columns[10]
        quote.PeRatio         = columns[11]
        quote.PeRatioX        = columns[12]
        quote.Yield           = columns[13]
        quote.MarketCap       = columns[14]
        quote.MarketCapX      = columns[15]
}

func (quotes Quotes) Format() string {
        str := time.Now().Format("3:04:05pm PST\n")

        for _, q := range quotes {
                str += fmt.Sprintf("%s - %s - %s - %s\n", q.Ticker, q.Ask, q.Change, q.ChangePercent)
        }
        return str
}

//
// http://query.yahooapis.com/v1/public/yql
// ?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in(%22ALU%22,%22AAPL%22)
// &env=http%3A%2F%2Fdatatables.org%2Falltables.env
// &format=json'
//
// ^IXIC NASDAQ composite
// ^GSPC S&P 500
//
// {
//   "query": {
//     "count": 2,
//     "created": "2013-06-28T03:28:19Z",
//     "lang": "en-US",
//     "results": {
//       "quote": [
//         {
//           "AfterHoursChangeRealtime": "N/A - N/A",
//           "AnnualizedGain": null,
//           "Ask": null,
//           "AskRealtime": "1.91",
//           "AverageDailyVolume": "11692300",
//           "Bid": null,
//           "BidRealtime": "1.86",
//           "BookValue": "1.249",
//           "Change": "+0.12",
//           "ChangeFromFiftydayMovingAverage": "+0.1626",
//           "ChangeFromTwoHundreddayMovingAverage": "+0.321",
//           "ChangeFromYearHigh": "-0.16",
//           "ChangeFromYearLow": "+0.94",
//           "ChangePercentRealtime": "N/A - +6.94%",
//           "ChangeRealtime": "+0.12",
//           "Change_PercentChange": "+0.12 - +6.94%",
//           "ChangeinPercent": "+6.94%",
//           "Commission": null,
//           "DaysHigh": "1.92",
//           "DaysLow": "1.79",
//           "DaysRange": "1.79 - 1.92",
//           "DaysRangeRealtime": "N/A - N/A",
//           "DaysValueChange": "- - +6.94%",
//           "DaysValueChangeRealtime": "N/A - N/A",
//           "DividendPayDate": "29-Jun-07",
//           "DividendShare": "0.00",
//           "DividendYield": null,
//           "EBITDA": "802.7M",
//           "EPSEstimateCurrentYear": "-0.30",
//           "EPSEstimateNextQuarter": "-0.05",
//           "EPSEstimateNextYear": "-0.07",
//           "EarningsShare": "-1.213",
//           "ErrorIndicationreturnedforsymbolchangedinvalid": null,
//           "ExDividendDate": "31-May-07",
//           "FiftydayMovingAverage": "1.6874",
//           "HighLimit": null,
//           "HoldingsGain": null,
//           "HoldingsGainPercent": "- - -",
//           "HoldingsGainPercentRealtime": "N/A - N/A",
//           "HoldingsGainRealtime": null,
//           "HoldingsValue": null,
//           "HoldingsValueRealtime": null,
//           "LastTradeDate": "6/27/2013",
//           "LastTradePriceOnly": "1.85",
//           "LastTradeRealtimeWithTime": "N/A - <b>1.85</b>",
//           "LastTradeTime": "4:00pm",
//           "LastTradeWithTime": "Jun 27 - <b>1.85</b>",
//           "LowLimit": null,
//           "MarketCapRealtime": null,
//           "MarketCapitalization": "4.202B",
//           "MoreInfo": "cnprmIed",
//           "Name": "Alcatel-Lucent Co",
//           "Notes": null,
//           "OneyrTargetPrice": "2.16",
//           "Open": "1.81",
//           "OrderBookRealtime": null,
//           "PEGRatio": "0.22",
//           "PERatio": null,
//           "PERatioRealtime": null,
//           "PercebtChangeFromYearHigh": "-7.96%",
//           "PercentChange": "+6.94%",
//           "PercentChangeFromFiftydayMovingAverage": "+9.63%",
//           "PercentChangeFromTwoHundreddayMovingAverage": "+20.99%",
//           "PercentChangeFromYearLow": "+103.30%",
//           "PreviousClose": "1.73",
//           "PriceBook": "1.39",
//           "PriceEPSEstimateCurrentYear": null,
//           "PriceEPSEstimateNextYear": null,
//           "PricePaid": null,
//           "PriceSales": "0.21",
//           "SharesOwned": null,
//           "ShortRatio": "0.90",
//           "StockExchange": "NYSE",
//           "Symbol": "ALU",
//           "TickerTrend": " +=-=+- ",
//           "TradeDate": null,
//           "TwoHundreddayMovingAverage": "1.529",
//           "Volume": "34193168",
//           "YearHigh": "2.01",
//           "YearLow": "0.91",
//           "YearRange": "0.91 - 2.01",
//           "symbol": "ALU"
//         },
//         {
//           "AfterHoursChangeRealtime": "N/A - N/A",
//           "AnnualizedGain": null,
//           "Ask": "393.45",
//           "AskRealtime": "393.45",
//           "AverageDailyVolume": "17939600",
//           "Bid": "393.32",
//           "BidRealtime": "393.32",
//           "BookValue": "144.124",
//           "Change": "-4.29",
//           "ChangeFromFiftydayMovingAverage": "-37.81",
//           "ChangeFromTwoHundreddayMovingAverage": "-111.877",
//           "ChangeFromYearHigh": "-311.29",
//           "ChangeFromYearLow": "+8.68",
//           "ChangePercentRealtime": "N/A - -1.08%",
//           "ChangeRealtime": "-4.29",
//           "Change_PercentChange": "-4.29 - -1.08%",
//           "ChangeinPercent": "-1.08%",
//           "Commission": null,
//           "DaysHigh": "401.39",
//           "DaysLow": "393.54",
//           "DaysRange": "393.54 - 401.39",
//           "DaysRangeRealtime": "N/A - N/A",
//           "DaysValueChange": "- - -1.08%",
//           "DaysValueChangeRealtime": "N/A - N/A",
//           "DividendPayDate": "May 16",
//           "DividendShare": "7.95",
//           "DividendYield": "2.00",
//           "EBITDA": "57.381B",
//           "EPSEstimateCurrentYear": "39.57",
//           "EPSEstimateNextQuarter": "8.21",
//           "EPSEstimateNextYear": "43.71",
//           "EarningsShare": "41.896",
//           "ErrorIndicationreturnedforsymbolchangedinvalid": null,
//           "ExDividendDate": "Feb  7",
//           "FiftydayMovingAverage": "431.59",
//           "HighLimit": null,
//           "HoldingsGain": null,
//           "HoldingsGainPercent": "- - -",
//           "HoldingsGainPercentRealtime": "N/A - N/A",
//           "HoldingsGainRealtime": null,
//           "HoldingsValue": null,
//           "HoldingsValueRealtime": null,
//           "LastTradeDate": "6/27/2013",
//           "LastTradePriceOnly": "393.78",
//           "LastTradeRealtimeWithTime": "N/A - <b>393.78</b>",
//           "LastTradeTime": "4:00pm",
//           "LastTradeWithTime": "Jun 27 - <b>393.78</b>",
//           "LowLimit": null,
//           "MarketCapRealtime": null,
//           "MarketCapitalization": "369.6B",
//           "MoreInfo": "cnsprmiIed",
//           "Name": "Apple Inc.",
//           "Notes": null,
//           "OneyrTargetPrice": "539.54",
//           "Open": "399.01",
//           "OrderBookRealtime": null,
//           "PEGRatio": "0.48",
//           "PERatio": "9.50",
//           "PERatioRealtime": null,
//           "PercebtChangeFromYearHigh": "-44.15%",
//           "PercentChange": "-1.08%",
//           "PercentChangeFromFiftydayMovingAverage": "-8.76%",
//           "PercentChangeFromTwoHundreddayMovingAverage": "-22.13%",
//           "PercentChangeFromYearLow": "+2.25%",
//           "PreviousClose": "398.07",
//           "PriceBook": "2.76",
//           "PriceEPSEstimateCurrentYear": "10.06",
//           "PriceEPSEstimateNextYear": "9.11",
//           "PricePaid": null,
//           "PriceSales": "2.21",
//           "SharesOwned": null,
//           "ShortRatio": "1.50",
//           "StockExchange": "NasdaqNM",
//           "Symbol": "AAPL",
//           "TickerTrend": " +--=== ",
//           "TradeDate": null,
//           "TwoHundreddayMovingAverage": "505.657",
//           "Volume": "12050007",
//           "YearHigh": "705.07",
//           "YearLow": "385.10",
//           "YearRange": "385.10 - 705.07",
//           "symbol": "AAPL"
//         }
//       ]
//     }
//   }
// }