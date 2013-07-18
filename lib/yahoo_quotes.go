// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	`bytes`
	`fmt`
	`io/ioutil`
	`net/http`
	`strings`
)

// See http://www.gummy-stuff.org/Yahoo-data.htm
// Also http://query.yahooapis.com/v1/public/yql
// ?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in(%22ALU%22,%22AAPL%22)
// &env=http%3A%2F%2Fdatatables.org%2Falltables.env
// &format=json'
//
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
		return `</green><green>`
	} else {
		return `` // `</red><red>`
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
