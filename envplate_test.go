package envplate

import (
	"fmt"
	"path/filepath"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.Verbose = true
}

func TestApplyNoGlobs(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{}

	_tmpdir(t, func(tempdir string) {

		globs := []string{
			tempdir,
			filepath.Join(tempdir, "*.not-here"),
			filepath.Join(tempdir, "*.not-there"),
		}

		_redirect_logs(func(logs func() string, raw func() string) {

			handler.Apply(globs)

			msg := fmt.Sprintf("[ ERROR ] Zero files matched passed globs '[%s %s %s]'",
				globs[0],
				globs[1],
				globs[2])

			assert.Contains(logs(), msg)

		})

	})

}

func TestCapture(t *testing.T) {

	assert := assert.New(t)

	var tt = []struct {
		in, e, v, s, d string
	}{
		{"${FOO}", "", "FOO", noDefaultDefined, noDefaultDefined},
		{"${FOO:-bar}", "", "FOO", ":-", "bar"},
		{"${FOO:-at the bar}", "", "FOO", ":-", "at the bar"},
		{"${FOO_3000:-near the bar}", "", "FOO_3000", ":-", "near the bar"},
		{"${FOO:--1}", "", "FOO", ":-", "-1"},
		{"${FOO:-http://www.example.com/bar/gar/war?a=b}", "", "FOO", ":-", "http://www.example.com/bar/gar/war?a=b"},
		{`\${FOO}`, `\`, "FOO", noDefaultDefined, noDefaultDefined},
		{`\\${FOO:-bar}`, `\\`, "FOO", ":-", "bar"},
		{`\\\${FOO:-bar}`, `\\\`, "FOO", ":-", "bar"},
		{`\\\\${FOO:-bar}`, `\\\\`, "FOO", ":-", "bar"},
		{"foo${FOO}", "", "FOO", noDefaultDefined, noDefaultDefined},
	}

	for _, tt := range tt {

		e, v, s, d := capture(tt.in)

		assert.Equal(tt.e, e)
		assert.Equal(tt.v, v)
		assert.Equal(tt.s, s)
		assert.Equal(tt.d, d)

	}

}

func TestDryRun(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		DryRun: true,
	}

	var (
		file = "test/template1.txt"
	)

	_redirect_logs(func(logs func() string, raw func() string) {

		err := handler.parse(file)
		assert.NoError(err)

		assert.Equal(`Database1=${DATABASE}
Mode=${MODE}
Null=${NULL}
Database2=${DATABASE}
Database3=$NOT_A_VARIABLE
Database4=${ANOTHER_DATABASE:-db2.example.com}
Database5=${DATABASE:-db2.example.com}
`, _read(t, file))

		assert.Equal(`Database1=db.example.com
Mode=debug
Null=
Database2=db.example.com
Database3=$NOT_A_VARIABLE
Database4=db2.example.com
Database5=db.example.com
`, raw())

	})

}

func TestEscape(t *testing.T) {

	assert := assert.New(t)

	var tt = []struct {
		in, e string
	}{
		{`foo`, notAnEscapeSequence},
		{`${FOO}`, notAnEscapeSequence},
		{`\${FOO}`, `${FOO}`},
		{`\\${FOO}`, notAnEscapeSequence},
		{`\\\${FOO:-bar}`, `\${FOO:-bar}`},
		{`\\\\${FOO}`, notAnEscapeSequence},
		{`\\\\\${FOO}`, `\\${FOO}`},
		{`\\\\\\${FOO}`, notAnEscapeSequence},
		{`\\\\\\\${FOO:-bar}`, `\\\${FOO:-bar}`},
	}

	for _, tt := range tt {
		esc := escape(tt.in)
		assert.Equal(tt.e, esc)
	}

}

func TestFullParse(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		Backup: true,
	}

	var (
		file = "test/template1.txt"
	)

	defer _restore(file)

	err := handler.parse(file)
	assert.NoError(err)

	assert.Equal(`Database1=db.example.com
Mode=debug
Null=
Database2=db.example.com
Database3=$NOT_A_VARIABLE
Database4=db2.example.com
Database5=db.example.com
`, _read(t, file))

}

func TestFullParseDefaults(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		Backup: true,
	}

	var (
		file = "test/template2.txt"
	)

	defer _restore(file)

	err := handler.parse(file)
	assert.NoError(err)

	assert.Equal(`Double1=db.example.com Double2=db.example.com
Double3=db.example.com Double4=db.example.com
DoubleDefault1=db2-example.com DoubleDefault2=db2-example.com
EmptyDefault=
`, _read(t, file))

}

func TestFullParseEscapes(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		Backup: true,
	}

	var (
		file = "test/template3.txt"
	)

	defer _restore(file)

	err := handler.parse(file)
	assert.NoError(err)

	assert.Equal(`db.example.com
Escaped1=${DATABASE} EscapedDefault1=${DATABASE:-db2.example.com} EscapedDefaultReplaced1=${ANOTHER_DATABASE:-db2.example.com}
NoEscape1=\db.example.com NoEscapeDefault1=\db.example.com NoEscapeDefaultReplaced1=\db2.example.com
Escaped2=\${DATABASE} EscapedDefault2=\${DATABASE:-db2.example.com} EscapedDefaultReplaced2=\${ANOTHER_DATABASE:-db2.example.com}
NoEscape2=\\db.example.com NoEscapeDefault2=\\db.example.com NoEscapeDefaultReplaced2=\\db2.example.com
Escaped3=\\${DATABASE} EscapedDefault3=\\${DATABASE:-db2.example.com} EscapedDefaultReplaced3=\\${ANOTHER_DATABASE:-db2.example.com}
NoEscape3=\\\db.example.com NoEscapeDefault3=\\\db.example.com NoEscapeDefaultReplaced3=\\\db2.example.com
`, _read(t, file))

}

func TestStrictParse(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		Backup: true,
		Strict: true,
	}

	var (
		file = "test/template4.txt"
	)

	err := handler.parse(file)
	assert.Error(err)

	assert.Equal("'test/template4.txt' requires undeclared environment variable 'ANOTHER_DATABASE', but cannot use default 'db2.example.com' (strict-mode)", err.Error())

}

func TestAbortOnParseErrors(t *testing.T) {

	assert := assert.New(t)
	handler := &Handler{
		Backup: true,
		Strict: true,
	}

	var (
		file     = "test/template4.txt"
		template = _read(t, file)
	)

	err := handler.parse(file)
	assert.Error(err)

	assert.Equal(template, _read(t, file))
	assert.NoFileExists(fmt.Sprintf("%s.bak", file))

}
