package main

import (
	"os"
	"testing"

	"github.com/cmj0121/argparse"
)

func ExampleIFace() {
	argparse.Stderr = os.Stdout
	argparse.ExitWhenCallback = false

	c := IFace{}
	parser := argparse.MustNew(&c)
	parser.Parse("-h")
	// Output:
	// usage: iface [OPTION] IFACE
	//
	// option:
	//          -h, --help                  show this message
	//              --iface IFACE           network interface
	//
	// argument:
	//     IFACE                            network interface
}

func TestIFace(t *testing.T) {
	c := IFace{}
	parser := argparse.MustNew(&c)
	if err := parser.Parse("--iface", ""); err == nil {
		t.Fatalf("expect --iface failure")
	} else if err := parser.Parse("--iface", "abc"); err == nil {
		t.Fatalf("expect --iface abc failure")
	}
}
