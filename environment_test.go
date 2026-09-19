package main

import (
	"testing"
)

func TestEnvironmentSetGet(t *testing.T) {
	env := NewEnvironment()

	value := Int(1712)

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

	old := Int(1712)

	new := Int(2412)

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

func TestEnvironmentDefine(t *testing.T) {
	env := NewEnvironment()

	name, value := "x", Int(1712)

	if err := env.Define(name, value); err != nil {
		t.Fatal(err)
	}

	got, err := env.Get(name)
	if err != nil {
		t.Fatal(err)
	}
	if got != value {
		t.Errorf("got %d, expected %d", got, value)
	}
}

func TestEnvironmentAlreadyDefined(t *testing.T) {
	env := NewEnvironment()

	name, value := "x", Int(1712)

	if err := env.Define(name, value); err != nil {
		t.Fatal(err)
	}

	if err := env.Define(name, Int(1205)); err == nil {
		t.Fatal("expected error")
	}
}

func TestEnvironmentAssign(t *testing.T) {
	env := NewEnvironment()

	name, value := "x", Bool(false)

	err := env.Define(name, value)
	if err != nil {
		t.Fatal(err)
	}

	newValue := Bool(false)

	err = env.Assign(name, newValue)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get(name)
	if err != nil {
		t.Fatal(err)
	}

	if got != newValue {
		t.Errorf("got %s, expected %s", got, newValue)
	}
}

func TestEnvironmentAssignUndefined(t *testing.T) {
	env := NewEnvironment()

	err := env.Assign("x", Int(1712))
	if err == nil {
		t.Error("expected undefined variable error")
	}
}

func TestEnvironmentAssignWrongType(t *testing.T) {
	env := NewEnvironment()

	err := env.Define("x", Int(1729))
	if err != nil {
		t.Fatal(err)
	}
	err = env.Assign("x", Bool(true))
	if err == nil {
		t.Error("expected wrong type assignment error")
	}
}
