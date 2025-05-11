package ast

type Statement interface {
	stmt()
}

type Expression interface {
	expr()
}

type VarDecl struct {
	Shared bool
	Type   string
	Name   string
	Value  Expression
}

func (vd *VarDecl) stmt() {}

type Assignment struct {
	Name  string
	Value Expression
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
