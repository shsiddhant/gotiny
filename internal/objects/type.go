package objects

import (
	"fmt"
	"strconv"
	"strings"
)

type Type interface {
	String() string
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=BaseType
type BaseType int

const (
	VoidType BaseType = iota
	IntType
	BoolType
)

type FunctionType struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (t FunctionType) String() string {
	params := []string{}
	for _, pt := range t.ParameterTypes {
		params = append(params, pt.String())
	}
	return fmt.Sprintf("fn(%s) %s", strings.Join(params, ", "), t.ReturnType)
}

type Value interface {
	Type() Type
	String() string
}

type Int int

func (i Int) Type() Type {
	return IntType
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

type Bool bool

func (b Bool) Type() Type {
	return BoolType
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

type ReturnValue struct {
	Value Value
}

func (rt *ReturnValue) Type() Type {
	return rt.Value.Type()
}

func (rt *ReturnValue) String() string {
	return rt.Value.String()
}
