package evaluator

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected objects.Value
	}{
		{"1;", objects.Int(1)},
		{"-1;", objects.Int(-1)},
		{"+1;", objects.Int(1)},
		{"--1;", objects.Int(1)},
		{"-+1;", objects.Int(-1)},
		{"2 + -5;", objects.Int(-3)},
		{"2 * -3;", objects.Int(-6)},
		{"-2 * 3;", objects.Int(-6)},
		{"-5 / 3;", objects.Int(-1)},
		{"-(1 + 2);", objects.Int(-3)},
		{"1 + 2;", objects.Int(3)},
		{"2 * 3;", objects.Int(6)},
		{"1 + 2 * 3;", objects.Int(7)},
		{"(1 + 2) * 3;", objects.Int(9)},
		{"20 / 5 / 2;", objects.Int(2)},
		{"10 - 3 - 2;", objects.Int(5)},
		{"1 > 2;", objects.Bool(false)},
		{"1 > 2 == true;", objects.Bool(false)},
		{"let x = 1712; let y = 1729; x <= y;", objects.Bool(true)},
		{"let x = 1712; x <= x;", objects.Bool(true)},
		{"1205 < 1205;", objects.Bool(false)},
		{"1 != 2;", objects.Bool(true)},
		{"1 != 2 == true != false;", objects.Bool(true)},
	}

	for _, tt := range tests {
		program, err := parser.Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: parse error: %v", tt.input, err)
		}

		env := NewEnvironment()
		typeEnv := NewTypeEnvironment()

		got, err := EvalProgram(program, env, typeEnv)
		if err != nil {
			t.Fatalf("%q: evaluation error: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Errorf("%q: got %s, expected %s", tt.input, got, tt.expected)
		}
	}
}

func TestEvalDivisionByZero(t *testing.T) {
	program, err := parser.Parse("10 / 0;")
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	_, err = EvalProgram(program, env, typeEnv)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestEvalVariable(t *testing.T) {
	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	value := objects.Int(1712)
	expected := objects.Int(-17)

	env.Set("x", value)
	typeEnv.Set("x", objects.IntType)

	program, err := parser.Parse("x + -1729;")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func TestEvalUndefinedVariable(t *testing.T) {
	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	program, err := parser.Parse("x;")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalProgram(program, env, typeEnv)
	if err == nil {
		t.Fatal("expected undefined variable error")
	}
}

func TestEvalLetStmt(t *testing.T) {
	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	program, err := parser.Parse("let x = 3;")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	program2, err := parser.Parse("x*x + 1;")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program2, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	if got != objects.Int(10) {
		t.Errorf("%q: got %d, expected %d", program2, got, 10)
	}
}

func TestEvalProgram(t *testing.T) {
	program, err := parser.Parse(`
	let x = 1712;
	let y = -1729;
	y = x + 2 * y;
	x + y;
	`)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	got, err := EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	if got != objects.Int(-34) {
		t.Errorf("got %d, expected -34", got)
	}

	x, err := env.Get("x")
	if err != nil {
		t.Fatal(err)
	}

	if x != objects.Int(1712) {
		t.Errorf("x = %d, expected 1712", x)
	}

	y, err := env.Get("y")
	if err != nil {
		t.Fatal(err)
	}

	if y != objects.Int(-1746) {
		t.Errorf("x = %d, expected -1746", y)
	}

}

func TestEvalProgramUsesPreviousStatements(t *testing.T) {
	program, err := parser.Parse(`
        let x = 1712;
        x + 1205;
        x - 1729;
    `)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	got, err := EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	if got != objects.Int(-17) {
		t.Errorf("got %d, expected -17", got)
	}
}

func TestEvalEmptyProgram(t *testing.T) {
	program, err := parser.Parse("")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program, NewEnvironment(), NewTypeEnvironment())
	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Errorf("got %v, expected nil", got)
	}
}

func TestEvalTypeErrors(t *testing.T) {
	tests := []string{
		"true + 1;",
		"1 + true;",
		"true + false;",
		"-true;",
		"+false;",
	}

	for _, input := range tests {
		program, err := parser.Parse(input)
		if err != nil {
			t.Fatalf("%q: parse error: %v", input, err)
		}

		_, err = EvalProgram(program, NewEnvironment(), NewTypeEnvironment())
		if err == nil {
			t.Errorf("%q: expected evaluation error", input)
		}
	}
}
