package app

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type fileReader struct {
	file string
}

func newFileReader(f string) *fileReader {
	return &fileReader{file: f}
}

func (r *fileReader) ReadExecutables() ([]string, error) {
	content, err := ioutil.ReadFile(r.file)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file "+r.file)
	}
	lines := strings.Split(string(content), "\n")
	lines = clean(lines)
	if len(lines) == 0 {
		return nil, errors.New("No executables in " + r.file)
	}
	return lines, nil
}

func clean(data []string) []string {
	res := make([]string, 0)
	for _, d := range data {
		d = strings.TrimSpace(d)
		if d != "" {
			res = append(res, d)
		}
	}
	return res
}
