// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import `github.com/nsf/termbox-go`

// ColumnEditor handles column sort order. When activated it highlights
// current column name in the header, then waits for arrow keys (choose
// another column), Enter (reverse sort order), or Esc (exit).
type ColumnEditor struct {
	screen  *Screen  // Pointer to Screen so we could use screen.Draw().
	quotes  *Quotes  // Pointer to Quotes to redraw them when the sort order changes.
	layout  *Layout  // Pointer to Layout to redraw stock quotes header.
	profile *Profile // Pointer to Profile where we save newly selected sort order.
}

// Returns new initialized ColumnEditor struct. As part of initialization it
// highlights current column name (as stored in Profile).
func NewColumnEditor(screen *Screen, quotes *Quotes) *ColumnEditor {
	editor := &ColumnEditor{
		screen:  screen,
		quotes:  quotes,
		layout:  screen.layout,
		profile: quotes.profile,
	}

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
	editor.profile.selectedColumn = editor.profile.SortColumn
	editor.redrawHeader()
	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectLeftColumn() *ColumnEditor {
	editor.profile.selectedColumn--
	if editor.profile.selectedColumn < 0 {
		editor.profile.selectedColumn = editor.layout.TotalColumns() - 1
	}
	return editor
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectRightColumn() *ColumnEditor {
	editor.profile.selectedColumn++
	if editor.profile.selectedColumn > editor.layout.TotalColumns()-1 {
		editor.profile.selectedColumn = 0
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
	editor.profile.selectedColumn = -1
	return true
}

//-----------------------------------------------------------------------------
func (editor *ColumnEditor) redrawHeader() {
	editor.screen.DrawLine(0, 4, editor.layout.Header(editor.profile))
	termbox.Flush()
}
