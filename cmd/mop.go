// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/termbox-go`
	`github.com/mop/util`
	`github.com/mop/view`
	`time`
)

const help = `Mop v0.1.0 -- Copyright (c) 2013 Michael Dvorkin. All Rights Reserved.
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
func mainLoop(screen *view.Screen, profile *util.Profile) {
	var lineEditor *view.LineEditor
	var columnEditor *view.ColumnEditor

	keyboardQueue := make(chan termbox.Event)
	timestampQueue := time.NewTicker(1 * time.Second)
	quotesQueue := time.NewTicker(5 * time.Second)
	//marketQueue := time.NewTicker(12 * time.Second)
	emailQueue := time.NewTicker(10 * time.Second)
	showingHelp := false
	paused := false

	go func() {
		for {
			keyboardQueue <- termbox.PollEvent() // user input
		}
	}()

	market := new(util.Market).Initialize()
	quotes := new(util.Quotes).Initialize(market, profile)
	screen.Draw(quotes)

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
						lineEditor = new(view.LineEditor).Initialize(screen, quotes)
						lineEditor.Prompt(event.Ch)
					} else if event.Ch == 'o' || event.Ch == 'O' {
						columnEditor = new(view.ColumnEditor).Initialize(screen, quotes)
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
					screen.Clear().Draw(quotes)
				}
			case termbox.EventResize:
				screen.Resize()
				if !showingHelp {
					screen.Draw(quotes)
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

		// case <-marketQueue.C:
		// 	if !showingHelp && !paused {
		// 		screen.Draw(market)
		// 	}
		case <-emailQueue.C:
			util.SendMail(screen.GetQuoteLayout(quotes))
		}
	}
}

//-----------------------------------------------------------------------------
func main() {
	screen := new(view.Screen).Initialize()
	defer screen.Close()

	profile := new(util.Profile).Initialize()
	mainLoop(screen, profile)
}
