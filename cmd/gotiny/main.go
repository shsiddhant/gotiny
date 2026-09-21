package main

import (
	"fmt"
	"os"
)

func main() {
	switch len(os.Args) {
	case 1:
		runREPL()
	case 2:
		if err := runFile(os.Args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: gotiny [file]")
		os.Exit(1)
	}

}
