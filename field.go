package argparse

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// the field in the argparse which may be 1) option, 2) argument and 3) sub-command
type FieldType int

const (
	OPTION FieldType = iota
	ARGUMENT
	SUBCOMMAND
)

// the field used in the argparse
type Field struct {
	// the passed value for the field, from the parser
	reflect.Value
	reflect.Type
	reflect.StructTag

	Subcommand *ArgParse

	// the field type
	FieldType

	// set flag
	BeenSet bool

	// the display field
	Name     string
	TypeHint string
	Shortcut rune
	Help     string

	Callback     string
	DefaultValue interface{}
	Choices      []string
}

func NewField(value reflect.Value, sfield reflect.StructField, ftyp FieldType) (field *Field, err error) {
	field = &Field{
		Value:     value,
		Type:      sfield.Type,
		StructTag: sfield.Tag,
		FieldType: ftyp,
	}

	if field.Name = strings.ToLower(sfield.Name); field.StructTag.Get(TAG_NAME) != "" {
		field.Name = strings.ToLower(field.StructTag.Get(TAG_NAME))
		field.Name = strings.TrimSpace(field.Name)
	}

	// customized pre-process by field type
	switch ftyp {
	case ARGUMENT:
		// set the display as the upper-case
		field.Name = strings.ToUpper(field.Name)
	case SUBCOMMAND:
		obj := reflect.New(field.Value.Type().Elem())
		if field.Subcommand, err = New(obj.Interface()); err != nil {
			// cannot set the sub-command
			return
		}

		if field.StructTag.Get(TAG_NAME) != "" {
			field.Subcommand.Name = strings.ToLower(field.StructTag.Get(TAG_NAME))
			field.Subcommand.Name = strings.TrimSpace(field.Name)
		}
	}

	if s := field.StructTag.Get(TAG_SHORTCUT); s != "" {
		shortcut := []rune(s)

		switch {
		case len(shortcut) > 1:
			err = fmt.Errorf("shortcut too large: %s", s)
			return
		case len(shortcut) == 1:
			field.Shortcut = shortcut[0]
		}
	}

	if help := field.StructTag.Get(TAG_HELP); help != "" {
		// set the help message
		field.Help = help
	}

	if callback := field.StructTag.Get(TAG_CALLBACK); callback != "" {
		// set the callback name
		field.Callback = callback
	}

	if field.Value.IsValid() && !field.Value.IsZero() {
		switch field.FieldType {
		case ARGUMENT:
			field.DefaultValue = field.Value.Elem().Interface()
		default:
			field.DefaultValue = field.Value.Interface()
		}
	}

	if c := field.StructTag.Get(TAG_CHOICES); c != "" {
		field.Choices = []string{}
		for _, choice := range strings.Split(c, TAG_CHOICES_SEP) {
			// add into the fix choices
			field.Choices = append(field.Choices, strings.TrimSpace(choice))
		}

		sort.Strings(field.Choices)
		// check the default in the choice or NOT
		if field.DefaultValue != nil {
			choice := fmt.Sprintf("%v", field.DefaultValue)
			if field.Choices[sort.SearchStrings(field.Choices, choice)] != choice {
				err = fmt.Errorf("%v not in the choices: %v", choice, field.Choices)
				return
			}
		}
	}

	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch field.Value.Kind() {
	case reflect.Int:
		field.TypeHint = TYPE_INT
	case reflect.String:
		field.TypeHint = TYPE_STRING
	}

	return
}

func (field *Field) String() (str string) {
	str = field.FormatString(4, 8, 18)
	return
}

// the format string for the field
// | margin | pending  | size | margin |      |
// |        | Shortcut | Name |        | Help |
func (field *Field) FormatString(margin, pending, size int) (str string) {
	option := field.Name

	switch field.FieldType {
	case OPTION:
		// --KEY TYPE
		option = fmt.Sprintf("%*v--%v %v", pending, "", field.Name, field.TypeHint)
		option = strings.TrimRight(option, " \t\n")

		// -SHORT TYPE, --KEY TYPE
		if field.Shortcut != rune(0) {
			shortcut := fmt.Sprintf("-%v %v", string(field.Shortcut), field.TypeHint)
			shortcut = fmt.Sprintf("%v, ", strings.TrimSpace(shortcut))
			shift := len(shortcut) - WidecharSize(shortcut)
			option = fmt.Sprintf("%*v--%v %v", pending-shift, shortcut, field.Name, field.TypeHint)
		}
	}

	help := fmt.Sprintf("%v", field.Help)
	if len(field.Choices) > 0 {
		choices := strings.Join(field.Choices, TAG_CHOICES_SEP)
		help = fmt.Sprintf("%v [%v]", help, choices)
	}

	if field.DefaultValue != nil {
		// set the default value
		switch field.FieldType {
		case ARGUMENT:
			help = fmt.Sprintf("%v (default: %v)", help, field.DefaultValue)
		default:
			help = fmt.Sprintf("%v (default: %v)", help, field.DefaultValue)
		}
	}

	shift := len(option) - WidecharSize(option)
	str = fmt.Sprintf("%*v%-*v%*v", margin, "", pending+size-shift, option, margin, strings.TrimSpace(help))
	str = strings.TrimRight(str, " \t\n")
	return
}

func (field *Field) SetValue(parser *ArgParse, args ...string) (size int, err error) {
	size = 1
	switch kind := field.Value.Kind(); kind {
	case reflect.Bool, reflect.Int, reflect.String:
		// the basic setter
		if size, err = field.setValue(field.Value, args...); err != nil {
			return
		}
	case reflect.Ptr:
		Log(DEBUG, "set pointer %v: %v", field.Name, field.Value)

		if field.Value.IsNil() {
			if field.Subcommand != nil {
				Log(INFO, "nil pointer, assign as sub-command")

				field.Value.Set(field.Subcommand.Value)
			} else {
				Log(INFO, "nil pointer, new instance: %v", field.Value.Type())

				obj := reflect.New(field.Value.Type().Elem())
				field.Value.Set(obj)
			}
		}

		if size, err = field.setValue(field.Value.Elem(), args...); err != nil {
			return
		}

		// HACK - override the used args as 1
		size = 1
	default:
		switch {
		case field.Value.Type() == reflect.TypeOf(time.Time{}):
			if len(args) == 0 {
				err = fmt.Errorf("should pass TIME: RFC-3339 (%v)", time.RFC3339)
				return
			}

			Log(INFO, "set time.Time as %v", args[0])
			var timestamp time.Time

			if timestamp, err = time.Parse(time.RFC3339, args[0]); err != nil {
				err = fmt.Errorf("should pass RFC-3339 (%v): %v: %v", time.RFC3339, args[0], err)
				Log(INFO, "should pass RFC-3339 (%v): %v: %v", time.RFC3339, args[0], err)
				return
			}
			field.Value.Set(reflect.ValueOf(timestamp))
		default:
			Log(WARN, "not implemented set field kind: %v (%v)", kind, field.Value.Type())
			err = fmt.Errorf("not support field: %v", field.Name)
			return
		}
	}

	if fn, ok := parser.callbacks[field.Callback]; ok {
		Log(DEBUG, "try execute %v", field.Callback)
		// trigger the callback
		if fn(parser) && parser.ExitOnCallback {
			// set the exit when return true
			Log(INFO, "exit when call %v", field.Callback)
			os.Exit(0)
		}
	}

	field.BeenSet = true
	Log(INFO, "set %v as %v", field.Name, field.Value)
	return
}

func (field *Field) setValue(value reflect.Value, args ...string) (size int, err error) {
	size = 1

	switch kind := value.Kind(); kind {
	case reflect.Bool:
		// toggle the boolean
		value.SetBool(!value.Interface().(bool))
	case reflect.Int:
		// override the integer
		if len(args) == 0 {
			err = fmt.Errorf("should pass %v", TYPE_INT)
			return
		}

		var val int
		if val, err = strconv.Atoi(args[0]); err != nil {
			err = fmt.Errorf("should pass %v: %v", TYPE_INT, args[0])
			return
		}

		if len(field.Choices) > 0 {
			idx := sort.SearchStrings(field.Choices, args[0])
			if idx == len(field.Choices) || field.Choices[idx] != args[0] {
				err = fmt.Errorf("%v should choice from %v", args[0], field.Choices)
				return
			}
		}

		value.SetInt(int64(val))
		size++
	case reflect.String:
		// override the string
		if len(args) == 0 {
			err = fmt.Errorf("should pass %v", TYPE_STRING)
			return
		}

		if len(field.Choices) > 0 {
			idx := sort.SearchStrings(field.Choices, args[0])
			if idx == len(field.Choices) || field.Choices[idx] != args[0] {
				err = fmt.Errorf("%v should choice from %v", args[0], field.Choices)
				return
			}
		}

		value.SetString(args[0])
		size++
	default:
		// execute sub-command
		if err = field.Subcommand.Parse(args...); err != nil {
			// only show the help message on the sub-command
			field.Subcommand.HelpMessage(err)
			os.Exit(1)
		}
		size = len(args)
	}

	Log(INFO, "success set %v", value)
	return
}

func WidecharSize(widechar string) (siz int) {
	for _, s := range widechar {
		siz++
		if len(string(s)) > 1 {
			// detect wide-char
			siz++
		}
	}
	return
}
