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

var printCommands = &Options{
	commands: []string{
		"today",
		"yesterday",
		"tomorrow",
	},
	descriptions: []string{
		"print today's entry",
		"print yesterday's entry",
		"print tomorrow's entry",
	},
}

func printEntry(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) == 0 {
		return fmt.Errorf("command required. Options are:\n%v\n%v", printCommands, dateHelp)
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
		r, err := Apropos(x[0], printCommands.commands)
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

	return printDate(path, datePath, x[1:])
}

func printDate(path string, datePath string, x []string) error {
	if len(x) != 0 {
		return fmt.Errorf("trailing commands: %v", x)
	}

	err := syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	f, err := os.Open(filepath.Join(datePath, entryName))
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, f)
	f.Close()

	return nil
}
