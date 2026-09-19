package main

import (
	"fmt"
)

type Environment struct {
	values map[string]Value
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]Value)}
}

func (e *Environment) Get(name string) (Value, error) {
	value, ok := e.values[name]

	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return value, nil
}

func (e *Environment) Set(name string, value Value) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value Value) error {
	current, ok := e.values[name]

	if !ok {
		return fmt.Errorf("undefined variable: %s", name)
	}
	if current.Type() != value.Type() {
		return fmt.Errorf("cannot assign %s value to %s variable", value.Type(), current.Type())
	}
	e.values[name] = value
	return nil
}

func (e *Environment) Define(name string, value Value) error {
	if _, exists := e.values[name]; exists {
		return fmt.Errorf("variable %s already defined", name)
	}
	e.values[name] = value
	return nil
}
