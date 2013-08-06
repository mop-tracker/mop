// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	`github.com/michaeldv/termbox-go`
)

type ColumnEditor struct {
	screen     *Screen
	layout     *Layout
	quotes     *Quotes
	profile    *Profile
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) Initialize(screen *Screen, quotes *Quotes) *ColumnEditor {
	self.screen = screen
	self.quotes = quotes
	self.profile = quotes.profile
	self.layout = new(Layout).Initialize()

	self.select_current_column()
	return self
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) Handle(ev termbox.Event) bool {
	defer self.redraw_header()

	switch ev.Key {
	case termbox.KeyEsc:
		return self.done()

	case termbox.KeyEnter:
		self.execute()

        case termbox.KeyArrowLeft:
		self.select_left_column()

	case termbox.KeyArrowRight:
		self.select_right_column()
	}

	return false
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) select_current_column() *ColumnEditor {
	self.profile.selected_column = self.profile.SortColumn
	self.redraw_header()
	return self
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) select_left_column() *ColumnEditor {
	self.profile.selected_column--
	if self.profile.selected_column < 0 {
		self.profile.selected_column = TotalColumns - 1
	}
	return self
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) select_right_column() *ColumnEditor {
	self.profile.selected_column++
	if self.profile.selected_column > TotalColumns - 1 {
		self.profile.selected_column = 0
	}
	return self
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) execute() *ColumnEditor {
	if self.profile.Reorder() == nil {
		self.screen.Draw(self.quotes)
	}

	return self
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) done() bool {
	self.profile.selected_column = -1
	return true
}

//-----------------------------------------------------------------------------
func (self *ColumnEditor) redraw_header() {
	self.screen.DrawLine(0, 4, self.layout.Header(self.profile))
	termbox.Flush()
}

