// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

//-----------------------------------------------------------------------------
func Format(message []Message) string {
	vars := struct {
		Now    string
		Stocks []Message
	}{
		time.Now().Format("3:04:05pm PST"),
		prettify(message),
	}

	markup :=
		`Hello<right>{{.Now}}</right>

Ticker     Last trade     Change   % Change   Dividend      Yield
{{range .Stocks}}{{.Color}}{{.Ticker}} {{.LastTrade}} {{.Change}} {{.ChangePercent}} {{.Dividend}} {{.Yield}}
{{end}}...`

	template, err := template.New("screen").Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, vars)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

func color(m Message) string {
	return "x"
}

func prettify(message []Message) []Message {
	pretty := make([]Message, len(message))
	for i, m := range message {
		pretty[i].Ticker = pad(m.Ticker, -10)
		pretty[i].LastTrade = pad(with_currency(m.LastTrade), 10)
		pretty[i].CurrentPrice = pad(with_currency(m.CurrentPrice), 10)
		pretty[i].Change = pad(with_currency(m.Change), 10)
		pretty[i].ChangePercent = pad(with_percent(m.ChangePercent), 10)
		// ExLastTrade         string `json:"el"`
		// ExCurrentPrice      string `json:"el_cur"`
		// ExLastTradeDateTime string `json:"elt"`
		// ExChange            string `json:"ec"`
		// ExChangePercentage  string `json:"ecp"`
		pretty[i].Dividend = pad(with_currency(nullify(m.Dividend)), 10)
		pretty[i].Yield = pad(with_currency(nullify(m.Yield)), 10)
	}
	// fmt.Printf("%q", pretty)
	return pretty
}

func nullify(str string) string {
	if len(str) > 0 {
		return str
	} else {
		return "-"
	}
}

func with_currency(str string) string {
	if str == "-" {
		return str
	} else {
		switch str[0:1] {
		case "+", "-":
			return str[0:1] + "$" + str[1:]
		default:
			return "$" + str
		}
	}
}

func with_percent(str string) string {
	if str == "-" {
		return str
	} else if str[0:1] != "-" {
		return "+" + str + "%"
	} else {
		return str + "%"
	}
}

func colorize(str string) string {
	if str == "-" {
		return str
	} else if str[0:1] == "-" {
		return "<red>" + str + "</red>"
	} else {
		return "<green>" + str + "</green>"
	}
}

func ticker(str string, change string) string {
	if change[0:1] == "-" {
		return "<red>" + str + "</red>"
	} else {
		return "<green>" + str + "</green>"
	}
}

func pad(str string, width int) string {
	return fmt.Sprintf("%*s", width, str)
}
