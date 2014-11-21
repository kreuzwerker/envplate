package envplate

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitArgs(t *testing.T) {

	assert := assert.New(t)

	var tests = []struct {
		in       []string
		flagArgs []string
		execArgs []string
	}{
		{
			[]string{"-a", "-b"},
			[]string{"-a", "-b"},
			[]string{},
		},
		{
			[]string{"-a", "-b", "--"},
			[]string{"-a", "-b"},
			[]string{},
		},
		{
			[]string{"-a", "-b", "--", "-c"},
			[]string{"-a", "-b"},
			[]string{"-c"},
		},
	}

	for _, test := range tests {

		os.Args = test.in

		flagArgs, execArgs := SplitArgs()

		msg := fmt.Sprintf("%v", test.in)

		assert.Equal(test.flagArgs, flagArgs, msg)
		assert.Equal(test.execArgs, execArgs, msg)

	}

}
