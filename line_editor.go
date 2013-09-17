// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`github.com/michaeldv/termbox-go`
	`regexp`
	`strings`
)

// LineEditor kicks in when user presses '+' or '-' to add or delete stock
// tickers. The data structure and methods are used to collect the input
// data and keep track of cursor movements (left, right, beginning of the
// line, end of the line, and backspace).
type LineEditor struct {
	command rune           // Keyboard command such as '+' or '-'.
	cursor  int            // Current cursor position within the input line.
	prompt  string         // Prompt string for the command.
	input   string         // User typed input string.
	screen  *Screen        // Pointer to Screen.
	quotes  *Quotes        // Pointer to Quotes.
	regex   *regexp.Regexp // Regex to split comma-delimited input string.
}

// Initialize sets internal pointers and compiles the regular expression.
func (editor *LineEditor) Initialize(screen *Screen, quotes *Quotes) *LineEditor {
	editor.screen = screen
	editor.quotes = quotes
	editor.regex = regexp.MustCompile(`[,\s]+`)

	return editor
}

// Prompt displays a prompt in response to '+' or '-' commands. Unknown commands
// are simply ignored. The prompt is displayed on the 3rd line (between the market
// data and the stock quotes).
func (editor *LineEditor) Prompt(command rune) *LineEditor {
	prompts := map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}
	if prompt, ok := prompts[command]; ok {
		editor.prompt = prompt
		editor.command = command

		editor.screen.DrawLine(0, 3, `<white>`+editor.prompt+`</>`)
		termbox.SetCursor(len(editor.prompt), 3)
		termbox.Flush()
	}

	return editor
}

// Handle takes over the keyboard events and dispatches them to appropriate
// line editor handlers. As user types or edits the text cursor movements
// are tracked in `editor.cursor` while the text itself is stored in
// `editor.input`. The method returns true when user presses Esc (discard)
// or Enter (process).
func (editor *LineEditor) Handle(ev termbox.Event) bool {
	defer termbox.Flush()

	switch ev.Key {
	case termbox.KeyEsc:
		return editor.done()

	case termbox.KeyEnter:
		return editor.execute().done()

	case termbox.KeyBackspace, termbox.KeyBackspace2:
		editor.deletePreviousCharacter()

	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		editor.moveLeft()

	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		editor.moveRight()

	case termbox.KeyCtrlA:
		editor.jumpToBeginning()

	case termbox.KeyCtrlE:
		editor.jumpToEnd()

	case termbox.KeySpace:
		editor.insertCharacter(' ')

	default:
		if ev.Ch != 0 {
			editor.insertCharacter(ev.Ch)
		}
	}

	return false
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) deletePreviousCharacter() *LineEditor {
	if editor.cursor > 0 {
		if editor.cursor < len(editor.input) {
			// Remove character in the middle of the input string.
			editor.input = editor.input[0:editor.cursor-1] + editor.input[editor.cursor:len(editor.input)]
		} else {
			// Remove last input character.
			editor.input = editor.input[:len(editor.input)-1]
		}
		editor.screen.DrawLine(len(editor.prompt), 3, editor.input+` `) // Erase last character.
		editor.moveLeft()
	}

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) insertCharacter(ch rune) *LineEditor {
	if editor.cursor < len(editor.input) {
		// Insert the character in the middle of the input string.
		editor.input = editor.input[0:editor.cursor] + string(ch) + editor.input[editor.cursor:len(editor.input)]
	} else {
		// Append the character to the end of the input string.
		editor.input += string(ch)
	}
	editor.screen.DrawLine(len(editor.prompt), 3, editor.input)
	editor.moveRight()

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) moveLeft() *LineEditor {
	if editor.cursor > 0 {
		editor.cursor--
		termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)
	}

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) moveRight() *LineEditor {
	if editor.cursor < len(editor.input) {
		editor.cursor++
		termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)
	}

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) jumpToBeginning() *LineEditor {
	editor.cursor = 0
	termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) jumpToEnd() *LineEditor {
	editor.cursor = len(editor.input)
	termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) execute() *LineEditor {
	switch editor.command {
	case '+':
		tickers := editor.tokenize()
		if len(tickers) > 0 {
			if added, _ := editor.quotes.AddTickers(tickers); added > 0 {
				editor.screen.Draw(editor.quotes)
			}
		}
	case '-':
		tickers := editor.tokenize()
		if len(tickers) > 0 {
			before := len(editor.quotes.profile.Tickers)
			if removed, _ := editor.quotes.RemoveTickers(tickers); removed > 0 {
				editor.screen.Draw(editor.quotes)

				// Clear the lines at the bottom of the list, if any.
				after := before - removed
				for i := before; i > after; i-- {
					editor.screen.ClearLine(0, i+4)
				}
			}
		}
	}

	return editor
}

//-----------------------------------------------------------------------------
func (editor *LineEditor) done() bool {
	editor.screen.ClearLine(0, 3)
	termbox.HideCursor()

	return true
}

// Split by whitespace/comma to convert a string to array of tickers. Make sure
// the string is trimmed to avoid empty tickers in the array.
func (editor *LineEditor) tokenize() []string {
	input := strings.ToUpper(strings.Trim(editor.input, `, `))
	return editor.regex.Split(input, -1)
}
