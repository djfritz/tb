package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var calendarCommands = &Options{
	commands: []string{
		"last",
		"next",
	},
	descriptions: []string{
		"show last month's calendar",
		"show next month's calendar",
	},
}

func calendar(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	var year, month int

	// default to this month if no argument provided
	if len(x) == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
		return showCalendar(path, year, month, x)
	} else if f := splitSlash(x[0]); len(f) == 2 {
		// check for a specific year/month first
		year, err = strconv.Atoi(f[0])
		if err != nil {
			return fmt.Errorf("invalid year: %v: %v", f[0], err)
		}
		month, err = strconv.Atoi(f[1])
		if err != nil {
			return fmt.Errorf("invalid month: %v: %v", f[1], err)
		}
	} else {
		r, err := Apropos(x[0], calendarCommands.commands)
		if err != nil {
			return fmt.Errorf("%w\n%v", err, monthHelp)
		}

		when := time.Now()
		switch r {
		case "last":
			when = when.AddDate(0, -1, 0)
		case "next":
			when = when.AddDate(0, 1, 0)
		}
		year = when.Year()
		month = int(when.Month())
	}

	// validate month
	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month: %v", month)
	}

	return showCalendar(path, year, month, x[1:])
}

const monthHelp = "year/month : Specific month"

func splitSlash(s string) []string {
	var parts []string
	var current string
	for _, c := range s {
		if c == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}

func showCalendar(path string, year, month int, x []string) error {
	if len(x) != 0 {
		return fmt.Errorf("trailing commands: %v", x)
	}

	err := syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// find which days have entries and files
	monthPath := filepath.Join(path, fmt.Sprintf("%v/%v", year, month))
	daysWithEntries := make(map[int]bool)
	daysWithFiles := make(map[int]bool)

	filepath.WalkDir(monthPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Name() == entryName && !d.IsDir() {
			// check if file has content
			info, err := d.Info()
			if err == nil && info.Size() > 0 {
				// extract day from path
				dir := filepath.Dir(p)
				dayStr := filepath.Base(dir)
				day, err := strconv.Atoi(dayStr)
				if err == nil {
					daysWithEntries[day] = true
				}
			}
		}
		return nil
	})

	// check for files in each day
	for day := 1; day <= daysIn(month, year); day++ {
		dayPath := filepath.Join(monthPath, strconv.Itoa(day))
		if hasFilesInDay(dayPath) {
			daysWithFiles[day] = true
		}
	}

	// render calendar
	renderCalendar(year, month, daysWithEntries, daysWithFiles)
	return nil
}

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"
	colorBold  = "\033[1m"
)

func renderCalendar(year, month int, entries map[int]bool, files map[int]bool) {
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthName := t.Month().String()

	// header - width is 36 to match the calendar body (7 cells × 5 chars + 1)
	header := fmt.Sprintf("%s %d", monthName, year)
	fmt.Printf("┌──────────────────────────────────┐\n")
	fmt.Printf("│%s│\n", centerString(header, 34))
	fmt.Printf("├────┬────┬────┬────┬────┬────┬────┤\n")
	fmt.Printf("│ Su │ Mo │ Tu │ We │ Th │ Fr │ Sa │\n")
	fmt.Printf("├────┼────┼────┼────┼────┼────┼────┤\n")

	// find first day of month and number of days
	firstWeekday := int(t.Weekday())
	daysInMonth := daysIn(month, year)

	// print calendar grid
	day := 1
	for week := 0; week < 6; week++ {
		if day > daysInMonth {
			break
		}
		fmt.Print("│")
		for weekday := 0; weekday < 7; weekday++ {
			if week == 0 && weekday < firstWeekday {
				fmt.Print("    │")
			} else if day > daysInMonth {
				fmt.Print("    │")
			} else {
				hasEntry := entries[day]
				hasFiles := files[day]
				if hasEntry && hasFiles {
					// both entry and files: green with * marker
					fmt.Printf(" %s%s%2d%s*%s│", colorBold, colorGreen, day, colorBlue, colorReset)
				} else if hasEntry {
					// entry only: green
					fmt.Printf(" %s%s%2d%s │", colorBold, colorGreen, day, colorReset)
				} else if hasFiles {
					// files only: blue with * marker
					fmt.Printf(" %s%s%2d*%s│", colorBold, colorBlue, day, colorReset)
				} else {
					fmt.Printf(" %2d │", day)
				}
				day++
			}
		}
		fmt.Println()
	}
	fmt.Printf("└────┴────┴────┴────┴────┴────┴────┘\n")
}

func daysIn(month, year int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func centerString(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	left := (width - len(s)) / 2
	right := width - len(s) - left
	result := ""
	for i := 0; i < left; i++ {
		result += " "
	}
	result += s
	for i := 0; i < right; i++ {
		result += " "
	}
	return result
}
