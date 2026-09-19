package main

import (
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"-1", -1},
		{"+1", 1},
		{"--1", 1},
		{"-+1", -1},
		{"2 + -5", -3},
		{"2 * -3", -6},
		{"-2 * 3", -6},
		{"-5 / 3", -1},
		{"-(1 + 2)", -3},
		{"1 + 2", 3},
		{"2 * 3", 6},
		{"1 + 2 * 3", 7},
		{"(1 + 2) * 3", 9},
		{"20 / 5 / 2", 2},
		{"10 - 3 - 2", 5},
	}

	for _, tt := range tests {
		stmt, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: parse error: %v", tt.input, err)
		}

		env := NewEnvironment()

		got, err := EvalStmt(stmt, env)
		if err != nil {
			t.Fatalf("%q: evaluation error: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Errorf("%q: got %d, expected %d", tt.input, got, tt.expected)
		}
	}
}

func TestEvalDivisionByZero(t *testing.T) {
	stmt, err := Parse("10 / 0")
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()

	_, err = EvalStmt(stmt, env)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestEvalVariable(t *testing.T) {
	env := NewEnvironment()

	value := 1712
	expected := -17

	env.Set("x", value)

	stmt, err := Parse("x + -1729")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalStmt(stmt, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func TestEvalUndefinedVariable(t *testing.T) {
	env := NewEnvironment()

	stmt, err := Parse("x")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalStmt(stmt, env)
	if err == nil {
		t.Fatal("expected undefined variable error")
	}
}

func TestEvalLetStmt(t *testing.T) {
	env := NewEnvironment()

	stmt, err := Parse("let x = 3")
	if err != nil {
		t.Fatal(err)
	}

	_, err = EvalStmt(stmt, env)
	if err != nil {
		t.Fatal(err)
	}

	stmt2, err := Parse("x*x + 1")
	if err != nil {
		t.Fatal(err)
	}

	got, err := EvalStmt(stmt2, env)
	if err != nil {
		t.Fatal(err)
	}

	if got != 10 {
		t.Errorf("%q: got %d, expected %d", stmt2, got, 10)
	}
}
