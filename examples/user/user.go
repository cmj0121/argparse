package main

import (
	"time"

	"github.com/cmj0121/argparse"
)

type AccessLog struct {
	argparse.Help

	Last  time.Time
	Login bool
}

// for the small example
type UserConf struct {
	argparse.Help

	Username string `short:"u" help:"username"`
	Password string `short:"p" help:"password"`

	*AccessLog `name:"log" help:"access log"`
}

func main() {
	c := UserConf{
		Username: "root",
		Password: "password",
		AccessLog: &AccessLog{
			Login: true,
		},
	}
	parser := argparse.MustNew(&c)
	parser.Run()
}
