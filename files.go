package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var filesCommands = &Options{
	commands: []string{
		"add",
		"list",
		"remove",
		"copy",
	},
	descriptions: []string{
		"add a file to a day: files add <date> <filepath>",
		"list files in a day: files list <date>",
		"remove a file: files remove <date> <filename>",
		"copy a file out: files copy <date> <filename> <destpath>",
	},
}

func files(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) == 0 {
		return fmt.Errorf("command required. Options are:\n%v", filesCommands)
	}

	r, err := Apropos(x[0], filesCommands.commands)
	if err != nil {
		return err
	}

	switch r {
	case "add":
		return filesAdd(path, x[1:])
	case "list":
		return filesList(path, x[1:])
	case "remove":
		return filesRemove(path, x[1:])
	case "copy":
		return filesCopy(path, x[1:])
	default:
		return fmt.Errorf("invalid command %v", r)
	}
}

func parseDateArg(x []string) (year, month, day int, rest []string, err error) {
	if len(x) == 0 {
		err = fmt.Errorf("date required (year/month/day, today, yesterday, tomorrow)")
		return
	}

	dateArg := x[0]
	rest = x[1:]

	// check for a specific date first
	if f := strings.Split(dateArg, "/"); len(f) == 3 {
		year, err = strconv.Atoi(f[0])
		if err != nil {
			err = fmt.Errorf("invalid year: %v: %v", f[0], err)
			return
		}
		month, err = strconv.Atoi(f[1])
		if err != nil {
			err = fmt.Errorf("invalid month: %v: %v", f[1], err)
			return
		}
		day, err = strconv.Atoi(f[2])
		if err != nil {
			err = fmt.Errorf("invalid day: %v: %v", f[2], err)
			return
		}
	} else {
		dateCommands := []string{"today", "yesterday", "tomorrow"}
		var r string
		r, err = Apropos(dateArg, dateCommands)
		if err != nil {
			err = fmt.Errorf("%w\n%v", err, dateHelp)
			return
		}

		when := time.Now()
		switch r {
		case "today":
		case "yesterday":
			when = when.Add(-24 * time.Hour)
		case "tomorrow":
			when = when.Add(24 * time.Hour)
		}
		year = when.Year()
		month = int(when.Month())
		day = when.Day()
	}

	// validate date
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	valid := t.Year() == year && t.Month() == time.Month(month) && t.Day() == day
	if !valid {
		err = fmt.Errorf("invalid date: %v/%v/%v", year, month, day)
	}
	return
}

func filesAdd(path string, x []string) error {
	year, month, day, rest, err := parseDateArg(x)
	if err != nil {
		return err
	}

	if len(rest) == 0 {
		return fmt.Errorf("file path required")
	}
	if len(rest) > 1 {
		return fmt.Errorf("trailing commands: %v", rest[1:])
	}

	srcPath := rest[0]

	// check source file exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("cannot access file: %v", err)
	}
	if srcInfo.IsDir() {
		return fmt.Errorf("cannot add directory: %v", srcPath)
	}

	err = syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	datePath := filepath.Join(path, fmt.Sprintf("%v/%v/%v", year, month, day))
	err = os.MkdirAll(datePath, 0755)
	if err != nil {
		return err
	}

	// copy file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := filepath.Join(datePath, filepath.Base(srcPath))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = syncPush(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
}

func filesList(path string, x []string) error {
	year, month, day, rest, err := parseDateArg(x)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return fmt.Errorf("trailing commands: %v", rest)
	}

	err = syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	datePath := filepath.Join(path, fmt.Sprintf("%v/%v/%v", year, month, day))
	files, err := listFilesInDay(datePath)
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f)
	}
	return nil
}

func filesRemove(path string, x []string) error {
	year, month, day, rest, err := parseDateArg(x)
	if err != nil {
		return err
	}

	if len(rest) == 0 {
		return fmt.Errorf("filename required")
	}
	if len(rest) > 1 {
		return fmt.Errorf("trailing commands: %v", rest[1:])
	}

	filename := rest[0]

	err = syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	datePath := filepath.Join(path, fmt.Sprintf("%v/%v/%v", year, month, day))
	filePath := filepath.Join(datePath, filename)

	// check file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %v", filename)
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	err = syncPush(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
}

func filesCopy(path string, x []string) error {
	year, month, day, rest, err := parseDateArg(x)
	if err != nil {
		return err
	}

	if len(rest) < 2 {
		return fmt.Errorf("usage: files copy <date> <filename> <destpath>")
	}
	if len(rest) > 2 {
		return fmt.Errorf("trailing commands: %v", rest[2:])
	}

	filename := rest[0]
	destPath := rest[1]

	err = syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	datePath := filepath.Join(path, fmt.Sprintf("%v/%v/%v", year, month, day))
	srcPath := filepath.Join(datePath, filename)

	// check source file exists
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("file not found: %v", filename)
	}

	// copy file out
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// listFilesInDay returns all files in a day's directory except "entry"
func listFilesInDay(datePath string) ([]string, error) {
	entries, err := os.ReadDir(datePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && e.Name() != entryName {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

// hasFilesInDay returns true if a day has any files other than "entry"
func hasFilesInDay(datePath string) bool {
	files, err := listFilesInDay(datePath)
	if err != nil {
		return false
	}
	return len(files) > 0
}
