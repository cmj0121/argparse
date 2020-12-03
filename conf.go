package argparse

// the general version info
const (
	PROJ_NAME = "argparse"
	MAJOR     = 0
	MINOR     = 5
	MACRO     = 0
)

// type hint of the field
const (
	TYPE_INT    = "INT"
	TYPE_STRING = "STR"
	TYPE_TIME   = "TIME"
)

// default key tag and value used when parse *Struct
const (
	TAG_IGNORE = "-"
	// the reserved key used in the structure
	TAG_RESERVED_KEY = "args"
	TAG_SHORTCUT     = "short"
	TAG_NAME         = "name"
	TAG_HELP         = "help"
	TAG_CALLBACK     = "callback"
	TAG_CHOICES      = "choices"
	TAG_CHOICES_SEP  = " "
	// the reserved key used in TAG_KEY
	KEY_PASSWORD = "password"
	// default callback KEY
	FN_HELP    = "_help"
	FN_VERSION = "_version"
)

// the default formatted string config
const (
	FMT_MARGIN  = 4
	FMT_PENDING = 8
	FMT_SIZE    = 24
)
