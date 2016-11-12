// Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`github.com/michaeldv/termbox-go`
	`strings`
	`time`
)

// Screen is thin wrapper aroung Termbox library to provide basic display
// capabilities as requied by Mop.
type Screen struct {
	width    int        // Current number of columns.
	height   int        // Current number of rows.
	cleared  bool       // True after the screens gets cleared.
	layout   *Layout    // Pointer to layout (gets created by screen).
	markup   *Markup    // Pointer to markup processor (gets created by screen).
	pausedAt *time.Time // Timestamp of the pause request or nil if none.
}

// Initializes Termbox, creates screen along with layout and markup, and
// calculates current screen dimensions. Once initialized the screen is
// ready for display.
func NewScreen() *Screen {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	screen := &Screen{}
	screen.layout = NewLayout()
	screen.markup = NewMarkup()

	return screen.Resize()
}

// Close gets called upon program termination to close the Termbox.
func (screen *Screen) Close() *Screen {
	termbox.Close()

	return screen
}

// Resize gets called when the screen is being resized. It recalculates screen
// dimensions and requests to clear the screen on next update.
func (screen *Screen) Resize() *Screen {
	screen.width, screen.height = termbox.Size()
	screen.cleared = false

	return screen
}

// Pause is a toggle function that either creates a timestamp of the pause
// request or resets it to nil.
func (screen *Screen) Pause(pause bool) *Screen {
	if pause {
		screen.pausedAt = new(time.Time)
		*screen.pausedAt = time.Now()
	} else {
		screen.pausedAt = nil
	}

	return screen
}

// Clear makes the entire screen blank using default background color.
func (screen *Screen) Clear() *Screen {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.cleared = true

	return screen
}

// ClearLine erases the contents of the line starting from (x,y) coordinate
// till the end of the line.
func (screen *Screen) ClearLine(x int, y int) *Screen {
	for i := x; i < screen.width; i++ {
		termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()

	return screen
}

// Draw accepts variable number of arguments and knows how to display the
// market data, stock quotes, current time, and an arbitrary string.
func (screen *Screen) Draw(objects ...interface{}) *Screen {
	zonename, _ := time.Now().In(time.Local).Zone()
	if screen.pausedAt != nil {
		defer screen.DrawLine(0, 0, `<right><r>`+screen.pausedAt.Format(`3:04:05pm ` + zonename)+`</r></right>`)
	}
	for _, ptr := range objects {
		switch ptr.(type) {
		case *Market:
			object := ptr.(*Market)
			screen.draw(screen.layout.Market(object.Fetch()))
		case *Quotes:
			object := ptr.(*Quotes)
			screen.draw(screen.layout.Quotes(object.Fetch()))
		case time.Time:
			timestamp := ptr.(time.Time).Format(`3:04:05pm ` + zonename)
			screen.DrawLine(0, 0, `<right>`+timestamp+`</right>`)
		default:
			screen.draw(ptr.(string))
		}
	}

	return screen
}

// DrawLine takes the incoming string, tokenizes it to extract markup
// elements, and displays it all starting at (x,y) location.
func (screen *Screen) DrawLine(x int, y int, str string) {
	start, column := 0, 0

	for _, token := range screen.markup.Tokenize(str) {
		// First check if it's a tag. Tags are eaten up and not displayed.
		if screen.markup.IsTag(token) {
			continue
		}

		// Here comes the actual text: display it one character at a time.
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
	termbox.Flush()
}

// Underlying workhorse function that takes multiline string, splits it into
// lines, and displays them row by row.
func (screen *Screen) draw(str string) {
	if !screen.cleared {
		screen.Clear()
	}
	for row, line := range strings.Split(str, "\n") {
		screen.DrawLine(0, row, line)
	}
}
