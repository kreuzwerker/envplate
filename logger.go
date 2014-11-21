package envplate

import (
	"fmt"
	"log"
)

type logLevel string

const (
	DEBUG logLevel = "DEBUG"
	INFO           = "INFO"
	ERROR          = "ERROR"
)

var (
	DebugAndInfoFunc = log.Printf
	ErrorFunc        = log.Fatalf
)

func Log(lvl logLevel, msg string, args ...interface{}) {

	if lvl == DEBUG && !Config.Verbose {
		return
	}

	msg = fmt.Sprintf("[ %s ] %s", lvl, msg)

	if lvl == ERROR {
		ErrorFunc(msg, args...)
	} else {
		DebugAndInfoFunc(msg, args...)
	}

}
