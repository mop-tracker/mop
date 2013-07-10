// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
        //"fmt"
        "bytes"
        "regexp"
        "strings"
        "net/http"
        "io/ioutil"
)

type Market struct {
        Dow       map[string]string
        Nasdaq    map[string]string
        Sp500     map[string]string
        Advances  map[string]string
        Declines  map[string]string
        Unchanged map[string]string
        Highs     map[string]string
        Lows      map[string]string
}

const yahoo_market_url = `http://finance.yahoo.com/marketupdate/overview`

func GetMarket() Market {
	response, err := http.Get(yahoo_market_url)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

        return extract(trim(body))
}

func trim(body []byte) []byte {
        start := bytes.Index(body, []byte("<table id=\"yfimktsumm\""))
        finish := bytes.LastIndex(body, []byte("<table id=\"yfimktsumm\""))
        snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
        snippet = bytes.Replace(snippet, []byte("&amp;"), []byte{'&'}, -1)
        
        return snippet
}

func extract(snippet []byte) Market {
        const any     = `\s*<.+?>`
        const some    = `<.+?`
        const space   = `\s*`
        const arrow   = `"(Up|Down)">\s*`
        const price   = `([\d\.,]+)`
        const percent = `\(([\d\.,%]+)\)`

        regex := []string{
                "(Dow)",       any, price, some, arrow, any, price, some, percent, any,
                "(Nasdaq)",    any, price, some, arrow, any, price, some, percent, any,
                "(S&P 500)",   any, price, some, arrow, any, price, some, percent, any,
                "(Advances)",  any, price, space, percent, any, price, space, percent, any,
                "(Declines)",  any, price, space, percent, any, price, space, percent, any,
                "(Unchanged)", any, price, space, percent, any, price, space, percent, any,
                "(New Hi's)",  any, price, any, price, any,
                "(New Lo's)",  any, price, any, price, any,
        }

        re := regexp.MustCompile(strings.Join(regex, ""))
        matches := re.FindAllStringSubmatch(string(snippet), -1)

        // if len(matches) > 0 {
        //         fmt.Printf("%d matches\n", len(matches[0]))
        //         for i, str := range matches[0][1:] {
        //                 fmt.Printf("%d) [%s]\n", i, str)
        //         }
        // } else {
        //         println("No matches")
        // }

        m := Market{
                Dow:       make(map[string]string),
                Nasdaq:    make(map[string]string),
                Sp500:     make(map[string]string),
                Advances:  make(map[string]string),
                Declines:  make(map[string]string),
                Unchanged: make(map[string]string),
                Highs:     make(map[string]string),
                Lows:      make(map[string]string),
        }
        m.Dow[`name`]          = matches[0][1]
        m.Dow[`latest`]        = matches[0][2]
        m.Dow[`change`]        = matches[0][4]
        if matches[0][3] == "Up" {
                m.Dow[`change`] = "+" + matches[0][4]
                m.Dow[`percent`] = "+" + matches[0][5]
        } else {
                m.Dow[`change`] = "-" + matches[0][4]
                m.Dow[`percent`] = "-" + matches[0][5]
        }

        m.Nasdaq[`name`]       = matches[0][6]
        m.Nasdaq[`latest`]     = matches[0][7]
        if matches[0][8] == "Up" {
                m.Nasdaq[`change`] = "+" + matches[0][9]
                m.Nasdaq[`percent`] = "+" + matches[0][10]
        } else {
                m.Nasdaq[`change`] = "-" + matches[0][9]
                m.Nasdaq[`percent`] = "-" + matches[0][10]
        }

        m.Sp500[`name`]        = matches[0][11]
        m.Sp500[`latest`]      = matches[0][12]
        if matches[0][13] == "Up" {
                m.Sp500[`change`] = "+" + matches[0][14]
                m.Sp500[`percent`] = "+" + matches[0][15]
        } else {
                m.Sp500[`change`] = "-" + matches[0][14]
                m.Sp500[`percent`] = "-" + matches[0][15]
        }

        m.Advances[`name`]     = matches[0][16]
        m.Advances[`nyse`]     = matches[0][17]
        m.Advances[`nysep`]    = matches[0][18]
        m.Advances[`nasdaq`]   = matches[0][19]
        m.Advances[`nasdaqp`]  = matches[0][20]

        m.Declines[`name`]     = matches[0][21]
        m.Declines[`nyse`]     = matches[0][22]
        m.Declines[`nysep`]    = matches[0][23]
        m.Declines[`nasdaq`]   = matches[0][24]
        m.Declines[`nasdaqp`]  = matches[0][25]

        m.Unchanged[`name`]    = matches[0][26]
        m.Unchanged[`nyse`]    = matches[0][27]
        m.Unchanged[`nysep`]   = matches[0][28]
        m.Unchanged[`nasdaq`]  = matches[0][29]
        m.Unchanged[`nasdaqp`] = matches[0][30]

        m.Highs[`name`]        = matches[0][31]
        m.Highs[`nyse`]        = matches[0][32]
        m.Highs[`nasdaq`]      = matches[0][33]
        m.Lows[`name`]         = matches[0][34]
        m.Lows[`nyse`]         = matches[0][35]
        m.Lows[`nasdaq`]       = matches[0][36]

        return m;
}
