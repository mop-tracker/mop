// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	// "fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	// "strings"
)

const real_time_url = "http://finance.google.com/finance/info?client=ig&q="

//     body := `
// // [
// {
// "id": "22144"
// ,"t" : "AAPL"
// ,"e" : "NASDAQ"
// ,"l" : "416.84"
// ,"l_cur" : "416.84"
// ,"s": "2"
// ,"ltt":"4:00PM EDT"
// ,"lt" : "Jun 20, 4:00PM EDT"
// ,"c" : "-6.16"
// ,"cp" : "-1.46"
// ,"ccol" : "chr"
// ,"el": "416.74"
// ,"el_cur": "416.74"
// ,"elt" : "Jun 20, 7:16PM EDT"
// ,"ec" : "-0.10"
// ,"ecp" : "-0.02"
// ,"eccol" : "chr"
// ,"div" : "3.05"
// ,"yld" : "2.93"
// }
// ]`

type Message struct {
	Ticker    string `json:"t"`
	LastTrade string `json:"l"`
	Change    string `json:"c"`
}

func Quote(ticker string) []Message {
	// Send the request.
	response, err := http.Get(real_time_url + ticker)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	// Parse JSON.
	var message []Message
	err = json.Unmarshal(sanitize(body), &message)
	if err != nil {
		panic(err)
	}
	return message
}

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
