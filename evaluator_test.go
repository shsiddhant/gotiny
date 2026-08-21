package main

import "testing"

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
		expr, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: parse error: %v", tt.input, err)
		}

		got, err := Eval(expr)
		if err != nil {
			t.Fatalf("%q: evaluation error: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Errorf("%q: got %d, expected %d", tt.input, got, tt.expected)
		}
	}
}

func TestEvalDivisionByZero(t *testing.T) {
	expr, err := Parse("10 / 0")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Eval(expr)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}
