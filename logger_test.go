package envplate

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _redirect_logs(f func(func() string)) {

	var buf bytes.Buffer

	log.SetOutput(&buf)

	logs := func() string {

		logs := buf.String()
		buf.Reset()

		return logs

	}

	f(logs)

	log.SetOutput(os.Stderr)

}

func TestLog(t *testing.T) {

	Config.Verbose = false

	ErrorFunc = log.Panicf

	_redirect_logs(func(logs func() string) {

		assert := assert.New(t)

		Log(DEBUG, "Hello debug world")
		assert.Empty(logs())

		Log(INFO, "Hello info world")
		assert.NotEmpty(logs())

		assert.Panics(func() { Log(ERROR, "Hello error world") })
		assert.NotEmpty(logs())

		Config.Verbose = true

		Log(DEBUG, "Hello debug world")
		assert.Regexp(".+ \\[ DEBUG \\] Hello debug world", logs())

		Log(DEBUG, "Hello %s %s", "debug", "world")
		assert.Regexp(".+ \\[ DEBUG \\] Hello debug world", logs())

		Log(INFO, "Hello info world")
		assert.Regexp(".+ \\[ INFO \\] Hello info world", logs())

		assert.Panics(func() { Log(ERROR, "Hello error world") })
		assert.Regexp(".+ \\[ ERROR \\] Hello error world", logs())

	})

}
