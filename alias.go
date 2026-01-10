package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const tagebuchAliases = "aliases"

var aliasCommands = &Options{
	commands: []string{
		"add",
		"remove",
	},
	descriptions: []string{
		"add an alias: alias add <name> <year/month/day>",
		"remove an alias by name",
	},
}

func alias(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) == 0 {
		// list aliases
		return aliasList(path)
	}

	r, err := Apropos(x[0], aliasCommands.commands)
	if err != nil {
		return err
	}

	switch r {
	case "add":
		return aliasAdd(path, x[1:])
	case "remove":
		return aliasRemove(path, x[1:])
	default:
		return fmt.Errorf("invalid command %v", r)
	}
}

type aliases struct {
	a map[string]string // name -> date (year/month/day)
}

func (a *aliases) String() string {
	var ret string
	for name, date := range a.a {
		ret += fmt.Sprintf("%v -> %v\n", name, date)
	}
	return strings.TrimSpace(ret)
}

func (a *aliases) save(path string) error {
	pt := filepath.Join(path, tagebuchAliases)
	f, err := os.OpenFile(pt, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for name, date := range a.a {
		_, err = f.WriteString(name + "=" + date + "\n")
		if err != nil {
			return err
		}
	}

	err = syncPush(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
}

func aliasLoad(path string) (*aliases, error) {
	err := syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	a := &aliases{a: make(map[string]string)}

	pt := filepath.Join(path, tagebuchAliases)
	f, err := os.Open(pt)
	if err != nil {
		if os.IsNotExist(err) {
			return a, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		parts := strings.SplitN(text, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid alias entry: %v", text)
		}
		a.a[parts[0]] = parts[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return a, nil
}

func aliasList(path string) error {
	a, err := aliasLoad(path)
	if err != nil {
		return err
	}
	if len(a.a) == 0 {
		return nil
	}
	fmt.Println(a.String())
	return nil
}

func aliasAdd(path string, x []string) error {
	if len(x) < 2 {
		return fmt.Errorf("usage: alias add <name> <year/month/day>")
	}

	if len(x) > 2 {
		return fmt.Errorf("trailing commands: %v", x[2:])
	}

	name := strings.TrimSpace(x[0])
	date := strings.TrimSpace(x[1])

	// validate date format
	parts := strings.Split(date, "/")
	if len(parts) != 3 {
		return fmt.Errorf("invalid date format: %v (expected year/month/day)", date)
	}

	a, err := aliasLoad(path)
	if err != nil {
		return err
	}

	a.a[name] = date
	return a.save(path)
}

func aliasRemove(path string, x []string) error {
	if len(x) == 0 {
		return fmt.Errorf("must provide alias name")
	}

	if len(x) > 1 {
		return fmt.Errorf("trailing commands: %v", x[1:])
	}

	name := strings.TrimSpace(x[0])

	a, err := aliasLoad(path)
	if err != nil {
		return err
	}

	if _, ok := a.a[name]; !ok {
		return fmt.Errorf("alias not found: %v", name)
	}

	delete(a.a, name)
	return a.save(path)
}

// aliasLookup returns the date for a given alias name, or empty string if not found
func aliasLookup(path, name string) (string, error) {
	a, err := aliasLoad(path)
	if err != nil {
		return "", err
	}
	return a.a[name], nil
}
