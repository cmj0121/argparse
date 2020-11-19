package argparse

import (
	"os"
)

type Inner struct {
}

type Foo struct {
	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch   bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count    int    `short:"C" help:"save as the integer"`
	Name     string `name:"user-name" help:"save the username"`
	Password string `args:"password"`

	// the pass argument
	Bind    *string `help:"pass the bind HOST:IP"`
	Timeout *int

	// the subcommand command
	*Inner `help:"sub-command"`
}

func ExampleArgParse() {
	foo := Foo{
		Switch: true,
		Count:  12,
	}
	parser := MustNew(&foo)
	// show the message on the STDOUT, for testing
	parser.Stderr = os.Stdout
	parser.HelpMessage(nil)
	// Output:
	// usage: foo
}
