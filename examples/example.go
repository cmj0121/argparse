package main

import (
	"fmt"
	"time"

	"github.com/cmj0121/argparse"
)

type Inner struct {
	time.Time
}

type SubInner struct {
	X   int
	Y   byte
	Now time.Time `help:"timestamp"`
}

type Foo struct {
	argparse.Model

	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch   bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count    int    `short:"計" help:"save as the integer"`
	Name     string `name:"姓名" help:"save the username"`
	Cases    string `short:"c" choices:"demo foo" help:"choice from fix possible"`
	Password string `args:"password"`

	// the pass argument
	Bind    *string `name:"綁定" help:"pass the bind HOST:IP"`
	Timeout *int    `choices:"1 12 123"`
	SubInner

	// the subcommand command
	*Inner `help:"sub-command"`
}

func main() {
	foo := Foo{
		Switch: true,
		Count:  12,
		Cases:  "foo",
	}
	if err := argparse.MustNew(&foo).Run(); err == nil {
		fmt.Printf("%#v\n", foo)
	}
}
