package envplate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/paulrosania/go-charset/charset"
	"github.com/yawn/envmap"
)

const (
	noDefaultDefined    = ""
	notAnEscapeSequence = ""
)

type Handler struct {
	Backup  bool
	DryRun  bool
	Strict  bool
	Charset string
}

var exp = regexp.MustCompile(`(\\*)\$\{(.+?)(?:(\:\-)(.*?))?\}`)

func (h *Handler) Apply(globs []string) error {

	matches := false

	for _, pattern := range globs {

		files, err := filepath.Glob(pattern)

		if err != nil {
			return err
		}

		for _, name := range files {

			if info, _ := os.Stat(name); info.IsDir() {
				continue
			}

			matches = true

			if err := h.parse(name); err != nil {
				return Log(ERROR, "Error while parsing '%s': %v", name, err)
			}

		}

	}

	if !matches {
		return Log(ERROR, "Zero files matched passed globs '%v'", globs)
	}

	return nil

}

func (h *Handler) parse(file string) error {

	env := envmap.Import()
	content, err := os.ReadFile(file)

	if err != nil {
		return Log(ERROR, "Cannot open %s: %v", file, err)
	}

	Log(DEBUG, "Parsing environment references in '%s'", file)

	var errors []error

	parsed := exp.ReplaceAllStringFunc(string(content), func(match string) string {

		var (
			esc, key, sep, def = capture(match)
			value, keyDefined  = env[key]
		)

		if len(esc)%2 == 1 {

			escaped := escape(match)

			if escaped == notAnEscapeSequence {
				errors = append(errors, Log(ERROR, "Tried to escape '%s', but was no escape sequence", content))
			}

			encodedValue, err := convertToCharset(escaped, h.Charset)
			if err != nil {
				errors = append(errors, Log(ERROR, "Tried to convert string '%s' to charset '%s'  but an error ocourred: %v", encodedValue, h.Charset, err))
				return value
			}
			return encodedValue

		}

		if !keyDefined {

			if sep == noDefaultDefined {
				errors = append(errors, Log(ERROR, "'%s' requires undeclared environment variable '%s', no default is given", file, key))
			} else {

				if h.Strict {
					errors = append(errors, Log(ERROR, "'%s' requires undeclared environment variable '%s', but cannot use default '%s' (strict-mode)", file, key, def))
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

		encodedValue, err := convertToCharset(value, h.Charset)
		if err != nil {
			errors = append(errors, Log(ERROR, "Tried to convert string '%s' to charset '%s'  but an error ocourred: %v", encodedValue, h.Charset, err))
			return value
		}
		return encodedValue

	})

	if len(errors) > 0 {
		return errors[0]
	}

	if h.DryRun {
		Log(DEBUG, "Expanding all references in '%s' without doing anything (dry-run)", file)
		Log(RAW, parsed)
	} else {

		if h.Backup {

			Log(DEBUG, "Creating backup of '%s'", file)

			if err := createBackup(file); err != nil {
				return err
			}

		}

		if err := saveFile(file, parsed, h.Charset); err != nil {
			return err
		}

	}

	return nil

}

func saveFile(file string, parsed string, cs string) error {
	mode, err := filemode(file)
	if err != nil {
		return err
	}

	if err := os.WriteFile(file, []byte(parsed), mode); err != nil {
		return err
	}

	return nil
}

func convertToCharset(data string, charSet string) (string, error) {
	if charSet == "" {
		return data, nil
	}

	buf := new(bytes.Buffer)
	w, err := charset.NewWriter(charSet, buf)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(w, data)
	w.Close()
	return buf.String(), nil
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
		return notAnEscapeSequence
	}

	bss := matches[1]

	if len(bss)%2 != 1 {
		return notAnEscapeSequence
	}

	parsedBss := bss[:len(bss)-1][:(len(bss)-1)/2]

	escaped = parsedBss + matches[2]

	Log(DEBUG, "Substituting escaped sequence '%s' with '%s'", s, escaped)

	return escaped

}
