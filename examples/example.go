package main

import (
	"fmt"

	"github.com/cmj0121/argparse"
)

type InnerX struct {
	InnerHide bool `-`
	Inner     bool
	// the pass argument
	Bind    *string `help:"pass the bind HOST:IP"`
	Timeout *int
}

type Foo struct {
	argparse.Help

	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count  int    `short:"C" help:"save as the integer"`
	Name   string `name:"user-name"`
	Cases  string `short:"c" choices:"demo foo" help:"choice from fix possible"`

	InnerX `help:"embedded and should not be display"`
}

func main() {
	foo := Foo{}
	if err := argparse.MustNew(&foo).Run(); err == nil {
		fmt.Printf("%#v\n", foo)
	}
}
