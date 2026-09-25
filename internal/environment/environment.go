package environment

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/objects"
)

// Runtime environment
type Environment struct {
	values map[string]objects.Value
	outer  *Environment
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]objects.Value)}
}

func (e *Environment) NewChild() *Environment {
	return &Environment{values: make(map[string]objects.Value), outer: e}
}

func (e *Environment) Get(name string) (objects.Value, error) {
	value, ok := e.values[name]

	if !ok {
		if e.outer != nil {
			return e.outer.Get(name)
		}
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return value, nil
}

func (e *Environment) Set(name string, value objects.Value) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value objects.Value) error {
	current, ok := e.values[name]

	if !ok {
		if e.outer != nil {
			return e.outer.Assign(name, value)
		}
		return fmt.Errorf("undefined variable: %s", name)
	}
	if current.Type() != value.Type() {
		return fmt.Errorf(
			"cannot assign %s value to %s variable",
			value.Type(),
			current.Type(),
		)
	}
	e.values[name] = value
	return nil
}

func (e *Environment) Define(name string, value objects.Value) error {
	if _, exists := e.values[name]; exists {
		return fmt.Errorf("variable %s already defined", name)
	}
	e.values[name] = value
	return nil
}
