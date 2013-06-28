// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const real_time_url = "http://finance.google.com/finance/info?client=ig&q="

// const body = `
// // [
// {
// "id": "22144"
// ,"t" : "AAPL"
// ,"e" : "NASDAQ"
// ,"l" : "393.78"
// ,"l_cur" : "393.78"
// ,"s": "2"
// ,"ltt":"4:00PM EDT"
// ,"lt" : "Jun 27, 4:00PM EDT"
// ,"c" : "-4.29"
// ,"cp" : "-1.08"
// ,"ccol" : "chr"
// ,"el": "393.40"
// ,"el_cur": "393.40"
// ,"elt" : "Jun 27, 5:04PM EDT"
// ,"ec" : "-0.38"
// ,"ecp" : "-0.10"
// ,"eccol" : "chr"
// ,"div" : "3.05"
// ,"yld" : "3.10"
// }
// ,{
// "id": "353353"
// ,"t" : "ATVI"
// ,"e" : "NASDAQ"
// ,"l" : "13.55"
// ,"l_cur" : "13.55"
// ,"s": "0"
// ,"ltt":"3:59PM EDT"
// ,"lt" : "Jun 21, 3:59PM EDT"
// ,"c" : "-0.33"
// ,"cp" : "-2.38"
// ,"ccol" : "chr"
// }
// ,{
// "id": "17154"
// ,"t" : "HPQ"
// ,"e" : "NYSE"
// ,"l" : "24.15"
// ,"l_cur" : "24.15"
// ,"s": "0"
// ,"ltt":"4:01PM EDT"
// ,"lt" : "Jun 21, 4:01PM EDT"
// ,"c" : "0.57"
// ,"cp" : "2.31"
// ,"ccol" : "chr"
// }
// ,{
// "id": "18241"
// ,"t" : "IBM"
// ,"e" : "NYSE"
// ,"l" : "195.46"
// ,"l_cur" : "195.46"
// ,"s": "0"
// ,"ltt":"4:02PM EDT"
// ,"lt" : "Jun 21, 4:02PM EDT"
// ,"c" : "-1.89"
// ,"cp" : "-0.96"
// ,"ccol" : "chr"
// }
// ]`

type Message struct {
	Ticker              string `json:"t"`
	Exchange            string `json:"e"`
	LastTrade           string `json:"l"`
	CurrentPrice        string `json:"l_cur"`
	LastTradeTime       string `json:"ltt"`
	LastTradeDateTime   string `json:"lt"`
	Change              string `json:"c"`
	ChangePercent       string `json:"cp"`
	ExLastTrade         string `json:"el"`
	ExCurrentPrice      string `json:"el_cur"`
	ExLastTradeDateTime string `json:"elt"`
	ExChange            string `json:"ec"`
	ExChangePercent     string `json:"ecp"`
	Dividend            string `json:"div"`
	Yield               string `json:"yld"`
}

var message []Message

func (m *Message) Color() string {
	if strings.Index(m.Change, "-") == -1 {
		return "</green><green>"
	} else {
		return "</red><red>"
	}
}

//-----------------------------------------------------------------------------
func Quote(ticker string) []Message {
	if len(message) > 0 && time.Now().Second()%5 != 0 { // Fetch quotes every 5 seconds.
		return message
	}

	// Send the request.
	response, err := http.Get(real_time_url + ticker)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	// Parse JSON.
	err = json.Unmarshal(sanitize(body), &message)

	// Parse JSON.
	// err := json.Unmarshal(sanitize([]byte(body)), &message)
	if err != nil {
		panic(err)
	}

	return message
}

//-----------------------------------------------------------------------------
func sanitize(ascii []byte) []byte {
	return bytes.Replace(ascii, []byte{'/'}, []byte{}, -1)
}

// func sanitize(str string) string {
//     r := strings.NewReplacer("//", "", "[", "", "]", "")
//     fmt.Printf("%s\n", []byte(r.Replace(str)))
//     return r.Replace(str)
// }
//
// func main() {
//
//     message := Quote("coh,atvi,hpq,ibm,xxx")
//     for _,m := range message {
//         fmt.Printf("%s, %s, %s\n", m.Ticker, m.LastTrade, m.Change)
//     }
// }
