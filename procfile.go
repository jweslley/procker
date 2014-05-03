package procker

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var procfileRegexp = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

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
			return nil, fmt.Errorf("procker: parse error: invalid line found: '%s'", line)
		}

		name, command := matches[1], matches[2]
		p[name] = NewProcess(name, command)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("procker: parse error: %s", err)
	}

	return p, nil
}
