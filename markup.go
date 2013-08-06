// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`github.com/michaeldv/termbox-go`
	`regexp`
	`strings`
)

type Markup struct {
	tags	      map[string]termbox.Attribute
	Foreground    termbox.Attribute
	Background    termbox.Attribute
	RightAligned  bool
}

//-----------------------------------------------------------------------------
func (self *Markup) Initialize() *Markup {
	self.tags = make(map[string]termbox.Attribute)
	self.tags[`/`]       = termbox.ColorDefault
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
	self.Foreground      = termbox.ColorDefault
	self.Background      = termbox.ColorDefault
	self.RightAligned    = false

	return self
}

//-----------------------------------------------------------------------------
func (self *Markup) Tokenize(str string) []string {
	matches := self.supported_tags().FindAllStringIndex(str, -1)
	strings := make([]string, 0, len(matches))

	head, tail := 0, 0
	for _, match := range matches {
		tail = match[0]
		if match[1] != 0 {
			if head != 0 || tail != 0 {
				strings = append(strings, str[head:tail]) // Apend text between tags.
			}
			strings = append(strings, str[match[0]:match[1]]) // Append tag.
		}
		head = match[1]
	}

	if head != len(str) && tail != len(str) {
		strings = append(strings, str[head:])
	}

	return strings
}

//-----------------------------------------------------------------------------
func (self *Markup) IsTag(str string) bool {
	tag, open := probe_tag(str)

	if tag == `` {
		return false
	}

	return self.process(tag, open)
}

//-----------------------------------------------------------------------------
func (self *Markup) process(tag string, open bool) bool {
	if attribute, ok := self.tags[tag]; ok {
		switch tag {
		case `right`:
			self.RightAligned = open
		default:
			if open {
				if attribute >= termbox.AttrBold {
					self.Foreground |= attribute
				} else {
					self.Foreground = attribute
				}
			} else {
				if attribute >= termbox.AttrBold {
					self.Foreground &= ^attribute
				} else {
					self.Foreground = termbox.ColorDefault
				}
			}
		}
	}

	return true
}

//
// Return regular expression that matches all possible tags, i.e.
// </?black>|</?red>| ... |</?white>
//-----------------------------------------------------------------------------
func (self *Markup) supported_tags() *regexp.Regexp {
	arr := []string{}

	for tag, _ := range self.tags {
		arr = append(arr, `</?` + tag + `>`)
	}

	return regexp.MustCompile(strings.Join(arr, `|`))
}

//-----------------------------------------------------------------------------
func probe_tag(str string) (string, bool) {
	if len(str) > 2 && str[0:1] == `<` && str[len(str)-1:] == `>` {
		return extract_tag_name(str), str[1:2] != `/`
	}

	return ``, false
}

//
// Extract tag name from the given tag, i.e. `<hello>` => `hello`
//-----------------------------------------------------------------------------
func extract_tag_name(str string) string {
	if len(str) < 3 {
		return ``
	} else if str[1:2] != `/` {
		return str[1 : len(str)-1]
	} else if len(str) > 3 {
		return str[2 : len(str)-1]
	}

	return `/`
}
