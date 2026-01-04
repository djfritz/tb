package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoMatches       = errors.New("no matching command")
	ErrMultipleMatches = errors.New("multiple matching commands")
)

func Apropos(input string, options []string) (string, error) {
	var results []string
	for _, v := range options {
		if strings.HasPrefix(v, input) {
			results = append(results, v)
		}
	}

	switch len(results) {
	case 0:
		return "", fmt.Errorf("%w: options are %v", ErrNoMatches, options)
	case 1:
		return results[0], nil
	default:
		return "", fmt.Errorf("%w: %v", ErrMultipleMatches, results)
	}
}
