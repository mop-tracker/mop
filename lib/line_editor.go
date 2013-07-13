// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
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
		     DrawLine(len(self.prompt), 3, self.input + " ")
		     termbox.SetCursor(len(self.prompt)+len(self.input), 3)
		     termbox.Flush()
             }
	default:
		if ev.Ch != 0 {
			self.input += string(ev.Ch)
			DrawLine(len(self.prompt), 3, self.input)
			termbox.SetCursor(len(self.prompt)+len(self.input), 3)
			termbox.Flush()
		}
	}
	return false
}
