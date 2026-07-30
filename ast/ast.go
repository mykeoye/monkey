package ast

import (
	"bytes"
	"monkey/token"
)

// A node representation in our AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statements as nodes which do not return values
type Statement interface {
	Node
	statementNode()
}

// Expressions as nodes that evaluate/return a value
type Expression interface {
	Node
	expressionNode()
}

// This struct will be the root node in our AST
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
	var outBuf bytes.Buffer
	for _, stmt := range p.Statements {
		outBuf.WriteString(stmt.String())
	}
	return outBuf.String()
}

// The let statement node. Form is let <identifier> = <expression>
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LetStatement) statementNode() {
}
func (l *LetStatement) String() string {
	var outBuf bytes.Buffer
	outBuf.WriteString(l.TokenLiteral() + " ")
	outBuf.WriteString(l.Name.String())
	outBuf.WriteString(" = ")

	if l.Value != nil {
		outBuf.WriteString(l.Value.String())
	}
	outBuf.WriteString(";")
	return outBuf.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}
func (ident *Identifier) expressionNode() {}
func (ident *Identifier) String() string {
	return ident.Value
}

// A return statement node. Form is return <expression>
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (ret *ReturnStatement) statementNode() {}
func (ret *ReturnStatement) TokenLiteral() string {
	return ret.Token.Literal
}
func (ret *ReturnStatement) String() string {
	var outBuf bytes.Buffer
	outBuf.WriteString(ret.TokenLiteral() + " ")
	if ret.ReturnValue != nil {
		outBuf.WriteString(ret.ReturnValue.String())
	}
	outBuf.WriteString(";")
	return outBuf.String()
}

// An expression statement, abstracting an expression which contains a token and
// one expression
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (exp *ExpressionStatement) statementNode() {}
func (exp *ExpressionStatement) TokenLiteral() string {
	return exp.Token.Literal
}
func (exp *ExpressionStatement) String() string {
	if exp.Expression != nil {
		return exp.Expression.String()
	}
	return ""
}
