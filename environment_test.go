package main

import (
	"testing"
)

func TestEnvironmentSetGet(t *testing.T) {
	env := NewEnvironment()

	value := 1712

	env.Set("x", value)

	got, err := env.Get("x")
	if err != nil {
		t.Fatal(err)
	}

	if got != value {
		t.Errorf("got %d, expected %d", got, value)
	}
}

func TestEnvironmentSetExisting(t *testing.T) {
	env := NewEnvironment()

	old := 1712

	new := 2412

	env.Set("x", old)
	env.Set("x", new)

	got, err := env.Get("x")
	if err != nil {
		t.Fatal(err)
	}

	if got != new {
		t.Errorf("got %d, expected %d", got, new)
	}
}
