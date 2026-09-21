package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shsiddhant/gotiny/internal/evaluator"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func runFile(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	source, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	program, err := parser.Parse(string(source))
	if err != nil {
		return err
	}

	env := evaluator.NewEnvironment()
	typeEnv := evaluator.NewTypeEnvironment()

	result, err := evaluator.EvalProgram(program, env, typeEnv)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
