package main

import (
	"bufio"
	"fmt"
	"os"
)

func RunREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	env := NewEnvironment()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		stmt, err := Parse(input)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		result, err := EvalStmt(stmt, env)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
}
