package main

import (
	"os"
	"testing"

	"github.com/cmj0121/argparse"
)

func ExampleFile() {
	argparse.Stderr = os.Stdout
	argparse.ExitWhenCallback = false

	c := File{}
	parser := argparse.MustNew(&c)
	parser.Parse("-h")
	// Output:
	// usage: file [OPTION] ACTION
	//
	// option:
	//          -h, --help                  show this message
	//          -v, --version               show argparse version
	//     -m PERM, --filemode PERM         file perm
	//     -c TIME, --created_at TIME       timestamp RFC-3339 (2006-01-02T15:04:05+07:00)
	//      -p STR, --path STR              file path list
	//
	// argument:
	//     ACTION                           action
	//
	// sub-command:
	//     fileaction                       sub-command 1
	//     sub                              sub-command 2
}

func TestFile(t *testing.T) {
	c := File{}
	parser := argparse.MustNew(&c)
	if err := parser.Parse("-c", "2020-01-02T11:22:33+07:00"); err != nil {
		t.Fatalf("cannot parse -c 2020-01-02T11:22:33+07:00: %v", err)
	}
}
