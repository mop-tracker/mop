// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package view

import (
	`github.com/michaeldv/termbox-go`
	`regexp`
	`strings`
)

// Markup implements some minimalistic text formatting conventions that
// get translated to Termbox colors and attributes. To colorize a string
// wrap it in <color-name>...</> tags. Unlike HTML each tag sets a new
// color whereas the </> tag changes color back to default. For example:
//
// <green>Hello, <red>world!</>
//
// The color tags could be combined with the attributes: <b>...</b> for
// bold, <u>...</u> for underline, and <r>...</r> for reverse. Unlike
// colors the attributes require matching closing tag.
//
// The <right>...</right> tag is used to right align the enclosed string
// (ex. when displaying current time in the upper right corner).
type Markup struct {
	Foreground     termbox.Attribute	     // Foreground color.
	Background     termbox.Attribute	     // Background color (so far always termbox.ColorDefault).
	RightAligned   bool			     // True when the string is right aligned.
	tags	       map[string]termbox.Attribute  // Tags to Termbox translation hash.
	regex	      *regexp.Regexp		     // Regex to identify the supported tag names.
}

// Initialize creates tags to Termbox translation hash and sets default
// colors and string alignment.
func (markup *Markup) Initialize() *Markup {
	markup.Foreground      = termbox.ColorDefault
	markup.Background      = termbox.ColorDefault
	markup.RightAligned    = false

	markup.tags = make(map[string]termbox.Attribute)
	markup.tags[`/`]       = termbox.ColorDefault
	markup.tags[`black`]   = termbox.ColorBlack
	markup.tags[`red`]     = termbox.ColorRed
	markup.tags[`green`]   = termbox.ColorGreen
	markup.tags[`yellow`]  = termbox.ColorYellow
	markup.tags[`blue`]    = termbox.ColorBlue
	markup.tags[`magenta`] = termbox.ColorMagenta
	markup.tags[`cyan`]    = termbox.ColorCyan
	markup.tags[`white`]   = termbox.ColorWhite
	markup.tags[`right`]   = termbox.ColorDefault	 // Termbox can combine attributes and a single color using bitwise OR.
	markup.tags[`b`]       = termbox.AttrBold	 // Attribute = 1 << (iota + 4)
	markup.tags[`u`]       = termbox.AttrUnderline
	markup.tags[`r`]       = termbox.AttrReverse
	markup.regex           = markup.supportedTags()  // Once we have the hash we could build the regex.

	return markup
}

// Tokenize works just like strings.Split() except the resulting array includes
// the delimiters. For example, the "<green>Hello, <red>world!</>" string when
// tokenized by tags produces the following:
//
//   [0] "<green>"
//   [1] "Hello, "
//   [2] "<red>"
//   [3] "world!"
//   [4] "</>"
//
func (markup *Markup) Tokenize(str string) []string {
	matches := markup.regex.FindAllStringIndex(str, -1)
	strings := make([]string, 0, len(matches))

	head, tail := 0, 0
	for _, match := range matches {
		tail = match[0]
		if match[1] != 0 {
			if head != 0 || tail != 0 {
				// Apend the text between tags.
				strings = append(strings, str[head:tail])
			}
			// Append the tag itmarkup.
			strings = append(strings, str[match[0]:match[1]])
		}
		head = match[1]
	}

	if head != len(str) && tail != len(str) {
		strings = append(strings, str[head:])
	}

	return strings
}

// IsTag returns true when the given string looks like markup tag. When the
// tag name matches one of the markup-supported tags it gets translated to
// relevant Termbox attributes and colors.
func (markup *Markup) IsTag(str string) bool {
	tag, open := probeForTag(str)

	if tag == `` {
		return false
	}

	return markup.process(tag, open)
}

//-----------------------------------------------------------------------------
func (markup *Markup) process(tag string, open bool) bool {
	if attribute, ok := markup.tags[tag]; ok {
		switch tag {
		case `right`:
			markup.RightAligned = open  // On for <right>, off for </right>.
		default:
			if open {
				if attribute >= termbox.AttrBold {
					markup.Foreground |= attribute   // Set the Termbox attribute.
				} else {
					markup.Foreground = attribute	 // Set the Termbox color.
				}
			} else {
				if attribute >= termbox.AttrBold {
					markup.Foreground &= ^attribute	 // Clear the Termbox attribute.
				} else {
					markup.Foreground = termbox.ColorDefault
				}
			}
		}
	}

	return true
}

// supportedTags returns regular expression that matches all possible tags
// supported by the markup, i.e. </?black>|</?red>| ... |<?b>| ... |</?right>
func (markup *Markup) supportedTags() *regexp.Regexp {
	arr := []string{}

	for tag := range markup.tags {
		arr = append(arr, `</?` + tag + `>`)
	}

	return regexp.MustCompile(strings.Join(arr, `|`))
}

//-----------------------------------------------------------------------------
func probeForTag(str string) (string, bool) {
	if len(str) > 2 && str[0:1] == `<` && str[len(str)-1:] == `>` {
		return extractTagName(str), str[1:2] != `/`
	}

	return ``, false
}

// Extract tag name from the given tag, i.e. `<hello>` => `hello`.
func extractTagName(str string) string {
	if len(str) < 3 {
		return ``
	} else if str[1:2] != `/` {
		return str[1 : len(str)-1]
	} else if len(str) > 3 {
		return str[2 : len(str)-1]
	}

	return `/`
}
