package token

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
	True
	False
	If
	Else

	// Operators

	// Arithmetic
	Plus  // +
	Minus // -
	Star  // *
	Slash // /

	// Assign
	Equal // =

	// Comparison
	Less         // <
	LessEqual    // <=
	Greater      // >
	GreaterEqual // >=

	// Equality
	EqualEqual // ==
	NotEqual   // !=

	// Not
	Not // !

	// Delimiters
	SemiColon       // ;
	LeftParen       // (
	RightParen      // )
	LeftCurlyBrace  // {
	RightCurlyBrace // }
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (token Token) String() string {
	return fmt.Sprintf("%d:%d:%s(%s)", token.Line, token.Column, token.Type, token.Value)
}
