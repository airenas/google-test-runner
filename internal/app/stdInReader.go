package app

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type stdInReader struct {
}

func newStdInReader(f string) *stdInReader {
	return &stdInReader{}
}

func (r *stdInReader) ReadExecutables() ([]string, error) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimSpace(s)
		if s != "" {
			lines = append(lines, s)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("No google executables from stdin")
	}
	return lines, nil
}
