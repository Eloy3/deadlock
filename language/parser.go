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
	token := p.peekN(0)
	switch token.Type {
	case KEYWORD:
		switch token.Value {
		case "shared":
			token = p.advance()
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
		token = p.advance()
	}
	vdecl := ast.VarDecl{
		Shared: shared,
		Type:   string(t[1].Value),
		Name:   t[2].Value,
		Value:  p.parseExpression(p.peekN(-1)), // parse element before ENDSTMT (;)
	}

	return vdecl
}

func (p *Parser) error(s string) {
	fmt.Printf("Error at line %d: %s\n", p.peekN(0).Line, s)
}

func (p *Parser) peekN(n int) Token {
	return p.tokens[p.current+n]
}

func (p *Parser) isAtEnd() bool {
	return p.peekN(0).Type == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.peekN(0)
}

func (p *Parser) parseExpression(token Token) ast.Expression {
	return p.parseBinaryOp(token, 0)
}

func (p *Parser) parsePrimaryExpression(token Token) ast.Expression {
	switch token.Type {
	case NUMBER:
		return ast.NumberLiteral{Value: token.Value}
	case IDENTIFIER:
		// Check if it's a function call
		if p.peekN(0).Type == SYMBOL && p.peekN(0).Value == "(" {
			return p.parseFunctionCall(token)
		}
		return ast.Identifier{Name: token.Value}
	default:
		p.error(fmt.Sprintf("Unexpected token in expression: %s", token.Value))
		return nil
	}
}

func (p *Parser) parseFunctionCall(token Token) ast.Expression {
	name := token.Value
	p.advance() // skip '('

	var args []ast.Expression
	for p.peekN(0).Type != SYMBOL || p.peekN(0).Value != ")" {
		args = append(args, p.parseExpression(p.peekN(0)))
		if p.peekN(0).Type == SYMBOL && p.peekN(0).Value == "," {
			p.advance() // skip ','
		}
	}
	p.advance() // skip ')'

	return ast.FunctionCall{Name: name, Args: args}
}

func (p *Parser) parseBinaryOp(token Token, minPrec int) ast.Expression {
	left := p.parsePrimaryExpression(token)

	for {
		op := p.peekN(0)
		if op.Type != SYMBOL || !isOperator(op.Value) || getPrecedence(op.Value) < minPrec {
			break
		}

		opStr := op.Value
		p.advance() // consume operator

		right := p.parseBinaryOp(p.peekN(0), getPrecedence(opStr)+1)
		left = ast.BinaryOp{Left: left, Operator: opStr, Right: right}
	}

	return left
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "="
}

func getPrecedence(op string) int {
	switch op {
	case "=":
		return 1
	case "+", "-":
		return 2
	case "*", "/":
		return 3
	default:
		return 0
	}
}
