package parser

import (
	"deadlock/language/ast"
	"deadlock/language/token"
	"testing"
)

// parseInput tokenizes and parses the input, returning the parsed program.
// It handles tokenization errors and parser errors, failing the test if either occur.
func parseInput(t *testing.T, input string) ast.Program {
	tokens, err := token.TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing program: %s", err.Error())
	}
	parser := NewParser(tokens)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	return program
}

// requireExactlyOneStatement asserts that the program contains exactly one statement.
// It fails the test if the count doesn't match.
func requireExactlyOneStatement(t *testing.T, program ast.Program) {
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
}
