package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const tagebuchMagic = ".tagebuch"

func initTagebuch(path string, x []string) error {
	if len(x) != 0 {
		return fmt.Errorf("trailing commands: %v", x)
	}

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("path %v exists", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	pf := filepath.Join(path, tagebuchMagic)
	f, err := os.OpenFile(pf, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
