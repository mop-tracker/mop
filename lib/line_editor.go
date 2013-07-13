// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"github.com/nsf/termbox-go"
)

// const (
//         add_prompt = "Add tickers: "
//         remove_prompt = "Remove tickers: "
// )

// const prompts = map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}

type LineEditor struct {
	command rune
	prompt  string
	cursor  int
	input   string
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Prompt(command rune) {
	prompts := map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}

	self.command = command
	switch self.command {
	case '+', '-':
		self.prompt = prompts[self.command]
		// if self.command == '+' {
		//         self.prompt = add_prompt
		// } else {
		//         self.prompt = remove_prompt
		// }
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
