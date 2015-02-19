package envplate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	NoDefaultDefined = ""
	NoKeyDefined     = ""
)

var exp = regexp.MustCompile(`(.?)\$\{(.+?)(?:\:\-(.+?))?\}`)

func Apply(globs []string) {

	matches := false

	for _, pattern := range globs {

		files, err := filepath.Glob(pattern)

		if err != nil {
			Log(ERROR, err.Error())
		}

		for _, name := range files {

			if info, _ := os.Stat(name); info.IsDir() {
				continue
			}

			matches = true

			if err := parse(name); err != nil {
				Log(ERROR, "Error while parsing '%s': %v", name, err)
			}

		}

	}

	if !matches {
		Log(ERROR, "Zero files matched passed globs '%v'", globs)
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

	if err := os.Chmod(target.Name(), filemode(source.Name())); err != nil {
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

	parsed := exp.ReplaceAllStringFunc(string(content), func(match string) string {

		var (
			pre, key, def = capture(match)
			value    = os.Getenv(key)
		)

                if pre == `\` {

                  return match[1:len(match)]

                }

		if value == NoKeyDefined {

			if def == NoDefaultDefined {
				Log(ERROR, "'%s' requires undeclared environment variable '%s', no default is given", file, key)
			} else {

				if Config.Strict {
					Log(ERROR, "'%s' requires undeclared environment variable '%s', but cannot use default '%s' (strict-mode)", file, key, def)
				} else {
					Log(DEBUG, "'%s' requires undeclared environment variable '%s', using default '%s'", file, key, def)
					value = def
				}

			}

		} else {
			Log(DEBUG, "Expanding reference to '%s' to value '%s'", key, value)
		}

		return pre + value

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

func capture(s string) (pre, key, def string) {

	matches := exp.FindStringSubmatch(s)

        pre = matches[1]
	key = matches[2]
	def = matches[3]

	return pre, key, def

}

func filemode(file string) os.FileMode {

	fileinfo, err := os.Stat(file)

	if err != nil {
		Log(ERROR, "Cannot stat '%s': %v", file, err)
	}

	return fileinfo.Mode()

}
