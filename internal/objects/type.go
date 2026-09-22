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
	IntType BaseType = iota
	BoolType
)

type FunctionType struct {
	ParameterTypes []Type
}

func (t FunctionType) String() string {
	params := []string{}
	for _, pt := range t.ParameterTypes {
		params = append(params, pt.String())
	}
	return fmt.Sprintf("fn(%s)", strings.Join(params, ", "))
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
