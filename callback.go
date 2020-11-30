package argparse

import (
	"fmt"
)

var (
	// the callback when option triggered
	callbacks = map[string]Callback{}
)

// execute the callback routine
type Callback func(parser *ArgParse) bool

func RegisterCallback(name string, fn Callback) (err error) {
	if _, ok := callbacks[name]; ok {
		err = fmt.Errorf("duplicated callback: %v", name)
		return
	}

	callbacks[name] = fn
	return
}
