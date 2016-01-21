package envplate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	NoDefaultDefined    = ""
	NotAnEscapeSequence = ""
	DefaultValueSyntax  = ":-"
)

var exp = regexp.MustCompile(`(\\*)\$\{(.+?)(?:(\:\-)(.*?))?\}`)

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

	env := envmap()
	content, err := ioutil.ReadFile(file)

	if err != nil {
		return fmt.Errorf("Cannot open %s: %v", file, err)
	}

	Log(DEBUG, "Parsing environment references in '%s'", file)

	parsed := exp.ReplaceAllStringFunc(string(content), func(match string) string {

		var (
			esc, key, sep, def     = capture(match)
			value, keyDefined = env[key]
		)

		if len(esc)%2 == 1 {

			escaped := escape(match)

			if escaped == NotAnEscapeSequence {

				Log(ERROR, "Tried to escape '%s', but was no escape sequence")

			}

			return escaped

		}

		if !keyDefined {

			if sep == NoDefaultDefined {
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

		if len(esc) > 0 {
			value = esc[:len(esc)/2] + value
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

func capture(s string) (esc, key, sep, def string) {

	matches := exp.FindStringSubmatch(s)

	esc = matches[1]
	key = matches[2]
	sep = matches[3]
	def = matches[4]

	return esc, key, sep, def

}

func escape(s string) (escaped string) {

	expEscaped := regexp.MustCompile(`(\\+)(.*)`)
	matches := expEscaped.FindStringSubmatch(s)

	if matches == nil {

		return NotAnEscapeSequence

	}

	bss := matches[1]

	if len(bss)%2 != 1 {

		return NotAnEscapeSequence

	}

	parsedBss := bss[:len(bss)-1][:(len(bss)-1)/2]

	escaped = parsedBss + matches[2]

	Log(DEBUG, "Substituting escaped sequence '%s' with '%s'", s, escaped)

	return escaped

}

func envmap() (m map[string]string) {

	m = make(map[string]string)

	for _, e := range os.Environ() {

		s := strings.Split(e, "=")

		key := s[0]
		val := strings.Join(s[1:], "=")

		m[key] = val

	}

	return

}

func filemode(file string) os.FileMode {

	fileinfo, err := os.Stat(file)

	if err != nil {
		Log(ERROR, "Cannot stat '%s': %v", file, err)
	}

	return fileinfo.Mode()

}
