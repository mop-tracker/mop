// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"bytes"
	"text/template"
)

//-----------------------------------------------------------------------------
func Format(message []Message) string {
	markup := `{{range .}}<green>{{.Ticker}}</green> ${{.LastTrade}} <red>{{.Change}}</red>
{{end}}...`

	template, err := template.New("screen").Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, message)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
