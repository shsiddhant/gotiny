package main

import "fmt"

func evalUnary(expr *UnaryExpr, env *Environment) (int, error) {
	operandEval, err := EvalExpr(expr.Operand, env)

	if err != nil {
		return 0, err
	}

	switch expr.Operator.Type {
	case Plus:
		return operandEval, nil
	case Minus:
		return -operandEval, nil
	default:
		return 0, fmt.Errorf("Invalid unary operator %s", expr.Operator)
	}

}

func evalBinary(expr *BinaryExpr, env *Environment) (int, error) {
	left, err := EvalExpr(expr.Left, env)

	if err != nil {
		return 0, err
	}

	right, err := EvalExpr(expr.Right, env)
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

func EvalExpr(expr Expr, env *Environment) (int, error) {
	switch expr := expr.(type) {
	case *LiteralExpr:
		return expr.Value, nil
	case *VariableExpr:
		return env.Get(expr.Name.Value)
	case *UnaryExpr:
		return evalUnary(expr, env)
	case *BinaryExpr:
		return evalBinary(expr, env)
	case *GroupExpr:
		return EvalExpr(expr.Expression, env)
	default:
		return 0, fmt.Errorf("unknown expression type %T", expr)

	}
}

func EvalStmt(stmt Stmt, env *Environment) (int, error) {
	switch stmt := stmt.(type) {
	case *ExprStmt:
		return EvalExpr(stmt.Expression, env)

	default:
		return 0, fmt.Errorf("unknown statement type %T", stmt)
	}
}
