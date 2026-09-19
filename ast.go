package main

import (
	"fmt"
)

type Expr interface {
	expr()
	String() string
}

type Stmt interface {
	stmt()
	String() string
}

type Program struct {
	Statements []Stmt
}

func (p Program) String() string {
	var s string
	for i, stmt := range p.Statements {
		if i != 0 {
			s = s + "\n"
		}
		s += stmt.String() + ";"

	}
	return s
}

type ExprStmt struct {
	Expression Expr
}

func (s ExprStmt) stmt() {}

func (s ExprStmt) String() string {
	return s.Expression.String()
}

type LetStmt struct {
	Name  string
	Value Expr
}

func (s LetStmt) stmt() {}

func (s LetStmt) String() string {
	return fmt.Sprintf("let %s = %s", s.Name, s.Value)
}

type LiteralExpr struct {
	Value Value
}

func (e LiteralExpr) expr() {}

func (e LiteralExpr) String() string {
	return e.Value.String()
}

type VariableExpr struct {
	Name Token
}

func (e VariableExpr) expr() {}

func (e VariableExpr) String() string {
	return e.Name.Value
}

type UnaryExpr struct {
	Operator Token
	Operand  Expr
}

func (e UnaryExpr) expr() {}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator.Value, e.Operand)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
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
