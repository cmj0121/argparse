package main

import (
	"net"
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
	//              --iface                 network interface
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
	} else if err := parser.Parse("--iface", "lo0"); err != nil {
		t.Fatalf("cannot parser --iface lo0: %v", err)
	} else {
		if c.IFace.Name != "lo0" {
			t.Errorf("expect c.IFace.Name = lo0: %v", c.IFace.Name)
		}

		if c.IFace.Flags != net.FlagLoopback|net.FlagUp|net.FlagMulticast {
			t.Errorf("expect c.IFace.Flags: %v", c.IFace.Flags)
		}
	}
}
