package evaluator

import (
	"fmt"
	"strings"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/environment"
	"github.com/shsiddhant/gotiny/internal/objects"
)

type Function struct {
	Declaration *ast.FnDeclareStmt
	Env         *environment.Environment
}

func (f *Function) Type() objects.Type {
	return &objects.FunctionType{
		ParameterTypes: ast.GetParameterTypes(f.Declaration.Parameters),
		ReturnType:     f.Declaration.ReturnType,
	}
}

func (f *Function) String() string {
	var params []string
	for _, param := range f.Declaration.Parameters {
		params = append(params, param.String())
	}
	return fmt.Sprintf(
		"%s(%s) %s",
		f.Declaration.Name.Value,
		strings.Join(params, ", "),
		objects.TypeString(f.Declaration.ReturnType),
	)
}

func evalFnDeclareStmt(
	stmt *ast.FnDeclareStmt,
	env *environment.Environment,
) (objects.Value, error) {
	value := &Function{
		Declaration: stmt,
		Env:         env,
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

func evalCallExpr(
	expr *ast.CallExpr,
	env *environment.Environment,
) (objects.Value, error) {
	function, err := env.Get(expr.Name.Value)
	if err != nil {
		return nil, err
	}
	fnValue := function.(*Function)
	args := make([]objects.Value, len(expr.Args))

	for i, arg := range expr.Args {
		value, err := EvalExpr(arg, env)
		if err != nil {
			return nil, err
		}
		args[i] = value
	}

	callEnv := fnValue.Env.NewChild()

	for i, param := range fnValue.Declaration.Parameters {
		callEnv.Set(param.Name.Value, args[i])
	}

	result, err := evalBlockStmt(fnValue.Declaration.Body, callEnv)
	if err != nil {
		return nil, err
	}

	if returnValue, ok := result.(*objects.ReturnValue); ok {
		return returnValue.Value, nil
	}

	return result, nil
}
