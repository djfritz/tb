package main

import (
	"fmt"
	"log"
	"strings"
)

type Options struct {
	commands     []string
	descriptions []string
}

func (o *Options) String() string {
	if len(o.commands) != len(o.descriptions) {
		log.Fatal("mismatched options lengths")
	}

	// Golang's tabwriter is still there, but it's deprecated, so we'll
	// just do it from scratch knowing exactly what/how we want to print.
	maxLength := 0
	for _, v := range o.commands {
		if len(v) > maxLength {
			maxLength = len(v)
		}
	}

	var output string
	for i, v := range o.commands {
		output += v
		for x := len(v); x < maxLength; x++ {
			output += " "
		}
		output = fmt.Sprintf("%v : %v\n", output, o.descriptions[i])
	}
	return strings.TrimSpace(output)
}
