# argparse #
The *argparse* is the Go-based command-line parser.


## Parser ##
The argparse is based on the reflect to generate the parser by pass the pointer of structure. The fields in
the structure may contains the tag which is the customized setting on the field. Without of the general,
the field in the structure may or may not the pointer. The general field, include the embedded structure
are treated as the *option*, and the pointer field will be treated as the argument.

The type of the field is used to control the pass data to the option and/or argument. For example, the boolean
type is used as the switch, and the integer will only allow to save the digest variable.

| type      | description                                          |
|-----------|------------------------------------------------------|
| bool      | the switch toggle without pass the extra varaiable   |
| int       | pass the valid degist and save as the int            |
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
| callback | the callback frunction and be triggered when pass the valid argument |
| choices  | fixed choice of the pass arguments, separated by the space           |

## Inner Log sub-system ##
The `Log` is the sub-system in the argparse which provide the simple logging system. It can be change the log
level by pass the environment *LOG_LEVEL* and change the level, and you can override the logger by `SetLogger`.
