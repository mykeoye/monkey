package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
	 	let y = 10;
		let foobar = 838383;`

	// Create a lexer for tokenizing the input text
	l := lexer.NewLexer(input)

	// Create a parser
	p := NewParser(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Check for errors while parsing
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements but found %d \n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		// Test statement to ensure it matches our expectation
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func TestReturnStatement(t *testing.T) {
	input := `return 5;
	 	return 10;
		return 838383;`

	// Create a lexer for tokenizing the input text
	l := lexer.NewLexer(input)

	// Create a parser
	p := NewParser(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Check for errors while parsing
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements but found %d \n", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected return statement but found %T", returnStmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Expected a 'return' literal but got=%s", returnStmt.TokenLiteral())
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func testLetStatement(t *testing.T, s ast.Statement, literal string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("Expected a let but got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("Expected a ast.LetStatement but got %T instead", letStmt)
		return false
	}

	if letStmt.Name.Value != literal {
		t.Fatalf("Expected %s in stmt.Name.Value but got %s \n", literal, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != literal {
		t.Fatalf("Expected %s in stmt.Name.TokenLiteral() but got %s \n", literal, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestIntegerLiterals(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)

	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}
	t.Errorf("Parser has %d errors\n", len(errs))
	for _, err := range errs {
		t.Errorf("Parser error %s \n", err)
	}
	t.FailNow()
}
