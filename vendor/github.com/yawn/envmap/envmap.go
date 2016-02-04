package envmap

import (
	"os"
	"strings"
)

const separator = "="

// Envmap is a mapping of environment keys to values
type Envmap map[string]string

// Export exports the keys and values defined in this Envmap to the
// actual environment
func (e Envmap) Export() {

	for k, v := range e {
		os.Setenv(k, v)
	}

}

// Import creates an Envmap from the actual environment.
func Import() Envmap {
	return ToMap(os.Environ())
}

// Join builds a environment variable declaration out of seperate
// key and value strings
func Join(k, v string) string {
	return strings.Join([]string{k, v}, separator)
}

// ToEnv converts a map of environment variables to a slice
// of key=value strings
func (e Envmap) ToEnv() (env []string) {

	for k, v := range e {
		env = append(env, Join(k, v))
	}

	return

}

// ToMap converts a slice of environment variables to a map
// of environment variables
func ToMap(env []string) (m Envmap) {

	m = make(map[string]string)

	for _, e := range env {

		s := strings.Split(e, separator)

		var (
			key = s[0]
			val = strings.Join(s[1:], separator)
		)

		m[key] = val

	}

	return

}
