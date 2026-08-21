package main

import "testing"

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1",
			expected: "1",
		},
		{
			input:    "1 + 2",
			expected: "(1 + 2)",
		},
		{
			input:    "1 + 2 * 3",
			expected: "(1 + (2 * 3))",
		},
		{
			input:    "(1 + 2) * 3",
			expected: "((group (1 + 2)) * 3)",
		},
		{
			input:    "20 / 5 / 2",
			expected: "((20 / 5) / 2)",
		},
	}

	for _, tt := range tests {
		expr, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", tt.input, err)
		}

		got := expr.String()

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
