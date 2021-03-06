package argparse

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"

	"github.com/cmj0121/logger"
)

var (
	Stderr           = os.Stderr
	ExitWhenCallback = true
	log              = logger.New(PROJ_NAME)
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
	log.Info("new %[1]T", in)

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
		Value: value,
		Name:  name,

		used_option:     map[string]*Field{},
		used_shortcut:   map[rune]*Field{},
		used_subcommand: map[string]*Field{},
	}

	// process the field
	typ := value.Elem().Type()
	log.Verbose("start process: %v", typ)
	for idx := 0; idx < typ.NumField(); idx++ {
		field := typ.Field(idx)
		log.Debug("#%d field: %-12v %v", idx, field.Name, field.Type)

		v := value.Elem().Field(idx)
		if !v.CanSet() || strings.TrimSpace(string(field.Tag)) == TAG_IGNORE {
			// the field will not be processed, skip
			log.Info("#%-2d field %-12v skip", idx, field.Name)
			continue
		}

		if err = parser.setField(v, field); err != nil {
			err = fmt.Errorf("cannot processed %v.%v: %v", typ.Name(), field.Name, err)
			return
		}
	}

	return
}

type ArgParse struct {
	// the raw value pass to parser
	reflect.Value

	// the program name for the parser, default is the name of passed structure as lowercase
	Name string
	// set exit when callback success

	// the field in the argparse
	options     []*Field
	arguments   []*Field
	subcommands []*Field

	// the cache for the used options
	used_option     map[string]*Field
	used_shortcut   map[rune]*Field
	used_subcommand map[string]*Field
}

func (parser *ArgParse) setField(val reflect.Value, field reflect.StructField) (err error) {
	var new_field *Field

	log.Debug("try set field: %v (%v) (%v)", val, val.Type(), field.Tag)

	switch field.Tag.Get(TAG_RESERVED_KEY) {
	case TAG_IGNORE:
		log.Info("skip field: %v (%v)", val, field.Tag)
		return
	default:
		switch {
		case field.Type.Kind() == reflect.Ptr: // argument or sub-command
			log.Debug("argument or sub-command: %v", field.Type.Elem().Kind())

			switch field.Tag.Get(TAG_RESERVED_KEY) {
			case TAG_OPTION:
				log.Info("force set as option: %[1]T", val.Interface())
				if new_field, err = NewField(val, field, OPTION); err != nil {
					return
				}

				if _, ok := parser.used_option["--"+new_field.Name]; ok {
					err = fmt.Errorf("duplicated option --%v", new_field.Name)
					return
				}
				parser.used_option["--"+new_field.Name] = new_field

				if new_field.Shortcut != rune(0) {
					if _, ok := parser.used_option["-"+string(new_field.Shortcut)]; ok {
						err = fmt.Errorf("duplicated option -%v", string(new_field.Shortcut))
						return
					}
					parser.used_option["-"+string(new_field.Shortcut)] = new_field
				}

				parser.options = append(parser.options, new_field)
			default:
				switch field.Type.Elem().Kind() {
				case reflect.Struct:
					switch val.Interface().(type) {
					case *net.Interface:
						if new_field, err = NewField(val, field, ARGUMENT); err != nil {
							return
						}
						parser.arguments = append(parser.arguments, new_field)
					default:
						if new_field, err = NewField(val, field, SUBCOMMAND); err != nil {
							return
						}

						if _, ok := parser.used_subcommand[new_field.Name]; ok {
							err = fmt.Errorf("duplicated subcommands %v", new_field.Name)
							return
						}
						parser.used_subcommand[new_field.Name] = new_field
						parser.subcommands = append(parser.subcommands, new_field)
					}
				default:
					if new_field, err = NewField(val, field, ARGUMENT); err != nil {
						return
					}
					parser.arguments = append(parser.arguments, new_field)
				}
			}
		case field.Type.Kind() == reflect.Struct && field.Anonymous: // embedded field
			log.Debug("embedded field: %T", val.Interface())

			switch val.Interface().(type) {
			default:
				for idx := 0; idx < field.Type.NumField(); idx++ {
					v := val.Field(idx)

					sub_field := field.Type.Field(idx)
					if !v.CanSet() || strings.TrimSpace(string(sub_field.Tag)) == TAG_IGNORE {
						// the field will not be processed, skip
						log.Info("#%-2d field %v.%v skip", idx, field.Name, sub_field.Name)
						continue
					}

					if err = parser.setField(v, sub_field); err != nil {
						err = fmt.Errorf("set %v.%v: %v", field.Name, sub_field.Name, err)
						return
					}
				}
			}
		default:
			if new_field, err = NewField(val, field, OPTION); err != nil {
				return
			}

			if _, ok := parser.used_option["--"+new_field.Name]; ok {
				err = fmt.Errorf("duplicated option --%v", new_field.Name)
				return
			}
			parser.used_option["--"+new_field.Name] = new_field

			if new_field.Shortcut != rune(0) {
				if _, ok := parser.used_option["-"+string(new_field.Shortcut)]; ok {
					err = fmt.Errorf("duplicated option -%v", string(new_field.Shortcut))
					return
				}
				parser.used_option["-"+string(new_field.Shortcut)] = new_field
			}

			parser.options = append(parser.options, new_field)
		}
	}

	if new_field != nil && new_field.Callback != "" && GetCallback(parser.Value, new_field.Callback) == nil {
		err = fmt.Errorf("callback %v not defined", new_field.Callback)
		return
	}

	log.Info("add new field: %v", new_field)
	return
}

func (parser *ArgParse) Run() (err error) {
	if err = parser.Parse(os.Args[1:]...); err != nil {
		// show the help message
		parser.HelpMessage(err)
	}
	return
}

func (parser *ArgParse) Parse(args ...string) (err error) {
	log.Info("parse %#v", args)

	for idx, size := 0, 0; idx < len(args); idx += size {
		token := args[idx]

		log.Info("%v parse #%-2d %v", parser.Name, idx, token)
	PROCESS_FIELD:
		switch {
		case len(token) > 2 && token[:2] == "--":
			log.Debug("optional: %v", token)

			for _, field := range parser.options {
				if field.Name == token[2:] {
					// set the value
					if size, err = field.SetValue(parser, args[idx+1:]...); err != nil {
						// cannot set the value, raise
						err = fmt.Errorf("%v %v", token, err)
						return
					}

					size++
					break PROCESS_FIELD
				}
			}

			log.Warn("unknown option: %v", token)
			err = fmt.Errorf("unknown option: %v", token)
			return
		case len(token) > 1 && token[:1] == "-":
			log.Debug("shortcut: %v (%d)", token[1:], WidecharSize(token[1:]))

			switch {
			case WidecharSize(token[1:]) == 1 || (WidecharSize(token[1:]) == 2 && len(token[1:]) > 2):
				shortcut := []rune(token[1:])[0]

				for _, field := range parser.options {
					if field.Shortcut == shortcut {
						if size, err = field.SetValue(parser, args[idx+1:]...); err != nil {
							// cannot set the value, raise
							err = fmt.Errorf("%v %v", token, err)
							return
						}

						size++
						break PROCESS_FIELD
					}
				}
			default:
				for _, shortcut := range token[1:] {
					found := false
					for _, field := range parser.options {
						if field.Shortcut == shortcut {
							if _, err = field.SetValue(parser); err != nil {
								// cannot set the value, raise
								err = fmt.Errorf("multi-shortcut %#v cannot set: %v", token, err)
								return
							}

							found = true
						}
					}

					if !found {
						err = fmt.Errorf("unknown option: -%v", string(shortcut))
						return
					}
				}

				// skip this option
				size++
				break PROCESS_FIELD
			}

			err = fmt.Errorf("unknown option: %v", token)
			return
		default:
			log.Debug("argument or sub-command: %v", token)

			// check the sub-command first
			for _, field := range parser.subcommands {
				if field.Name == token {
					log.Info("set sub-command %v", field.Name)
					if _, err = field.SetValue(parser, args[idx+1:]...); err != nil {
						// cannot set the value, raise
						err = fmt.Errorf("%v %v", field.Name, err)
						return
					}

					// always return when process sub-command
					return
				}
			}

			for _, field := range parser.arguments {
				if field.BeenSet {
					log.Info("field %v already set %v, skip", field.Name, field.Value)
					continue
				}

				if size, err = field.SetValue(parser, args[idx:]...); err != nil {
					// cannot set the value, raise
					err = fmt.Errorf("%v %v", field.Name, err)
					return
				}

				break PROCESS_FIELD
			}

			log.Warn("unknown argument: %v", token)
			err = fmt.Errorf("unknown argument: %v", token)
			return
		}
	}
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
			log.Debug("format string m:%d, p:%d, s:%d", margin, pending, siz)
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
			log.Debug("format string m:%d, p:%d, s:%d", margin, pending, siz)
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
			log.Debug("format string m:%d, p:%d, s:%d", margin, pending, siz)
			msgs = append(msgs, field.FormatString(margin, pending, siz))
		}
	}

	msg := strings.Join(msgs, "\n") + "\n"
	Stderr.WriteString(msg)
}

func (parser *ArgParse) usage() (str string) {
	str = fmt.Sprintf("usage: %v", parser.Name)

	if len(parser.options) > 0 {
		// add the option
		str = fmt.Sprintf("%v [OPTION]", str)
	}

	// add the command
	for _, field := range parser.arguments {
		if field.FieldType == ARGUMENT {
			str = fmt.Sprintf("%v %v", str, strings.ToUpper(field.Name))
		}
	}

	return
}
