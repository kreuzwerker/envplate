package envplate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _delete(t *testing.T, name string) {

	if err := os.Remove(name); err != nil {
		t.Fatalf("Error while deleting '%s': %v", name, err)
	}

}

func _exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func _read(t *testing.T, name string) string {

	content, err := ioutil.ReadFile(name)

	if err != nil {
		t.Fatalf("Error while reading '%s': %v", name, err)
	}

	return string(content)

}

func _template(t *testing.T) (string, string) {

	tpl := `Database1=${DATABASE}
  Mode=${MODE}
  Database2=${DATABASE}
  Database4=${DATABASE1:-FOOBAR}
  Database5=${DATABASE:-FOOBAR}
  Database3=$FOOBAR`

	return _write(t, "parse.txt", tpl, 0644), tpl

}

func _write(t *testing.T, name, content string, mode os.FileMode) string {

	file, err := ioutil.TempFile("", name)

	if err != nil {
		t.Fatalf("Error while opening '%s': %v", name, err)
	}

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("Error while writing to '%s': %v", name, err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Error while closing '%s': %v", name, err)
	}

	if err := os.Chmod(file.Name(), mode); err != nil {
		t.Fatalf("Error while chmod '%s': %v", name, err)
	}

	return file.Name()

}

func TestCreateBackup(t *testing.T) {

	assert := assert.New(t)

	file := _write(t, "create-backup.txt", "hello world", 0644)
	defer _delete(t, file)

	backup := fmt.Sprintf("%s.bak", file)

	assert.False(_exists(backup))

	err := createBackup(file)
	defer _delete(t, backup)

	assert.NoError(err)

	content := _read(t, backup)

	assert.Equal("hello world", content)
	assert.Equal(filemode(file).String(), filemode(backup).String())

}

func TestApplyNoGlobs(t *testing.T) {

	var buf bytes.Buffer

	ErrorFunc = func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}

	globs := []string{
		"*.not-here",
		"*.not-there",
	}

	Apply(globs)

	assert.Equal(t, "[ ERROR ] Zero files matched passed globs '[*.not-here *.not-there]'", buf.String())

}

func TestFullParse(t *testing.T) {

	Config.Backup = true
	Config.DryRun = false
	Config.Verbose = true

	ErrorFunc = log.Panicf

	assert := assert.New(t)

	file, _ := _template(t)
	defer _delete(t, file)

	backup := fmt.Sprintf("%s.bak", file)

	err := parse(file)

	assert.NoError(err)
	assert.True(_exists(backup))
	assert.Equal(`Database1=db.example.com
  Mode=debug
  Database2=db.example.com
  Database4=FOOBAR
  Database5=db.example.com
  Database3=$FOOBAR`, _read(t, file))

}

func TestFilemode(t *testing.T) {

	file := _write(t, "filemode.text", "", 0654)
	defer _delete(t, file)

	mode := filemode(file)

	assert.Equal(t, 0654, mode)

}
