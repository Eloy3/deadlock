package parser

import (
	"deadlock/language/token"

	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	tokens, errors := token.TokenizeProgram(input)

	if errors != nil {
		t.Fatalf("Error tokenizing program")
	}

	if len(tokens) < 1 {
		t.Fatalf("No tokens detected")
	}

	parser := NewParser(tokens)
	statements := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(statements))
	}
}
