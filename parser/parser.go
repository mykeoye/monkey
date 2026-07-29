package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	int = iota
	_
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// The parser for the program which will read tokens and build our AST
type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Advance the cursor twice so we fill both the cur and peek tokens
	p.advance()
	p.advance()

	// Register parsing functions for prefix and infix operations here
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixParseFn(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixParseFn(token.INT, p.parseIntegerLiteral)
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Appends peek errors to the parser's internal error slice
func (p *Parser) appendPeekError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected %s but got=%s", t, p.curToken.Type))
}

func (p *Parser) advance() {
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
		p.advance()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Function to parse an expression
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expr := &ast.ExpressionStatement{Token: p.curToken}
	expr.Expression = p.parseExpression(LOWEST)

	if p.isPeekTokenEq(token.SEMICOLON) {
		p.advance()
	}
	return expr
}

// Parses an expression following a defined precedence
func (p *Parser) parseExpression(precedence int32) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExpr := prefix()
	return leftExpr
}

// Parses tokens to form a let statemtent AST representation. A lets statement takes
// the form let <identifie> = <expression>
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
		p.advance()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.curToken}
	p.advance() // Advance the cursor to the next token

	// Skip parsing expressions for now
	for !p.expectPeekThenAdvance(token.SEMICOLON) {
		p.advance()
	}
	return returnStmt
}

// ----------------------------------------------------------------------- //
//
//	Parsing Functions for Expressions									   //
//
// ----------------------------------------------------------------------- //
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	exp := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, errMsg)
		return nil
	}
	exp.Value = value
	return exp
}

// ----------------------------------------------------------------------- //

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
		p.advance() // advance the parser to then next token
		return true
	}
	p.appendPeekError(expectedToken)
	return false
}

func (p *Parser) registerPrefixParseFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
