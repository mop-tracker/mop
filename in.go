// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main
import (
        `fmt`
	`reflect`
)

type Formatter struct {
	entity interface{}
}

type Dog struct {
	name string
}

func (self *Formatter) Initialize(e interface{}) *Formatter {
	self.entity = e
	fmt.Printf("[%v]\n", reflect.TypeOf(e).String())
	return self
}

func main() {
	str := `hello`
	f1 := new(Formatter).Initialize(str)
	dog := new(Dog)
	dog.name = `Google`
	f2 := new(Formatter).Initialize(dog)

	fmt.Printf("[%v] [%v]\n", f1, f2)
}
