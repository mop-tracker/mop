// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`regexp`
	`strings`
	`github.com/michaeldv/termbox-go`
)

type LineEditor struct {
	command  rune		// Keyboard command such as '+' or '-'.
	cursor   int 		// Current cursor position within the input line.
	prompt   string		// Prompt string for the command.
	input    string		// User typed input string.
	screen  *Screen		// Pointer to Screen.
	quotes  *Quotes		// Pointer to Quotes.
	regex   *regexp.Regexp	// Regex to split comma-delimited input string.
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Initialize(screen *Screen, quotes *Quotes) *LineEditor {
	self.screen = screen
	self.quotes = quotes
	self.regex = regexp.MustCompile(`[,\s]+`)

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Prompt(command rune) *LineEditor {
	prompts := map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}
	if prompt, ok := prompts[command]; ok {
		self.prompt = prompt
		self.command = command

		self.screen.DrawLine(0, 3, `<white>` + self.prompt + `</>`)
		termbox.SetCursor(len(self.prompt), 3)
		termbox.Flush()
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Handle(ev termbox.Event) bool {
	defer termbox.Flush()

	switch ev.Key {
	case termbox.KeyEsc:
		return self.done()

	case termbox.KeyEnter:
		return self.execute().done()

        case termbox.KeyBackspace, termbox.KeyBackspace2:
		self.delete_previous_character()

	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		self.move_left()

	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		self.move_right()

	case termbox.KeyCtrlA:
		self.jump_to_beginning()

	case termbox.KeyCtrlE:
		self.jump_to_end()

	case termbox.KeySpace:
		self.insert_character(' ')

	default:
		if ev.Ch != 0 {
			self.insert_character(ev.Ch)
		}
	}

	return false
}

//-----------------------------------------------------------------------------
func (self *LineEditor) delete_previous_character() *LineEditor {
	if self.cursor > 0 {
		if self.cursor < len(self.input) {
			// Remove character in the middle of the input string.
			self.input = self.input[0 : self.cursor-1] + self.input[self.cursor : len(self.input)]
		} else {
			// Remove last input character.
			self.input = self.input[ : len(self.input)-1]
		}
		self.screen.DrawLine(len(self.prompt), 3, self.input + ` `) // Erase last character.
		self.move_left()
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) insert_character(ch rune) *LineEditor {
	if self.cursor < len(self.input) {
		// Insert the character in the middle of the input string.
		self.input = self.input[0 : self.cursor] + string(ch) + self.input[self.cursor : len(self.input)]
	} else {
		// Append the character to the end of the input string.
		self.input += string(ch)
	}
	self.screen.DrawLine(len(self.prompt), 3, self.input)
	self.move_right()

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) move_left() *LineEditor {
	if self.cursor > 0 {
		self.cursor--
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) move_right() *LineEditor {
	if self.cursor < len(self.input) {
		self.cursor++
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) jump_to_beginning() *LineEditor {
	self.cursor = 0
	termbox.SetCursor(len(self.prompt) + self.cursor, 3)

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) jump_to_end() *LineEditor {
	self.cursor = len(self.input)
	termbox.SetCursor(len(self.prompt) + self.cursor, 3)

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) execute() *LineEditor {
	switch self.command {
	case '+':
		tickers := self.tokenize()
		if len(tickers) > 0 {
			if added,_ := self.quotes.AddTickers(tickers); added > 0 {
				self.screen.Draw(self.quotes)
			}
		}
	case '-':
		tickers := self.tokenize()
		if len(tickers) > 0 {
			before := len(self.quotes.profile.Tickers)
			if removed,_ := self.quotes.RemoveTickers(tickers); removed > 0 {
				self.screen.Draw(self.quotes)
				after := before - removed
				for i := before; i > after; i-- {
					self.screen.ClearLine(0, i + 4)
				}
			}
		}
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *LineEditor) done() bool {
	self.screen.ClearLine(0, 3)
	termbox.HideCursor()

	return true
}

//
// Split by whitespace/comma to convert a string to array of tickers. Make sure
// the string is trimmed to avoid empty tickers in the array.
//
func (self *LineEditor) tokenize() []string {
	input := strings.ToUpper(strings.Trim(self.input, `, `))
	return self.regex.Split(input, -1)
}
