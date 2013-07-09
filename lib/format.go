// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"fmt"
	"time"
	"bytes"
        "regexp"
        "strings"
	"text/template"
)

//-----------------------------------------------------------------------------
func Format(quotes Quotes) string {
	vars := struct {
		Now    string
                Header string
		Stocks Quotes
	}{
		time.Now().Format("3:04:05pm PST"),
                header(),
		prettify(quotes),
	}

	markup :=
		`Hello<right><white>{{.Now}}</white></right>

{{.Header}}
{{range .Stocks}}{{.Color}}{{.Ticker}} {{.LastTrade}} {{.Change}} {{.ChangePercent}} {{.Open}} {{.Low}} {{.High}} {{.Low52}} {{.High52}} {{.Volume}} {{.AvgVolume}} {{.PeRatio}} {{.Dividend}} {{.Yield}} {{.MarketCap}}
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

func header() string {
        str := fmt.Sprintf("%-7s ", "Ticker")
        str += fmt.Sprintf("%9s ",  "Last")
        str += fmt.Sprintf("%9s ",  "Change")
        str += fmt.Sprintf("%9s ",  "%Change")
        str += fmt.Sprintf("%9s ",  "Open")
        str += fmt.Sprintf("%9s ",  "Low")
        str += fmt.Sprintf("%9s ",  "High")
        str += fmt.Sprintf("%9s ",  "52w Low")
        str += fmt.Sprintf("%9s ",  "52w High")
        str += fmt.Sprintf("%10s ", "Volume")
        str += fmt.Sprintf("%10s ", "AvgVolume")
        str += fmt.Sprintf("%9s ",  "P/E")
        str += fmt.Sprintf("%9s ",  "Dividend")
        str += fmt.Sprintf("%9s ",  "Yield")
        str += fmt.Sprintf("%10s",  "MktCap")

        return str
}

func prettify(quotes Quotes) Quotes {
	pretty := make(Quotes, len(quotes))
	for i, q := range quotes {
		pretty[i].Ticker        = pad(q.Ticker, -7)
		pretty[i].LastTrade     = pad(with_currency(q.LastTrade), 9)
		pretty[i].Change        = pad(with_currency(q.Change), 9)
		pretty[i].ChangePercent = pad(last_of_pair(q.ChangePercent), 9)
		pretty[i].Open          = pad(with_currency(q.Open), 9)
		pretty[i].Low           = pad(with_currency(q.Low), 9)
		pretty[i].High          = pad(with_currency(q.High), 9)
		pretty[i].Low52         = pad(with_currency(q.Low52), 9)
		pretty[i].High52        = pad(with_currency(q.High52), 9)
		pretty[i].Volume        = pad(q.Volume, 10)
		pretty[i].AvgVolume     = pad(q.AvgVolume, 10)
		pretty[i].PeRatio       = pad(nullify(q.PeRatioX), 9)
		pretty[i].Dividend      = pad(with_currency(q.Dividend), 9)
		pretty[i].Yield         = pad(with_percent(q.Yield), 9)
		pretty[i].MarketCap     = pad(with_currency(q.MarketCapX), 10)
	}
	return pretty
}

func nullify(str string) string {
        if len(str) == 3 && str[0:3] == "N/A" {
                return "-"
        } else {
                return str
        }
}

func last_of_pair(str string) string {
        if len(str) >= 6 && str[0:6] != "N/A - " {
		return str
        } else {
                return str[6:]
        }
}

func with_currency(str string) string {
	if str == "N/A" || str == "0.00" {
		return "-"
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
	if str == "N/A" {
		return "-"
	} else {
		return str + "%"
	}
}

func colorize(str string) string {
	if str == "N/A" {
		return "-"
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
        re := regexp.MustCompile(`(\.\d+)[MB]?$`)
        match := re.FindStringSubmatch(str)
        if len(match) > 0 {
                switch len(match[1]) {
                case 2: str = strings.Replace(str, match[1], match[1] + "0", 1)
                case 4, 5: str = strings.Replace(str, match[1], match[1][0:3], 1)
                }
        }

	return fmt.Sprintf("%*s", width, str)
}
