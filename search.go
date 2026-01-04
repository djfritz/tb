package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var searchTerm string

func search(path string, x []string) error {
	// no need to validate because we support any path

	if len(x) == 0 {
		return fmt.Errorf("search term required")
	}

	if len(x) != 1 {
		return fmt.Errorf("trailing commands: %v", x[2:])
	}

	searchTerm = x[0]

	return filepath.WalkDir(path, searchFunc)
}

func searchFunc(path string, d fs.DirEntry, err error) error {
	base := filepath.Base(path)
	if base == entryName {
		return searchEntry(path)
	}
	return nil
}

func searchEntry(path string) error {
	// grep!

	cmd := exec.Command("grep", "-H", searchTerm, path)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ExitCode() == 1 {
				// ignore, this just means that no lines were cu if exiterr, ok := err.(*exec.ExitError); ok {
				return nil
			}
		}
		return err
	}
	return nil
}
