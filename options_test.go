package main

import "testing"

func TestOptions(t *testing.T) {
	o := &Options{
		commands:     []string{"foo", "foobar"},
		descriptions: []string{"it's foo", "it's foobar"},
	}

	expected := `foo    : it's foo
foobar : it's foobar`
	if o.String() != expected {
		t.Fatal("invalid output:", o.String())
	}
}
