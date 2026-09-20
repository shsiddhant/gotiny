package parser

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/ast"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1;",
			expected: "1",
		},
		{
			input:    "-1;",
			expected: "(-1)",
		},
		{
			input:    "+1;",
			expected: "(+1)",
		},
		{
			input:    "--1;",
			expected: "(-(-1))",
		},
		{
			input:    "2 + -5;",
			expected: "(2 + (-5))",
		},
		{
			input:    "2 * -3;",
			expected: "(2 * (-3))",
		},
		{
			input:    "-2 * 3;",
			expected: "((-2) * 3)",
		},
		{
			input:    "-(1 + 2);",
			expected: "(-(group (1 + 2)))",
		},
		{
			input:    "1 + 2;",
			expected: "(1 + 2)",
		},
		{
			input:    "1 + 2 * 3;",
			expected: "(1 + (2 * 3))",
		},
		{
			input:    "(1 + 2) * 3;",
			expected: "((group (1 + 2)) * 3)",
		},
		{
			input:    "20 / 5 / 2;",
			expected: "((20 / 5) / 2)",
		},
		{
			input:    "-x + 2;",
			expected: "((-x) + 2)",
		},
		{
			input:    "let x = y + 2;",
			expected: "let x = (y + 2)",
		},
		{
			input:    "let x = -y + 2;",
			expected: "let x = ((-y) + 2)",
		},
		{
			input:    "let x = true;",
			expected: "let x = true",
		},
		{
			input:    "x = true;",
			expected: "x = true",
		},
	}

	for _, tt := range tests {
		program, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", tt.input, err)
		}

		got := program.Statements[0].String()

		if got != tt.expected {
			t.Errorf("%q: got %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestParserRejectsTrailingTokens(t *testing.T) {
	_, err := Parse("1 2")

	if err == nil {
		t.Fatal("expected error")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError: %+v", parseErr)
	}

	if parseErr.Token.Column != 3 {
		t.Fatalf("expected parse error at column 3, got %d", parseErr.Token.Column)
	}
	if parseErr.Token.Line != 1 {
		t.Fatalf("expected parse error at line: 1, got %d", parseErr.Token.Line)
	}
	if parseErr.Token.Value != "2" {
		t.Fatalf("expected parse error for token value 2, got %s", parseErr.Token.Value)
	}

}

func TestParserRejectsIncompleteExpression(t *testing.T) {
	_, err := Parse("1 +")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRejectsUnclosedGrouping(t *testing.T) {
	_, err := Parse("(1 + 2")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRejectsNestedLet(t *testing.T) {
	_, err := Parse("let x = let y = 2")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionWithLet(t *testing.T) {
	_, err := Parse("1712 + let x = 1729")

	if err == nil {
		t.Fatal(err)
	}
}

func TestParserProgram(t *testing.T) {

	progString := `
	let x =1712;
	let y =-1729;
	x-y;
	`

	expected := []string{
		"let x = 1712",
		"let y = (-1729)",
		"(x - y)",
	}

	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("got %d statements, expected 3", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		got := stmt.String()
		if got != expected[i] {
			t.Errorf("got %q, expected %q", got, expected[i])
		}
	}
}

func TestParserIfStmt(t *testing.T) {
	progString := `
	if x {
		x = 1712;
		x;
	} else {
		x = 1729;
		x - 1712;
	}
	`
	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatal("expected an if statement")
	}
	if ifStmt.Cond.String() != "x" {
		t.Fatalf("got %s, expected x", ifStmt.Cond.String())
	}
	body := ifStmt.Body
	if len(body.Statements) != 2 {
		t.Fatalf("got %d statements, expected 2", len(body.Statements))
	}
	expectedBody := []string{
		"x = 1712",
		"x",
	}
	for i, got := range body.Statements {
		if got.String() != expectedBody[i] {
			t.Fatalf("got %q, expected %q", got, expectedBody[i])
		}
	}
	elseBlock := ifStmt.Else
	if elseBlock == nil {
		t.Fatal("expected else block, got nil")
	}
	expectedElse := []string{
		"x = 1729",
		"(x - 1712)",
	}
	if len(elseBlock.Statements) != 2 {
		t.Fatalf("got %d statements, expected 2", len(elseBlock.Statements))
	}
	for i, got := range elseBlock.Statements {
		if got.String() != expectedElse[i] {
			t.Fatalf("got %q, expected %q", got, expectedElse[i])
		}
	}
}

func TestParserNestedIfStmt(t *testing.T) {
	progString := `
    if x {
        if y {
            x = 1;
        }
    }
    `
	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatal("expected an if statement")
	}
	if ifStmt.Cond.String() != "x" {
		t.Fatalf("got %s, expected x", ifStmt.Cond.String())
	}
	if len(ifStmt.Body.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(ifStmt.Body.Statements))
	}

	innerIf, ok := ifStmt.Body.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatal("expected nested if statement")
	}
	if innerIf.Cond.String() != "y" {
		t.Fatalf("got %s, expected y", innerIf.Cond.String())
	}
	if len(innerIf.Body.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(innerIf.Body.Statements))
	}
	if innerIf.Body.Statements[0].String() != "x = 1" {
		t.Fatalf("got %q, expected %q", innerIf.Body.Statements[0], "x = 1")
	}
}

func TestParserEmptyIfBlock(t *testing.T) {
	progString := `
    if x {
    } else {
    }
    `
	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatal("expected an if statement")
	}
	if ifStmt.Cond.String() != "x" {
		t.Fatalf("got %s, expected x", ifStmt.Cond.String())
	}

	if len(ifStmt.Body.Statements) != 0 {
		t.Fatalf("got %d statements in body, expected 0", len(ifStmt.Body.Statements))
	}

	elseBlock := ifStmt.Else
	if elseBlock == nil {
		t.Fatal("expected else block, got nil")
	}
	if len(elseBlock.Statements) != 0 {
		t.Fatalf("got %d statements in else block, expected 0", len(elseBlock.Statements))
	}
}

func TestParserIfStmtNoElse(t *testing.T) {
	progString := `
    if x {
		x = x + 1;
    }
    `
	program, err := Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatal("expected an if statement")
	}
	if ifStmt.Cond.String() != "x" {
		t.Fatalf("got %s, expected x", ifStmt.Cond.String())
	}

	if len(ifStmt.Body.Statements) != 1 {
		t.Fatalf("got %d statements in body, expected 1", len(ifStmt.Body.Statements))
	}

	elseBlock := ifStmt.Else
	if elseBlock != nil {
		t.Fatal("expected nil else block")
	}
}

func TestParserComparison(t *testing.T) {

	program, err := Parse("x>=1;")
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected an expression statement, got %T", program.Statements[0])
	}
	binaryExpr, ok := exprStmt.Expression.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected a binary expression, got %T", exprStmt.Expression)
	}
	if binaryExpr.String() != "(x >= 1)" {
		t.Errorf("expected expression (x >= 1), got %s", binaryExpr)
	}

}

func TestParserComparisonMalformed(t *testing.T) {
	_, err := Parse("x> = 1;")

	if err == nil {
		t.Fatal("expected error")
	}
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError: %+v", parseErr)
	}

	if parseErr.Token.Column != 4 {
		t.Fatalf("expected parse error at column 4, got %d", parseErr.Token.Column)
	}
	if parseErr.Token.Line != 1 {
		t.Fatalf("expected parse error at line: 1, got %d", parseErr.Token.Line)
	}
	if parseErr.Token.Value != "=" {
		t.Fatalf("expected parse error for token value \"=\", got %q", parseErr.Token.Value)
	}

}

func TestParserEquality(t *testing.T) {
	expected := "(((x + 2) == 1) != true)"
	program, err := Parse("x + 2 == 1 != true;")
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements, expected 1", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected an expression statement, got %T", program.Statements[0])
	}
	binaryExpr, ok := exprStmt.Expression.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected a binary expression, got %T", exprStmt.Expression)
	}
	if binaryExpr.String() != expected {
		t.Errorf("expected expression %s, got %s", expected, binaryExpr)
	}
}

func TestParserMalformedEquality(t *testing.T) {
	_, err := Parse("x === 1;")
	if err == nil {
		t.Fatal("expected error")
	}
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError: %+v", parseErr)
	}

	if parseErr.Token.Column != 5 {
		t.Fatalf("expected parse error at column 5, got %d", parseErr.Token.Column)
	}
	if parseErr.Token.Line != 1 {
		t.Fatalf("expected parse error at line: 1, got %d", parseErr.Token.Line)
	}
	if parseErr.Token.Value != "=" {
		t.Fatalf("expected parse error for token value \"=\", got %q", parseErr.Token.Value)
	}

}

func TestParserMalformedInEquality(t *testing.T) {
	_, err := Parse("x !=== 1;")
	if err == nil {
		t.Fatal("expected error")
	}
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError: %+v", parseErr)
	}

	if parseErr.Token.Column != 5 {
		t.Fatalf("expected parse error at column 5, got %d", parseErr.Token.Column)
	}
	if parseErr.Token.Line != 1 {
		t.Fatalf("expected parse error at line: 1, got %d", parseErr.Token.Line)
	}
	if parseErr.Token.Value != "==" {
		t.Fatalf("expected parse error for token value \"==\", got %q", parseErr.Token.Value)
	}

}
