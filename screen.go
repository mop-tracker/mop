// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

// Screen is thin wrapper around Termbox library to provide basic display
// capabilities as required by Mop.
type Screen struct {
	width      int        // Current number of columns.
	height     int        // Current number of rows.
	cleared    bool       // True after the screens gets cleared.
	layout     *Layout    // Pointer to layout (gets created by screen).
	markup     *Markup    // Pointer to markup processor (gets created by screen).
	pausedAt   *time.Time // Timestamp of the pause request or nil if none.
	offset     int        // Offset for scolling
	headerLine int        // Line number of header for scroll feature
	max        int        // highest offset
}

// Initializes Termbox, creates screen along with layout and markup, and
// calculates current screen dimensions. Once initialized the screen is
// ready for display.
func NewScreen(profile *Profile) *Screen {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	screen := &Screen{}
	screen.layout = NewLayout()
	screen.markup = NewMarkup(profile)
	screen.offset = 0

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

// Increase the offset for scrolling feature by n
// Takes number of tickers as max, so not scrolling down forever
func (screen *Screen) IncreaseOffset(n int) {
	if screen.offset+n <= screen.max {
		screen.offset += n
	} else if screen.max > screen.height {
		screen.offset = screen.max
	}
}

// Decrease the offset for scrolling feature by n
func (screen *Screen) DecreaseOffset(n int) {
	if screen.offset > n {
		screen.offset -= n
	} else {
		screen.offset = 0
	}
}

func (screen *Screen) ScrollTop() {
	screen.offset = 0
}

func (screen *Screen) ScrollBottom() {
	if screen.max > screen.height {
		screen.offset = screen.max
	}
}

func (screen *Screen) DrawOldQuotes(quotes *Quotes) {
	screen.draw(screen.layout.Quotes(quotes), true)
	termbox.Flush()
}

func (screen *Screen) DrawOldMarket(market *Market) {
	screen.draw(screen.layout.Market(market), false)
	termbox.Flush()
}

// Draw accepts variable number of arguments and knows how to display the
// market data, stock quotes, current time, and an arbitrary string.
func (screen *Screen) Draw(objects ...interface{}) *Screen {
	zonename, _ := time.Now().In(time.Local).Zone()
	if screen.pausedAt != nil {
		defer screen.DrawLine(0, 0, `<right><r>`+screen.pausedAt.Format(`3:04:05pm `+zonename)+`</r></right>`)
	}
	for _, ptr := range objects {
		switch ptr.(type) {
		case *Market:
			object := ptr.(*Market)
			screen.draw(screen.layout.Market(object.Fetch()), false)
		case *Quotes:
			object := ptr.(*Quotes)
			screen.draw(screen.layout.Quotes(object.Fetch()), true)
		case time.Time:
			timestamp := ptr.(time.Time).Format(`3:04:05pm ` + zonename)
			screen.DrawLineInverted(0, 0, `<right><time>`+timestamp+`</></right>`)
		default:
			screen.draw(ptr.(string), false)
		}
	}

	termbox.Flush()

	return screen
}

// DrawLine takes the incoming string, tokenizes it to extract markup
// elements, and displays it all starting at (x,y) location.

// DrawLineFlush gives the option to flush screen after drawing

// wrapper for DrawLineFlush
func (screen *Screen) DrawLine(x int, y int, str string) {
	screen.DrawLineFlush(x, y, str, true)
}

func (screen *Screen) DrawLineInverted(x int, y int, str string) {
	screen.DrawLineFlushInverted(x, y, str, true)
}

func (screen *Screen) DrawLineFlush(x int, y int, str string, flush bool) {
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
	if flush {
		termbox.Flush()
	}
}

func (screen *Screen) DrawLineFlushInverted(x int, y int, str string, flush bool) {
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
			termbox.SetCell(start, y, char, screen.markup.tags[`black`], screen.markup.Foreground)
		}
	}
	if flush {
		termbox.Flush()
	}
}

// Underlying workhorse function that takes multiline string, splits it into
// lines, and displays them row by row.
func (screen *Screen) draw(str string, offset bool) {
	if !screen.cleared {
		screen.Clear()
	}
	var allLines []string
	drewHeading := false

	screen.width, screen.height = termbox.Size()

	tempFormat := "%" + strconv.Itoa(screen.width) + "s"
	blankLine := fmt.Sprintf(tempFormat, "")
	allLines = strings.Split(str, "\n")

	if offset {
		screen.max = len(allLines) - screen.height + screen.headerLine
	}

	// Write the lines being updated.
	for row := 0; row < len(allLines); row++ {
		if offset {
			// Did we draw the underlined heading row?  This is a crude
			// check, but--see comments below...
			// --- Heading row only appears for quotes, so offset is true
			if !drewHeading {
				if strings.Contains(allLines[row], "Ticker") &&
					strings.Contains(allLines[row], "Last") &&
					strings.Contains(allLines[row], "Change") {
					drewHeading = true
					screen.headerLine = row
					screen.DrawLine(0, row, allLines[row])
					// move on to the point to offset to
					row += screen.offset
				}
			} else {
				// only write the necessary lines
				if row <= len(allLines) &&
					row > screen.headerLine {
					screen.DrawLineFlush(0, row-screen.offset, allLines[row], false)
				} else if row > len(allLines)+1 {
					row = len(allLines)
				}
			}
		} else {
			screen.DrawLineFlush(0, row, allLines[row], false)
		}
	}
	// If the quotes lines in this cycle are shorter than in the previous
	// cycles, e.g., because a filter was just applied, then one or more
	// lines from the previous cycles will not be cleared.  Since the
	// incoming lines don't mark explicitly whether they are part of the
	// market summary or quotes, we can't check whether quotes were updated
	// in a way that is robust for code changes.  This is a simple test: if
	// we drew the heading row ("Ticker Last Change..."), then we are
	// updating the quotes section in this cycle, and we should pad the
	// quotes section with blank lines.  If we didn't draw the heading row,
	// then we probably only updated the market summary at the top in this
	// cycle.  In that case, padding with blank lines would overwrite the
	// stocks list.)
	if drewHeading {
		for i := len(allLines) - 1 - screen.offset; i < screen.height; i++ {
			if i > screen.headerLine {
				screen.DrawLine(0, i, blankLine)
			}
		}
	}
}
