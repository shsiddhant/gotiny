package evaluator

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func TestEvalIfElseMutation(t *testing.T) {
	progString := `
        let x = 2000;
        if x > 1729 {
            x = x - 1729;
        } else {
            x = 1729 - x;
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()
	typEnv := NewTypeEnvironment()

	_, err = EvalProgram(program, env, typEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(271) {
		t.Errorf("got %d, expected 271", got)
	}
}

func TestEvalVariableShadowing(t *testing.T) {
	progString := `
        let x = 1712;
        if true {
            let x = -1729;
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	env := NewEnvironment()
	typEnv := NewTypeEnvironment()

	_, err = EvalProgram(program, env, typEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(1712) {
		t.Errorf("expected outer x to stay 1712, got %v", got)
	}
}

func TestEvalNestedMutation(t *testing.T) {
	progString := `
        let x = 1712;
        if true {
            if true {
                x = 1205;
            }
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnvironment()
	typeEnv := NewTypeEnvironment()

	_, err = EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(1205) {
		t.Errorf("got %v, expected 1205", got)
	}
}
