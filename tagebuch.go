package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func tagebuch(x []string) error {
	if len(x) == 0 {
		return fmt.Errorf("must provide tagebuch path")
	}

	p := filepath.Join(baseDir, x[0])

	return base(p, x[1:])
}

func validate(path string) error {
	_, err := os.Stat(filepath.Join(path, tagebuchMagic))
	if err != nil {
		return fmt.Errorf("invalid tagebuch: %v: %v", path, err)
	}
	return nil
}
