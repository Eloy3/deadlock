package ast

import (
	"bytes"
	"deadlock/language/token"
	"strings"
)

// safeExprString safely converts an Expression to string, returning "" if nil
func safeExprString(expr Expression) string {
	if expr == nil {
		return ""
	}
	return expr.String()
}

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	stmt()
}

type Expression interface {
	Node
	expr()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var s strings.Builder

	for _, stmt := range p.Statements {
		s.WriteString(stmt.String())
	}

	return s.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) TokenLiteral() string { return id.Token.Literal }
func (id *Identifier) expr()                {}
func (id *Identifier) String() string       { return id.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expr()                {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FunctionCall struct {
	Name string
	Args []Expression
}

func (fc FunctionCall) expr() {}

type VariableDeclaration struct {
	Token  token.Token
	Type   string
	Shared bool
	Name   *Identifier
	Value  Expression
}

// String implements Statement.
// let shared int x = 1
func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString("let ")
	if vd.Shared {
		out.WriteString("shared ")
	}
	out.WriteString(vd.Type)
	if vd.Name != nil {
		out.WriteString(vd.Name.String())
	}
	out.WriteString(" = ")
	out.WriteString(safeExprString(vd.Value))

	return out.String()
}

func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableDeclaration) stmt()                {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expr()                {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(safeExprString(pe.Right))
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expr()                {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(safeExprString(oe.Left))
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(safeExprString(oe.Right))
	out.WriteString(")")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) stmt()                {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expr()                {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type Assignment struct {
	Name   string
	Value  Expression
	Shared bool
}

type LockStmt struct {
	LockName string
}

type PrintStmt struct {
	Value Expression
}

type ThreadDecl struct {
	Name string
	Body []Statement
}
