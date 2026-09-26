package lexer

import (
	"fmt"

	"github.com/shsiddhant/gotiny/internal/token"
)

type Lexer struct {
	source  []rune
	start   int
	current int
	line    int
	column  int
}

func NewLexer(source string) *Lexer {
	return &Lexer{source: []rune(source), line: 1, column: 1}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() rune {
	ch := l.source[l.current]
	l.current++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
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

func (l *Lexer) number(line int, column int) token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	return token.Token{
		Type:   token.Number,
		Value:  string(l.source[l.start:l.current]),
		Line:   line,
		Column: column,
	}
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentifierPart(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func (l *Lexer) identifierOrKeyword(line int, column int) token.Token {
	for isIdentifierPart(l.peek()) {
		l.advance()
	}

	value := string(l.source[l.start:l.current])

	if tokenType, ok := keywords[value]; ok {
		return token.Token{
			Type:   tokenType,
			Value:  value,
			Line:   line,
			Column: column,
		}
	}

	return token.Token{
		Type:   token.Identifier,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

func (l *Lexer) skipComment() {
	for !l.isAtEnd() {
		switch l.peek() {
		case '\n':
			return
		}
		l.advance()
	}
}

func (l *Lexer) Next() (token.Token, error) {
	l.skipWhitespace()
	l.start = l.current
	line, column := l.line, l.column

	if l.isAtEnd() {
		return token.Token{Type: token.EOF, Line: line, Column: column}, nil
	}

	c := l.advance()

	switch c {
	case '#':
		// Comments
		l.skipComment()
		return l.Next()
	// Operators
	case '+':
		return token.Token{
			Type:   token.Plus,
			Value:  "+",
			Line:   line,
			Column: column,
		}, nil
	case '-':
		return token.Token{
			Type:   token.Minus,
			Value:  "-",
			Line:   line,
			Column: column,
		}, nil
	case '*':
		return token.Token{
			Type:   token.Star,
			Value:  "*",
			Line:   line,
			Column: column,
		}, nil
	case '/':
		return token.Token{
			Type:   token.Slash,
			Value:  "/",
			Line:   line,
			Column: column,
		}, nil
	case '=':
		if l.peek() == '=' {
			l.advance()
			return token.Token{
				Type:   token.EqualEqual,
				Value:  "==",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{
			Type:   token.Equal,
			Value:  "=",
			Line:   line,
			Column: column,
		}, nil
	case '>':
		if l.peek() == '=' {
			l.advance()
			return token.Token{
				Type:   token.GreaterEqual,
				Value:  ">=",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{
			Type:   token.Greater,
			Value:  ">",
			Line:   line,
			Column: column,
		}, nil
	case '<':
		if l.peek() == '=' {
			l.advance()
			return token.Token{
				Type:   token.LessEqual,
				Value:  "<=",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{
			Type:   token.Less,
			Value:  "<",
			Line:   line,
			Column: column,
		}, nil
	case '!':
		if l.peek() == '=' {
			l.advance()
			return token.Token{
				Type:   token.NotEqual,
				Value:  "!=",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{Type: token.Not, Value: "!", Line: line, Column: column}, nil
	case '&':
		peek := l.peek()
		if l.peek() == '&' {
			l.advance()
			return token.Token{
				Type:   token.And,
				Value:  "&&",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{
			Line:   line,
			Column: column,
		}, fmt.Errorf(
			"expected '&', got '%c'",
			peek,
		)
	case '|':
		peek := l.peek()
		if l.peek() == '|' {
			l.advance()
			return token.Token{
				Type:   token.Or,
				Value:  "||",
				Line:   line,
				Column: column,
			}, nil
		}
		return token.Token{
			Line:   line,
			Column: column,
		}, fmt.Errorf(
			"expected '|', got '%c'",
			peek,
		)
	// Delimiters
	case ';':
		return token.Token{
			Type:   token.SemiColon,
			Value:  ";",
			Line:   line,
			Column: column,
		}, nil
	case '(':
		return token.Token{
			Type:   token.LeftParen,
			Value:  "(",
			Line:   line,
			Column: column,
		}, nil
	case ')':
		return token.Token{
			Type:   token.RightParen,
			Value:  ")",
			Line:   line,
			Column: column,
		}, nil
	case '{':
		return token.Token{
			Type:   token.LeftCurlyBrace,
			Value:  "{",
			Line:   line,
			Column: column,
		}, nil
	case '}':
		return token.Token{
			Type:   token.RightCurlyBrace,
			Value:  "}",
			Line:   line,
			Column: column,
		}, nil
	case ',':
		return token.Token{
			Type:   token.Comma,
			Value:  ",",
			Line:   line,
			Column: column,
		}, nil
	}

	if isDigit(c) {
		return l.number(line, column), nil
	}

	if isLetter(c) {
		return l.identifierOrKeyword(line, column), nil
	}

	return token.Token{
		Line:   line,
		Column: column,
	}, fmt.Errorf(
		"unexpected character: \"%c\"",
		c,
	)
}

var keywords = map[string]token.TokenType{
	"let":    token.Let,
	"true":   token.True,
	"false":  token.False,
	"if":     token.If,
	"else":   token.Else,
	"fn":     token.Fn,
	"Int":    token.Int,
	"Bool":   token.Bool,
	"return": token.Return,
	"print":  token.Print,
}
