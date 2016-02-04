package envplate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {

	reset := func(v bool) {
		Logger.Verbose = v
	}

	defer reset(Logger.Verbose)
	reset(false)

	_redirect_logs(func(logs func() string, raw func() string) {

		assert := assert.New(t)

		err := Log(RAW, "Hello raw world")
		assert.NoError(err)
		assert.Empty(logs())
		assert.Equal(raw(), "Hello raw world")

		err = Log(DEBUG, "Hello debug world")
		assert.NoError(err)
		assert.Empty(logs())

		err = Log(INFO, "Hello info world")
		assert.NoError(err)
		assert.NotEmpty(logs())

		err = Log(ERROR, "Hello error world")
		assert.Error(err)
		assert.Equal(errors.New("Hello error world"), err)
		assert.NotEmpty(logs())

		Logger.Verbose = true

		err = Log(DEBUG, "Hello debug world")
		assert.NoError(err)
		assert.Regexp(".+ \\[ DEBUG \\] Hello debug world", logs())

		err = Log(DEBUG, "Hello %s %s", "debug", "world")
		assert.NoError(err)
		assert.Regexp(".+ \\[ DEBUG \\] Hello debug world", logs())

		err = Log(INFO, "Hello info world")
		assert.NoError(err)
		assert.Regexp(".+ \\[ INFO \\] Hello info world", logs())

		err = Log(ERROR, "Hello error world")
		assert.Error(err)
		assert.Equal(errors.New("Hello error world"), err)
		assert.Regexp(".+ \\[ ERROR \\] Hello error world", logs())

	})

}
