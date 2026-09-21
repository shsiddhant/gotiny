package evaluator

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/parser"
)

func TestTypeCheckerErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"1 + true;"},
		{"!1;"},
		{"-true;"},
		{"false && 5;"},
		{"undeclaredVar + 1;"},
		{"let x = 1; let x = 2;"},
		{"let x = 1; x = true;"},
	}

	for _, tt := range tests {
		program, err := parser.Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: unexpected parse error: %v", tt.input, err)
		}

		typeEnv := NewTypeEnvironment()
		err = CheckProgram(program, typeEnv)
		if err == nil {
			t.Errorf("%q: expected type check error, got nil", tt.input)
		}
	}
}
