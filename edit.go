package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var editCommands = &Options{
	commands: []string{
		"today",
		"yesterday",
		"tomorrow",
	},
	descriptions: []string{
		"edit today's entry",
		"edit yesterday's entry",
		"edit tomorrow's entry",
	},
}

const (
	dateHelp  = "year/month/day : Specific date"
	entryName = "entry"
)

func edit(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) == 0 {
		return fmt.Errorf("command required. Options are:\n%v\n%v", editCommands, dateHelp)
	}

	var year, month, day int

	// check for a specific date first
	if f := strings.Split(x[0], "/"); len(f) == 3 {
		year, err = strconv.Atoi(f[0])
		if err != nil {
			return fmt.Errorf("invalid year: %v: %v", f[0], err)
		}
		month, err = strconv.Atoi(f[1])
		if err != nil {
			return fmt.Errorf("invalid month: %v: %v", f[1], err)
		}
		day, err = strconv.Atoi(f[2])
		if err != nil {
			return fmt.Errorf("invalid day: %v: %v", f[2], err)
		}

	} else {
		r, err := Apropos(x[0], editCommands.commands)
		if err != nil {
			return fmt.Errorf("%w\n%v", err, dateHelp)
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

	// validate
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	valid := t.Year() == year && t.Month() == time.Month(month) && t.Day() == day
	if !valid {
		return fmt.Errorf("invalid date: %v/%v/%v", year, month, day)
	}

	datePath := filepath.Join(path, fmt.Sprintf("%v/%v/%v", year, month, day))

	return editDate(path, datePath, x[1:])
}

func editDate(path, datePath string, x []string) error {
	if len(x) != 0 {
		return fmt.Errorf("trailing commands: %v", x)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("$EDITOR not set")
	}

	err := syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = os.MkdirAll(datePath, 0755)
	if err != nil {
		return err
	}

	filename := filepath.Join(datePath, entryName)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = syncPush(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
}
