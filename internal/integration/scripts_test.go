package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shsiddhant/gotiny/internal/evaluator"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func evalScript(name string, t *testing.T) (objects.Value, error) {
	t.Helper()

	rootPath := "../../"
	filePath := filepath.Join(rootPath, "scripts/", name)

	source, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	program, err := parser.Parse(string(source))
	if err != nil {
		return nil, err
	}

	env := evaluator.NewEnvironment()
	typeEnv := evaluator.NewTypeEnvironment()

	return evaluator.EvalProgram(program, env, typeEnv)
}

func TestScripts(t *testing.T) {
	tests := []struct {
		name     string
		expected objects.Value
	}{
		{"v0.2.0/greenbutterfly.gt", objects.Bool(true)},
		{"v0.2.0/negation.gt", objects.Int(1729)},
		{"v0.2.0/precedence.gt", objects.Int(-1729)},
		{"v0.2.0/scope.gt", objects.Int(524)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalScript(tt.name, t)
			if err != nil {
				t.Fatal(err)
			}
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
