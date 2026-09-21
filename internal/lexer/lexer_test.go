package lexer

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/token"
)

func TestLexer(t *testing.T) {
	input := "let x = 123 + 4 * (2 - 1) -  _tgb;"

	lexer := NewLexer(input)

	expected := []token.Token{
		{Type: token.Let, Value: "let", Line: 1, Column: 1},
		{Type: token.Identifier, Value: "x", Line: 1, Column: 5},
		{Type: token.Equal, Value: "=", Line: 1, Column: 7},
		{Type: token.Number, Value: "123", Line: 1, Column: 9},
		{Type: token.Plus, Value: "+", Line: 1, Column: 13},
		{Type: token.Number, Value: "4", Line: 1, Column: 15},
		{Type: token.Star, Value: "*", Line: 1, Column: 17},
		{Type: token.LeftParen, Value: "(", Line: 1, Column: 19},
		{Type: token.Number, Value: "2", Line: 1, Column: 20},
		{Type: token.Minus, Value: "-", Line: 1, Column: 22},
		{Type: token.Number, Value: "1", Line: 1, Column: 24},
		{Type: token.RightParen, Value: ")", Line: 1, Column: 25},
		{Type: token.Minus, Value: "-", Line: 1, Column: 27},
		{Type: token.Identifier, Value: "_tgb", Line: 1, Column: 30},
		{Type: token.SemiColon, Value: ";", Line: 1, Column: 34},
		{Type: token.EOF, Line: 1, Column: 35},
	}

	for i := 0; ; i++ {
		got, err := lexer.Next()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("lexer produced unexpected token: %v", got)
		}

		if got != expected[i] {
			t.Errorf("token %d: got %v, expected %v", i, got, expected[i])
		}

		if got.Type == token.EOF {
			if i != len(expected)-1 {
				t.Errorf("lexer reached EOF early: got %d tokens, expected %d", i+1, len(expected))
			}
			break
		}
	}
}

func TestLexerUnexpectedCharacter(t *testing.T) {
	lexer := NewLexer("1 @ 2")

	_, err := lexer.Next()
	if err != nil {
		t.Fatal(err)
	}

	_, err = lexer.Next()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLexerKeywords(t *testing.T) {
	lexer := NewLexer("let x = true;")

	expected := []token.Token{
		{Type: token.Let, Value: "let", Line: 1, Column: 1},
		{Type: token.Identifier, Value: "x", Line: 1, Column: 5},
		{Type: token.Equal, Value: "=", Line: 1, Column: 7},
		{Type: token.True, Value: "true", Line: 1, Column: 9},
		{Type: token.SemiColon, Value: ";", Line: 1, Column: 13},
		{Type: token.EOF, Line: 1, Column: 14},
	}

	for i := 0; ; i++ {
		got, err := lexer.Next()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("lexer produced unexpected token: %v", got)
		}

		if got != expected[i] {
			t.Errorf("token %d: got %v, expected %v", i, got, expected[i])
		}

		if got.Type == token.EOF {
			if i != len(expected)-1 {
				t.Errorf("lexer reached EOF early: got %d tokens, expected %d", i+1, len(expected))
			}
			break
		}
	}

}

func TestLexerPositionMultiLine(t *testing.T) {
	source := `let x = 1729;
let y = 1712;
y - x;`
	lexer := NewLexer(source)

	expected := []token.Token{
		{Type: token.Let, Value: "let", Line: 1, Column: 1},
		{Type: token.Identifier, Value: "x", Line: 1, Column: 5},
		{Type: token.Equal, Value: "=", Line: 1, Column: 7},
		{Type: token.Number, Value: "1729", Line: 1, Column: 9},
		{Type: token.SemiColon, Value: ";", Line: 1, Column: 13},

		{Type: token.Let, Value: "let", Line: 2, Column: 1},
		{Type: token.Identifier, Value: "y", Line: 2, Column: 5},
		{Type: token.Equal, Value: "=", Line: 2, Column: 7},
		{Type: token.Number, Value: "1712", Line: 2, Column: 9},
		{Type: token.SemiColon, Value: ";", Line: 2, Column: 13},

		{Type: token.Identifier, Value: "y", Line: 3, Column: 1},
		{Type: token.Minus, Value: "-", Line: 3, Column: 3},
		{Type: token.Identifier, Value: "x", Line: 3, Column: 5},
		{Type: token.SemiColon, Value: ";", Line: 3, Column: 6},
		{Type: token.EOF, Value: "", Line: 3, Column: 7},
	}

	for i := 0; ; i++ {
		got, err := lexer.Next()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("lexer produced unexpected token: %v", got)
		}

		if got != expected[i] {
			t.Errorf("token %d: got %v, expected %v", i, got, expected[i])
		}

		if got.Type == token.EOF {
			if i != len(expected)-1 {
				t.Errorf("lexer reached EOF early: got %d tokens, expected %d", i+1, len(expected))
			}
			break
		}
	}
}

func TestLexerComparison(t *testing.T) {
	lexer := NewLexer(`x>=1712
y===1205
1 != 3
!==
a > b
> =`)

	expected := []token.Token{
		{Type: token.Identifier, Value: "x", Line: 1, Column: 1},
		{Type: token.GreaterEqual, Value: ">=", Line: 1, Column: 2},
		{Type: token.Number, Value: "1712", Line: 1, Column: 4},
		{Type: token.Identifier, Value: "y", Line: 2, Column: 1},
		{Type: token.EqualEqual, Value: "==", Line: 2, Column: 2},
		{Type: token.Equal, Value: "=", Line: 2, Column: 4},
		{Type: token.Number, Value: "1205", Line: 2, Column: 5},
		{Type: token.Number, Value: "1", Line: 3, Column: 1},
		{Type: token.NotEqual, Value: "!=", Line: 3, Column: 3},
		{Type: token.Number, Value: "3", Line: 3, Column: 6},
		{Type: token.NotEqual, Value: "!=", Line: 4, Column: 1},
		{Type: token.Equal, Value: "=", Line: 4, Column: 3},
		{Type: token.Identifier, Value: "a", Line: 5, Column: 1},
		{Type: token.Greater, Value: ">", Line: 5, Column: 3},
		{Type: token.Identifier, Value: "b", Line: 5, Column: 5},
		{Type: token.Greater, Value: ">", Line: 6, Column: 1},
		{Type: token.Equal, Value: "=", Line: 6, Column: 3},
		{Type: token.EOF, Line: 6, Column: 4},
	}

	for i := 0; ; i++ {
		got, err := lexer.Next()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("lexer produced unexpected token: %v", got)
		}

		if got != expected[i] {
			t.Errorf("token %d: got %v, expected %v", i, got, expected[i])
		}

		if got.Type == token.EOF {
			if i != len(expected)-1 {
				t.Errorf("lexer reached EOF early: got %d tokens, expected %d", i+1, len(expected))
			}
			break
		}
	}
}

func TestLexerBoolean(t *testing.T) {
	lexer := NewLexer(`x == 1 && y`)

	expected := []token.Token{
		{Type: token.Identifier, Value: "x", Line: 1, Column: 1},
		{Type: token.EqualEqual, Value: "==", Line: 1, Column: 3},
		{Type: token.Number, Value: "1", Line: 1, Column: 6},
		{Type: token.And, Value: "&&", Line: 1, Column: 8},
		{Type: token.Identifier, Value: "y", Line: 1, Column: 11},
		{Type: token.EOF, Line: 1, Column: 12},
	}

	for i := 0; ; i++ {
		got, err := lexer.Next()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("lexer produced unexpected token: %v", got)
		}

		if got != expected[i] {
			t.Errorf("token %d: got %v, expected %v", i, got, expected[i])
		}

		if got.Type == token.EOF {
			if i != len(expected)-1 {
				t.Errorf("lexer reached EOF early: got %d tokens, expected %d", i+1, len(expected))
			}
			break
		}
	}
}
