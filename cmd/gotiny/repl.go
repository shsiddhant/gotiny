package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/shsiddhant/gotiny/internal/checker"
	"github.com/shsiddhant/gotiny/internal/environment"
	"github.com/shsiddhant/gotiny/internal/evaluator"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	env := environment.NewEnvironment()
	typeEnv := checker.NewTypeEnvironment()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		program, err := parser.Parse(input)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		result, err := evaluator.EvalProgram(program, env, typeEnv)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if result != nil {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
}
