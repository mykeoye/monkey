package ast

import (
	"monkey/token"
	"testing"
)

func TestAstNodes(t *testing.T) {
	var expected = "let a = 5;"

	ast := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "a"},
					Value: "a",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: "5",
				},
			},
		},
	}

	if ast.String() != expected {
		t.Errorf("The generated AST %s doesn't match the expected program string %s \n", ast.String(), expected)
	}
}
