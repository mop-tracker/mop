// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
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
		message,
	}

	markup := `Hello<right>{{.Now}}</right>
{{range .Stocks}}<green>{{.Ticker}}</green> ${{.LastTrade}} <red>{{.Change}}</red>
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
