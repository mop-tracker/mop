// Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/mop`
	`github.com/michaeldv/termbox-go`
	`time`
)

const help = `Mop v0.2.0 -- Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
NO WARRANTIES OF ANY KIND WHATSOEVER. SEE THE LICENSE FILE FOR DETAILS.

<u>Command</u>    <u>Description                                </u>
   +       Add stocks to the list.
   -       Remove stocks from the list.
   ?       Display this help screen.
   g       Group stocks by advancing/declining issues.
   o       Change column sort order.
   p       Pause market data and stock updates.
   q       Quit mop.
  esc      Ditto.

Enter comma-delimited list of stock tickers when prompted.

<r> Press any key to continue </r>
`

//-----------------------------------------------------------------------------
func mainLoop(screen *mop.Screen, profile *mop.Profile) {
	var lineEditor *mop.LineEditor
	var columnEditor *mop.ColumnEditor

	keyboardQueue := make(chan termbox.Event)
	timestampQueue := time.NewTicker(1 * time.Second)
	quotesQueue := time.NewTicker(5 * time.Second)
	marketQueue := time.NewTicker(12 * time.Second)
	showingHelp := false
	paused := false

	go func() {
		for {
			keyboardQueue <- termbox.PollEvent()
		}
	}()

	market := mop.NewMarket()
	quotes := mop.NewQuotes(market, profile)
	screen.Draw(market, quotes)

loop:
	for {
		select {
		case event := <-keyboardQueue:
			switch event.Type {
			case termbox.EventKey:
				if lineEditor == nil && columnEditor == nil && !showingHelp {
					if event.Key == termbox.KeyEsc || event.Ch == 'q' || event.Ch == 'Q' {
						break loop
					} else if event.Ch == '+' || event.Ch == '-' {
						lineEditor = mop.NewLineEditor(screen, quotes)
						lineEditor.Prompt(event.Ch)
					} else if event.Ch == 'o' || event.Ch == 'O' {
						columnEditor = mop.NewColumnEditor(screen, quotes)
					} else if event.Ch == 'g' || event.Ch == 'G' {
						if profile.Regroup() == nil {
							screen.Draw(quotes)
						}
					} else if event.Ch == 'p' || event.Ch == 'P' {
						paused = !paused
						screen.Pause(paused).Draw(time.Now())
					} else if event.Ch == '?' || event.Ch == 'h' || event.Ch == 'H' {
						showingHelp = true
						screen.Clear().Draw(help)
					}
				} else if lineEditor != nil {
					if done := lineEditor.Handle(event); done {
						lineEditor = nil
					}
				} else if columnEditor != nil {
					if done := columnEditor.Handle(event); done {
						columnEditor = nil
					}
				} else if showingHelp {
					showingHelp = false
					screen.Clear().Draw(market, quotes)
				}
			case termbox.EventResize:
				screen.Resize()
				if !showingHelp {
					screen.Draw(market, quotes)
				} else {
					screen.Draw(help)
				}
			}

		case <-timestampQueue.C:
			if !showingHelp && !paused {
				screen.Draw(time.Now())
			}

		case <-quotesQueue.C:
			if !showingHelp && !paused {
				screen.Draw(quotes)
			}

		case <-marketQueue.C:
			if !showingHelp && !paused {
				screen.Draw(market)
			}
		}
	}
}

//-----------------------------------------------------------------------------
func main() {
	screen := mop.NewScreen()
	defer screen.Close()

	profile := mop.NewProfile()
	mainLoop(screen, profile)
}
