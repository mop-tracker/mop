// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	"bytes"
	"fmt"
	"github.com/michaeldv/just"
	"github.com/michaeldv/mop/lib"
	"github.com/nsf/termbox-go"
	"regexp"
	"strings"
	"text/template"
)

// Can combine attributes and a single color using bitwise OR.
//
// AttrBold Attribute = 1 << (iota + 4)
// AttrUnderline
// AttrReverse
//
var colors = map[string]termbox.Attribute{
	"black":   termbox.ColorBlack,
	"red":     termbox.ColorRed,
	"green":   termbox.ColorGreen,
	"yellow":  termbox.ColorYellow,
	"blue":    termbox.ColorBlue,
	"magenta": termbox.ColorMagenta,
	"cyan":    termbox.ColorCyan,
	"white":   termbox.ColorWhite,
}

//
// Return regular expression that matches all possible color tags, i.e.
// </?black>|</?red>| ... |</?white>
//
func color_tags() *regexp.Regexp {
	tags := []string{}

	for color, _ := range colors {
		tags = append(tags, "</?"+color+">")
	}

	return regexp.MustCompile(strings.Join(tags, "|"))
}

//
// Return true if a string looks like a tag.
//
func is_tag(str string) (is bool, open bool) {
	is = (str[0:1] == "<" && str[len(str)-1:] == ">")
	open = (is && str[1:2] != "/")
	return
}

//
// Extract tag name from the given tag, i.e. "<hello>" => "hello"
//
func tag_name(str string) string {
	if str[1:2] != "/" {
		return str[1 : len(str)-1]
	} else {
		return str[2 : len(str)-1]
	}
}

func draw_color_line(x int, y int, str string) {
	column := 0
	foreground, background := termbox.ColorDefault, termbox.ColorDefault

	for _, token := range just.Split(color_tags(), str) {
		if tag, open := is_tag(token); tag {
			if color, ok := colors[tag_name(token)]; ok {
				token = ""
				if open {
					foreground = color
				} else {
					foreground = termbox.ColorDefault
				}
			}
		}

		for _, char := range token {
			termbox.SetCell(x+column, y, char, foreground, background)
			column += 1
		}
	}
}

func draw_screen(str string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for row, line := range strings.Split(str, "\n") {
		draw_color_line(0, row, line)
	}
	termbox.Flush()
}

func main() {
	message := mop.Quote("coh,atvi,hpq,ibm,xxx")
	for _, m := range message {
		fmt.Printf("%s, %s, %s\n", m.Ticker, m.LastTrade, m.Change)
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	markup := `<green>line 1</green>: <white>Hello</white> world
<green>line 2</green>: <white>Hello</white> again
<green>line 3</green>: <white>Hello</white> one more time :-)`

	template, err := template.New("screen").Parse(markup)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, nil)
	if err != nil {
		panic(err)
	}
	draw_screen(buffer.String())

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break loop
			}
		case termbox.EventResize:
			//x, y := termbox.Size()
			str := fmt.Sprintf("(%d:%d)", ev.Width, ev.Height)
			draw_screen(str + ": <red>Hello world</red>, how <white>are</white> <blue>you?</blue>")
		}
	}
}
