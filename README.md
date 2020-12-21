# argparse #
![Go test](https://github.com/cmj0121/argparse/workflows/test/badge.svg)
The *argparse* is the Go-based command-line parser.

## Example ##

```go
type Simple struct {
	argparse.Model

	// the ignore field that will not be processed
	Ignore bool `-`
	ignore bool
	_      [8]byte `pending byte`

	// the option can be set repeatedly by-default
	Switch bool   `short:"s" name:"toggle" help:"toggle the boolean value"`
	Count  int    `short:"C" help:"save as the integer"`
	Name   string `name:"user-name"`
	Cases  string `short:"c" choices:"demo foo" help:"choice from fix possible"`
}

```

## Types ##
In the argparse it support several built-in type. The type of the field is used to control the pass data to the option and/or
the argument. For example, the boolean type is used as the switch, and the integer will only allow to save the as digest. It
is implemented in the `field.setValue`:

| type      | description                                          |
|-----------|------------------------------------------------------|
| bool      | the switch toggle without pass the extra variable    |
| int       | pass the valid gigital and save as the int            |
| string    | pass any string, include empty string or binary data |

### Syntax-Sugar ###
The argparse supports few types that can be easily parse and used in the command-line.

| type          | hint  | description                           |
|---------------|-------|---------------------------------------|
| os.FileMode   | PERM  | the file permission in the system     |
| time.Time     | TIME  | the timestamp noted by the RFC-3339   |
| net.Interface | IFACE | the interface in the system           |
| net.IP        | IP    | the IP format string                  |
| net.IPNet     | CIDR  | the IP with mask (CIDR) format string |

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
| args     | force set as the option (value: option)                              |

### Callback ##
You can define the **callback** when you have to execute some specified method when set the valid option or argument.
There are two methods when define the callback: 1) global callback and 2) the method in your structure. When call the
**RegisterCallback** the parser will register callback in the global scope, and it can be used on other parser. Also
you can define the method as the same type of `Callback`, correct defined in the tag and it will be executed when set
the valid value.

The `GetCallback` will find the customized callback first, and then try the global callback. It may return **nil** 
when no valid callback found.

