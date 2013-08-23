// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`github.com/michaeldv/termbox-go`
	`strings`
	`time`
)

// Screen ...
type Screen struct {
	width	  int
	height	  int
	cleared   bool
	layout   *Layout
	markup   *Markup
}

// Initialize ...
func (screen *Screen) Initialize() *Screen {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	screen.layout = new(Layout).Initialize()
	screen.markup = new(Markup).Initialize()

	return screen.Resize()
}

// Close ...
func (screen *Screen) Close() *Screen {
	termbox.Close()

	return screen
}

// Resize ...
func (screen *Screen) Resize() *Screen {
	screen.width, screen.height = termbox.Size()
	screen.cleared = false

	return screen
}

// Clear ...
func (screen *Screen) Clear() *Screen {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.cleared = true

	return screen
}

// ClearLine ...
func (screen *Screen) ClearLine(x int, y int) {
	for i := x; i < screen.width; i++ {
		termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}

// Draw ...
func (screen *Screen) Draw(objects ...interface{}) *Screen {
	for _, ptr := range objects {
		switch ptr.(type) {
		case *Market:
			object := ptr.(*Market)
			screen.draw(screen.layout.Market(object.Fetch()))
		case *Quotes:
			object := ptr.(*Quotes)
			screen.draw(screen.layout.Quotes(object.Fetch()))
		default:
			screen.draw(ptr.(string))
		}
	}

	return screen
}

// DrawLine ...
func (screen *Screen) DrawLine(x int, y int, str string) {
	start, column := 0, 0

	for _, token := range screen.markup.Tokenize(str) {
		if !screen.markup.IsTag(token) {
			for i, char := range token {
				if !screen.markup.RightAligned {
					start = x + column
					column++
				} else {
					start = screen.width - len(token) + i
				}
				termbox.SetCell(start, y, char, screen.markup.Foreground, screen.markup.Background)
			}
		}
	}
	termbox.Flush()
}

// DrawTime ...
func (screen *Screen) DrawTime() {
	now := time.Now().Format(`3:04:05pm PST`)
	screen.DrawLine(0, 0, `<right>` + now + `</right>`)
}


//-----------------------------------------------------------------------------
func (screen *Screen) draw(str string) {
	if !screen.cleared {
		screen.Clear()
	}
	for row, line := range strings.Split(str, "\n") {
		screen.DrawLine(0, row, line)
	}
}
