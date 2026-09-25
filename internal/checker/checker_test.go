package checker

import (
	"strings"
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

func TestTypeChecker(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		errorSubstr string
	}{
		{
			name: "Shadowing should be allowed",
			source: `
				let x = 1712;
				if x < 1729 {
					let x = -x;
				}
			`,
			expectError: false,
		},
		{
			name: "If statement else block shouldn't see body env",
			source: `
				if true {
    				let x = 1712;
				} else {
				    x;
				}
			`,
			expectError: true,
			errorSubstr: "undefined name: x",
		},
		{
			name: "If statement body shouldn't see else block env",
			source: `
				if true {
    				y = 1205;
				} else {
				    let y = 1013;
				}
			`,
			expectError: true,
			errorSubstr: "undefined name: y",
		},
		{
			name: "Valid sequential return",
			source: `
				fn square(x Int) Int {
					return x * x;
				}
			`,
			expectError: false,
		},
		{
			name: "Valid conditional paths (both sides return)",
			source: `
				fn max(a Int, b Int) Int {
					if a > b {
						return a;
					} else {
						return b;
					}
				}
			`,
			expectError: false,
		},
		{
			name: "Missing return path (no else branch)",
			source: `
				fn abs(x Int) Int {
					if x > 0 {
						return x;
					}
				}
			`,
			expectError: true,
			errorSubstr: `missing return statement at end of function "abs"`,
		},
		{
			name: "Missing return path (nested statement missing return)",
			source: `
				fn badMax(a Int, b Int) Int {
					if a > b {
						return a;
					} else {
						let x = b;
					}
				}
			`,
			expectError: true,
			errorSubstr: `missing return statement at end of function "badMax"`,
		},
		{
			name: "Implicit Void flow allowed to omit return",
			source: `
				fn greenbutterfly(us Bool) {
					if us {
						let theMonster = -1712;
					}
				}
			`,
			expectError: false,
		},
		{
			name: "Return outside function block should fail",
			source: `
				let x = 5;
				return x;
			`,
			expectError: true,
			errorSubstr: "return not allowed outside function",
		},
		{
			name: "Mismatched return type",
			source: `
				fn checkFlag(flag Bool) Int {
					return flag;
				}
			`,
			expectError: true,
			errorSubstr: "cannot return BoolType from function expecting IntType",
		},
		{
			name: "Duplicate parameter names should fail",
			source: `
				fn add(x Int, x Int) Int {
					return x + x;
				}
			`,
			expectError: true,
			errorSubstr: `duplicate parameter name: x`,
		},
		{
			name: "Multiple different duplicate parameters should fail",
			source: `
				fn detectiveConan(val Int, flag Bool, val Bool) {
					let x = 1;
				}
			`,
			expectError: true,
			errorSubstr: `duplicate parameter name: val`,
		},
		{
			name: "Parameters shadowing outer variables is legal",
			source: `
				let x = 10;
				fn shadowTest(x Bool) Bool {
					return x;
				}
			`,
			expectError: false,
		},
		{
			name: "Function name same as already existing variable",
			source: `
				let carissasWierd = 1116;
				fn carissasWierd() Int {
					return 1116;
				}
			`,
			expectError: true,
			errorSubstr: `name already defined: carissasWierd`,
		},
		{
			name: "Valid function call with matching arguments",
			source: `
				fn add(a Int, b Int) Int {
					return a + b;
				}
				let result = add(5, 10);
			`,
			expectError: false,
		},
		{
			name: "Function call with wrong argument count should fail",
			source: `
				fn greet(id Int, status Bool) {
					let x = id;
				}
				greet(101);
			`,
			expectError: true,
			errorSubstr: `expected 2 args, got 1`,
		},
		{
			name: "Function call with mismatched argument type should fail",
			source: `
				fn toggle(state Bool) Bool {
					return !state;
				}
				toggle(1729);
			`,
			expectError: true,
			errorSubstr: `expected BoolType arg, got IntType`,
		},
		{
			name: "Recursive functions should be allowed",
			source: `
				fn factorial(n Int) Int {
					if n < 0 {
						return 0;
					}
					if n == 0 {
						return 1;
					} else {
						return n * factorial(n-1);
					}

				}
			`,
			expectError: false,
		},
		{
			name: "Redeclaration of parameter inside immediate function scope not allowed",
			source: `
				fn f(x Int) Int {
    				let x = 100;
    				return x;
				}
			`,
			expectError: true,
			errorSubstr: "name already defined: x",
		},
		{
			name: "Shadowing works as expected inside nested block in functions.",
			source: `
				fn f(x Int) Bool {
					if x > 0 {
						let x = true;
						return x;
					}
					return false;
				}
			`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parser.Parse(tt.source)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}

			typeEnv := NewTypeEnvironment()

			err = CheckProgram(program, typeEnv)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected a static check error, but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf(
						"expected error to contain %q, got total message: %q",
						tt.errorSubstr,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected checker error: %v", err)
				}
			}
		})
	}
}
