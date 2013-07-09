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
	keyboard_queue := make(chan termbox.Event)
	quotes_queue := time.NewTicker(5 * time.Second)
	timestamp_queue := time.NewTicker(1 * time.Second)

	go func() {
		for {
			keyboard_queue <- termbox.PollEvent()
		}
	}()

	mop.Draw(profile)
loop:
	for {
		select {
		case event := <-keyboard_queue:
			switch event.Type {
			case termbox.EventKey:
				if event.Key == termbox.KeyEsc {
					break loop
				}
			case termbox.EventResize:
				mop.Draw(profile)
			}

		case <-quotes_queue.C:
			mop.Draw(profile)

		case <-timestamp_queue.C:
			mop.DrawTime()
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
