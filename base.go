package main

import "fmt"

var baseCommands = &Options{
	commands: []string{
		"init",
		"edit",
		"print",
		"todo",
		"search",
		"calendar",
		"list",
		"sync",
		"alias",
		"files",
	},
	descriptions: []string{
		"initialize a new tagebuch",
		"edit an entry",
		"print an entry",
		"interact with todos",
		"search within a tagebuch",
		"show calendar of entries",
		"list all days with entries",
		"sync with git remote",
		"manage named aliases to dates",
		"manage files attached to entries",
	},
}

func base(path string, x []string) error {
	if len(x) == 0 {
		return fmt.Errorf("command required. Options are:\n%v", baseCommands)
	}

	r, err := Apropos(x[0], baseCommands.commands)
	if err != nil {
		return err
	}

	switch r {
	case "init":
		return initTagebuch(path, x[1:])
	case "edit":
		return edit(path, x[1:])
	case "print":
		return printEntry(path, x[1:])
	case "todo":
		return todo(path, x[1:])
	case "search":
		return search(path, x[1:])
	case "calendar":
		return calendar(path, x[1:])
	case "list":
		return list(path, x[1:])
	case "sync":
		return sync(path, x[1:])
	case "alias":
		return alias(path, x[1:])
	case "files":
		return files(path, x[1:])
	default:
		return fmt.Errorf("invalid command %v", r)
	}
}
