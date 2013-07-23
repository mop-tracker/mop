// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	`github.com/michaeldv/mop/lib`
	`github.com/nsf/termbox-go`
	`time`
)

//-----------------------------------------------------------------------------
func mainLoop(screen *mop.Screen, profile *mop.Profile) {
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

	market := new(mop.Market).Initialize().Fetch()
	quotes := new(mop.Quotes).Initialize(market, profile)
	screen.Draw(market)
	screen.Draw(quotes)

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
						line_editor = new(mop.LineEditor).Initialize(screen, quotes)
						line_editor.Prompt(event.Ch)
					}
				} else {
					done := line_editor.Handle(event)
					if done {
						line_editor = nil
					}
				}
			case termbox.EventResize:
				screen.Resize()
				screen.Draw(market)
				screen.Draw(quotes)
			}

		case <-timestamp_queue.C:
			screen.DrawTime()

		case <-quotes_queue.C:
			screen.Draw(quotes)

		case <-market_queue.C:
			screen.Draw(market)
		}
	}
}

//-----------------------------------------------------------------------------
func main() {
	screen := new(mop.Screen).Initialize()
	defer screen.Close()

	profile := new(mop.Profile).Initialize()
	mainLoop(screen, profile)
}
