package main

import (
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1;", 1},
		{"-1;", -1},
		{"+1;", 1},
		{"--1;", 1},
		{"-+1;", -1},
		{"2 + -5;", -3},
		{"2 * -3;", -6},
		{"-2 * 3;", -6},
		{"-5 / 3;", -1},
		{"-(1 + 2);", -3},
		{"1 + 2;", 3},
		{"2 * 3;", 6},
		{"1 + 2 * 3;", 7},
		{"(1 + 2) * 3;", 9},
		{"20 / 5 / 2;", 2},
		{"10 - 3 - 2;", 5},
	}

	for _, tt := range tests {
		program, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: parse error: %v", tt.input, err)
		}

		env := NewEnvironment()

		got, err := EvalProgram(program, env)
		if err != nil {
			t.Fatalf("%q: evaluation error: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Errorf("%q: got %d, expected %d", tt.input, got, tt.expected)
		}
	}
}

func TestEvalDivisionByZero(t *testing.T) {
	program, err := Parse("10 / 0;")
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()

	_, err = EvalProgram(program, env)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestEvalVariable(t *testing.T) {
	env := NewEnvironment()

	value := 1712
	expected := -17

	env.Set("x", value)

	program, err := Parse("x + -1729;")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func TestEvalUndefinedVariable(t *testing.T) {
	env := NewEnvironment()

	program, err := Parse("x;")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalProgram(program, env)
	if err == nil {
		t.Fatal("expected undefined variable error")
	}
}

func TestEvalLetStmt(t *testing.T) {
	env := NewEnvironment()

	program, err := Parse("let x = 3;")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalProgram(program, env)
	if err != nil {
		t.Fatal(err)
	}

	program2, err := Parse("x*x + 1;")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program2, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != 10 {
		t.Errorf("%q: got %d, expected %d", program2, got, 10)
	}
}

func TestEvalProgram(t *testing.T) {
	program, err := Parse(`
	let x = 1712;
	let y = -1729;
	x + y;
	`)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()

	got, err := EvalProgram(program, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != -17 {
		t.Errorf("got %d, expected -17", got)
	}

	x, err := env.Get("x")
	if err != nil {
		t.Fatal(err)
	}

	if x != 1712 {
		t.Errorf("x = %d, expected 1712", x)
	}

	y, err := env.Get("y")
	if err != nil {
		t.Fatal(err)
	}

	if y != -1729 {
		t.Errorf("x = %d, expected -1729", y)
	}

}

func TestEvalProgramUsesPreviousStatements(t *testing.T) {
	program, err := Parse(`
        let x = 1712;
        x + 1205;
        x - 1729;
    `)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()

	got, err := EvalProgram(program, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != -17 {
		t.Errorf("got %d, expected -17", got)
	}
}

func TestEvalEmptyProgram(t *testing.T) {
	program, err := Parse("")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalProgram(program, NewEnvironment())
	if err != nil {
		t.Fatal(err)
	}

	if got != 0 {
		t.Errorf("got %d, expected 0", got)
	}
}
