package evaluator

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/objects"
)

type EvalError struct {
	Line    int
	Column  int
	Message string
}

func (e *EvalError) Error() string {
	return fmt.Sprintf("evaluation error at line %d, column %d: %s",
		e.Line,
		e.Column,
		e.Message,
	)
}

func evalAssignStmt(stmt *ast.AssignStmt, env *Environment) (objects.Value, error) {
	value, err := EvalExpr(stmt.Value, env)

	if err != nil {
		return nil, err
	}

	if err := env.Assign(stmt.Name.Value, value); err != nil {
		return nil, &EvalError{
			Line:    stmt.Name.Line,
			Column:  stmt.Name.Column,
			Message: err.Error(),
		}
	}
	return nil, nil
}

func evalLetStmt(stmt *ast.LetStmt, env *Environment) (objects.Value, error) {
	value, err := EvalExpr(stmt.Value, env)

	if err != nil {
		return nil, err
	}

	if err := env.Define(stmt.Name.Value, value); err != nil {
		return nil, &EvalError{
			Line:    stmt.Name.Line,
			Column:  stmt.Name.Column,
			Message: err.Error(),
		}
	}

	return nil, nil
}

func evalStmt(stmt ast.Stmt, env *Environment) (objects.Value, error) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return EvalExpr(stmt.Expression, env)
	case *ast.AssignStmt:
		return evalAssignStmt(stmt, env)
	case *ast.LetStmt:
		return evalLetStmt(stmt, env)
	default:
		return nil, fmt.Errorf("unknown statement type %T", stmt)
	}
}

func EvalProgram(program *ast.Program, env *Environment, typeEnv *TypeEnvironment) (objects.Value, error) {
	if err := CheckProgram(program, typeEnv); err != nil {
		return nil, err
	}
	var result objects.Value
	var err error

	for _, stmt := range program.Statements {
		result, err = evalStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
