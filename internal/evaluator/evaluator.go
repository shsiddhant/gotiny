package evaluator

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/objects"
)

func evalAssignStmt(stmt *ast.AssignStmt, env *Environment) (objects.Value, error) {
	value, err := EvalExpr(stmt.Value, env)

	if err != nil {
		return nil, err
	}

	return nil, env.Assign(stmt.Name.Value, value)
}

func evalLetStmt(stmt *ast.LetStmt, env *Environment) (objects.Value, error) {
	value, err := EvalExpr(stmt.Value, env)

	if err != nil {
		return nil, err
	}

	return nil, env.Define(stmt.Name.Value, value)
}

func EvalStmt(stmt ast.Stmt, env *Environment) (objects.Value, error) {
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

func EvalProgram(program *ast.Program, env *Environment) (objects.Value, error) {
	var result objects.Value
	var err error

	for _, stmt := range program.Statements {
		result, err = EvalStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
