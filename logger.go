package envplate

import (
	"errors"
	"fmt"
	"io"
	_log "log"
	"os"
)

const (

	// RAW indicates a log message should be send directly to stdout
	RAW logLevel = "RAW"

	// DEBUG indicates log message that are only visible when the verbose flag it set
	DEBUG = "DEBUG"

	// INFO indicates regular log messages
	INFO = "INFO"

	// ERROR indicates error log messages
	ERROR = "ERROR"
)

type logLevel string

type logger struct {
	Out     io.Writer
	Verbose bool
}

// Log exposes a logger with a custom envplate formatting syntax
var Logger = &logger{Out: os.Stdout}

// Log emits log messages and filters based on the set verbosity - if an error
// is logged, the error msg is returned as error object
func Log(lvl logLevel, msg string, args ...interface{}) error {
	return Logger.log(lvl, msg, args...)
}

func (l *logger) log(lvl logLevel, msg string, args ...interface{}) error {

	if lvl == DEBUG && !l.Verbose {
		return nil
	}

	msg = fmt.Sprintf(msg, args...)

	if lvl == RAW {
		fmt.Fprintf(l.Out, msg)
		return nil
	}

	_log.Printf("[ %s ] %s", lvl, msg)

	if lvl == ERROR {
		return errors.New(msg)
	}

	return nil

}
