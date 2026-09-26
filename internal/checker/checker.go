package checker

import (
	"fmt"
	"slices"

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

type CheckContext struct {
	ReturnType objects.Type
}

func CheckProgram(program *ast.Program, env *TypeEnvironment) error {
	var err error

	for _, stmt := range program.Statements {
		err = checkStmt(stmt, env, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkStmt(stmt ast.Stmt, env *TypeEnvironment, ctx *CheckContext) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		_, err := typeOfExpr(stmt.Expression, env)
		return err
	case *ast.AssignStmt:
		return checkAssignStmt(stmt, env)
	case *ast.LetStmt:
		return checkLetStmt(stmt, env)
	case *ast.IfStmt:
		return checkIfStmt(stmt, env, ctx)
	case *ast.WhileStmt:
		return checkWhileStmt(stmt, env, ctx)
	case *ast.ReturnStmt:
		return checkReturnStmt(stmt, env, ctx)
	case *ast.FnDeclareStmt:
		return checkFnDeclareStmt(stmt, env)
	case *ast.PrintStmt:
		return checkPrintStmt(stmt, env)
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

func checkReturnStmt(
	stmt *ast.ReturnStmt,
	env *TypeEnvironment,
	ctx *CheckContext,
) error {
	if ctx == nil {
		return &CheckError{
			Token:   stmt.Expr.LocToken(),
			Message: "return not allowed outside function",
		}
	}
	returnType, err := typeOfExpr(stmt.Expr, env)
	if err != nil {
		return err
	}

	if !objects.SameType(returnType, ctx.ReturnType) {
		return &CheckError{
			Token: stmt.Expr.LocToken(),
			Message: fmt.Sprintf(
				"cannot return %s from function expecting %s",
				returnType,
				ctx.ReturnType,
			),
		}
	}

	return nil
}

func checkPrintStmt(stmt *ast.PrintStmt, env *TypeEnvironment) error {
	_, err := typeOfExpr(stmt.Expr, env)
	if err != nil {
		return err
	}
	return nil
}

func checkBlockStmt(
	stmt *ast.BlockStmt,
	env *TypeEnvironment,
	ctx *CheckContext,
) error {
	for _, childStmt := range stmt.Statements {
		if err := checkStmt(childStmt, env, ctx); err != nil {
			return err
		}
	}
	return nil
}

func checkIfStmt(stmt *ast.IfStmt, env *TypeEnvironment, ctx *CheckContext) error {
	condType, err := typeOfExpr(stmt.Cond, env)
	if err != nil {
		return err
	}
	if !objects.SameType(condType, objects.BoolType) {
		return &CheckError{
			Token:   stmt.Cond.LocToken(),
			Message: fmt.Sprintf("condition must be BoolType, got %s", condType),
		}
	}

	bodyEnv := env.NewChild()
	if err := checkBlockStmt(stmt.Body, bodyEnv, ctx); err != nil {
		return err
	}
	if stmt.Else != nil {
		elseEnv := env.NewChild()
		if err := checkBlockStmt(stmt.Else, elseEnv, ctx); err != nil {
			return err
		}
	}
	return nil
}

func checkWhileStmt(
	stmt *ast.WhileStmt,
	env *TypeEnvironment,
	ctx *CheckContext,
) error {
	condType, err := typeOfExpr(stmt.Cond, env)
	if err != nil {
		return err
	}
	if !objects.SameType(condType, objects.BoolType) {
		return &CheckError{
			Token:   stmt.Cond.LocToken(),
			Message: fmt.Sprintf("condition must be BoolType, got %s", condType),
		}
	}

	bodyEnv := env.NewChild()
	if err := checkBlockStmt(stmt.Body, bodyEnv, ctx); err != nil {
		return err
	}
	return nil
}

func checkFnDeclareStmt(stmt *ast.FnDeclareStmt, env *TypeEnvironment) error {
	var paramTypes []objects.Type

	if _, err := env.Get(stmt.Name.Value); err == nil {
		return &CheckError{
			Token:   stmt.Name,
			Message: fmt.Sprintf("name already defined: %s", stmt.Name.Value),
		}
	}

	for _, param := range stmt.Parameters {
		paramTypes = append(paramTypes, param.Type)
	}

	fnType := &objects.FunctionType{
		ParameterTypes: paramTypes,
		ReturnType:     stmt.ReturnType,
	}

	fnTypeEnv := env.NewChild()

	if err := fnTypeEnv.Define(stmt.Name.Value, fnType); err != nil {
		return err
	}

	for _, param := range stmt.Parameters {
		if err := fnTypeEnv.Define(param.Name.Value, param.Type); err != nil {
			return &CheckError{
				Token:   param.Name,
				Message: fmt.Sprintf("duplicate parameter name: %s", param.Name.Value),
			}
		}
	}

	ctx := &CheckContext{stmt.ReturnType}

	if err := checkBlockStmt(stmt.Body, fnTypeEnv, ctx); err != nil {
		return err
	}

	if !objects.SameType(stmt.ReturnType, objects.VoidType) &&
		!alwaysReturns(stmt.Body) {
		return &CheckError{
			Token: stmt.Name,
			Message: fmt.Sprintf(
				"missing return statement at end of function %q",
				stmt.Name.Value,
			),
		}
	}

	if err := env.Define(stmt.Name.Value, fnType); err != nil {
		return &CheckError{
			Token:   stmt.Name,
			Message: err.Error(),
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
	case *ast.CallExpr:
		return typeOfCallExpr(expr, env)
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
		if !objects.SameType(operandType, objects.IntType) {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"unary operator %q cannot be applied to %s",
					expr.Operator.Value,
					operandType,
				),
			}
		}
		return operandType, nil
	case token.Not:
		if !objects.SameType(operandType, objects.BoolType) {
			return nil, &CheckError{
				Token: expr.LocToken(),
				Message: fmt.Sprintf(
					"unary operator %q cannot be applied to %s",
					expr.Operator.Value,
					operandType,
				),
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
		if !objects.SameType(leftType, objects.IntType) ||
			!objects.SameType(rightType, objects.IntType) {
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
		if !objects.SameType(leftType, objects.IntType) ||
			!objects.SameType(rightType, objects.IntType) {
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
		if !objects.SameType(leftType, rightType) {
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
		if !objects.SameType(leftType, objects.BoolType) ||
			!objects.SameType(rightType, objects.BoolType) {
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

func typeOfCallExpr(expr *ast.CallExpr, env *TypeEnvironment) (objects.Type, error) {
	exprType, err := env.Get(expr.Name.Value)
	if err != nil {
		return nil, &CheckError{
			Token:   expr.LocToken(),
			Message: err.Error(),
		}
	}

	fnType, ok := exprType.(*objects.FunctionType)
	if !ok {
		return nil, &CheckError{
			Token:   expr.LocToken(),
			Message: fmt.Sprintf("expected function type, got %s", exprType),
		}
	}

	if len(expr.Args) != len(fnType.ParameterTypes) {
		return nil, &CheckError{
			Token: expr.LocToken(),
			Message: fmt.Sprintf(
				"expected %d args, got %d",
				len(fnType.ParameterTypes),
				len(expr.Args),
			),
		}
	}

	for i, arg := range expr.Args {
		argType, err := typeOfExpr(arg, env)
		if err != nil {
			return nil, err
		}
		if !objects.SameType(argType, fnType.ParameterTypes[i]) {
			return nil, &CheckError{
				Token: arg.LocToken(),
				Message: fmt.Sprintf(
					"expected %s arg, got %s",
					fnType.ParameterTypes[i],
					argType,
				),
			}
		}
	}

	return fnType.ReturnType, nil
}

func alwaysReturns(stmt ast.Stmt) bool {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BlockStmt:
		return slices.ContainsFunc(s.Statements, alwaysReturns)
	case *ast.IfStmt:
		return s.Else != nil && alwaysReturns(s.Body) && alwaysReturns(s.Else)
	default:
		return false
	}
}
