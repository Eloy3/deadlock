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
	switch p.peek().Type {
	case KEYWORD:
		switch p.peek().Value {
		case "shared":
			return p.parseVarDecl()
		default:
			p.error("Unexpected keyword")
		}
	default:
		p.error("Unexpected token")
	}
	return nil
}

func (p *Parser) parseVarDecl() ast.Statement {
	panic("unimplemented")
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
