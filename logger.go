// the customized log system
package argparse

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// the log level
type LogLevel int

const (
	CRIT LogLevel = iota
	WARN
	INFO
	DEBUG
	VERBOSE
)

const ENV_LOG_LEVEL = "LOG_LEVEL"

var (
	log_level = CRIT
	logger    = log.New(os.Stderr, "", log.Lshortfile)
)

func init() {
	logs := map[string]LogLevel{
		"CRIT":    CRIT,
		"WARN":    WARN,
		"INFO":    INFO,
		"DEBUG":   DEBUG,
		"VERBOSE": VERBOSE,
	}

	lv := strings.ToUpper(os.Getenv(ENV_LOG_LEVEL))
	if level, ok := logs[lv]; ok {
		// override the log level by ENV
		log_level = level
	}
}

func Log(lv LogLevel, msg string, args ...interface{}) {
	if lv <= log_level {
		logger.Output(2, fmt.Sprintf(msg, args...))
	}
}

func SetLogger(log *log.Logger) {
	logger = log
	return
}
