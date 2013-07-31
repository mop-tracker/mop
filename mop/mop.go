// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	`github.com/michaeldv/mop`
	`github.com/nsf/termbox-go`
	`time`
)

const help = `Mop v0.1.0 -- Copyright (c) 2013 Michael Dvorkin. All Rights Reserved.
NO WARRANTIES OF ANY KIND WHATSOEVER. USE AT YOUR OWN DISCRETION.

<u>Command</u>    <u>Description                                </u>
   +       Add stocks to the list.
   -       Remove stocks from the list.
   o       Change default sorting order.
   g       Group stocks by advancing/declining issues.
   ?       Display this help screen.
  esc      Quit mop.

<r> Press any key to continue </r>
`

//-----------------------------------------------------------------------------
func main_loop(screen *mop.Screen, profile *mop.Profile) {
	var line_editor *mop.LineEditor
	var column_editor *mop.ColumnEditor

	keyboard_queue := make(chan termbox.Event)
	timestamp_queue := time.NewTicker(1 * time.Second)
	quotes_queue := time.NewTicker(5 * time.Second)
	market_queue := time.NewTicker(12 * time.Second)
	showing_help := false

	go func() {
		for {
			keyboard_queue <- termbox.PollEvent()
		}
	}()

	market := new(mop.Market).Initialize()
	quotes := new(mop.Quotes).Initialize(market, profile)
	screen.Draw(market, quotes)

loop:
	for {
		select {
		case event := <-keyboard_queue:
			switch event.Type {
			case termbox.EventKey:
				if line_editor == nil && column_editor == nil && !showing_help {
					if event.Key == termbox.KeyEsc {
						break loop
					} else if event.Ch == '+' || event.Ch == '-' {
						line_editor = new(mop.LineEditor).Initialize(screen, quotes)
						line_editor.Prompt(event.Ch)
					} else if event.Ch == 'o' || event.Ch == 'O' {
						column_editor = new(mop.ColumnEditor).Initialize(screen, quotes)
					} else if event.Ch == 'g' || event.Ch == 'G' {
						profile.Regroup()
						screen.Draw(quotes)
					} else if event.Ch == '?' || event.Ch == 'h' || event.Ch == 'H' {
						showing_help = true
						screen.Clear().Draw(help)
					}
				} else if line_editor != nil {
					done := line_editor.Handle(event)
					if done {
						line_editor = nil
					}
				} else if column_editor != nil {
					done := column_editor.Handle(event)
					if done {
						column_editor = nil
					}
				} else if showing_help {
					showing_help = false
					screen.Clear().Draw(market, quotes)
				}
			case termbox.EventResize:
				screen.Resize()
				if !showing_help {
					screen.Draw(market, quotes)
				} else {
					screen.Draw(help)
				}
			}

		case <-timestamp_queue.C:
			if !showing_help {
				screen.DrawTime()
			}

		case <-quotes_queue.C:
			if !showing_help {
				screen.Draw(quotes)
			}

		case <-market_queue.C:
			if !showing_help {
				screen.Draw(market)
			}
		}
	}
}

//-----------------------------------------------------------------------------
func main() {
	screen := new(mop.Screen).Initialize()
	defer screen.Close()

	profile := new(mop.Profile).Initialize()
	main_loop(screen, profile)
}
