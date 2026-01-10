package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func list(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) != 0 {
		return fmt.Errorf("trailing commands: %v", x)
	}

	err = syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var entries []string

	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Name() == entryName && !d.IsDir() {
			info, err := d.Info()
			if err == nil && info.Size() > 0 {
				// extract year/month/day from path
				rel, err := filepath.Rel(path, p)
				if err != nil {
					return nil
				}
				// rel is like "2026/1/6/entry"
				dir := filepath.Dir(rel)
				parts := splitSlash(dir)
				if len(parts) == 3 {
					year, err1 := strconv.Atoi(parts[0])
					month, err2 := strconv.Atoi(parts[1])
					day, err3 := strconv.Atoi(parts[2])
					if err1 == nil && err2 == nil && err3 == nil {
						entries = append(entries, fmt.Sprintf("%d/%d/%d", year, month, day))
					}
				}
			}
		}
		return nil
	})

	// sort chronologically
	sort.Slice(entries, func(i, j int) bool {
		return compareDates(entries[i], entries[j])
	})

	for _, e := range entries {
		fmt.Println(e)
	}

	return nil
}

func compareDates(a, b string) bool {
	pa := splitSlash(a)
	pb := splitSlash(b)
	ya, _ := strconv.Atoi(pa[0])
	yb, _ := strconv.Atoi(pb[0])
	if ya != yb {
		return ya < yb
	}
	ma, _ := strconv.Atoi(pa[1])
	mb, _ := strconv.Atoi(pb[1])
	if ma != mb {
		return ma < mb
	}
	da, _ := strconv.Atoi(pa[2])
	db, _ := strconv.Atoi(pb[2])
	return da < db
}
