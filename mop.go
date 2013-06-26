// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main

import (
	"github.com/michaeldv/mop/lib"
	"github.com/nsf/termbox-go"
)

//-----------------------------------------------------------------------------
func main() {
	profile := mop.LoadProfile()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

        mop.Draw(profile)
        mop.Refresh(profile)
}
