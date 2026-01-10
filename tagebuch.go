package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func tagebuch(x []string) error {
	if len(x) == 0 {
		return listJournals()
	}

	p := filepath.Join(baseDir, x[0])

	return base(p, x[1:])
}

func listJournals() error {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		fmt.Println("no journals found")
		return nil
	}

	var journals []string
	filepath.WalkDir(baseDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Name() == tagebuchMagic && !d.IsDir() {
			rel, err := filepath.Rel(baseDir, filepath.Dir(p))
			if err == nil {
				journals = append(journals, rel)
			}
		}
		return nil
	})

	if len(journals) == 0 {
		fmt.Println("no journals found")
		return nil
	}

	fmt.Println("available journals:")
	for _, j := range journals {
		fmt.Println("  " + j)
	}
	return nil
}

func validate(path string) error {
	_, err := os.Stat(filepath.Join(path, tagebuchMagic))
	if err != nil {
		return fmt.Errorf("invalid tagebuch: %v: %v", path, err)
	}
	return nil
}
