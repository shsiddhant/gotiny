package main

import "testing"

func TestLexer(t *testing.T) {
	input := "let x = 123 + 4 * (2 - 1) - _tgb"

	lexer := NewLexer(input)

	expected := []Token{
		{Type: Let, Value: "let"},
		{Type: Identifier, Value: "x"},
		{Type: Equal, Value: "="},
		{Type: Number, Value: "123"},
		{Type: Plus, Value: "+"},
		{Type: Number, Value: "4"},
		{Type: Star, Value: "*"},
		{Type: LeftParen, Value: "("},
		{Type: Number, Value: "2"},
		{Type: Minus, Value: "-"},
		{Type: Number, Value: "1"},
		{Type: RightParen, Value: ")"},
		{Type: Minus, Value: "-"},
		{Type: Identifier, Value: "_tgb"},
		{Type: EOF},
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

		if got.Type == EOF {
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
