package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/cmj0121/argparse"
)

func ExampleSimple() {
	argparse.Stderr = os.Stdout
	argparse.ExitWhenCallback = false

	c := Simple{}
	parser := argparse.MustNew(&c)
	parser.Parse("-h")
	// Output:
	// usage: simple [OPTION] PATH
	//
	// option:
	//          -h, --help                  show this message
	//          -v, --version               show argparse version
	//          -s, --toggle                toggle the boolean value
	//      -C INT, --count INT             save as the integer
	//              --user-name STR
	//      -c STR, --cases STR             choice from fix possible [demo foo]
	//
	// argument:
	//     PATH                             multi-argument
}

func ExampleSimpleDefault() {
	argparse.Stderr = os.Stdout
	argparse.ExitWhenCallback = false

	c := Simple{
		Ignore: true,
		ignore: true,

		Switch: false,
		Count:  123,
		Name:   "simple",
		Cases:  "demo",
	}
	parser := argparse.MustNew(&c)
	parser.Parse("-h")
	// Output:
	// usage: simple [OPTION] PATH
	//
	// option:
	//          -h, --help                  show this message
	//          -v, --version               show argparse version
	//          -s, --toggle                toggle the boolean value
	//      -C INT, --count INT             save as the integer (default: 123)
	//              --user-name STR         (default: simple)
	//      -c STR, --cases STR             choice from fix possible [demo foo] (default: demo)
	//
	// argument:
	//     PATH                             multi-argument
}

func TestSimple(t *testing.T) {
	c := Simple{
		Ignore: true,
		ignore: true,

		Switch: false,
		Count:  123,
		Name:   "simple",
		Cases:  "demo",
	}
	parser := argparse.MustNew(&c)
	if err := parser.Parse("-s"); err != nil {
		t.Fatalf("cannot parse -s: %v", err)
	} else {
		if c.Switch == false {
			t.Errorf("parse -s should change: %v", c.Switch)
		}
	}

	if err := parser.Parse("-C", "123", "--user-name", "username"); err != nil {
		t.Fatalf("cannot parse -C 123 --user-name username: %v", err)
	} else {
		if c.Count != 123 {
			t.Errorf("parse -C 123: %v", c.Count)
		}
		if c.Name != "username" {
			t.Errorf("parse --user-name username: %v", c.Name)
		}
	}

	if err := parser.Parse("-c", "abc"); err == nil {
		t.Fatalf("expe t -c abc should failure")
	} else if err := parser.Parse("-c", "foo"); err != nil {
		t.Fatalf("cannot parse -c abc: %v", err)
	} else {
		if c.Cases != "foo" {
			t.Errorf("parse -c foo: %v", c.Cases)
		}
	}

	if err := parser.Parse("x", "y", "z"); err != nil {
		t.Fatalf("cannot parse x y x : %v", err)
	} else {
		if ans := []string{"x", "y", "z"}; !reflect.DeepEqual(*c.Path, ans) {
			t.Errorf("parse x y z: %#v", c.Path)
		}
	}
}
