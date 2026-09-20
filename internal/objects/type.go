package objects

import (
	"strconv"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Type
type Type int

const (
	IntType Type = iota
	BoolType
)

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
