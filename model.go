package argparse

import (
	"fmt"
	"os"
)

// the basic info for the default model
type Model struct {
	Help
	Version
}

type Help struct {
	// show the default help message
	ShowHelp bool `short:"h" name:"help" help:"show this message" callback:"_help"`
}

type Version struct {
	// show the version
	ShowVersion bool `short:"v" name:"version" help:"show argparse version" callback:"_version"`
}

// show the help message and exit
func (parser *ArgParse) defaultHelpMessage(in *ArgParse) (exit bool) {
	parser.HelpMessage(nil)
	exit = true
	return
}

func (parser *ArgParse) defaultVersionMessage(in *ArgParse) (exit bool) {
	os.Stdout.WriteString(fmt.Sprintf("%v (v%d.%d.%d)\n", PROJ_NAME, MAJOR, MINOR, MACRO))
	exit = true
	return
}
