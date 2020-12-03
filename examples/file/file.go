package main

import (
	"os"

	"github.com/cmj0121/argparse"
)

type File struct {
	*os.File
	os.FileMode
}

func main() {
	file := File{}
	parser := argparse.MustNew(&file)
	parser.Run()
}
