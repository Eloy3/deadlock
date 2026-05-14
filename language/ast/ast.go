package ast

import (
	"bytes"
	"deadlock/language/token"
	"fmt"
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

// PrintTree prints the AST as a tree structure for debugging
func (p *Program) PrintTree() {
	fmt.Println("Program")
	for i, stmt := range p.Statements {
		isLast := i == len(p.Statements)-1
		printNodeTree(stmt, "", isLast)
	}
}

func printNodeTree(node Node, indent string, isLast bool) {
	var prefix string
	if isLast {
		prefix = indent + "└── "
		indent = indent + "    "
	} else {
		prefix = indent + "├── "
		indent = indent + "│   "
	}

	switch n := node.(type) {
	case *VariableDeclaration:
		fmt.Printf("%s%T (name=%s)\n", prefix, n, n.Name.Value)
		if n.Value != nil {
			printExprTree(n.Value, indent, true)
		}

	case *MutexStatement:
		fmt.Printf("%s%T (name=%s)\n", prefix, n, n.Name.Value)

	case *ExpressionStatement:
		fmt.Printf("%s%T\n", prefix, n)
		if n.Expression != nil {
			printExprTree(n.Expression, indent, true)
		}

	case *BlockStatement:
		fmt.Printf("%sBlockStatement (%d statements)\n", prefix, len(n.Statements))
		for i, stmt := range n.Statements {
			printNodeTree(stmt, indent, i == len(n.Statements)-1)
		}

	default:
		fmt.Printf("%s%T\n", prefix, n)
	}
}

func printExprTree(expr Expression, indent string, isLast bool) {
	var prefix string
	if isLast {
		prefix = indent + "└── "
		indent = indent + "    "
	} else {
		prefix = indent + "├── "
		indent = indent + "│   "
	}

	switch e := expr.(type) {
	case *Identifier:
		fmt.Printf("%sIdentifier (value=%s)\n", prefix, e.Value)

	case *IntegerLiteral:
		fmt.Printf("%sIntegerLiteral (value=%d)\n", prefix, e.Value)

	case *StringLiteral:
		fmt.Printf("%sStringLiteral (value=%s)\n", prefix, e.Value)

	case *Boolean:
		fmt.Printf("%sBoolean (value=%v)\n", prefix, e.Value)

	case *PrefixExpression:
		fmt.Printf("%sPrefixExpression (op=%s)\n", prefix, e.Operator)
		if e.Right != nil {
			printExprTree(e.Right, indent, true)
		}

	case *InfixExpression:
		fmt.Printf("%sInfixExpression (op=%s)\n", prefix, e.Operator)
		if e.Left != nil {
			printExprTree(e.Left, indent, false)
		}
		if e.Right != nil {
			printExprTree(e.Right, indent, true)
		}

	case *IfExpression:
		fmt.Printf("%sIfExpression\n", prefix)
		if e.Condition != nil {
			fmt.Printf("%s├── Condition:\n", indent)
			printExprTree(e.Condition, indent+"│   ", false)
		}
		if e.Consequence != nil {
			fmt.Printf("%s├── Consequence:\n", indent)
			printNodeTree(e.Consequence, indent+"│   ", e.Alternative == nil)
		}
		if e.Alternative != nil {
			fmt.Printf("%s└── Alternative:\n", indent)
			printNodeTree(e.Alternative, indent+"    ", true)
		}

	case *AssignmentExpression:
		fmt.Printf("%sAssignmentExpression\n", prefix)
		if e.Left != nil {
			printExprTree(e.Left, indent, false)
		}
		if e.Value != nil {
			printExprTree(e.Value, indent, true)
		}

	case *IndexExpression:
		fmt.Printf("%sIndexExpression\n", prefix)
		if e.Left != nil {
			printExprTree(e.Left, indent, false)
		}
		if e.Index != nil {
			printExprTree(e.Index, indent, true)
		}

	case *SharedBlock:
		fmt.Printf("%sSharedBlock (%d declarations)\n", prefix, len(e.Declarations))
		for i, decl := range e.Declarations {
			printNodeTree(decl, indent, i == len(e.Declarations)-1)
		}

	case *ThreadExpression:
		fmt.Printf("%sThreadExpression (name=%s)\n", prefix, e.Name.Value)
		if e.Body != nil {
			printNodeTree(e.Body, indent, true)
		}

	default:
		fmt.Printf("%s%T\n", prefix, e)
	}
}
