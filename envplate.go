package envplate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	NoKeyDefined = ""
)

func Apply(globs []string) {

	for _, pattern := range globs {

		files, err := filepath.Glob(pattern)

		if err != nil {
			Log(ERROR, err.Error())
		}

		for _, name := range files {

			if err := parse(name); err != nil {
				Log(ERROR, "Error while parsing '%s': %v", name, err)
			}

		}

	}

}

func createBackup(file string) error {

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

func parse(file string) error {

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return fmt.Errorf("Cannot open %s: %v", file, err)
	}

	Log(DEBUG, "Parsing environment references in '%s'", file)

	parsed := os.Expand(string(content), func(key string) string {

		value := os.Getenv(key)

		Log(DEBUG, "Expanding reference to '%s' to value '%s'", key, value)

		if value == NoKeyDefined {
			Log(ERROR, "'%s' requires undeclared environment variable '%s'", file, key)
		}

		return value

	})

	if Config.DryRun {
		Log(INFO, "Expanding all references in '%s' would look like this:\n%s", file, parsed)
	} else {

		if Config.Backup {

			Log(DEBUG, "Creating backup of '%s'", file)

			if err := createBackup(file); err != nil {
				return err
			}

		}

		return ioutil.WriteFile(file, []byte(parsed), filemode(file))

	}

	return nil

}

func filemode(file string) os.FileMode {

	fileinfo, err := os.Stat(file)

	if err != nil {
		Log(ERROR, "Cannot state '%s': %v", file, err)
	}

	return fileinfo.Mode()

}
