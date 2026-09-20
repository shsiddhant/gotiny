package lexer

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/token"
)

func TestLexer(t *testing.T) {
	input := "let x = 123 + 4 * (2 - 1) - _tgb;"

	lexer := NewLexer(input)

	expected := []token.Token{
		{Type: token.Let, Value: "let"},
		{Type: token.Identifier, Value: "x"},
		{Type: token.Equal, Value: "="},
		{Type: token.Number, Value: "123"},
		{Type: token.Plus, Value: "+"},
		{Type: token.Number, Value: "4"},
		{Type: token.Star, Value: "*"},
		{Type: token.LeftParen, Value: "("},
		{Type: token.Number, Value: "2"},
		{Type: token.Minus, Value: "-"},
		{Type: token.Number, Value: "1"},
		{Type: token.RightParen, Value: ")"},
		{Type: token.Minus, Value: "-"},
		{Type: token.Identifier, Value: "_tgb"},
		{Type: token.SemiColon, Value: ";"},
		{Type: token.EOF},
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
		{Type: token.Let, Value: "let"},
		{Type: token.Identifier, Value: "x"},
		{Type: token.Equal, Value: "="},
		{Type: token.True, Value: "true"},
		{Type: token.SemiColon, Value: ";"},
		{Type: token.EOF},
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
