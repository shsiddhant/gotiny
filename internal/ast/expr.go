package ast

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

type Expr interface {
	expr()
	String() string
}

type LiteralExpr struct {
	Value objects.Value
}

func (e LiteralExpr) expr() {}

func (e LiteralExpr) String() string {
	return e.Value.String()
}

type VariableExpr struct {
	Name token.Token
}

func (e VariableExpr) expr() {}

func (e VariableExpr) String() string {
	return e.Name.Value
}

type UnaryExpr struct {
	Operator token.Token
	Operand  Expr
}

func (e UnaryExpr) expr() {}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator.Value, e.Operand)
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

type GroupExpr struct {
	Expression Expr
}

func (e GroupExpr) expr() {}

func (e GroupExpr) String() string {
	return fmt.Sprintf("(group %s)", e.Expression)
}
