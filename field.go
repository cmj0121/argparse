package argparse

import (
	"fmt"
	"net"
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

func (ftyp FieldType) String() (str string) {
	ftyps := []string{
		"OPTION",
		"ARGUMENT",
		"SUB-COMMAND",
	}
	str = ftyps[ftyp]
	return
}

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
	log.Verbose("new field: %T (tag: `%v`) (%v)", value.Interface(), sfield.Tag, ftyp)

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
		var obj reflect.Value

		if obj = field.Value; field.Value.IsNil() {
			// nil sub-command, new instance
			obj = reflect.New(field.Value.Type().Elem())
		}
		if field.Subcommand, err = New(obj.Interface()); err != nil {
			// cannot set the sub-command
			return
		}

		if field.StructTag.Get(TAG_NAME) != "" {
			field.Subcommand.Name = strings.ToLower(field.StructTag.Get(TAG_NAME))
			field.Subcommand.Name = strings.TrimSpace(field.Name)
		}
	}

	// set the type hint
	field.setTypeHint(value.Type())

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
		case SUBCOMMAND:
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

	// HACK - set the default setting
	switch field.Value.Interface().(type) {
	case os.FileMode, *os.FileMode:
		log.Info("HACK - set the os.FileMode default settinig")
		if field.Help == "" {
			// set the default help message
			field.Help = fmt.Sprintf("file perm")
		}
	case time.Time, *time.Time:
		log.Info("HACK - set the time.Time default settinig")
		if field.Help == "" {
			// set the default help message
			field.Help = fmt.Sprintf("timestamp RFC-3339 (2006-01-02T15:04:05+07:00)")
		}
	case net.Interface, *net.Interface:
		log.Info("HACK - set the net.Interface default settinig")
		if field.Help == "" {
			// set the default help message
			field.Help = "network interface"
		}
	}

	return
}

func (field *Field) String() (str string) {
	str = field.FormatString(4, 8, 18)
	return
}

func (field *Field) setTypeHint(typ reflect.Type) {
	switch typ.Kind() {
	case reflect.Int:
		field.TypeHint = TYPE_INT
	case reflect.String:
		field.TypeHint = TYPE_STRING
	case reflect.Slice:
		field.setTypeHint(typ.Elem())
	default:
		switch field.Value.Interface().(type) {
		case os.FileMode:
			field.TypeHint = TYPE_PERM
		case time.Time:
			field.TypeHint = TYPE_TIME
		case net.Interface, *net.Interface:
			field.TypeHint = TYPE_IFACE
		case net.IP, *net.IP:
			field.TypeHint = TYPE_IP
		case net.IPNet, *net.IPNet:
			field.TypeHint = TYPE_CIDR
		}
	}
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

// pre-process the field setting, include new instance
func (field *Field) SetValue(parser *ArgParse, args ...string) (size int, err error) {
	size = 1
	// the basic setter
	if size, err = field.setValue(field.Value, args...); err != nil {
		return
	}

	if fn := GetCallback(parser.Value, field.Callback); fn != nil {
		log.Debug("try execute %v", field.Callback)
		// trigger the callback, exit when callback return true
		if fn(parser) && ExitWhenCallback {
			log.Info("execute callback %v, and exit 0", field.Callback)
			os.Exit(0)
		}
	}

	field.BeenSet = true
	if field.Value.Kind() == reflect.Ptr && field.Value.Elem().Kind() == reflect.Slice {
		// can set repeat
		field.BeenSet = false
	}
	log.Info("set %v as %v (%d)", field.Name, field.Value, size)
	return
}

// the exactly set the value to the field
func (field *Field) setValue(value reflect.Value, args ...string) (size int, err error) {
	log.Debug("try set value %[1]T (%#v)", value.Interface(), args)

	switch value.Interface().(type) {
	case bool:
		// toggle the boolean
		value.SetBool(!value.Interface().(bool))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
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

		switch field.Value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value.SetInt(int64(val))
		default:
			value.SetUint(uint64(val))
		}

		size++
	case string:
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
	case os.FileMode:
		if len(args) == 0 {
			err = fmt.Errorf("should pass file perm")
			return
		}

		log.Info("set os.FileMode as %v", args[0])
		var perm int
		if perm, err = strconv.Atoi(args[0]); err != nil || uint64(perm)&uint64(0xFFFFFFFF00000000) != 0 {
			log.Info("cannot set os.FileMode %#v: %v", args[0], err)
			err = fmt.Errorf("cannot set os.FileMode %#v: %v", args[0], err)
			return
		}

		value.Set(reflect.ValueOf(os.FileMode(uint32(perm))))
		size++
	case time.Time:
		if len(args) == 0 {
			err = fmt.Errorf("should pass TIME: RFC-3339 (2006-01-02T15:04:05+07:00)")
			return
		}

		log.Info("set time.Time as %v", args[0])
		var timestamp time.Time

		if timestamp, err = time.Parse(time.RFC3339, args[0]); err != nil {
			log.Info("should pass RFC-3339 (2006-01-02T15:04:05+07:00): %v: %v", args[0], err)
			err = fmt.Errorf("should pass RFC-3339 (2006-01-02T15:04:05+07:00): %v: %v", args[0], err)
			return
		}
		value.Set(reflect.ValueOf(timestamp))
		size++
	case net.Interface:
		if len(args) == 0 {
			err = fmt.Errorf("should pass IFACE")
			return
		}

		var iface *net.Interface

		if iface, err = net.InterfaceByName(args[0]); err != nil {
			err = fmt.Errorf("invalid IFACE %#v: %v", args[0], err)
			return
		}

		value.Set(reflect.ValueOf(*iface))
		size++
	case net.IP:
		if len(args) == 0 {
			err = fmt.Errorf("should pass IP")
			return
		}

		ip := net.ParseIP(args[0])
		if ip == nil {
			err = fmt.Errorf("invalid IP: %v", args[0])
			return
		}

		value.Set(reflect.ValueOf(ip))
		size++
	case net.IPNet:
		if len(args) == 0 {
			err = fmt.Errorf("should pass CIDR")
			return
		}

		var inet *net.IPNet
		_, inet, err = net.ParseCIDR(args[0])
		if err != nil {
			err = fmt.Errorf("invalid CIDR: %v", args[0])
			return
		}

		value.Set(reflect.ValueOf(*inet))
		size++
	default:
		switch value.Kind() {
		case reflect.Struct:
			// execute sub-command
			if err = field.Subcommand.Parse(args...); err != nil {
				// only show the help message on the sub-command
				field.Subcommand.HelpMessage(err)
				os.Exit(1)
			}
		case reflect.Ptr:
			log.Debug("set pointer %v: %v", field.Name, value)

			if value.IsNil() {
				if field.Subcommand != nil {
					log.Info("nil pointer, assign as sub-command")

					value.Set(field.Subcommand.Value)
				} else {
					log.Info("nil pointer, new instance: %v", value.Type())

					obj := reflect.New(value.Type().Elem())
					value.Set(obj)
				}
			}

			if size, err = field.setValue(value.Elem(), args...); err != nil {
				return
			}
		case reflect.Slice:
			elem := reflect.New(value.Type().Elem()).Elem()
			if size, err = field.setValue(elem, args...); err != nil {
				err = fmt.Errorf("cannot set %#v: %v", value.Type(), err)
				return
			}

			// append to the slice
			value.Set(reflect.Append(value, elem))
		default:
			log.Warn("not implemented set value: %[1]v (%[1]T)", value.Interface())
			err = fmt.Errorf("not implemented set value: %[1]v (%[1]T)", value.Interface())
			return
		}
	}

	log.Debug("success set %v (%d)", value, size)
	return
}

// calculate the multiple-char size
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
