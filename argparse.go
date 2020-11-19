package argparse

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func MustNew(in interface{}) (parser *ArgParse) {
	var err error

	if parser, err = New(in); err != nil {
		// cannot new from pass interface, panic
		panic(err)
	}

	return
}

func New(in interface{}) (parser *ArgParse, err error) {
	Log(INFO, "new %[1]T", in)

	value := reflect.ValueOf(in)
	if value.Kind() != reflect.Ptr || !value.IsValid() {
		// invalid pass value
		err = fmt.Errorf("should pass *Struct: %T", in)
		return
	}

	// set the default program name as the pass structure with lowercase
	name := value.Elem().Type().Name()
	name = strings.ToLower(name)

	parser = &ArgParse{
		Value:  value,
		Name:   name,
		Stderr: os.Stderr,
	}

	return
}

type ArgParse struct {
	// the raw value pass to parser
	reflect.Value

	// the program name for the parser, default is the name of passed structure as lowercase
	Name string

	// IO for show the help message
	Stderr io.StringWriter
}

func (parser *ArgParse) Run() (err error) {
	if err = parser.Parse(os.Args[1:]...); err != nil {
		// show the help message
		parser.HelpMessage(err)
	}
	return
}

func (parser *ArgParse) Parse(args ...string) (err error) {
	Log(INFO, "parse %#v", args)
	return
}

func (parser *ArgParse) HelpMessage(err error) {
	msgs := []string{}

	if err != nil {
		msg := fmt.Sprintf("error: %v", err)
		msgs = append(msgs, msg)
	}

	msgs = append(msgs, fmt.Sprintf("usage: %v", parser.Name))

	msg := strings.Join(msgs, "\n") + "\n"
	parser.Stderr.WriteString(msg)
}
