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

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentifierPart(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func (l *Lexer) identifierOrKeyword() Token {
	for isIdentifierPart(l.peek()) {
		l.advance()
	}

	value := string(l.source[l.start:l.current])

	if tokenType, ok := keywords[value]; ok {
		return Token{
			Type:  tokenType,
			Value: value,
		}
	}

	return Token{
		Type:  Identifier,
		Value: value,
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
	//Operators
	case '+':
		return Token{Type: Plus, Value: "+"}, nil
	case '-':
		return Token{Type: Minus, Value: "-"}, nil
	case '*':
		return Token{Type: Star, Value: "*"}, nil
	case '/':
		return Token{Type: Slash, Value: "/"}, nil
	case '=':
		return Token{Type: Equal, Value: "="}, nil

	// Delimiters
	case ';':
		return Token{Type: SemiColon, Value: ";"}, nil
	case '(':
		return Token{Type: LeftParen, Value: "("}, nil
	case ')':
		return Token{Type: RightParen, Value: ")"}, nil
	}

	if isDigit(c) {
		return l.number(), nil
	}

	if isLetter(c) {
		return l.identifierOrKeyword(), nil
	}

	return Token{}, fmt.Errorf("unexpected character: %q", c)
}
