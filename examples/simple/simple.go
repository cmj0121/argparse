package main

import (
	"time"

	"github.com/cmj0121/argparse"
)

type Simple struct {
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
	Now    time.Time
}

func main() {
	c := Simple{
		Ignore: false,
		ignore: true,
	}
	parser := argparse.MustNew(&c)
	parser.Run()
}