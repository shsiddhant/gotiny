package main

import "fmt"

//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
type TokenType int

const (
	InvalidToken TokenType = iota
	EOF
	Number
	Identifier

	// Keywords
	Let

	// Operators
	Plus
	Minus
	Star
	Slash
	Equal

	// Delimiters
	SemiColon
	LeftParen
	RightParen
)

type Token struct {
	Type  TokenType
	Value string
}

func (token Token) String() string {
	return fmt.Sprintf("%s(%s)", token.Type, token.Value)
}

var keywords = map[string]TokenType{
	"let": Let,
}
