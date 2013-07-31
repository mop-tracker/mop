// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	`github.com/nsf/termbox-go`
	`strings`
	`time`
)

type Screen struct {
	width	 int
	height	 int
	cleared  bool
	markup   *Markup
}

//-----------------------------------------------------------------------------
func (self *Screen) Initialize() *Screen {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	self.markup = new(Markup).Initialize()

	return self.Resize()
}

//-----------------------------------------------------------------------------
func (self *Screen) Resize() *Screen {
	self.width, self.height = termbox.Size()
	self.cleared = false

	return self
}

//-----------------------------------------------------------------------------
func (self *Screen) Clear() *Screen {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	self.cleared = true

	return self
}

//-----------------------------------------------------------------------------
func (self *Screen) Close() *Screen {
	termbox.Close()

	return self
}

//-----------------------------------------------------------------------------
func (self *Screen) Draw(objects ...interface{}) {
	for _, ptr := range objects {
		switch ptr.(type) {
		case *Market:
			object := ptr.(*Market)
			self.draw(object.Fetch().Format())
		case *Quotes:
			object := ptr.(*Quotes)
			self.draw(object.Fetch().Format())
		default:
			self.draw(ptr.(string))
		}
	}
}

//-----------------------------------------------------------------------------
func (self *Screen) DrawTime() {
	now := time.Now().Format(`3:04:05pm PST`)
	self.DrawLine(0, 0, `<right>` + now + `</right>`)
}

//-----------------------------------------------------------------------------
func (self *Screen) ClearLine(x int, y int) {
	for i := x; i < self.width; i++ {
		termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}

//-----------------------------------------------------------------------------
func (self *Screen) DrawLine(x int, y int, str string) {
	start, column := 0, 0

	for _, token := range self.markup.Tokenize(str) {
		if !self.markup.IsTag(token) {
			for i, char := range token {
				if !self.markup.RightAligned {
					start = x + column
					column++
				} else {
					start = self.width - len(token) + i
				}
				termbox.SetCell(start, y, char, self.markup.Foreground, self.markup.Background)
			}
		}
	}
	termbox.Flush()
}

// private
//-----------------------------------------------------------------------------
func (self *Screen) draw(str string) {
	if !self.cleared {
		self.Clear()
	}
	for row, line := range strings.Split(str, "\n") {
		self.DrawLine(0, row, line)
	}
}
