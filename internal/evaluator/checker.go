package evaluator

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

type CheckError struct {
	Token   token.Token
	Message string
}

func (e *CheckError) Error() string {
	return fmt.Sprintf("check error at line %d, column %d: %s",
		e.Token.Line,
		e.Token.Column,
		e.Message,
	)
}

func CheckProgram(program *ast.Program, env *TypeEnvironment) error {
	var err error

	for _, stmt := range program.Statements {
		err = checkStmt(stmt, env)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkStmt(stmt ast.Stmt, env *TypeEnvironment) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		_, err := typeOfExpr(stmt.Expression, env)
		return err
	case *ast.AssignStmt:
		return checkAssignStmt(stmt, env)
	case *ast.LetStmt:
		return checkLetStmt(stmt, env)
	case *ast.IfStmt:
		return checkIfStmt(stmt, env)
	default:
		return fmt.Errorf("unknown statement type %T", stmt)
	}
}

func checkAssignStmt(stmt *ast.AssignStmt, env *TypeEnvironment) error {
	_, err := env.Get(stmt.Name.Value)
	if err != nil {
		return &CheckError{
			Token:   stmt.Name,
			Message: err.Error(),
		}
	}

	typ, err := typeOfExpr(stmt.Value, env)
	if err != nil {
		return err
	}
	if err := env.Assign(stmt.Name.Value, typ); err != nil {
		return &CheckError{
			Token:   stmt.Name,
			Message: err.Error(),
		}
	}
	return nil

}

func checkLetStmt(stmt *ast.LetStmt, env *TypeEnvironment) error {
	typ, err := typeOfExpr(stmt.Value, env)
	if err != nil {
		return err
	}
	if err := env.Define(stmt.Name.Value, typ); err != nil {
		return &CheckError{
			Token:   stmt.Name,
			Message: err.Error(),
		}
	}
	return nil
}

func checkBlockStmt(stmt *ast.BlockStmt, env *TypeEnvironment) error {
	blockEnv := env.NewChild()

	for _, childStmt := range stmt.Statements {
		if err := checkStmt(childStmt, blockEnv); err != nil {
			return err
		}
	}
	return nil
}

func checkIfStmt(stmt *ast.IfStmt, env *TypeEnvironment) error {
	condType, err := typeOfExpr(stmt.Cond, env)
	if err != nil {
		return err
	}
	if condType != objects.BoolType {
		return &CheckError{
			Token:   stmt.Cond.LocToken(),
			Message: fmt.Sprintf("condition must be BoolType, got %s", condType),
		}
	}

	if err := checkBlockStmt(stmt.Body, env); err != nil {
		return err
	}
	if stmt.Else != nil {
		if err := checkBlockStmt(stmt.Else, env); err != nil {
			return err
		}
	}
	return nil
}

func typeOfExpr(expr ast.Expr, env *TypeEnvironment) (objects.Type, error) {
	switch expr := expr.(type) {
	case *ast.LiteralExpr:
		return expr.Value.Type(), nil
	case *ast.VariableExpr:
		typ, err := env.Get(expr.Name.Value)
		if err != nil {
			return nil, &CheckError{
				Token:   expr.LocToken(),
				Message: err.Error(),
			}
		}
		return typ, nil
	case *ast.UnaryExpr:
		return typeOfUnary(expr, env)
	case *ast.BinaryExpr:
		return typeOfBinary(expr, env)
	case *ast.GroupExpr:
		return typeOfExpr(expr.Expression, env)
	default:
		return nil, fmt.Errorf("unknown expression type %T", expr)
	}
}

func typeOfUnary(expr *ast.UnaryExpr, env *TypeEnvironment) (objects.Type, error) {
	operandType, err := typeOfExpr(expr.Operand, env)

	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Plus, token.Minus:
		if operandType != objects.IntType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"unary operator %q cannot be applied to %s", expr.Operator.Value, operandType),
			}
		}
		return operandType, nil
	case token.Not:
		if operandType != objects.BoolType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"unary operator %q cannot be applied to %s", expr.Operator.Value, operandType),
			}
		}
		return objects.BoolType, nil
	default:
		return nil, &CheckError{
			Token:   expr.LocToken(),
			Message: fmt.Sprintf("invalid unary operator %s", expr.Operator),
		}
	}

}

func typeOfBinary(expr *ast.BinaryExpr, env *TypeEnvironment) (objects.Type, error) {
	leftType, err := typeOfExpr(expr.Left, env)
	if err != nil {
		return nil, err
	}

	rightType, err := typeOfExpr(expr.Right, env)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Plus, token.Minus, token.Star, token.Slash:
		if leftType != objects.IntType || rightType != objects.IntType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"binary operator %q cannot be applied to %s and %s",
					expr.Operator.Value,
					leftType,
					rightType,
				),
			}
		}
		return objects.IntType, nil
	case token.Less, token.LessEqual, token.Greater, token.GreaterEqual:
		if leftType != objects.IntType || rightType != objects.IntType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"binary operator %q cannot be applied to %s and %s",
					expr.Operator.Value,
					leftType,
					rightType,
				),
			}
		}
		return objects.BoolType, nil
	case token.EqualEqual, token.NotEqual:
		if leftType != rightType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"cannot compare %s and %s",
					leftType,
					rightType,
				),
			}
		}
		return objects.BoolType, nil
	case token.And, token.Or:
		if leftType != objects.BoolType || rightType != objects.BoolType {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"binary operator %q cannot be applied to %s and %s",
					expr.Operator.Value,
					leftType,
					rightType,
				),
			}
		}
		return objects.BoolType, nil
	default:
		return nil, &CheckError{
			Token:   expr.LocToken(),
			Message: fmt.Sprintf("invalid binary operator %s", expr.Operator),
		}
	}
}
