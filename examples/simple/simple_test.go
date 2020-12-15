package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/cmj0121/argparse"
)

func ExampleSimple() {
	argparse.Stderr = os.Stdout

	c := Simple{}
	parser := argparse.MustNew(&c)
	parser.Parse("-h")
	// Output:
	// usage: simple [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//         -s, --toggle                toggle the boolean value
	//     -C INT, --count INT             save as the integer
	//             --user-name STR
	//     -c STR, --cases STR             choice from fix possible [demo foo]
	//             --now TIME
	//             --opt                   multiple option and save as array
}

func ExampleSimpleDefault() {
	argparse.Stderr = os.Stdout

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
	// usage: simple [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//         -s, --toggle                toggle the boolean value
	//     -C INT, --count INT             save as the integer (default: 123)
	//             --user-name STR         (default: simple)
	//     -c STR, --cases STR             choice from fix possible [demo foo] (default: demo)
	//             --now TIME
	//             --opt                   multiple option and save as array
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

	if err := parser.Parse("--opt", "a"); err != nil {
		t.Fatalf("cannot parse --opt a: %v", err)
	} else {
		if ans := []string{"a"}; !reflect.DeepEqual(ans, c.Optional) {
			t.Errorf("parse --opt a: %#v <> %#v", ans, c.Optional)
		}
	}

	if err := parser.Parse("--opt", "b", "--opt", "c"); err != nil {
		t.Fatalf("cannot parse --opt a: %v", err)
	} else {
		if ans := []string{"a", "b", "c"}; !reflect.DeepEqual(ans, c.Optional) {
			t.Errorf("parse --opt a: %#v <> %#v", ans, c.Optional)
		}
	}
}
