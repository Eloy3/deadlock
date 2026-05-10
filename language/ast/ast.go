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
	Token token.Token
	Local bool
	Name  *Identifier
	Value Expression
}

func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer
	if vd.Local {
		out.WriteString("local ")
	}
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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expr()                {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) stmt()                {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type SharedBlock struct {
	Token        token.Token
	Declarations []*VariableDeclaration
}

func (sb *SharedBlock) expr()                {}
func (sb *SharedBlock) TokenLiteral() string { return sb.Token.Literal }
func (sb *SharedBlock) String() string {
	var s strings.Builder
	s.WriteString("shared {\n")
	for _, decl := range sb.Declarations {
		s.WriteString("\t")
		s.WriteString(decl.String())
		s.WriteString("\n")
	}
	s.WriteString("}")
	return s.String()
}

type MutexStatement struct {
	Token token.Token
	Name  *Identifier
}

func (md *MutexStatement) stmt()                {}
func (md *MutexStatement) TokenLiteral() string { return md.Token.Literal }
func (md *MutexStatement) String() string {
	var out bytes.Buffer
	out.WriteString("mutex ")
	if md.Name != nil {
		out.WriteString(md.Name.String())
	}
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // for else
}

func (is *IfExpression) expr()                {}
func (is *IfExpression) TokenLiteral() string { return is.Token.Literal }
func (is *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(safeExprString(is.Condition))
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

type AssignmentExpression struct {
	Token token.Token // The = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expr() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(ae.TokenLiteral())
	out.WriteString(ae.Value.String())

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expr() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type LockStatement struct {
	Token    token.Token
	Argument *Identifier
}

func (ls *LockStatement) stmt()                {}
func (ls *LockStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LockStatement) String() string       { return "lock(" + ls.Argument.Value + ")" }

type UnlockStatement struct {
	Token    token.Token
	Argument *Identifier
}

func (us *UnlockStatement) stmt()                {}
func (us *UnlockStatement) TokenLiteral() string { return us.Token.Literal }
func (us *UnlockStatement) String() string       { return "unlock(" + us.Argument.Value + ")" }

type PrintStmt struct {
	Token token.Token
	Value Expression
}

func (ps *PrintStmt) stmt()                {}
func (ps *PrintStmt) TokenLiteral() string { return ps.Token.Literal }
func (ps *PrintStmt) String() string       { return "print(" + safeExprString(ps.Value) + ")" }

type ThreadExpression struct {
	Token token.Token
	Name  *Identifier
	Body  *BlockStatement
}

func (td *ThreadExpression) expr()                {}
func (td *ThreadExpression) TokenLiteral() string { return td.Token.Literal }
func (td *ThreadExpression) String() string {
	var s strings.Builder
	s.WriteString("thread ")
	s.WriteString(td.Name.Value)
	s.WriteString(" {\n")
	s.WriteString(td.Body.String())
	s.WriteString("}")
	return s.String()
}
