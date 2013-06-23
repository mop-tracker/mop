// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
	"encoding/json"
	// "io/ioutil"
	// "net/http"
)

const real_time_url = "http://finance.google.com/finance/info?client=ig&q="
const body = `
// [
{
"id": "665300"
,"t" : "COH"
,"e" : "NYSE"
,"l" : "56.54"
,"l_cur" : "56.54"
,"s": "0"
,"ltt":"4:01PM EDT"
,"lt" : "Jun 21, 4:01PM EDT"
,"c" : "-0.75"
,"cp" : "-1.31"
,"ccol" : "chr"
}
,{
"id": "353353"
,"t" : "ATVI"
,"e" : "NASDAQ"
,"l" : "13.55"
,"l_cur" : "13.55"
,"s": "0"
,"ltt":"3:59PM EDT"
,"lt" : "Jun 21, 3:59PM EDT"
,"c" : "-0.33"
,"cp" : "-2.38"
,"ccol" : "chr"
}
,{
"id": "17154"
,"t" : "HPQ"
,"e" : "NYSE"
,"l" : "24.15"
,"l_cur" : "24.15"
,"s": "0"
,"ltt":"4:01PM EDT"
,"lt" : "Jun 21, 4:01PM EDT"
,"c" : "-0.57"
,"cp" : "-2.31"
,"ccol" : "chr"
}
,{
"id": "18241"
,"t" : "IBM"
,"e" : "NYSE"
,"l" : "195.46"
,"l_cur" : "195.46"
,"s": "0"
,"ltt":"4:02PM EDT"
,"lt" : "Jun 21, 4:02PM EDT"
,"c" : "-1.89"
,"cp" : "-0.96"
,"ccol" : "chr"
}
]`

type Message struct {
	Ticker    string `json:"t"`
	LastTrade string `json:"l"`
	Change    string `json:"c"`
}

func Quote(ticker string) []Message {
	// Send the request.
    // response, err := http.Get(real_time_url + ticker)
    // if err != nil {
    //     panic(err)
    // }
    //
    // // Fetch response and get its body.
    // defer response.Body.Close()
    // body, err := ioutil.ReadAll(response.Body)
    //
    // // Parse JSON.
    // var message []Message
    // err = json.Unmarshal(sanitize(body, &message)

	// Parse JSON.
	var message []Message
	err := json.Unmarshal(sanitize([]byte(body)), &message)
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
