package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	stmt()
}

type Expression interface {
	expr()
}

type NumberLiteral struct {
	Value string
}

func (nl NumberLiteral) expr() {}

type Identifier struct {
	Name string
}

func (id Identifier) expr() {}

type BinaryOp struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (bo BinaryOp) expr() {}

type FunctionCall struct {
	Name string
	Args []Expression
}

func (fc FunctionCall) expr() {}

type VarDecl struct {
	Shared bool
	Type   string
	Name   string
	Value  Expression
}

func (vd VarDecl) TokenLiteral() {}
func (vd VarDecl) stmt()         {}

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
