# argparse #
![Go test](https://github.com/cmj0121/argparse/workflows/test/badge.svg)
The *argparse* is the Go-based command-line parser.

## Example ##
The following is the sample structure that can be generate the command-line parse by the argparse. The `argparse.Help`.
Call `Run` or `Parse` and the argparse will parse the input argument (default is os.Args) and then set the field in the
struct

```go
type Simple struct {
	argparse.Help

	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count  int    `short:"C" help:"save as the integer"`
	Name   string `name:"user-name"`
	Cases  string `short:"c" choices:"demo foo" help:"choice from fix possible"`
	Now    time.Time

	Optional []string  `name:"opt" help:"multiple option and save as array"`
	Args     *[]string `help:"arbitrary argument"`
}
```

## Parser ##
The argparse is based on the reflect to generate the parser by pass the pointer of structure. The fields in
the structure may contains the tag which is the customized setting on the field. Without of the general,
the field in the structure may or may not the pointer. The general field, include the embedded structure
are treated as the *option*, and the pointer field will be treated as the argument.

The type of the field is used to control the pass data to the option and/or argument. For example, the boolean
type is used as the switch, and the integer will only allow to save the digest variable.

| type      | description                                          |
|-----------|------------------------------------------------------|
| bool      | the switch toggle without pass the extra variable    |
| int       | pass the valid gigital and save as the int            |
| string    | pass any string, include empty string or binary data |
| time.Time | pass the valid RFC-3339 time format string           |

### tags ###
There are few tags use for the customized field setting

| tag      | description                                                          |
|----------|----------------------------------------------------------------------|
| -        | ignore this field                                                    |
| name     | replace the field name, and will only treated as the lowercase       |
| short    | the shortcut of option, should be one and only one rune              |
| help     | the help message of the option or argument                           |
| callback | the callback function and be triggered when pass the valid argument  |
| choices  | fixed choice of the pass arguments, separated by the space           |

### Callback ##
You can define the **callback** when you have to execute some specified method when set the valid option or argument.
There are two methods when define the callback: 1) global callback and 2) the method in your structure. When call the
**RegisterCallback** the parser will register callback in the global scope, and it can be used on other parser. Also
you can define the method as the same type of `Callback`, correct defined in the tag and it will be executed when set
the valid value.

The `GetCallback` will find the customized callback first, and then try the global callback. It may return **nil** 
when no valid callback found.

## Inner Log sub-system ##
The `Log` is the sub-system in the argparse which provide the simple logging system. It can be change the log
level by pass the environment *LOG_LEVEL* and change the level, and you can override the logger by `SetLogger`.
