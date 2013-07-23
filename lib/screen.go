// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	`github.com/michaeldv/just`
	`github.com/nsf/termbox-go`
	`regexp`
	`strings`
	`time`
)

type Screen struct {
	width	int
	height	int
	cleared bool
	tags	map[string]termbox.Attribute
}

//-----------------------------------------------------------------------------
func (self *Screen) Initialize() *Screen {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	self.Resize()
	self.tags = make(map[string]termbox.Attribute)
	self.tags[`black`]   = termbox.ColorBlack
	self.tags[`red`]     = termbox.ColorRed
	self.tags[`green`]   = termbox.ColorGreen
	self.tags[`yellow`]  = termbox.ColorYellow
	self.tags[`blue`]    = termbox.ColorBlue
	self.tags[`magenta`] = termbox.ColorMagenta
	self.tags[`cyan`]    = termbox.ColorCyan
	self.tags[`white`]   = termbox.ColorWhite
	self.tags[`right`]   = termbox.ColorDefault	// Termbox can combine attributes and a single color using bitwise OR.
	self.tags[`b`]       = termbox.AttrBold		// Attribute = 1 << (iota + 4)
	self.tags[`u`]       = termbox.AttrUnderline
	self.tags[`r`]       = termbox.AttrReverse

	return self
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
func (self *Screen) Close() {
	termbox.Close()
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
	column, right := 0, false
	foreground, background := termbox.ColorDefault, termbox.ColorDefault

	for _, token := range just.Split(self.possible_tags(), str) {
		if tag, open := self.is_tag(token); tag {
			key := self.tag_name(token)
			if value, ok := self.tags[key]; ok {
				token = ``
				switch key {
				case `right`:
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
				termbox.SetCell(self.width-len(token)+i, y, char, foreground, background)
			}
			column += 1
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

//
// Return regular expression that matches all possible tags, i.e.
// </?black>|</?red>| ... |</?white>
//-----------------------------------------------------------------------------
func (self *Screen) possible_tags() *regexp.Regexp {
	arr := []string{}

	for tag, _ := range self.tags {
		arr = append(arr, `</?` + tag + `>`)
	}

	return regexp.MustCompile(strings.Join(arr, `|`))
}

//
// Return true if a string looks like a tag.
//-----------------------------------------------------------------------------
func (self *Screen) is_tag(str string) (is bool, open bool) {
	is = (len(str) > 2 && str[0:1] == `<` && str[len(str)-1:] == `>`)
	open = (is && str[1:2] != `/`)
	return
}

//
// Extract tag name from the given tag, i.e. `<hello>` => `hello`
//-----------------------------------------------------------------------------
func (self *Screen) tag_name(str string) string {
	if len(str) < 3 {
		return ``
	} else if str[1:2] != `/` {
		return str[1 : len(str)-1]
	} else {
		return str[2 : len(str)-1]
	}
}
