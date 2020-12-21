package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cmj0121/argparse"
)

type FileAction struct {
	Action *string
}

type File struct {
	argparse.Model

	os.FileMode `short:"m"`
	CreatedAt   time.Time `short:"c" name:"created_at"`
	Path        []string  `short:"p" name:"path" help:"file path list"`

	Action *string `help:"action"`

	*FileAction `help:"sub-command 1"`
	Sub         *FileAction `help:"sub-command 2"`
}

func main() {
	c := File{}
	parser := argparse.MustNew(&c)
	if err := parser.Run(); err == nil {
		data, _ := json.MarshalIndent(c, "", "    ")
		fmt.Println(string(data))
	}
}
