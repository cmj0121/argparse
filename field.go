package argparse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

	// the field type
	FieldType

	// the display field
	Name     string
	TypeHint string
	Shortcut rune
	Help     string

	Callback     string
	DefaultValue interface{}
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

	if field.Value.IsValid() && !field.Value.IsZero() {
		switch field.FieldType {
		case ARGUMENT:
			field.DefaultValue = field.Value.Elem().Interface()
		default:
			field.DefaultValue = field.Value.Interface()
		}
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
	if field.DefaultValue != nil {
		// set the default value
		switch field.FieldType {
		case ARGUMENT:
			help = fmt.Sprintf("%v (default: %v)", field.Help, field.DefaultValue)
		default:
			help = fmt.Sprintf("%v (default: %v)", field.Help, field.DefaultValue)
		}
	}

	shift := len(option) - WidecharSize(option)
	str = fmt.Sprintf("%*v%-*v%*v", margin, "", pending+size-shift, option, margin, strings.TrimSpace(help))
	str = strings.TrimRight(str, " \t\n")
	return
}

func (field *Field) SetValue(args ...string) (size int, err error) {
	size = 1
	switch field.Value.Kind() {
	case reflect.Bool:
		// toggle the boolean
		field.Value.SetBool(!field.Value.Interface().(bool))
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

		field.Value.SetInt(int64(val))
		size++
	case reflect.String:
		// override the string
		if len(args) == 0 {
			err = fmt.Errorf("should pass %v", TYPE_STRING)
			return
		}

		field.Value.SetString(args[0])
		size++
	}

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
