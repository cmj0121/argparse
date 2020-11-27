package argparse

// the basic info for the default model
type Help struct {
	// show the default help message
	ShowHelp bool `short:"h" name:"help" help:"show this message" callback:"_help"`
}

// show the help message and exit
func (parser *ArgParse) defaultHelpMessage(in *ArgParse) (exit bool) {
	parser.HelpMessage(nil)
	exit = true
	return
}
