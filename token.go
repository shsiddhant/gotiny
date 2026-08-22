package main

import "fmt"

type TokenType int

func (t TokenType) String() string {
	switch t {
	case InvalidToken:
		return "Invalid"
	case EOF:
		return "EOF"
	case Number:
		return "Number"
	case Identifier:
		return "Identifier"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	default:
		return "Unknown"
	}
}

const (
	InvalidToken TokenType = iota
	EOF
	Number
	Identifier
	Plus
	Minus
	Star
	Slash
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
