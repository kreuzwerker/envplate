package envplate

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	NoKeyDefined = ""
)

type Envplate struct {
	Backup, Debug, Verbose *bool
}

func (e *Envplate) Apply(globs []string) {

	for _, pattern := range globs {

		files, err := filepath.Glob(pattern)

		if err != nil {
			e.log(ERROR, err.Error())
		}

		for _, name := range files {

			if err := e.parse(name); err != nil {
				e.log(ERROR, "Error while parsing '%s': %v", name, err)
			}

		}

	}

}

func (e *Envplate) createBackup(file string) error {

	source, err := os.Open(file)

	if err != nil {
		return err
	}

	defer source.Close()

	target, err := os.Create(fmt.Sprintf("%s.bak", file))

	if err != nil {
		return err
	}

	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return err
	}

	return nil

}

func (e *Envplate) parse(file string) error {

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return fmt.Errorf("Cannot open %s: %v", file, err)
	}

	e.log(DEBUG, "Parsing environment references in '%s'", file)

	parsed := os.Expand(string(content), func(key string) string {

		value := os.Getenv(key)

		e.log(DEBUG, "Expanding reference to '%s' to value '%s'", key, value)

		if value == NoKeyDefined {
			e.log(ERROR, "'%s' requires undeclared environment variable '%s'", file, key)
		}

		return value

	})

	if *e.Debug {
		e.log(INFO, "Expanding all references in '%s' would look like this:\n%s", file, parsed)
	} else {

		if *e.Backup {

			e.log(DEBUG, "Creating backup of '%s'", file)

			if err := e.createBackup(file); err != nil {
				return err
			}

		}

		// TODO: take over permissions from original file
		return ioutil.WriteFile(file, []byte(parsed), 0644)

	}

	return nil

}

type logLevel string

const (
	DEBUG logLevel = "DEBUG"
	INFO           = "INFO"
	ERROR          = "ERROR"
)

func (e *Envplate) log(lvl logLevel, msg string, args ...interface{}) {

	if lvl == DEBUG && !*e.Verbose {
		return
	}

	msg = fmt.Sprintf("[ %s ] %s", lvl, msg)

	if lvl == ERROR {
		log.Fatalf(msg, args...)
	} else {
		log.Printf(msg, args...)
	}

}
