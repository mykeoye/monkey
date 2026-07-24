package ast

import (
	"monkey/token"
)

// A node representation in our AST
type Node interface {
	TokenLiteral() string
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

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

func (ident *Identifier) expressionNode() {}

// A return statement node. Form is return <expression>
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (ret *ReturnStatement) statementNode() {}
func (ret *ReturnStatement) TokenLiteral() string {
	return ret.Token.Literal
}
