package local

import (
	"fmt"
	"os"
)

var (
	ErrEmptyPath = fmt.Errorf(".Path is empty")
)

func New(path string) (*os.File, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", path)
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
