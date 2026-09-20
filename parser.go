package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
}

func (p *Parser) advance() error {
	p.current = p.peek

	token, err := p.lexer.Next()
	if err != nil {
		return err
	}

	p.peek = token
	return nil
}

func NewParser(lexer *Lexer) (*Parser, error) {
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

func (p *Parser) consume(kind TokenType) error {
	if p.current.Type != kind {
		return fmt.Errorf("expected %s, got %s", kind, p.current.Type)
	}
	return p.advance()
}

func (p *Parser) expression() (Expr, error) {
	return p.addition()
}

func (p *Parser) primary() (Expr, error) {
	switch p.current.Type {
	case Number:
		value, err := strconv.Atoi(p.current.Value)
		if err != nil {
			return nil, err
		}

		expr := &LiteralExpr{Value: Int(value)}

		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case True:
		expr := &LiteralExpr{Value: Bool(true)}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case False:
		expr := &LiteralExpr{Value: Bool(false)}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case Identifier:
		expr := &VariableExpr{Name: p.current}

		if err := p.advance(); err != nil {
			return nil, err
		}
		return expr, nil

	case LeftParen:
		if err := p.advance(); err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if err := p.consume(RightParen); err != nil {
			return nil, err
		}

		return &GroupExpr{Expression: expr}, nil
	}

	return nil, fmt.Errorf("expected expression, got %s", p.current.Type)
}

func (p *Parser) unary() (Expr, error) {
	if p.current.Type == Plus || p.current.Type == Minus {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		operand, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &UnaryExpr{Operator: operator, Operand: operand}, nil
	}
	return p.primary()
}

func (p *Parser) multiplication() (Expr, error) {
	expr, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.current.Type == Star || p.current.Type == Slash {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) addition() (Expr, error) {
	expr, err := p.multiplication()

	if err != nil {
		return nil, err
	}
	for p.current.Type == Plus || p.current.Type == Minus {
		operator := p.current

		if err := p.advance(); err != nil {
			return nil, err
		}

		right, err := p.multiplication()

		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) exprStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &ExprStmt{Expression: expr}, nil
}

func (p *Parser) assignStatement() (Stmt, error) {
	name := p.current.Value
	if err := p.consume(Identifier); err != nil {
		return nil, err
	}
	if err := p.consume(Equal); err != nil {
		return nil, err
	}
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &AssignStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) letStatement() (Stmt, error) {
	if err := p.consume(Let); err != nil {
		return nil, err
	}
	name := p.current.Value
	if err := p.consume(Identifier); err != nil {
		return nil, err
	}
	if err := p.consume(Equal); err != nil {
		return nil, err
	}
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &LetStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) block() (*BlockStmt, error) {
	if err := p.consume(LeftCurlyBrace); err != nil {
		return nil, err
	}
	block := &BlockStmt{}

	for p.current.Type != RightCurlyBrace {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
	}
	if err := p.consume(RightCurlyBrace); err != nil {
		return nil, err
	}
	return block, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	if err := p.consume(If); err != nil {
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
	var elseblock *BlockStmt
	if p.current.Type == Else {
		if err := p.advance(); err != nil {
			return nil, err
		}
		elseblock, err = p.block()
		if err != nil {
			return nil, err
		}
	}
	return &IfStmt{
		Cond: cond,
		Body: body,
		Else: elseblock,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	var stmt Stmt
	var err error

	switch {
	case p.current.Type == Let:
		stmt, err = p.letStatement()
	case p.current.Type == If:
		return p.ifStatement()
	case p.current.Type == Identifier && p.peek.Type == Equal:
		stmt, err = p.assignStatement()
	default:
		stmt, err = p.exprStatement()
	}
	if err != nil {
		return nil, err
	}
	if err := p.consume(SemiColon); err != nil {
		return nil, err
	}

	return stmt, err
}

func (p *Parser) Parse() (*Program, error) {

	program := &Program{}

	for p.current.Type != EOF {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}

	return program, nil
}

func Parse(source string) (*Program, error) {
	lexer := NewLexer(source)

	parser, err := NewParser(lexer)

	if err != nil {
		return nil, err
	}
	return parser.Parse()
}
