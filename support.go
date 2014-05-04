package procker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// PrefixedWriter implements prefixed output for an io.Writer object.
type PrefixedWriter struct {
	Prefix string
	writer io.Writer
	inline bool
}

// NewPrefixedWriter creates a PrefixedWriter
func NewPrefixedWriter(w io.Writer, prefix string) io.Writer {
	return &PrefixedWriter{Prefix: prefix, writer: w}
}

// Writes a Prefix string before writing to the underlying writer.
func (w *PrefixedWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if !w.inline {
			io.WriteString(w.writer, w.Prefix)
		}
		w.writer.Write([]byte{b})
		w.inline = b != '\n'
	}
	return len(p), nil
}

var procfileRegexp = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

// ParseProcfile parses io.Reader into a Process's map.
// Read more about Procfiles: https://devcenter.heroku.com/articles/procfile
func ParseProcfile(r io.Reader) (map[string]*Process, error) {
	p := make(map[string]*Process)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		matches := procfileRegexp.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("procker: parse procfile error: invalid line found: '%s'", line)
		}

		name, command := matches[1], matches[2]
		p[name] = NewProcess(name, command)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("procker: parse procfile error: %s", err)
	}

	return p, nil
}

// ParseEnv parses io.Reader into an arrays of strings
// representing the environment, in the form "key=value".
func ParseEnv(r io.Reader) ([]string, error) {
	sysenv := env2Map(os.Environ())
	localenv := make(map[string]string)
	mapping := func(key string) string {
		value, ok := localenv[key]
		if !ok {
			value, ok = sysenv[key]
		}
		return value
	}

	env := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry := strings.TrimSpace(scanner.Text())
		if len(entry) == 0 {
			continue
		}

		pair := strings.SplitN(entry, "=", 2)
		key := pair[0]
		value := os.Expand(pair[1], mapping)
		localenv[key] = value
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("procker: parse env error: %s", err)
	}

	return env, nil
}

func env2Map(env []string) map[string]string {
	m := make(map[string]string)
	for _, value := range env {
		pair := strings.SplitN(value, "=", 2)
		if len(pair) == 2 {
			m[pair[0]] = pair[1]
		}
	}
	return m
}
