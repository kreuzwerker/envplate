package envplate

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func _delete(t *testing.T, name string) {

	if err := os.Remove(name); err != nil {
		t.Fatal(fmt.Errorf("error while deleting '%s': %v", name, err))
	}

}

func _exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func _read(t *testing.T, name string) string {

	content, err := ioutil.ReadFile(name)

	if err != nil {
		t.Fatal(fmt.Errorf("error while reading '%s': %v", name, err))
	}

	return string(content)

}

func _redirect_logs(f func(func() string, func() string)) {

	var (
		logBuf bytes.Buffer
		rawBuf bytes.Buffer
	)

	defer func(w io.Writer) {
		Logger.Out = w
	}(Logger.Out)

	defer func(w io.Writer) {
		log.SetOutput(w)
	}(os.Stderr)

	log.SetOutput(&logBuf)
	Logger.Out = &rawBuf

	logs := func(buf *bytes.Buffer) func() string {

		return func() string {

			logs := buf.String()
			buf.Reset()

			return logs

		}

	}

	f(logs(&logBuf), logs(&rawBuf))

}

func _restore(file string) {

	backup := fmt.Sprintf("%s.bak", file)

	os.Remove(file)
	os.Link(backup, file)
	os.Remove(backup)

}

func _tmpdir(t *testing.T, f func(string)) {

	dir, err := ioutil.TempDir("", "")

	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	f(dir)

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
