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

// Stores a map of each token type to its precedence (order of evaluation)
var PRECEDENCE_TABLE = map[token.TokenType]uint64{
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
}

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
	p.registerPrefixParseFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixParseFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParseFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixParseFn(token.EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfixParseFn(token.GREATER_THAN, p.parseInfixExpression)

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
func (p *Parser) parseExpression(precedence uint64) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.prefixParseFnNotFound(p.curToken.Type)
		return nil
	}
	leftExpr := prefix()

	for !p.isPeekTokenEq(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.advance()

		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) prefixParseFnNotFound(tokenType token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("No prefix parse function found for token %s \n", tokenType))
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

// ----------------------------------------------------------------- //
// 			    Parsing Functions For Expressions           	     //
// ----------------------------------------------------------------- //

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

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// Move the cursor to the next token to parse
	p.advance()

	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infx := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// Get the current precedence for the operator, as we will need it to parse
	// the right side of the expression
	precedence := p.currentPrecendence()

	// Advance the read pointer to the next token
	p.advance()

	infx.Right = p.parseExpression(precedence)
	return infx
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

// ----------------------------------------------------------------- //
// Precedence functions 						                     //
// ----------------------------------------------------------------- //

// Gets the precedence of the successor token if there is an any. Returns the
// LOWEST precedence if no precedence is found for the peek token.
func (p *Parser) peekPrecedence() uint64 {
	if pred, ok := PRECEDENCE_TABLE[p.peekToken.Type]; ok {
		return pred
	}
	return LOWEST
}

// Gets the precedence of the current token if there is any, returns the LOWEST
// precedende if none is found for the current token
func (p *Parser) currentPrecendence() uint64 {
	if pred, ok := PRECEDENCE_TABLE[p.curToken.Type]; ok {
		return pred
	}
	return LOWEST
}

// -------------------------------------------------------------------- //
