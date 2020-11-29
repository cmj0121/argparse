package argparse

import (
	"os"
	"testing"
)

type Inner struct {
}

type Foo struct {
	Help

	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch   bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count    int    `short:"C" help:"save as the integer"`
	Name     string `name:"user-name"`
	Password string `args:"password"`

	// the pass argument
	Bind    *string `help:"pass the bind HOST:IP"`
	Timeout *int

	// the subcommand command
	*Inner `help:"sub-command"`
}

func TestArgParse(t *testing.T) {
	foo := Foo{
		Count: 12,
		Name:  "user",
	}
	parser := MustNew(&foo)

	// test toggle boolean
	if parser.Parse("-s"); foo.Switch != true {
		t.Fatalf("-s not work: %v", foo.Switch)
	} else if parser.Parse("-sss"); foo.Switch != false {
		t.Fatalf("-sss not work: %v", foo.Switch)
	} else if parser.Parse("-s", "-ss", "--toggle"); foo.Switch != false {
		t.Fatalf("-s -ss --toggle not work: %v", foo.Switch)
	}

	// test integer
	if err := parser.Parse("-C", "123"); err != nil || foo.Count != 123 {
		t.Fatalf("-C 123 not work: %v (%v)", foo.Count, err)
	} else if err := parser.Parse("--count", "22"); err != nil || foo.Count != 22 {
		t.Fatalf("--count not work: %v (%v)", foo.Count, err)
	} else if err := parser.Parse("--count", "333", "-C", "44"); err != nil || foo.Count != 44 {
		t.Fatalf("--count 333 -C 44 not work: %v (%v)", foo.Count, err)
	}

	// test string
	if err := parser.Parse("--user-name", "username II"); err != nil || foo.Name != "username II" {
		t.Fatalf("--username 'username II' not work: %v (%v)", foo.Name, err)
	}

	// test argument
	if err := parser.Parse(":9999"); err != nil || foo.Bind == nil || *foo.Bind != ":9999" {
		t.Fatalf(":9999 not work: %v (%v)", foo.Bind, err)
	} else if err := parser.Parse("98765"); err != nil || foo.Timeout == nil || *foo.Timeout != 98765 {
		t.Fatalf("98765 not work: %v (%v)", foo.Timeout, err)
	}
}

func Example() {
	foo := Foo{
		Count: 12,
		Name:  "user",
	}
	parser := MustNew(&foo)
	// show the message on the STDOUT, for testing
	parser.Stderr = os.Stdout
	parser.HelpMessage(nil)
	// Output:
	// usage: foo [OPTION] [BIND] [TIMEOUT]
	//
	// option:
	//         -h, --help                  show this message
	//         -s, --toggle                toggle the boolean value
	//     -C INT, --count INT             save as the integer (default: 12)
	//             --user-name STR         (default: user)
	//             --password STR
	//
	// argument:
	//     BIND                            pass the bind HOST:IP
	//     TIMEOUT
	//
	// sub-command:
	//     inner                           sub-command
}
