package argparse

import (
	"fmt"
)

// execute the callback routine
type Callback func(parser *ArgParse) bool

func (parser *ArgParse) RegisterCallback(name string, fn Callback) (err error) {
	if _, ok := parser.callbacks[name]; ok {
		err = fmt.Errorf("duplicated callback: %v", name)
		return
	}

	parser.callbacks[name] = fn
	return
}
