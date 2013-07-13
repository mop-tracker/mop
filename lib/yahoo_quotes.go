// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// See http://www.gummy-stuff.org/Yahoo-data.htm
// Current, Change, Open, High, Low, 52-W High, 52-W Low, Volume, AvgVolume, P/E, Yield, Market Cap.
// l1: last trade
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
// d: dividend/share
// y: wield
// j3: market cap rt
// j1: market cap

const yahoo_quotes_url = `http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=,l1c6k2oghjkva2r2rdyj3j1`

// "AAPL", 417.42, "-3.38", "N/A - -0.80%", 420.33, 415.35,   423.29, 385.10, 705.07, 9788680, 15181900, N/A, 10.04, 11.00, 2.61, N/A, 391.8B
// "ALU",    1.83, "+0.07", "N/A - +3.98%",   1.77,   1.75,     1.83,   0.91,   2.01, 7957103, 11640700, N/A,   N/A,  0.00,  N/A, N/A,   4.156B
// "IBM",  194.93, "+1.68", "N/A - +0.87%", 192.83, 192.3501, 195.16, 181.85, 215.90, 2407971,  4376120, N/A, 13.33,  3.50, 1.81, N/A, 216.1B
// "TSLA", 120.09, "+4.85", "N/A - +4.21%", 118.75, 115.70,   120.28,  25.52, 121.89, 6827497,  9464530, N/A,   N/A,  0.00,  N/A, N/A, 13.877B

type Quote struct {
	Ticker        string
	LastTrade     string
	Change        string
	ChangePercent string
	Open          string
	Low           string
	High          string
	Low52         string
	High52        string
	Volume        string
	AvgVolume     string
	PeRatio       string
	PeRatioX      string
	Dividend      string
	Yield         string
	MarketCap     string
	MarketCapX    string
}
type Quotes []Quote

func GetQuotes(tickers string) Quotes {

	// Format the URL and send the request.
	// url := fmt.Sprintf(yahoo_quotes_url, strings.Join(tickers, "+"))
	url := fmt.Sprintf(yahoo_quotes_url, tickers)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return parse(sanitize(body))
}

func (q *Quote) Color() string {
	if strings.Index(q.Change, "-") == -1 {
		return "</green><green>"
	} else {
		return "" // "</red><red>"
	}
}

func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}

func parse(body []byte) Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes := make(Quotes, len(lines))

	for i, line := range lines {
		// fmt.Printf("\n\n{%d} -> [%s]\n\n", i, string(line))
		parse_line(line, &quotes[i])
	}

	return quotes
}

func parse_line(line []byte, quote *Quote) {
	columns := bytes.Split(bytes.TrimSpace(line), []byte{','})

	quote.Ticker = string(columns[0])
	quote.LastTrade = string(columns[1])
	quote.Change = string(columns[2])
	quote.ChangePercent = string(columns[3])
	quote.Open = string(columns[4])
	quote.Low = string(columns[5])
	quote.High = string(columns[6])
	quote.Low52 = string(columns[7])
	quote.High52 = string(columns[8])
	quote.Volume = string(columns[9])
	quote.AvgVolume = string(columns[10])
	quote.PeRatio = string(columns[11])
	quote.PeRatioX = string(columns[12])
	quote.Dividend = string(columns[13])
	quote.Yield = string(columns[14])
	quote.MarketCap = string(columns[15])
	quote.MarketCapX = string(columns[16])
}

// func (quotes Quotes) Format() string {
//         str := time.Now().Format("3:04:05pm PST\n")
//
//         for _, q := range quotes {
//                 str += fmt.Sprintf("%s - %s - %s - %s\n", q.Ticker, q.Ask, q.Change, q.ChangePercent)
//         }
//         return str
// }

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
