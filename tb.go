package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	fBase = flag.String("b", "~/.tb/", "path to tagebuch journals")

	baseDir string
)

func main() {
	flag.Parse()

	if strings.HasPrefix(*fBase, "~") {
		hd, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		baseDir = filepath.Join(hd, strings.TrimPrefix(*fBase, "~"))
	} else {
		baseDir = *fBase
	}

	args := os.Args[1:]

	err := tagebuch(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
