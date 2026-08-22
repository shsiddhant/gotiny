package main

import "fmt"

type Environment struct {
	values map[string]int
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]int)}
}

func (e *Environment) Get(name string) (int, error) {
	value, ok := e.values[name]

	if !ok {
		return 0, fmt.Errorf("undefined variable: %s", name)
	}
	return value, nil
}

func (e *Environment) Set(name string, value int) {
	e.values[name] = value
}
