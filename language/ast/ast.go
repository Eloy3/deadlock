package ast

import (
	"bytes"
	"deadlock/language/token"
)

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
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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
// shared int x = 1
func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer

	if vd.Shared {
		out.WriteString("shared ")
	}
	out.WriteString(vd.Type)
	out.WriteString(vd.Name.String())
	out.WriteString(" = ")
	out.WriteString(vd.Value.String())

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
	out.WriteString(pe.Right.String())
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
