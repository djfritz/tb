package main

import (
	"errors"
	"testing"
)

func TestApropos(t *testing.T) {
	input := "f"
	options := []string{"foo", "bar"}

	x, err := Apropos(input, options)
	if err != nil {
		t.Fatal(err)
	}
	if x != "foo" {
		t.Fatal("invalid response")
	}
}

func TestAproposNone(t *testing.T) {
	input := "fool"
	options := []string{"foo", "bar"}

	_, err := Apropos(input, options)
	if err.Error() != "no matching command: options are [foo bar]" {
		t.Fatal("invalid or missing error", err)
	}
}

func TestAproposMulti(t *testing.T) {
	input := "foo"
	options := []string{"foo", "foobar", "bar"}

	_, err := Apropos(input, options)
	if !errors.Is(err, ErrMultipleMatches) {
		t.Fatal("invalid or missing error")
	}
	if err.Error() != "multiple matching commands: [foo foobar]" {
		t.Fatal("invalid error", err)
	}
}
