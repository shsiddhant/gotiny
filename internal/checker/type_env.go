package checker

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/objects"
)

// Type environment
type TypeEnvironment struct {
	types map[string]objects.Type
	outer *TypeEnvironment
}

func NewTypeEnvironment() *TypeEnvironment {
	return &TypeEnvironment{types: make(map[string]objects.Type)}
}

func (e *TypeEnvironment) NewChild() *TypeEnvironment {
	return &TypeEnvironment{types: make(map[string]objects.Type), outer: e}
}

func (e *TypeEnvironment) Get(name string) (objects.Type, error) {
	typ, ok := e.types[name]

	if !ok {
		if e.outer != nil {
			return e.outer.Get(name)
		}
		return nil, fmt.Errorf("undefined name: %s", name)
	}
	return typ, nil
}

func (e *TypeEnvironment) Set(name string, typ objects.Type) {
	e.types[name] = typ
}

func (e *TypeEnvironment) Assign(name string, typ objects.Type) error {
	current, ok := e.types[name]

	if !ok {
		if e.outer != nil {
			return e.outer.Assign(name, typ)
		}
		return fmt.Errorf("undefined name: %s", name)
	}
	if current != typ {
		return fmt.Errorf("cannot assign %s value to %s variable", typ, current)
	}
	e.types[name] = typ
	return nil
}

func (e *TypeEnvironment) Define(name string, typ objects.Type) error {
	if _, exists := e.types[name]; exists {
		return fmt.Errorf("name already defined: %s", name)
	}
	e.types[name] = typ
	return nil
}
