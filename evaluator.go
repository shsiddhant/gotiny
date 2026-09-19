package main

import "fmt"

func evalUnary(expr *UnaryExpr, env *Environment) (Value, error) {
	operandEval, err := EvalExpr(expr.Operand, env)

	if err != nil {
		return nil, err
	}

	value, ok := operandEval.(Int)
	if !ok {
		return nil, fmt.Errorf(
			"unary operator %q cannot be applied to %s", expr.Operator.Value, operandEval.Type())
	}

	switch expr.Operator.Type {
	case Plus:
		return value, nil
	case Minus:
		return -value, nil
	default:
		return nil, fmt.Errorf("Invalid unary operator %s", expr.Operator)
	}

}

func evalBinary(expr *BinaryExpr, env *Environment) (Value, error) {
	left, err := EvalExpr(expr.Left, env)

	if err != nil {
		return nil, err
	}

	right, err := EvalExpr(expr.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != IntType || right.Type() != IntType {
		return nil, fmt.Errorf(
			"binary operator %q cannot be applied to %s and %s",
			expr.Operator.Value,
			left.Type(),
			right.Type(),
		)
	}

	leftInt, rightInt := left.(Int), right.(Int)

	switch expr.Operator.Type {
	case Plus:
		return leftInt + rightInt, nil
	case Minus:
		return leftInt - rightInt, nil
	case Star:
		return leftInt * rightInt, nil
	case Slash:
		if right == Int(0) {
			return nil, fmt.Errorf("division by zero")
		}
		return leftInt / rightInt, nil
	default:
		return nil, fmt.Errorf("Invalid binary operator %s", expr.Operator)
	}
}

func EvalExpr(expr Expr, env *Environment) (Value, error) {
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
		return nil, fmt.Errorf("unknown expression type %T", expr)

	}
}

func evalLetStmt(stmt *LetStmt, env *Environment) (Value, error) {
	value, err := EvalExpr(stmt.Value, env)

	if err != nil {
		return nil, err
	}

	return nil, env.Define(stmt.Name, value)
}

func EvalStmt(stmt Stmt, env *Environment) (Value, error) {
	switch stmt := stmt.(type) {
	case *ExprStmt:
		return EvalExpr(stmt.Expression, env)
	case *LetStmt:
		return evalLetStmt(stmt, env)
	default:
		return nil, fmt.Errorf("unknown statement type %T", stmt)
	}
}

func EvalProgram(program *Program, env *Environment) (Value, error) {
	var result Value
	var err error

	for _, stmt := range program.Statements {
		result, err = EvalStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
