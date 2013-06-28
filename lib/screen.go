// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"github.com/michaeldv/just"
	"github.com/nsf/termbox-go"
	"regexp"
	"strings"
)

// Can combine attributes and a single color using bitwise OR.
//
// AttrBold Attribute = 1 << (iota + 4)
// AttrUnderline
// AttrReverse
//
var tags = map[string]termbox.Attribute{
	"black":   termbox.ColorBlack,
	"red":     termbox.ColorRed,
	"green":   termbox.ColorGreen,
	"yellow":  termbox.ColorYellow,
	"blue":    termbox.ColorBlue,
	"magenta": termbox.ColorMagenta,
	"cyan":    termbox.ColorCyan,
	"white":   termbox.ColorWhite,
	"right":   termbox.ColorDefault,
}

//-----------------------------------------------------------------------------
func Draw(stocks string) {
	message := Quote(stocks)

	// for _, m := range message {
	//         fmt.Printf("%s, %s, %s\n", m.Ticker, m.LastTrade, m.Change)
	// }
	// fmt.Printf("%s\n", Format(message))

	drawScreen(Format(message))
}

//
// Return regular expression that matches all possible tags, i.e.
// </?black>|</?red>| ... |</?white>
//-----------------------------------------------------------------------------
func tagsRegexp() *regexp.Regexp {
	arr := []string{}

	for tag, _ := range tags {
		arr = append(arr, "</?"+tag+">")
	}

	return regexp.MustCompile(strings.Join(arr, "|"))
}

//
// Return true if a string looks like a tag.
//-----------------------------------------------------------------------------
func isTag(str string) (is bool, open bool) {
	is = (len(str) > 3 && str[0:1] == "<" && str[len(str)-1:] == ">")
	open = (is && str[1:2] != "/")
	return
}

//
// Extract tag name from the given tag, i.e. "<hello>" => "hello"
//-----------------------------------------------------------------------------
func tagName(str string) string {
	if len(str) < 3 {
		return ""
	} else if str[1:2] != "/" {
		return str[1 : len(str)-1]
	} else {
		return str[2 : len(str)-1]
	}
}

//-----------------------------------------------------------------------------
func drawLine(x int, y int, str string) {
	column, right := 0, false
	foreground, background := termbox.ColorDefault, termbox.ColorDefault

	for _, token := range just.Split(tagsRegexp(), str) {
		if tag, open := isTag(token); tag {
			key := tagName(token)
			if value, ok := tags[key]; ok {
				token = ""
				switch key {
				case "right":
					right = open
				default:
					if open {
						foreground = value
					} else {
						foreground = termbox.ColorDefault
					}
				}
			}
		}

		for i, char := range token {
			if !right {
				termbox.SetCell(x+column, y, char, foreground, background)
			} else {
				width, _ := termbox.Size()
				termbox.SetCell(width-len(token)+i, y, char, foreground, background)
			}
			column += 1
		}
	}
}

//-----------------------------------------------------------------------------
func drawScreen(str string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for row, line := range strings.Split(str, "\n") {
		drawLine(0, row, line)
	}
	termbox.Flush()
}

func DrawScreen(str string) {
	drawScreen(str)
}
