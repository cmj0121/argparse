package argparse

import (
	"os"
	"strings"

	log "github.com/cmj0121/logger"
)

var (
	logger = log.New("argparse")
)

const (
	CRIT    = log.CRIT
	WARN    = log.WARN
	INFO    = log.INFO
	DEBUG   = log.DEBUG
	VERBOSE = log.VERBOSE
)

func init() {
	lv := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	logger.SetLevel(lv)
}

func Log(lv log.LogLevel, msg string, args ...interface{}) {
	logger.Log(lv, msg, args...)
	return
}
