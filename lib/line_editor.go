// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
"fmt"
	"github.com/nsf/termbox-go"
)

type LineEditor struct {
	command rune
	prompt  string
	cursor  int
	input   string
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Prompt(command rune) {
	prompts := map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}
	if prompt, ok := prompts[command]; ok {
		self.prompt = prompt
		self.command = command

		DrawLine(0, 3, "<white>"+self.prompt+"</white>")
		termbox.SetCursor(len(self.prompt), 3)
		termbox.Flush()
	}
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Handle(ev termbox.Event) bool {
	switch ev.Key {
	case termbox.KeyEsc:
		ClearLine(0, 3)
		termbox.HideCursor()
		termbox.Flush()
		return true
	case termbox.KeyEnter:
		ClearLine(0, 3)
		termbox.HideCursor()
		termbox.Flush()
		return true
        case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(self.input) > 0 {
			self.input = self.input[:len(self.input)-1]
			self.cursor -= 1
			DrawLine(len(self.prompt), 3, self.input + " ")
			termbox.SetCursor(len(self.prompt) + self.cursor, 3)
			termbox.Flush()
		}
	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		if self.cursor > 0 {
			self.cursor -= 1
			termbox.SetCursor(len(self.prompt) + self.cursor, 3)
			termbox.Flush()
		}
	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		if self.cursor < len(self.input) {
			self.cursor += 1
			termbox.SetCursor(len(self.prompt) + self.cursor, 3)
			termbox.Flush()
		}
	case termbox.KeyCtrlA: // Jump to the beginning of line.
		self.cursor = 0
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
		termbox.Flush()
	case termbox.KeyCtrlE: // Jump to the end of line.
		self.cursor = len(self.input)
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
		termbox.Flush()
	case termbox.KeySpace:
		self.append_character(' ')
	default:
		if ev.Ch != 0 {
			self.append_character(ev.Ch)
		}
	}
	DrawLine(20,20, fmt.Sprintf("cursor: %02d [%s] %08d", self.cursor, self.input, ev.Ch))
	return false
}

//-----------------------------------------------------------------------------
func (self *LineEditor) append_character(ch rune) {
	self.input += string(ch)
	self.cursor += 1
	DrawLine(len(self.prompt), 3, self.input)
	termbox.SetCursor(len(self.prompt) + self.cursor, 3)
	termbox.Flush()
}

