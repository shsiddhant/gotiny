package ast

import (
	"fmt"
	"strings"

	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

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
	Name  token.Token
	Value Expr
}

func (s AssignStmt) stmt() {}

func (s AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", s.Name.Value, s.Value)
}

// Let Statement for variable declaration.
//
// Example:
//
//	`let x = 1712;`
type LetStmt struct {
	Name  token.Token
	Value Expr
}

func (s LetStmt) stmt() {}

func (s LetStmt) String() string {
	return fmt.Sprintf("let %s = %s", s.Name.Value, s.Value)
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
	if len(s.Statements) == 0 {
		return "{}"
	}
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

// WhileStmt represents a while loop.
type WhileStmt struct {
	Cond Expr
	Body *BlockStmt
}

func (s WhileStmt) stmt() {}

func (s WhileStmt) String() string {
	return fmt.Sprintf("while %s %s ", s.Cond, s.Body)
}

type Parameter struct {
	Name token.Token
	Type objects.Type
}

func (param Parameter) String() string {
	return fmt.Sprintf("%s %s", param.Name.Value, objects.TypeString(param.Type))
}

func GetParameterTypes(params []Parameter) []objects.Type {
	var paramTypes []objects.Type

	for _, param := range params {
		paramTypes = append(paramTypes, param.Type)
	}
	return paramTypes
}

// Function Declaration Statement
type FnDeclareStmt struct {
	Name       token.Token
	Parameters []Parameter
	ReturnType objects.Type
	Body       *BlockStmt
}

func (stmt FnDeclareStmt) stmt() {}

func (stmt FnDeclareStmt) String() string {
	var params []string
	for _, param := range stmt.Parameters {
		params = append(params, param.String())
	}
	return fmt.Sprintf(
		"fn %s(%s) %s %s",
		stmt.Name.Value,
		strings.Join(params, ", "),
		objects.TypeString(stmt.ReturnType),
		stmt.Body,
	)
}

// Return statement
type ReturnStmt struct {
	Expr Expr
}

func (stmt ReturnStmt) stmt() {}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %s", stmt.Expr)
}

// Print statement
type PrintStmt struct {
	Expr Expr
}

func (stmt PrintStmt) stmt() {}

func (stmt PrintStmt) String() string {
	return fmt.Sprintf("print %s", stmt.Expr)
}
