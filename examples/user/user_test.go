package main

import (
	"os"
	"testing"

	"github.com/cmj0121/argparse"
)

func ExampleUserConf() {
	argparse.Stderr = os.Stdout

	c := UserConf{
		Username: "root",
		Password: "password",
	}
	parser := argparse.MustNew(&c)
	parser.Name = "config"
	parser.Parse("-h")
	// Output:
	// usage: config [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//     -u STR, --username STR          username (default: root)
	//     -p STR, --password STR          password (default: password)
	//
	// sub-command:
	//     log                             access log
}

func ExampleUserConfLog() {
	argparse.Stderr = os.Stdout

	c := UserConf{
		Username: "root",
		Password: "password",
	}
	parser := argparse.MustNew(&c)
	parser.Name = "config"
	parser.Parse("log", "-h")
	// Output:
	// usage: log [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//             --last TIME
	//             --login
}

func ExampleUserConfDefault() {
	argparse.Stderr = os.Stdout

	c := UserConf{
		Username: "root",
		Password: "password",
		AccessLog: &AccessLog{
			Login: true,
		},
	}
	parser := argparse.MustNew(&c)
	parser.Name = "config"
	parser.Parse("-h")
	// Output:
	// usage: config [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//     -u STR, --username STR          username (default: root)
	//     -p STR, --password STR          password (default: password)
	//
	// sub-command:
	//     log                             access log
}

func ExampleUserConfLogDefault() {
	argparse.Stderr = os.Stdout

	c := UserConf{
		Username: "root",
		Password: "password",
		AccessLog: &AccessLog{
			Login: true,
		},
	}
	parser := argparse.MustNew(&c)
	parser.Name = "config"
	parser.Parse("log", "-h")
	// Output:
	// usage: log [OPTION]
	//
	// option:
	//         -h, --help                  show this message
	//             --last TIME
	//             --login                 (default: true)
}

func TestUserConf(t *testing.T) {
	c := UserConf{
		Username: "root",
		Password: "password",
		AccessLog: &AccessLog{
			Login: true,
		},
	}
	parser := argparse.MustNew(&c)
	if err := parser.Parse("-u", "abc", "-p", "qwertyu"); err != nil {
		t.Fatalf("cannot parse -u abc -p qwertyu: %v", err)
	} else {
		if c.Username != "abc" {
			t.Errorf("parse -u abc should change: %v", c.Username)
		}
		if c.Password != "qwertyu" {
			t.Errorf("parse -p qwertyu should change: %v", c.Password)
		}
	}

	if err := parser.Parse("log", "--last", "2020-12-31T23:59:59Z", "--login"); err != nil {
		t.Fatalf("parse log --last 2020-12-31T23:59:59Z --login: %v", err)
	} else {
		if c.AccessLog == nil {
			t.Errorf("missing the default value")
		} else {
			if c.AccessLog.Last.String() != "2020-12-31 23:59:59 +0000 UTC" {
				t.Errorf("parse --last 2020-12-31T23:59:59Z: %v", c.AccessLog.Last)
			}

			if c.AccessLog.Login == true {
				t.Errorf("parse --login should change: %v", c.AccessLog.Login)
			}
		}
	}
}
