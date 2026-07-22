package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// The parser for the program which will read tokens and build our AST
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Advance the cursor twice so we fill both the cur and peek tokens
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Appends peek errors to the parser's internal error slice
func (p *Parser) appendPeekError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected %s but got=%s", t, p.curToken.Type))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Loop until we get to the EOF token, parsing tokens into parts of our AST
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// Parses tokens to form a let statemtent AST representation. A lets statement takes
// the form let <identifie> =  <expression>
func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{Token: p.curToken}

	// Peek ahead to the next token, it should be an indentifier
	if !p.expectPeekThenAdvance(token.IDENTIFIER) {
		return nil
	}

	// The parser has advanced and the curToken is now an identifier
	letStatement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Next we expect an assignemt operator as one other compoent of the let
	// statement.
	if !p.expectPeekThenAdvance(token.ASSIGNMENT) {
		return nil
	}

	// For now we don't parse expressions so we skip them until we get to a semi-colon
	for !p.isCurrentTokenEq(token.SEMICOLON) {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) isCurrentTokenEq(expectedTokenType token.TokenType) bool {
	return p.curToken.Type == expectedTokenType
}

func (p *Parser) isPeekTokenEq(expectedTokenType token.TokenType) bool {
	return p.peekToken.Type == expectedTokenType
}

// This function checks to see if the next token is the token we expect
// if it is, we advance the cursor and return true, or false otherwise
func (p *Parser) expectPeekThenAdvance(expectedToken token.TokenType) bool {
	if p.isPeekTokenEq(expectedToken) {
		p.nextToken() // advance the parser to then next token
		return true
	}
	p.appendPeekError(expectedToken)
	return false
}
