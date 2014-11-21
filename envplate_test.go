package envplate

import (
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

	tpl := `
  Database1=${DATABASE}
  Mode=${MODE}
  Database2=${DATABASE}
  `

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
	assert.Equal(filemode(file), filemode(backup))

}

func TestParse(t *testing.T) {

	Config.Backup = true
	Config.DryRun = false
	Config.Verbose = true

	ErrorFunc = log.Panicf

	assert := assert.New(t)

	file, tpl := _template(t)
	defer _delete(t, file)

	backup := fmt.Sprintf("%s.bak", file)
	fmt.Println(backup)

	//	time.Sleep(5 * time.Second)

	err := parse(file)

	assert.NoError(err)

	assert.True(_exists(backup))
	assert.NotEqual(tpl, _read(t, file), "content unchanged")

}

func TestFilemode(t *testing.T) {

	file := _write(t, "filemode.text", "", 0654)
	defer _delete(t, file)

	mode := filemode(file)

	assert.Equal(t, 0654, mode)

}
