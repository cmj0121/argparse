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

func init() {
	// set the default callback
	RegisterCallback(FN_HELP, defaultHelpMessage)
	RegisterCallback(FN_VERSION, defaultVersionMessage)
}

// show the help message and exit
func defaultHelpMessage(in *ArgParse) (exit bool) {
	in.HelpMessage(nil)
	exit = true
	return
}

func defaultVersionMessage(in *ArgParse) (exit bool) {
	os.Stdout.WriteString(fmt.Sprintf("%v (v%d.%d.%d)\n", PROJ_NAME, MAJOR, MINOR, MACRO))
	exit = true
	return
}
