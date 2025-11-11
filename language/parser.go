package language

import (
	"deadlock/language/ast"
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) ParseProgram() []ast.Statement {
	var statements []ast.Statement
	for !p.isAtEnd() {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *Parser) parseStatement() ast.Statement {
	token := p.peek()
	switch token.Type {
	case KEYWORD:
		switch token.Value {
		case "shared":
			p.current++
			return p.parseVarDecl(token, true)
		default:
			p.error("Unexpected keyword")
		}
	default:
		p.error("Unexpected token")
	}
	return nil
}

func (p *Parser) parseVarDecl(token Token, shared bool) ast.Statement {
	var t []Token
	for token.Type != ENDSTMT {
		t = append(t, token)
		token = p.peek()
		p.current++
	}
	vdecl := ast.VarDecl{
		Shared: shared,
		Type:   string(t[1].Value),
		Name:   t[2].Value,
		Value:  p.parseExpression(),
	}

	return vdecl
}

func (p *Parser) error(s string) {
	fmt.Printf("Error at line %d: %s\n", p.peek().Line, s)
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.peek()
}

func (p *Parser) parseExpression() ast.Expression {
	// Placeholder for now: return literal expression if the next token is a number
	tok := p.advance()
	if tok.Type == NUMBER {
		//return &ast.Literal{Value: tok.Value}
	}
	p.error("Expected expression")
	return nil
}
