// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	`github.com/michaeldv/mop/lib`
	`github.com/nsf/termbox-go`
	`time`
)

//-----------------------------------------------------------------------------
func initTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

//-----------------------------------------------------------------------------
func mainLoop(profile *mop.Profile) {
	var line_editor *mop.LineEditor
	keyboard_queue := make(chan termbox.Event)
	timestamp_queue := time.NewTicker(1 * time.Second)
	quotes_queue := time.NewTicker(5 * time.Second)
	market_queue := time.NewTicker(12 * time.Second)

	go func() {
		for {
			keyboard_queue <- termbox.PollEvent()
		}
	}()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	mop.DrawMarket()
	mop.DrawQuotes(profile.Quotes())
loop:
	for {
		select {
		case event := <-keyboard_queue:
			switch event.Type {
			case termbox.EventKey:
				if line_editor == nil {
					if event.Key == termbox.KeyEsc {
						break loop
					} else if event.Ch == '+' || event.Ch == '-' {
						line_editor = new(mop.LineEditor)
						line_editor.Prompt(event.Ch, profile)
					}
				} else {
					done := line_editor.Handle(event)
					if done {
						line_editor = nil
					}
				}
			case termbox.EventResize:
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				mop.DrawMarket()
				mop.DrawQuotes(profile.Quotes())
			}

		case <-timestamp_queue.C:
			mop.DrawTime()

		case <-quotes_queue.C:
			mop.DrawQuotes(profile.Quotes())

		case <-market_queue.C:
			mop.DrawMarket()
		}
	}
}

//-----------------------------------------------------------------------------
func main() {

	initTermbox()
	defer termbox.Close()

	profile := new(mop.Profile).Initialize()
	mainLoop(profile)
}
