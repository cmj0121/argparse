package argparse

import (
	"reflect"
)

var (
	// the callback when option triggered
	callbacks = map[string]Callback{}
)

// execute the callback routine
type Callback func(parser *ArgParse) bool

func RegisterCallback(name string, fn Callback) {
	if _, ok := callbacks[name]; ok {
		// show the alert
		Log(WARN, "duplicated callback %v, override", name)
	}
	callbacks[name] = fn
	return
}

func GetCallback(value reflect.Value, name string) (fn Callback) {
	var ok bool

	if fn_val := value.MethodByName(name); fn_val.IsValid() && !fn_val.IsZero() {
		// find the pass method of name
		// HACK - the type of the callback should be same as Callback, but only can convert as func(*ArgParse) bool
		if fn, ok = fn_val.Interface().(func(*ArgParse) bool); ok {
			return
		}
	}
	// try the global callback
	fn, _ = callbacks[name]
	return
}
