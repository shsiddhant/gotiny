package evaluator

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

func evalUnary(expr *ast.UnaryExpr, env *Environment) (objects.Value, error) {
	operandEval, err := EvalExpr(expr.Operand, env)

	if err != nil {
		return nil, err
	}

	value, ok := operandEval.(objects.Int)
	if !ok {
		return nil, &EvalError{
			Line:   expr.Operator.Line,
			Column: expr.Operator.Column,
			Message: fmt.Sprintf(
				"unary operator %q cannot be applied to %s", expr.Operator.Value, operandEval.Type()),
		}
	}

	switch expr.Operator.Type {
	case token.Plus:
		return value, nil
	case token.Minus:
		return -value, nil
	default:
		return nil, &EvalError{
			Line:    expr.Operator.Line,
			Column:  expr.Operator.Column,
			Message: fmt.Sprintf("Invalid unary operator %s", expr.Operator),
		}
	}

}

func evalBinary(expr *ast.BinaryExpr, env *Environment) (objects.Value, error) {
	left, err := EvalExpr(expr.Left, env)

	if err != nil {
		return nil, err
	}

	right, err := EvalExpr(expr.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != objects.IntType || right.Type() != objects.IntType {
		return nil, &EvalError{
			Line:   expr.Operator.Line,
			Column: expr.Operator.Column,
			Message: fmt.Sprintf(
				"binary operator %q cannot be applied to %s and %s",
				expr.Operator.Value,
				left.Type(),
				right.Type(),
			),
		}
	}

	leftInt, rightInt := left.(objects.Int), right.(objects.Int)

	switch expr.Operator.Type {
	case token.Plus:
		return leftInt + rightInt, nil
	case token.Minus:
		return leftInt - rightInt, nil
	case token.Star:
		return leftInt * rightInt, nil
	case token.Slash:
		if right == objects.Int(0) {
			return nil, &EvalError{
				Line:    expr.Operator.Line,
				Column:  expr.Operator.Column,
				Message: "division by zero",
			}
		}
		return leftInt / rightInt, nil
	default:
		return nil, &EvalError{
			Line:    expr.Operator.Line,
			Column:  expr.Operator.Column,
			Message: fmt.Sprintf("Invalid binary operator %s", expr.Operator),
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
