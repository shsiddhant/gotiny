package main

import "fmt"

type Expr interface {
	String() string
}

type LiteralExpr struct {
	Value int
}

func (e LiteralExpr) String() string {
	return fmt.Sprintf("%d", e.Value)
}

type VariableExpr struct {
	Name Token
}

func (e VariableExpr) String() string {
	return e.Name.Value
}

type UnaryExpr struct {
	Operator Token
	Operand  Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator.Value, e.Operand)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator.Value, e.Right)
}

type GroupExpr struct {
	Expression Expr
}

func (e GroupExpr) String() string {
	return fmt.Sprintf("(group %s)", e.Expression)
}
