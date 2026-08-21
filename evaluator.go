package main

import "fmt"

func evalBinary(expr *BinaryExpr) (int, error) {
	left, err := Eval(expr.Left)

	if err != nil {
		return 0, err
	}

	right, err := Eval(expr.Right)
	if err != nil {
		return 0, err
	}

	switch expr.Operator.Type {
	case Plus:
		return left + right, nil
	case Minus:
		return left - right, nil
	case Star:
		return left * right, nil
	case Slash:
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("Invalid binary operator %s", expr.Operator)
	}
}

func Eval(expr Expr) (int, error) {
	switch expr := expr.(type) {
	case *LiteralExpr:
		return expr.Value, nil
	case *BinaryExpr:
		return evalBinary(expr)
	case *GroupExpr:
		return Eval(expr.Expression)
	default:
		return 0, fmt.Errorf("unknown expression type %T", expr)

	}
}
