// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	"github.com/michaeldv/mop/lib"
	"github.com/nsf/termbox-go"
	"time"
)

//-----------------------------------------------------------------------------
func initTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

//-----------------------------------------------------------------------------
func mainLoop(profile string) {
	event_queue := make(chan termbox.Event)
	event_tick := time.NewTicker(1 * time.Second)

	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	mop.Draw(profile)
loop:
	for {
		select {
		case event := <-event_queue:
			switch event.Type {
			case termbox.EventKey:
				if event.Key == termbox.KeyEsc {
					break loop
				}
			case termbox.EventResize:
				mop.Draw(profile)
			}
		case <-event_tick.C:
			mop.Draw(profile)
		}
	}
}

//-----------------------------------------------------------------------------
func main() {

	initTermbox()
	defer termbox.Close()

	profile := mop.LoadProfile()
	mainLoop(profile)
}
