package main

import "fmt"

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

func (l *Lexer) number() Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	return Token{
		Type:  Number,
		Value: string(l.source[l.start:l.current]),
	}
}

func (l *Lexer) Next() (Token, error) {
	l.skipWhitespace()
	l.start = l.current

	if l.isAtEnd() {
		return Token{Type: EOF}, nil
	}

	c := l.advance()

	switch c {
	case '+':
		return Token{Type: Plus, Value: "+"}, nil
	case '-':
		return Token{Type: Minus, Value: "-"}, nil
	case '*':
		return Token{Type: Star, Value: "*"}, nil
	case '/':
		return Token{Type: Slash, Value: "/"}, nil
	case '(':
		return Token{Type: LeftParen, Value: "("}, nil
	case ')':
		return Token{Type: RightParen, Value: ")"}, nil
	}

	if isDigit(c) {
		return l.number(), nil
	}

	return Token{}, fmt.Errorf("unexpected character: %q", c)
}
