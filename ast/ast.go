package ast

import (
	"bytes"
	"strings"

	"pidgin-lang/token"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string // Returns the literal value of the token (for debugging)
	String() string       // Returns a string representation of the node
}

// Statement represents a statement node (doesn't produce a value)
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node (produces a value)
type Expression interface {
	Node
	expressionNode()
}

// =============================================================================
// Program - The root node of every AST
// =============================================================================

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// =============================================================================
// Statements
// =============================================================================

// MakeStatement represents: make x be 5
type MakeStatement struct {
	Token token.Token // the 'make' token
	Name  *Identifier // variable name
	Value Expression  // the value being assigned
}

func (ms *MakeStatement) statementNode()       {}
func (ms *MakeStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *MakeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ms.TokenLiteral() + " ")
	out.WriteString(ms.Name.String())
	out.WriteString(" be ")
	if ms.Value != nil {
		out.WriteString(ms.Value.String())
	}
	return out.String()
}

// BringStatement represents: bring x (return statement)
type BringStatement struct {
	Token       token.Token // the 'bring' token
	ReturnValue Expression
}

func (bs *BringStatement) statementNode()       {}
func (bs *BringStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BringStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral() + " ")
	if bs.ReturnValue != nil {
		out.WriteString(bs.ReturnValue.String())
	}
	return out.String()
}

// ExpressionStatement represents a statement consisting of a single expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block of statements: { ... }
type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// =============================================================================
// Expressions
// =============================================================================

// Identifier represents a variable name
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer: 5, 42, 1000
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral represents a string: "How far!"
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// Boolean represents: tru or lie
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// NothingLiteral represents: nothing (null)
type NothingLiteral struct {
	Token token.Token
}

func (nl *NothingLiteral) expressionNode()      {}
func (nl *NothingLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NothingLiteral) String() string       { return "nothing" }

// PrefixExpression represents: -5, !tru, no be x
type PrefixExpression struct {
	Token    token.Token // the prefix token (!, -, no)
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents: 5 + 5, x big pass y, a and b
type InfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// SupposeExpression represents: suppose x big pass 5 { ... } abi { ... }
type SupposeExpression struct {
	Token       token.Token // the 'suppose' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // optional 'abi' block
}

func (se *SupposeExpression) expressionNode()      {}
func (se *SupposeExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SupposeExpression) String() string {
	var out bytes.Buffer
	out.WriteString("suppose ")
	out.WriteString(se.Condition.String())
	out.WriteString(" ")
	out.WriteString(se.Consequence.String())
	if se.Alternative != nil {
		out.WriteString(" abi ")
		out.WriteString(se.Alternative.String())
	}
	return out.String()
}

// WhileExpression represents: dey do while x no reach 10 { ... }
type WhileExpression struct {
	Token     token.Token // the 'dey' token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("dey do while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}

// DoExpression represents function definition: do add(a, b) { bring a + b }
type DoExpression struct {
	Token      token.Token   // the 'do' token
	Name       *Identifier   // function name (optional for anonymous functions)
	Parameters []*Identifier // function parameters
	Body       *BlockStatement
}

func (de *DoExpression) expressionNode()      {}
func (de *DoExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DoExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range de.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("do ")
	if de.Name != nil {
		out.WriteString(de.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(de.Body.String())
	return out.String()
}

// CallExpression represents function call: add(1, 2) or yarn("hello")
type CallExpression struct {
	Token     token.Token  // the '(' token
	Function  Expression   // Identifier or DoExpression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}