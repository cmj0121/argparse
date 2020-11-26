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
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct || !value.IsValid() {
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

	// process the field
	typ := value.Elem().Type()
	Log(VERBOSE, "start process: %v", typ)
	for idx := 0; idx < typ.NumField(); idx++ {
		field := typ.Field(idx)
		Log(DEBUG, "#%d field: %-12v %v", idx, field.Name, field.Type)

		v := value.Elem().Field(idx)
		if !v.CanSet() || strings.TrimSpace(string(field.Tag)) == TAG_IGNORE {
			// the field will not be processed, skip
			Log(INFO, "#%-2d field %-12v skip", idx, field.Name)
			continue
		}

		var new_field *Field
		switch {
		case field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct:
			Log(DEBUG, "#%-2d field %-12v is sub-command", idx, field.Name)
			if new_field, err = NewField(v, field, SUBCOMMAND); err != nil {
				err = fmt.Errorf("cannot process %v.%v: %v", typ.Name(), field.Name, err)
				return
			}
			parser.subcommands = append(parser.subcommands, new_field)
		case field.Type.Kind() == reflect.Ptr:
			Log(DEBUG, "#%-2d field %-12v is argument", idx, field.Name)
			if new_field, err = NewField(v, field, ARGUMENT); err != nil {
				err = fmt.Errorf("cannot process %v.%v: %v", typ.Name(), field.Name, err)
				return
			}
			parser.arguments = append(parser.arguments, new_field)
		default:
			Log(DEBUG, "#%-2d field %-12v is option", idx, field.Name)
			if new_field, err = NewField(v, field, OPTION); err != nil {
				err = fmt.Errorf("cannot process %v.%v: %v", typ.Name(), field.Name, err)
				return
			}
			parser.options = append(parser.options, new_field)
		}

		Log(INFO, "add new field: %v", new_field)
	}

	return
}

type ArgParse struct {
	// the raw value pass to parser
	reflect.Value

	// the program name for the parser, default is the name of passed structure as lowercase
	Name string

	// the field in the argparse
	options     []*Field
	arguments   []*Field
	subcommands []*Field

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

	msgs = append(msgs, parser.usage())

	if len(parser.options) > 0 {
		margin, pending, siz := FMT_MARGIN, FMT_PENDING, FMT_SIZE
		msgs = append(msgs, []string{"", "option:"}...)

		for _, field := range parser.options {
			if field.Shortcut != rune(0) {
				if p := WidecharSize(string(field.Shortcut)) + WidecharSize(field.TypeHint) + 4; p > pending {
					// override the pending
					pending = p
				}
			}

			if s := WidecharSize(field.Name) + WidecharSize(field.TypeHint) + 6; s > siz {
				// override the size
				siz = s
			}
		}

		for _, field := range parser.options {
			msgs = append(msgs, field.FormatString(margin, pending, siz))
		}
	}

	if len(parser.arguments) > 0 {
		margin, pending, siz := FMT_MARGIN, FMT_PENDING, FMT_SIZE
		msgs = append(msgs, []string{"", "argument:"}...)

		for _, field := range parser.arguments {
			if s := WidecharSize(field.Name) + WidecharSize(field.TypeHint) + 4; s > siz {
				// override the size
				siz = s
			}
		}

		for _, field := range parser.arguments {
			msgs = append(msgs, field.FormatString(margin, pending, siz))
		}
	}

	if len(parser.subcommands) > 0 {
		margin, pending, siz := FMT_MARGIN, FMT_PENDING, FMT_SIZE
		msgs = append(msgs, []string{"", "sub-command:"}...)

		for _, field := range parser.subcommands {
			if s := WidecharSize(field.Name) + WidecharSize(field.TypeHint) + 4; s > siz {
				// override the size
				siz = s
			}
		}

		for _, field := range parser.subcommands {
			msgs = append(msgs, field.FormatString(margin, pending, siz))
		}
	}

	msg := strings.Join(msgs, "\n") + "\n"
	parser.Stderr.WriteString(msg)
}

func (parser *ArgParse) usage() (str string) {
	str = fmt.Sprintf("usage: %v", parser.Name)

	if len(parser.options) > 0 {
		// add the option
		str = fmt.Sprintf("%v [OPTION]", str)
	}

	if len(parser.arguments) > 0 || len(parser.subcommands) > 0 {
		// add the command
		str = fmt.Sprintf("%v ARGUMENT", str)
	}

	return
}
