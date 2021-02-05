package argparse

// the general version info
const (
	PROJ_NAME = "argparse"

	MAJOR = 0
	MINOR = 6
	MACRO = 6
)

// type hint of the field
const (
	TYPE_INT    = "INT"
	TYPE_STRING = "STR"
	TYPE_PERM   = "PERM"
	TYPE_TIME   = "TIME"
	TYPE_IFACE  = "IFACE"
	TYPE_IP     = "IP"
	TYPE_CIDR   = "CIDR"
	TYPE_FILE   = "FILE"
)

// default key tag and value used when parse *Struct
const (
	TAG_IGNORE = "-"
	// the reserved key used in the structure
	TAG_RESERVED_KEY = "args"
	TAG_OPTION       = "option"

	TAG_SHORTCUT    = "short"
	TAG_NAME        = "name"
	TAG_HELP        = "help"
	TAG_CALLBACK    = "callback"
	TAG_CHOICES     = "choices"
	TAG_CHOICES_SEP = " "

	TAG_DEFAULT_KEY = "default"
	// the reserved key used in TAG_KEY
	KEY_PASSWORD = "password"
	// default callback KEY
	FN_HELP    = "_help"
	FN_VERSION = "_version"
)

// the default formatted string config
const (
	FMT_MARGIN  = 4
	FMT_PENDING = 9
	FMT_SIZE    = 24
)
