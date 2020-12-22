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
	//              --ip IP                 IP address
	//              --inet CIDR             IP with mask
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

	if err := parser.Parse("--ip", "127.0.0.1"); err != nil {
		t.Fatalf("cannot parse --ip 127.0.0.1: %v", err)
	}

	if err := parser.Parse("--ip", "github.com"); err != nil {
		t.Fatalf("cannot parse --ip github.com: %v", err)
	}

	if err := parser.Parse("--inet", "192.168.1.2/24"); err != nil {
		t.Fatalf("cannot parse --inet 192.168.1.2/24: %v", err)
	}

	if err := parser.Parse("--inet", "github.com"); err != nil {
		t.Fatalf("cannot parse --inet github.com: %v", err)
	}

	if err := parser.Parse("--inet", "github.com/16"); err != nil {
		t.Fatalf("cannot parse --inet github.com/16: %v", err)
	}
}
