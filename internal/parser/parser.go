package parser

import (
	"fmt"
	"strconv"

	"github.com/shsiddhant/gotiny/internal/ast"
	"github.com/shsiddhant/gotiny/internal/lexer"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/token"
)

type Parser struct {
	lexer   *lexer.Lexer
	current token.Token
	peek    token.Token
}

type ParseError struct {
	Token   token.Token
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s",
		e.Token.Line,
		e.Token.Column,
		e.Message,
	)
}

func (p *Parser) advance() error {
	p.current = p.peek

	token, err := p.lexer.Next()
	if err != nil {
		return &ParseError{
			Token:   token,
			Message: err.Error(),
		}
	}

	p.peek = token
	return nil
}

func NewParser(lexer *lexer.Lexer) (*Parser, error) {
	p := &Parser{lexer: lexer}

	// advance twice to initialize both current and peek.
	if err := p.advance(); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) consume(kind token.TokenType) error {
	if p.current.Type != kind {
		return &ParseError{
			Token:   p.current,
			Message: fmt.Sprintf("expected %s, got %s", kind, p.current.Type),
		}
	}
	return p.advance()
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.booleanOr()
}

func (p *Parser) primary() (ast.Expr, error) {
	switch p.current.Type {
	case token.Number:
		value, err := strconv.Atoi(p.current.Value)
		if err != nil {
			return nil, err
		}

		expr := &ast.LiteralExpr{Token: p.current, Value: objects.Int(value)}

		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case token.True:
		expr := &ast.LiteralExpr{Token: p.current, Value: objects.Bool(true)}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case token.False:
		expr := &ast.LiteralExpr{Token: p.current, Value: objects.Bool(false)}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case token.Identifier:
		if p.peek.Type == token.LeftParen {
			expr, err := p.callExpression()
			if err != nil {
				return nil, err
			}
			return expr, nil
		}
		expr := &ast.VariableExpr{Name: p.current}

		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case token.LeftParen:
		if err := p.advance(); err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if err := p.consume(token.RightParen); err != nil {
			return nil, err
		}

		return &ast.GroupExpr{Expression: expr}, nil
	}

	return nil, &ParseError{
		Token:   p.current,
		Message: fmt.Sprintf("expected expression, got %s", p.current.Type),
	}
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.current.Type == token.Plus || p.current.Type == token.Minus || p.current.Type == token.Not {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		operand, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{Operator: operator, Operand: operand}, nil
	}
	return p.primary()
}

func (p *Parser) multiplication() (ast.Expr, error) {
	expr, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.current.Type == token.Star || p.current.Type == token.Slash {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) addition() (ast.Expr, error) {
	expr, err := p.multiplication()

	if err != nil {
		return nil, err
	}
	for p.current.Type == token.Plus || p.current.Type == token.Minus {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.multiplication()

		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}
	for p.current.Type == token.Less ||
		p.current.Type == token.LessEqual ||
		p.current.Type == token.Greater ||
		p.current.Type == token.GreaterEqual {

		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.addition()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.current.Type == token.EqualEqual || p.current.Type == token.NotEqual {

		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) booleanAnd() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.current.Type == token.And {

		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) booleanOr() (ast.Expr, error) {
	expr, err := p.booleanAnd()
	if err != nil {
		return nil, err
	}

	for p.current.Type == token.Or {

		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.booleanAnd()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) callArgs() ([]ast.Expr, error) {
	arg, err := p.expression()
	if err != nil {
		return nil, err
	}

	args := []ast.Expr{arg}

	for p.current.Type == token.Comma {
		if err := p.advance(); err != nil {
			return nil, err
		}
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

func (p *Parser) callExpression() (ast.Expr, error) {
	name := p.current
	if err := p.advance(); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil {
		return nil, err
	}

	var args []ast.Expr
	var err error

	if p.current.Type != token.RightParen {
		args, err = p.callArgs()
		if err != nil {
			return nil, err
		}
	}
	if err := p.consume(token.RightParen); err != nil {
		return nil, err
	}

	return &ast.CallExpr{Name: name, Args: args}, nil
}

func (p *Parser) exprStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &ast.ExprStmt{Expression: expr}, nil
}

func (p *Parser) assignStatement() (ast.Stmt, error) {
	name := p.current
	if err := p.consume(token.Identifier); err != nil {
		return nil, err
	}
	if err := p.consume(token.Equal); err != nil {
		return nil, err
	}
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &ast.AssignStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) letStatement() (ast.Stmt, error) {
	if err := p.consume(token.Let); err != nil {
		return nil, err
	}
	name := p.current
	if err := p.consume(token.Identifier); err != nil {
		return nil, err
	}
	if err := p.consume(token.Equal); err != nil {
		return nil, err
	}
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &ast.LetStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) block() (*ast.BlockStmt, error) {
	if err := p.consume(token.LeftCurlyBrace); err != nil {
		return nil, err
	}
	block := &ast.BlockStmt{}

	for p.current.Type != token.RightCurlyBrace {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
	}
	if err := p.consume(token.RightCurlyBrace); err != nil {
		return nil, err
	}
	return block, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	if err := p.consume(token.If); err != nil {
		return nil, err
	}
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	var elseblock *ast.BlockStmt
	if p.current.Type == token.Else {
		if err := p.advance(); err != nil {
			return nil, err
		}
		elseblock, err = p.block()
		if err != nil {
			return nil, err
		}
	}
	return &ast.IfStmt{
		Cond: cond,
		Body: body,
		Else: elseblock,
	}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	if err := p.consume(token.Return); err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{Expr: expr}, nil
}

func (p *Parser) fnParam() (ast.Parameter, error) {
	name := p.current
	if err := p.consume(token.Identifier); err != nil {
		return ast.Parameter{}, &ParseError{
			Token:   p.current,
			Message: err.Error(),
		}
	}
	switch p.current.Type {
	case token.Int:
		param := ast.Parameter{Name: name, Type: objects.IntType}
		if err := p.advance(); err != nil {
			return ast.Parameter{}, &ParseError{
				Token:   p.current,
				Message: err.Error(),
			}
		}
		return param, nil

	case token.Bool:
		param := ast.Parameter{Name: name, Type: objects.BoolType}
		if err := p.advance(); err != nil {
			return ast.Parameter{}, &ParseError{
				Token:   p.current,
				Message: err.Error(),
			}
		}
		return param, nil
	}
	return ast.Parameter{}, &ParseError{
		Token:   p.current,
		Message: fmt.Sprintf("expected parameter type, got %s", p.current.Type),
	}
}

func (p *Parser) fnParams() ([]ast.Parameter, error) {
	param, err := p.fnParam()
	if err != nil {
		return nil, err
	}
	params := []ast.Parameter{param}

	for p.current.Type == token.Comma {
		if err := p.advance(); err != nil {
			return nil, err
		}
		param, err := p.fnParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

func (p *Parser) fnDeclareStatement() (ast.Stmt, error) {
	if err := p.consume(token.Fn); err != nil {
		return nil, err
	}

	name := p.current

	if err := p.consume(token.Identifier); err != nil {
		return nil, err
	}
	if err := p.consume(token.LeftParen); err != nil {
		return nil, err
	}

	var params []ast.Parameter
	var returnType objects.Type
	var err error

	if p.current.Type != token.RightParen {
		params, err = p.fnParams()
		if err != nil {
			return nil, err
		}
	}
	if err := p.consume(token.RightParen); err != nil {
		return nil, err
	}

	switch p.current.Type {
	case token.Int:
		returnType = objects.IntType
		if err := p.advance(); err != nil {
			return nil, err
		}
	case token.Bool:
		returnType = objects.BoolType
		if err := p.advance(); err != nil {
			return nil, err
		}
	default:
		returnType = objects.VoidType
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FnDeclareStmt{Name: name, Parameters: params, ReturnType: returnType, Body: body}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	var stmt ast.Stmt
	var err error

	switch {
	case p.current.Type == token.Let:
		stmt, err = p.letStatement()
	case p.current.Type == token.If:
		return p.ifStatement()
	case p.current.Type == token.Identifier && p.peek.Type == token.Equal:
		stmt, err = p.assignStatement()
	case p.current.Type == token.Fn:
		return p.fnDeclareStatement()
	case p.current.Type == token.Return:
		stmt, err = p.returnStatement()
	default:
		stmt, err = p.exprStatement()
	}
	if err != nil {
		return nil, err
	}
	if err := p.consume(token.SemiColon); err != nil {
		return nil, err
	}

	return stmt, err
}

func (p *Parser) Parse() (*ast.Program, error) {

	program := &ast.Program{}

	for p.current.Type != token.EOF {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}

	return program, nil
}

func Parse(source string) (*ast.Program, error) {
	newLexer := lexer.NewLexer(source)

	parser, err := NewParser(newLexer)

	if err != nil {
		return nil, err
	}
	return parser.Parse()
}
