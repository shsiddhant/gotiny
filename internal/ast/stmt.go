package ast

import "fmt"

// Statement interface
type Stmt interface {
	stmt()
	String() string
}

// Expression Statement
//
// Example:
//
//	`x + y * 1712;`
type ExprStmt struct {
	Expression Expr
}

func (s ExprStmt) stmt() {}

func (s ExprStmt) String() string {
	return s.Expression.String()
}

// Assignment statement
//
// Example:
//
//	`z = x + 1729 * y;`
type AssignStmt struct {
	Name  string
	Value Expr
}

func (s AssignStmt) stmt() {}

func (s AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", s.Name, s.Value)
}

// Let Statement for variable declaration.
//
// Example:
//
//	`let x = 1712;`
type LetStmt struct {
	Name  string
	Value Expr
}

func (s LetStmt) stmt() {}

func (s LetStmt) String() string {
	return fmt.Sprintf("let %s = %s", s.Name, s.Value)
}

// BlockStmt represents a block of statement inside curly braces.
//
// Example:
//
//	`{
//		let x = 1712;
//		let y = 1729;
//		x - y;
//	 }`
type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) stmt() {}

func (s BlockStmt) String() string {
	str := "{\n"
	for _, stmt := range s.Statements {
		if stmt != nil {
			str = str + "  " + stmt.String() + ";\n"
		}
	}
	return str + "}"
}

// IfStmt represents and if..else statement, with a condition,
// body block, and an optional else block
//
// Example:
//
//	`if x {
//		y = x;
//	 } else {
//		y = -x;
//	 }`
type IfStmt struct {
	Cond Expr
	Body *BlockStmt
	Else *BlockStmt
}

func (s IfStmt) stmt() {}

func (s IfStmt) String() string {
	str := "if " + s.Cond.String() + " " + s.Body.String()
	if s.Else != nil {
		str += " else" + s.Else.String()
	}
	return str
}
