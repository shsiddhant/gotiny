package lexer

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/token"
)

type Lexer struct {
	source  []rune
	start   int
	current int
}

func NewLexer(source string) *Lexer {
	return &Lexer{source: []rune(source)}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() rune {
	ch := l.source[l.current]
	l.current++
	return ch
}

func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) skipWhitespace() {
	for !l.isAtEnd() {
		switch l.peek() {
		case ' ', '\t', '\n', '\r':
			l.advance()
		default:
			return
		}
	}
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (l *Lexer) number() token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	return token.Token{
		Type:  token.Number,
		Value: string(l.source[l.start:l.current]),
	}
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentifierPart(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func (l *Lexer) identifierOrKeyword() token.Token {
	for isIdentifierPart(l.peek()) {
		l.advance()
	}

	value := string(l.source[l.start:l.current])

	if tokenType, ok := keywords[value]; ok {
		return token.Token{
			Type:  tokenType,
			Value: value,
		}
	}

	return token.Token{
		Type:  token.Identifier,
		Value: value,
	}
}

func (l *Lexer) Next() (token.Token, error) {
	l.skipWhitespace()
	l.start = l.current

	if l.isAtEnd() {
		return token.Token{Type: token.EOF}, nil
	}

	c := l.advance()

	switch c {
	//Operators
	case '+':
		return token.Token{Type: token.Plus, Value: "+"}, nil
	case '-':
		return token.Token{Type: token.Minus, Value: "-"}, nil
	case '*':
		return token.Token{Type: token.Star, Value: "*"}, nil
	case '/':
		return token.Token{Type: token.Slash, Value: "/"}, nil
	case '=':
		return token.Token{Type: token.Equal, Value: "="}, nil

	// Delimiters
	case ';':
		return token.Token{Type: token.SemiColon, Value: ";"}, nil
	case '(':
		return token.Token{Type: token.LeftParen, Value: "("}, nil
	case ')':
		return token.Token{Type: token.RightParen, Value: ")"}, nil
	case '{':
		return token.Token{Type: token.LeftCurlyBrace, Value: "{"}, nil
	case '}':
		return token.Token{Type: token.RightCurlyBrace, Value: "}"}, nil
	}

	if isDigit(c) {
		return l.number(), nil
	}

	if isLetter(c) {
		return l.identifierOrKeyword(), nil
	}

	return token.Token{}, fmt.Errorf("unexpected character: %q", c)
}

var keywords = map[string]token.TokenType{
	"let":   token.Let,
	"true":  token.True,
	"false": token.False,
	"if":    token.If,
	"else":  token.Else,
}
