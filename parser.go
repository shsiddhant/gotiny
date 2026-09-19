package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer   *Lexer
	current Token
}

func (p *Parser) advance() error {
	token, err := p.lexer.Next()
	if err != nil {
		return err
	}

	p.current = token
	return nil
}

func NewParser(lexer *Lexer) (*Parser, error) {
	p := &Parser{lexer: lexer}

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

		expr := &LiteralExpr{Value: value}

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

func (p *Parser) Parse() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.current.Type != EOF {
		return nil, fmt.Errorf("unexpected token: %s", p.current)
	}

	return &ExprStmt{Expression: expr}, nil
}

func Parse(source string) (Stmt, error) {
	lexer := NewLexer(source)

	parser, err := NewParser(lexer)

	if err != nil {
		return nil, err
	}
	return parser.Parse()
}
