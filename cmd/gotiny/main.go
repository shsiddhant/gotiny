package main

import (
	"fmt"
	"os"

	"github.com/shsiddhant/gotiny/internal/checker"
	"github.com/shsiddhant/gotiny/internal/evaluator"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func main() {
	switch len(os.Args) {
	case 1:
		runREPL()
	case 2:
		if err := runFile(os.Args[1]); err != nil {
			fmt.Fprintln(os.Stderr, formatError(err, os.Args[1]))
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: gotiny [file]")
		os.Exit(1)
	}

}

func formatError(err error, filePath string) string {
	switch e := err.(type) {
	case *parser.ParseError:
		return fmt.Sprintf("%s:%d:%d: %s", filePath, e.Token.Line, e.Token.Column, e.Message)
	case *checker.CheckError:
		return fmt.Sprintf("%s:%d:%d: %s", filePath, e.Token.Line, e.Token.Column, e.Message)
	case *evaluator.EvalError:
		return fmt.Sprintf("%s:%d:%d: %s", filePath, e.Line, e.Column, e.Message)
	}
	return err.Error()
}
