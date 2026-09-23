package checker

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/parser"
)

func TestCheckIfStatementValid(t *testing.T) {
	progString := `
        let x = 10;
        if true {
            let y = x + 5;
        }
    `
	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	typeEnv := NewTypeEnvironment()
	err = CheckProgram(program, typeEnv)
	if err != nil {
		t.Error(err)
	}
}

func TestIfElseScoping(t *testing.T) {
	progString := `
if false {
	let a = 1712;
} else {
    let a = 1729;
}
let c = a;`
	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	typeEnv := NewTypeEnvironment()
	err = CheckProgram(program, typeEnv)
	if err == nil {
		t.Error("expected undefined variable error")
	}
}

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
		{`if 1712 {
            let x = 1205;
        }`},
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
