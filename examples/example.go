package main

import (
	"time"

	"github.com/cmj0121/argparse"
)

type Inner struct {
	time.Time
}

type SubInner struct {
	X int
	Y byte
}

type Foo struct {
	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch                        bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count                         int    `short:"C" help:"save as the integer"`
	Name                          string `name:"user-name" help:"save the username"`
	Password                      string `args:"password"`
	TheFuckingLog0123456789ABCDEF int

	// the pass argument
	Bind    *string `help:"pass the bind HOST:IP"`
	Timeout *int
	SubInner

	// the subcommand command
	*Inner `help:"sub-command"`
}

func main() {
	foo := Foo{
		Switch: true,
		Count:  12,
	}
	argparse.MustNew(&foo).HelpMessage(nil)
}
