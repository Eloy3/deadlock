package parser

import (
	"deadlock/language/ast"
	"deadlock/language/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	NOT
	IN
	ASSIGN       // := or =
	EQUALS       // ==
	LESSGREATER  // > or <
	BitwiseOR    // |
	BitwiseXOR   // ^
	BitwiseAND   // &
	BitwiseShift // << or >>
	SUM          // + or -
	PRODUCT      // * / or %
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // array[index]
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.OR:              OR,
	token.AND:             AND,
	token.NOT:             NOT,
	token.ASSIGN:          ASSIGN,
	token.EQUALS:          EQUALS,
	token.NOT_EQUALS:      EQUALS,
	token.GREATER:         LESSGREATER,
	token.GREATER_EQ:      LESSGREATER,
	token.LESS_GREATER:    LESSGREATER,
	token.LESS_GREATER_EQ: LESSGREATER,
	//token.BitwiseOR:       BitwiseOR,
	//token.BitwiseXOR:      BitwiseXOR,
	//token.BitwiseAND:      BitwiseAND,
	//token.LeftShift:       BitwiseShift,
	//token.RightShift:      BitwiseShift,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIVIDE:   PRODUCT,
	token.MULTIPLY: PRODUCT,
	token.MODULE:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.DOT:      INDEX,
}

type Parser struct {
	tokens  []token.Token
	errors  []string
	current int

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(tokens []token.Token) *Parser {
	parser := &Parser{
		tokens:  tokens,
		current: 0,
		errors:  []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.NUMBER, parser.parseIntegerLiteral)
	parser.registerPrefix(token.NOT, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.DIVIDE, parser.parseInfixExpression)
	parser.registerInfix(token.MULTIPLY, parser.parseInfixExpression)
	parser.registerInfix(token.EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.LESS_GREATER, parser.parseInfixExpression)
	parser.registerInfix(token.GREATER, parser.parseInfixExpression)

	return parser
}

func (p *Parser) ParseProgram() ast.Program {
	var statements []ast.Statement
	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.advance()
	}
	return ast.Program{Statements: statements}
}

func (p *Parser) parseStatement() ast.Statement {
	tok := p.peekN(0)
	switch tok.Type {
	case token.SHARED:
		tok = p.advance()
		return p.parseVarDecl(true)
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.peekN(0)}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekNtokenIs(1, token.SEMICOLON) {
		p.advance()
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var curToken = p.peekN(0)
	prefix := p.prefixParseFns[curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekNtokenIs(1, token.SEMICOLON) && precedence < p.peekPrecedence(1) {
		infix := p.infixParseFns[p.peekN(1).Type]
		if infix == nil {
			return leftExp
		}

		p.advance()

		leftExp = infix(leftExp)
	}

	p.advance()
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseVarDecl(shared bool) *ast.VariableDeclaration {
	stmt := &ast.VariableDeclaration{Shared: shared}

	if !p.peekNtokenIs(1, token.IDENTIFIER) {
		return nil
	}

	tok := p.peekN(0)
	stmt.Name = &ast.Identifier{Token: tok, Value: tok.Literal}

	if !p.peekNtokenIs(1, token.EQUALS) {
		return nil
	}

	for !p.peekNtokenIs(0, token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	var identifier ast.Identifier
	identifier.Token = p.peekN(0)
	identifier.Value = p.peekN(0).Literal
	return &identifier
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	curToken := p.peekN(0)

	lit := &ast.IntegerLiteral{Token: curToken}

	value, err := strconv.ParseInt(curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	curToken := p.peekN(0)
	expression := &ast.PrefixExpression{
		Token:    curToken,
		Operator: curToken.Literal,
	}
	p.advance()

	right := p.parseExpression(PREFIX)
	if right == nil {
		msg := fmt.Sprintf("could not parse expression after %q", curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	expression.Right = right

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	token := p.peekN(0)

	expression := &ast.InfixExpression{
		Token:    token,
		Operator: token.Literal,
		Left:     left,
	}

	precedence := p.peekPrecedence(1)
	p.advance()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()

	exp := p.parseExpression(LOWEST)

	if !p.peekNtokenIs(1, token.RPAREN) {
		return nil
	}

	p.advance()
	return exp
}

func (p *Parser) error(s string) {
	fmt.Printf("Error at line %d: %s\n", p.peekN(0).Line, s)
}

func (p *Parser) peekN(n int) token.Token {
	index := p.current + n
	if index >= len(p.tokens) {
		return token.Token{Type: token.EOF, Literal: "EOF", Line: 0}
	}
	return p.tokens[index]
}

func (p *Parser) peekNtokenIs(n int, t token.TokenType) bool {
	return p.peekN(p.current+n).Type == t
}

func (p *Parser) peekPrecedence(n int) int {
	if p, ok := precedences[p.peekN(n).Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) isAtEnd() bool {
	return p.peekN(0).Type == token.EOF
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.peekN(0)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}
