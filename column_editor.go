// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import `github.com/michaeldv/termbox-go`

// ColumnEditor handles column sort order. When activated it highlights
// current column name in the header, then waits for arrow keys (choose
// another column), Enter (reverse sort order), or Esc (exit).
type ColumnEditor struct {
	screen   *Screen	// Pointer to Screen so we could use screen.Draw().
	quotes   *Quotes	// Pointer to Quotes to redraw them when the sort order changes.
	profile	 *Profile	// Pointer to Profile where we save newly selected sort order.
}

// Initialize sets internal variables and highlights current column name
// (as stored in Profile).
func (editor *ColumnEditor) Initialize(screen *Screen, quotes *Quotes) *ColumnEditor {
	editor.screen = screen
	editor.quotes = quotes
	editor.profile = quotes.profile

	editor.selectCurrentColumn()

	return editor
}

// Handle takes over the keyboard events and dispatches them to appropriate
// column editor handlers. It returns true when user presses Esc.
func (editor *ColumnEditor) Handle(event termbox.Event) bool {
	defer editor.redrawHeader()

	switch event.Key {
	case termbox.KeyEsc:
		return editor.done()

	case termbox.KeyEnter:
		editor.execute()

        case termbox.KeyArrowLeft:
		editor.selectLeftColumn()

	case termbox.KeyArrowRight:
		editor.selectRightColumn()
	}

	return false
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectCurrentColumn() *ColumnEditor {
	editor.profile.selected_column = editor.profile.SortColumn
	editor.redrawHeader()
	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectLeftColumn() *ColumnEditor {
	editor.profile.selected_column--
	if editor.profile.selected_column < 0 {
		editor.profile.selected_column = TotalColumns - 1
	}
	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectRightColumn() *ColumnEditor {
	editor.profile.selected_column++
	if editor.profile.selected_column > TotalColumns - 1 {
		editor.profile.selected_column = 0
	}
	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) execute() *ColumnEditor {
	if editor.profile.Reorder() == nil {
		editor.screen.Draw(editor.quotes)
	}

	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) done() bool {
	editor.profile.selected_column = -1
	return true
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) redrawHeader() {
	editor.screen.DrawLine(0, 4, editor.screen.layout.Header(editor.profile))
	termbox.Flush()
}

