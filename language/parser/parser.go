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
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)

	return parser
}

func (p *Parser) ParseProgram() []ast.Statement {
	var statements []ast.Statement
	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt == nil {
			return nil
		}
		statements = append(statements, stmt)
	}
	return statements
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
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var currentToken = p.peekN(0)
	prefix := p.prefixParseFns[currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
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

func (p *Parser) error(s string) {
	fmt.Printf("Error at line %d: %s\n", p.peekN(0).Line, s)
}

func (p *Parser) peekN(n int) token.Token {
	return p.tokens[p.current+n]
}

func (p *Parser) peekNtokenIs(n int, t token.TokenType) bool {
	return p.peekN(n).Type == t
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
