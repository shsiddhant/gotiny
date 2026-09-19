package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1;",
			expected: "1",
		},
		{
			input:    "-1;",
			expected: "(-1)",
		},
		{
			input:    "+1;",
			expected: "(+1)",
		},
		{
			input:    "--1;",
			expected: "(-(-1))",
		},
		{
			input:    "2 + -5;",
			expected: "(2 + (-5))",
		},
		{
			input:    "2 * -3;",
			expected: "(2 * (-3))",
		},
		{
			input:    "-2 * 3;",
			expected: "((-2) * 3)",
		},
		{
			input:    "-(1 + 2);",
			expected: "(-(group (1 + 2)))",
		},
		{
			input:    "1 + 2;",
			expected: "(1 + 2)",
		},
		{
			input:    "1 + 2 * 3;",
			expected: "(1 + (2 * 3))",
		},
		{
			input:    "(1 + 2) * 3;",
			expected: "((group (1 + 2)) * 3)",
		},
		{
			input:    "20 / 5 / 2;",
			expected: "((20 / 5) / 2)",
		},
		{
			input:    "-x + 2;",
			expected: "((-x) + 2)",
		},
		{
			input:    "let x = y + 2;",
			expected: "let x = (y + 2)",
		},
		{
			input:    "let x = -y + 2;",
			expected: "let x = ((-y) + 2)",
		},
		{
			input:    "let x = true;",
			expected: "let x = true",
		},
	}

	for _, tt := range tests {
		program, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", tt.input, err)
		}

		got := program.Statements[0].String()

		if got != tt.expected {
			t.Errorf("%q: got %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestParserRejectsTrailingTokens(t *testing.T) {
	_, err := Parse("1 2")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRejectsIncompleteExpression(t *testing.T) {
	_, err := Parse("1 +")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRejectsUnclosedGrouping(t *testing.T) {
	_, err := Parse("(1 + 2")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRejectsNestedLet(t *testing.T) {
	_, err := Parse("let x = let y = 2")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionWithLet(t *testing.T) {
	_, err := Parse("1712 + let x = 1729")

	if err == nil {
		t.Fatal(err)
	}
}

func TestParserProgram(t *testing.T) {

	progString := `
	let x =1712;
	let y =-1729;
	x-y;
	`

	expected := []string{
		"let x = 1712",
		"let y = (-1729)",
		"(x - y)",
	}

	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("got %d statements, expected 3", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		got := stmt.String()
		if got != expected[i] {
			t.Errorf("got %q, expected %q", got, expected[i])
		}
	}
}
