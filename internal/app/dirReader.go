package app

import (
	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
)

type dirReader struct {
	dir    string
	filter string
}

func newDirReader(d string, f string) *dirReader {
	return &dirReader{dir: d, filter: f}
}

func (r *dirReader) ReadExecutables() ([]string, error) {
	lines, err := filepathx.Glob(r.filter)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.Errorf("No google executables from filter '%s'", r.filter)
	}
	return lines, nil
}
