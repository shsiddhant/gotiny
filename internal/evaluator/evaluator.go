package evaluator

import (
	"fmt"
	"os"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/checker"
	"github.com/shsiddhant/gotiny/internal/environment"
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

func evalAssignStmt(
	stmt *ast.AssignStmt,
	env *environment.Environment,
) (objects.Value, error) {
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

func evalLetStmt(
	stmt *ast.LetStmt,
	env *environment.Environment,
) (objects.Value, error) {
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

func evalPrintStmt(
	stmt *ast.PrintStmt,
	env *environment.Environment,
) (objects.Value, error) {
	value, err := EvalExpr(stmt.Expr, env)

	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, value)
	return nil, nil
}

func evalReturnStmt(
	stmt *ast.ReturnStmt,
	env *environment.Environment,
) (objects.Value, error) {
	value, err := EvalExpr(stmt.Expr, env)
	if err != nil {
		return nil, err
	}

	return &objects.ReturnValue{Value: value}, nil
}

func evalBlockStmt(
	stmt *ast.BlockStmt,
	env *environment.Environment,
) (objects.Value, error) {
	blockEnv := env.NewChild()

	var result objects.Value
	var err error

	for _, childStmt := range stmt.Statements {
		result, err = evalStmt(childStmt, blockEnv)
		if err != nil {
			return nil, err
		}
		if _, ok := result.(*objects.ReturnValue); ok {
			return result, nil
		}
	}
	return result, nil
}

func evalIfStmt(stmt *ast.IfStmt, env *environment.Environment) (objects.Value, error) {
	condValue, err := EvalExpr(stmt.Cond, env)
	if err != nil {
		return nil, err
	}

	boolValue := condValue.(objects.Bool)

	if boolValue {
		return evalBlockStmt(stmt.Body, env)
	} else if stmt.Else != nil {
		return evalBlockStmt(stmt.Else, env)
	}
	return nil, nil
}

func evalStmt(stmt ast.Stmt, env *environment.Environment) (objects.Value, error) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return EvalExpr(stmt.Expression, env)
	case *ast.AssignStmt:
		return evalAssignStmt(stmt, env)
	case *ast.LetStmt:
		return evalLetStmt(stmt, env)
	case *ast.IfStmt:
		return evalIfStmt(stmt, env)
	case *ast.ReturnStmt:
		return evalReturnStmt(stmt, env)
	case *ast.FnDeclareStmt:
		return evalFnDeclareStmt(stmt, env)
	case *ast.PrintStmt:
		return evalPrintStmt(stmt, env)
	default:
		return nil, fmt.Errorf("unknown statement type %T", stmt)
	}
}

func EvalProgram(
	program *ast.Program,
	env *environment.Environment,
	typeEnv *checker.TypeEnvironment,
) (objects.Value, error) {
	if err := checker.CheckProgram(program, typeEnv); err != nil {
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
