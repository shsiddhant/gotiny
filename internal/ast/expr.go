package ast

import (
	"fmt"
	"strings"

	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

type Expr interface {
	expr()
	String() string
	LocToken() token.Token // token for location
}

type LiteralExpr struct {
	Token token.Token
	Value objects.Value
}

func (e LiteralExpr) expr() {}

func (e LiteralExpr) String() string {
	return e.Value.String()
}

func (e LiteralExpr) LocToken() token.Token {
	return e.Token
}

type VariableExpr struct {
	Name token.Token
}

func (e VariableExpr) expr() {}

func (e VariableExpr) String() string {
	return e.Name.Value
}
func (e VariableExpr) LocToken() token.Token {
	return e.Name
}

func (e VariableExpr) Column() int {
	return e.Name.Column
}

type UnaryExpr struct {
	Operator token.Token
	Operand  Expr
}

func (e UnaryExpr) expr() {}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator.Value, e.Operand)
}

func (e UnaryExpr) LocToken() token.Token {
	return e.Operator
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e BinaryExpr) expr() {}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator.Value, e.Right)
}

func (e BinaryExpr) LocToken() token.Token {
	return e.Operator
}

type GroupExpr struct {
	Expression Expr
}

func (e GroupExpr) expr() {}

func (e GroupExpr) String() string {
	return fmt.Sprintf("(group %s)", e.Expression)
}

func (e GroupExpr) LocToken() token.Token {
	return e.Expression.LocToken()
}

// Call Expression
type CallExpr struct {
	Name token.Token
	Args []Expr
}

func (e CallExpr) expr() {}

func (e CallExpr) String() string {
	var args []string
	for _, arg := range e.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf(
		"%s(%s)",
		e.Name.Value,
		strings.Join(args, ", "),
	)
}

func (e CallExpr) LocToken() token.Token {
	return e.Name
}
