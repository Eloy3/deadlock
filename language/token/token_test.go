package token

import (
	"testing"
)

func TestSharedBlockTokenization(t *testing.T) {
	input := `shared {
				counter = 0;
				ready = false;
			  }`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		SHARED,
		LBRACE,
		IDENTIFIER, // counter
		ASSIGN,
		NUMBER,
		SEMICOLON,
		IDENTIFIER, // ready
		ASSIGN,
		FALSE,
		SEMICOLON,
		RBRACE,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestMutexDeclarationTokenization(t *testing.T) {
	input := "mutex m;"
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		MUTEX,
		IDENTIFIER,
		SEMICOLON,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestThreadDeclarationTokenization(t *testing.T) {
	input := `thread incrementer {
local temp = 0;
lock(m);
temp = counter;
temp = temp + 1;
counter = temp;
unlock(m);
}`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		THREAD,
		IDENTIFIER, // incrementer
		LBRACE,
		LOCAL,
		IDENTIFIER, // temp
		ASSIGN,
		NUMBER,
		SEMICOLON,
		LOCK,
		LPAREN,
		IDENTIFIER, // m
		RPAREN,
		SEMICOLON,
		IDENTIFIER, // temp
		ASSIGN,
		IDENTIFIER, // counter
		SEMICOLON,
		IDENTIFIER, // temp
		ASSIGN,
		IDENTIFIER, // temp
		PLUS,
		NUMBER,
		SEMICOLON,
		IDENTIFIER, // counter
		ASSIGN,
		IDENTIFIER, // temp
		SEMICOLON,
		UNLOCK,
		LPAREN,
		IDENTIFIER, // m
		RPAREN,
		SEMICOLON,
		RBRACE,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestLocalVariableDeclarationTokenization(t *testing.T) {
	input := "local temp = 0;"
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		LOCAL,
		IDENTIFIER,
		ASSIGN,
		NUMBER,
		SEMICOLON,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestIfStatementTokenization(t *testing.T) {
	input := `if counter > 0 {
print("counter changed");
}`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		IF,
		IDENTIFIER, // counter
		GREATER,
		NUMBER,
		LBRACE,
		PRINT,
		LPAREN,
		STRING,
		RPAREN,
		SEMICOLON,
		RBRACE,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestStringLiteralTokenization(t *testing.T) {
	input := `"hello world";`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		STRING,
		SEMICOLON,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestLockUnlockStatementsTokenization(t *testing.T) {
	input := `lock(m);
unlock(m);`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		LOCK,
		LPAREN,
		IDENTIFIER,
		RPAREN,
		SEMICOLON,
		UNLOCK,
		LPAREN,
		IDENTIFIER,
		RPAREN,
		SEMICOLON,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestPrintStatementTokenization(t *testing.T) {
	input := `print("message");`
	tokens, err := TokenizeProgram(input)
	if err != nil {
		t.Fatalf("Error tokenizing: %s", err)
	}
	expected := []TokenType{
		PRINT,
		LPAREN,
		STRING,
		RPAREN,
		SEMICOLON,
		EOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}
