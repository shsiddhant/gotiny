package evaluator

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

// Note: eval* functions now assume type checks have already been done.

func evalUnary(expr *ast.UnaryExpr, env *Environment) (objects.Value, error) {
	operandEval, err := EvalExpr(expr.Operand, env)

	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Plus:
		return operandEval.(objects.Int), nil
	case token.Minus:
		return -operandEval.(objects.Int), nil
	case token.Not:
		return !operandEval.(objects.Bool), nil
	default:
		return nil, &EvalError{
			Line:    expr.Operator.Line,
			Column:  expr.Operator.Column,
			Message: fmt.Sprintf("invalid unary operator %s", expr.Operator),
		}
	}

}

func evalIntBinary(left, right objects.Value, operator token.Token) (objects.Value, error) {

	leftInt, rightInt := left.(objects.Int), right.(objects.Int)
	switch operator.Type {
	case token.Plus:
		return leftInt + rightInt, nil
	case token.Minus:
		return leftInt - rightInt, nil
	case token.Star:
		return leftInt * rightInt, nil
	case token.Slash:
		if right == objects.Int(0) {
			return nil, &EvalError{
				Line:    operator.Line,
				Column:  operator.Column,
				Message: "division by zero",
			}
		}
		return leftInt / rightInt, nil
	case token.Less:
		return objects.Bool(leftInt < rightInt), nil
	case token.LessEqual:
		return objects.Bool(leftInt <= rightInt), nil
	case token.Greater:
		return objects.Bool(leftInt > rightInt), nil
	case token.GreaterEqual:
		return objects.Bool(leftInt >= rightInt), nil
	default:
		return nil, &EvalError{
			Line:    operator.Line,
			Column:  operator.Column,
			Message: fmt.Sprintf("invalid Int binary operator %q", operator.Value),
		}
	}
}

func evalEquality(left, right objects.Value, operator token.Token) (objects.Value, error) {
	switch operator.Type {
	case token.EqualEqual:
		return objects.Bool(left == right), nil
	case token.NotEqual:
		return objects.Bool(left != right), nil
	default:
		return nil, &EvalError{
			Line:    operator.Line,
			Column:  operator.Column,
			Message: fmt.Sprintf("invalid comparison operator %q", operator.Value),
		}
	}
}

func evalBinary(expr *ast.BinaryExpr, env *Environment) (objects.Value, error) {
	left, err := EvalExpr(expr.Left, env)

	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Plus, token.Minus, token.Star, token.Slash, token.Less, token.LessEqual, token.Greater, token.GreaterEqual:
		right, err := EvalExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return evalIntBinary(left, right, expr.Operator)
	case token.EqualEqual, token.NotEqual:
		right, err := EvalExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return evalEquality(left, right, expr.Operator)
	case token.And:
		if !left.(objects.Bool) {
			return objects.Bool(false), nil
		}
		right, err := EvalExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return right.(objects.Bool), nil
	case token.Or:
		if left.(objects.Bool) {
			return objects.Bool(true), nil
		}
		right, err := EvalExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return right.(objects.Bool), nil
	default:
		return nil, &EvalError{
			Line:    expr.Operator.Line,
			Column:  expr.Operator.Column,
			Message: fmt.Sprintf("invalid binary operator %s", expr.Operator),
		}

	}
}

func EvalExpr(expr ast.Expr, env *Environment) (objects.Value, error) {
	switch expr := expr.(type) {
	case *ast.LiteralExpr:
		return expr.Value, nil
	case *ast.VariableExpr:
		return env.Get(expr.Name.Value)
	case *ast.UnaryExpr:
		return evalUnary(expr, env)
	case *ast.BinaryExpr:
		return evalBinary(expr, env)
	case *ast.GroupExpr:
		return EvalExpr(expr.Expression, env)
	default:
		return nil, fmt.Errorf("unknown expression type %T", expr)

	}
}
