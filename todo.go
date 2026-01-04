package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var todoCommands = &Options{
	commands: []string{
		"add",
		"complete",
	},
	descriptions: []string{
		"add a todo item",
		"complete a todo item by number",
	},
}

func todo(path string, x []string) error {
	err := validate(path)
	if err != nil {
		return err
	}

	if len(x) == 0 {
		// print todos
		return todoPrint(path)
	}

	r, err := Apropos(x[0], todoCommands.commands)
	if err != nil {
		return err
	}

	switch r {
	case "add":
		return todoAdd(path, x[1:])
	case "complete":
		return todoComplete(path, x[1:])
	default:
		return fmt.Errorf("invalid command %v", r)
	}
}

type todos struct {
	t []string
}

func (t *todos) String() string {
	var ret string
	for i, v := range t.t {
		ret += fmt.Sprintf("%v: %v\n", i, v)
	}
	return strings.TrimSpace(ret)
}

func (t *todos) save(path string) error {
	pt := filepath.Join(path, tagebuchTodo)
	f, err := os.OpenFile(pt, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range t.t {
		_, err = f.WriteString(v + "\n")
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

func todoLoad(path string) (*todos, error) {
	err := syncPull(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// base path is already validated
	pt := filepath.Join(path, tagebuchTodo)
	f, err := os.Open(pt)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &todos{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			t.t = append(t.t, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return t, nil
}

func todoPrint(path string) error {
	t, err := todoLoad(path)
	if err != nil {
		return err
	}
	fmt.Println(t.String())
	return nil
}

func todoAdd(path string, x []string) error {
	if len(x) == 0 {
		return fmt.Errorf("must provide todo text")
	}

	if len(x) != 1 {
		return fmt.Errorf("trailing commands: %v", x[2:])
	}

	text := strings.TrimSpace(x[0])

	t, err := todoLoad(path)
	if err != nil {
		return err
	}

	// deduplicate
	for _, v := range t.t {
		if v == text {
			return nil
		}
	}

	t.t = append(t.t, text)
	return t.save(path)
}

func todoComplete(path string, x []string) error {
	if len(x) == 0 {
		return fmt.Errorf("must provide todo item number")
	}

	if len(x) != 1 {
		return fmt.Errorf("trailing commands: %v", x[2:])
	}

	item := strings.TrimSpace(x[0])
	no, err := strconv.Atoi(item)
	if err != nil {
		return err
	}

	t, err := todoLoad(path)
	if err != nil {
		return err
	}

	if no >= len(t.t) {
		return fmt.Errorf("invalid index %v", no)
	}

	t.t = slices.Delete(t.t, no, no+1)
	return t.save(path)
}
