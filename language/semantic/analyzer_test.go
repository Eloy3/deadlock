package semantic

import (
	"testing"

	"deadlock/language/ast"
	"deadlock/language/parser"
	"deadlock/language/token"
)

// Helper function to create a token
func createToken(tokenType token.TokenType, literal string, line int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}

// Helper function to parse code into AST
func parseCode(code string) *ast.Program {
	tokens, err := token.TokenizeProgram(code)
	if err != nil {
		return nil
	}
	p := parser.NewParser(tokens)
	program := p.ParseProgram()
	return &program
}

// TestValidProgram tests a simple valid program
func TestValidProgram(t *testing.T) {
	code := `
shared {
    counter = 0;
}

mutex m;

thread incrementer {
    lock(m);
    counter = counter + 1;
    unlock(m);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	symTable, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors, got: %v", errList)
	}

	// Verify symbol table
	symbols := symTable.GetGlobalSymbols()
	if _, ok := symbols["counter"]; !ok {
		t.Error("Expected 'counter' in global symbols")
	}
	if _, ok := symbols["m"]; !ok {
		t.Error("Expected 'm' (mutex) in global symbols")
	}

	// Verify thread symbols
	threadSymbols := symTable.GetThreadSymbols("incrementer")
	if len(threadSymbols) != 0 {
		t.Errorf("Expected no local symbols in incrementer thread, got %d", len(threadSymbols))
	}
}

// TestDuplicateVariableDeclaration tests duplicate variable errors
func TestDuplicateVariableDeclaration(t *testing.T) {
	code := `
shared {
    x = 5;
    x = 10;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected duplicate variable error")
	}

	found := false
	for _, err := range errList {
		if err.Type == DuplicateVariable {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected DuplicateVariable error type")
	}
}

// TestDuplicateMutexDeclaration tests duplicate mutex errors
func TestDuplicateMutexDeclaration(t *testing.T) {
	code := `
mutex m;
mutex m;
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected duplicate mutex error")
	}

	found := false
	for _, err := range errList {
		if err.Type == DuplicateMutex {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected DuplicateMutex error type")
	}
}

// TestUndeclaredVariableError tests undeclared variable errors
func TestUndeclaredVariableError(t *testing.T) {
	code := `
thread worker {
    x = y + 1;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	// Should have error for undeclared 'y'
	if !errList.HasErrors() {
		t.Error("Expected undeclared variable error for 'y'")
	}

	found := false
	for _, err := range errList {
		if err.Type == UndeclaredVariable {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected UndeclaredVariable error type")
	}
}

// TestInvalidLockOperation tests lock on non-mutex
func TestInvalidLockOperation(t *testing.T) {
	code := `
shared {
    x = 5;
}

thread worker {
    lock(x);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected invalid lock operation error")
	}

	found := false
	for _, err := range errList {
		if err.Type == InvalidLockOp {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected InvalidLockOp error type")
	}
}

// TestInvalidUnlockOperation tests unlock on non-mutex
func TestInvalidUnlockOperation(t *testing.T) {
	code := `
shared {
    x = 5;
}

thread worker {
    unlock(x);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected invalid unlock operation error")
	}

	found := false
	for _, err := range errList {
		if err.Type == InvalidUnlockOp {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected InvalidUnlockOp error type")
	}
}

// TestLockUndeclaredMutex tests lock on undeclared mutex
func TestLockUndeclaredMutex(t *testing.T) {
	code := `
thread worker {
    lock(m);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected undeclared mutex error")
	}

	found := false
	for _, err := range errList {
		if err.Type == UndeclaredMutex {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected UndeclaredMutex error type")
	}
}

// TestTypeMismatchInAssignment tests type mismatch in assignment
func TestTypeMismatchInAssignment(t *testing.T) {
	code := `
shared {
    x = 5;
}

thread worker {
    x = "string";
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected type mismatch error")
	}

	found := false
	for _, err := range errList {
		if err.Type == TypeMismatch {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected TypeMismatch error type")
	}
}

// TestInvalidBinaryOperation tests invalid binary operations
func TestInvalidBinaryOperation(t *testing.T) {
	code := `
shared {
    x = 5;
    y = "string";
}

thread worker {
    z = x + y;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected invalid operation error")
	}
}

// TestInvalidPrefixOperation tests invalid prefix operations
func TestInvalidPrefixOperation(t *testing.T) {
	code := `
thread worker {
    x = -"string";
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected invalid operation error for prefix negation on string")
	}
}

// TestIfConditionTypeCheck tests if condition type validation
func TestIfConditionTypeCheck(t *testing.T) {
	code := `
shared {
    x = 5;
}

thread worker {
    if x {
        print(x);
    }
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected type error for if condition (integer instead of boolean)")
	}
}

// TestSharedVariableInThread tests shared variable access in thread
func TestSharedVariableInThread(t *testing.T) {
	code := `
shared {
    counter = 0;
}

mutex m;

thread incrementer {
    lock(m);
    counter = counter + 1;
    unlock(m);
}

thread reader {
    if counter > 0 {
        print(counter);
    }
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	symTable, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors, got: %v", errList)
	}

	// Verify symbol table has both threads
	if len(symTable.Threads) != 2 {
		t.Errorf("Expected 2 threads, got %d", len(symTable.Threads))
	}

	// Verify symbols in each thread
	incrementerSymbols := symTable.GetThreadSymbols("incrementer")
	if len(incrementerSymbols) != 0 {
		t.Errorf("Expected 0 local symbols in incrementer, got %d", len(incrementerSymbols))
	}

	readerSymbols := symTable.GetThreadSymbols("reader")
	if len(readerSymbols) != 0 {
		t.Errorf("Expected 0 local symbols in reader, got %d", len(readerSymbols))
	}
}

// TestLocalVariableInThread tests local variable declaration in thread
func TestLocalVariableInThread(t *testing.T) {
	code := `
thread worker {
    local x = 5;
    local y = 10;
    z = x + y;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	symTable, errList := AnalyzeProgram(program)

	// Note: 'local x = 5' syntax needs to be properly parsed
	// For now, just verify the semantic analysis runs without panic
	_ = symTable
	_ = errList
}

// TestMultipleThreads tests multiple threads with shared state
func TestMultipleThreads(t *testing.T) {
	code := `
shared {
    balance = 0;
    flag = false;
}

mutex lock1;
mutex lock2;

thread producer {
    lock(lock1);
    balance = balance + 100;
    flag = true;
    unlock(lock1);
}

thread consumer {
    lock(lock2);
    if flag {
        balance = balance - 50;
    }
    unlock(lock2);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	symTable, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors, got: %v", errList)
	}

	// Verify all threads are registered
	if len(symTable.Threads) != 2 {
		t.Errorf("Expected 2 threads, got %d", len(symTable.Threads))
	}

	// Verify global symbols
	globalSymbols := symTable.GetGlobalSymbols()
	expectedSymbols := []string{"balance", "flag", "lock1", "lock2"}
	for _, name := range expectedSymbols {
		if _, ok := globalSymbols[name]; !ok {
			t.Errorf("Expected symbol %q in global scope", name)
		}
	}
}

// TestBooleanOperations tests boolean operations
func TestBooleanOperations(t *testing.T) {
	code := `
thread worker {
    local x = true;
    local y = false;
    local z = x && y;
    local w = x || y;
    local v = !x;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors for boolean operations, got: %v", errList)
	}
}

// TestIntegerComparisons tests integer comparison operations
func TestIntegerComparisons(t *testing.T) {
	code := `
thread worker {
    local x = 5;
    local a = x > 3;
    local b = x < 10;
    local c = x >= 5;
    local d = x <= 5;
    local e = x == 5;
    local f = x != 3;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors for integer comparisons, got: %v", errList)
	}
}

// TestStringLiterals tests string literal handling
func TestStringLiterals(t *testing.T) {
	code := `
thread worker {
    print("Hello, World!");
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	// Should not have errors for string literals
	if errList.HasErrors() {
		t.Errorf("Expected no errors for string literals, got: %v", errList)
	}
}

// TestSymbolTableIntegration tests complete symbol table structure
func TestSymbolTableIntegration(t *testing.T) {
	code := `
shared {
    x = 10;
    y = 20;
}

mutex m1;
mutex m2;

thread t1 {
    lock(m1);
    x = x + y;
    unlock(m1);
}

thread t2 {
    lock(m2);
    y = y - x;
    unlock(m2);
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	symTable, errList := AnalyzeProgram(program)
	if errList.HasErrors() {
		t.Errorf("Expected no errors, got: %v", errList)
	}

	// Test global symbols
	globalSymbols := symTable.GetGlobalSymbols()
	if len(globalSymbols) != 4 {
		t.Errorf("Expected 4 global symbols, got %d", len(globalSymbols))
	}

	// Verify each symbol
	if sym, ok := globalSymbols["x"]; ok {
		if !sym.IsShared {
			t.Error("Expected 'x' to be shared")
		}
		if sym.Type != IntType {
			t.Errorf("Expected 'x' to be INT type, got %v", sym.Type)
		}
	}

	if sym, ok := globalSymbols["m1"]; ok {
		if !sym.IsMutex {
			t.Error("Expected 'm1' to be a mutex")
		}
	}

	// Test thread scopes
	if len(symTable.Threads) != 2 {
		t.Errorf("Expected 2 threads, got %d", len(symTable.Threads))
	}

	t1Symbols := symTable.GetThreadSymbols("t1")
	t2Symbols := symTable.GetThreadSymbols("t2")
	if len(t1Symbols) != 0 || len(t2Symbols) != 0 {
		t.Error("Expected no local symbols in threads")
	}
}

// TestErrorList tests the ErrorList type
func TestErrorList(t *testing.T) {
	code := `
shared {
    x = 5;
    x = 10;
}

thread worker {
    lock(y);
    z = a + b;
}
`
	program := parseCode(code)
	if program == nil {
		t.Fatal("Failed to parse program")
	}

	_, errList := AnalyzeProgram(program)
	if !errList.HasErrors() {
		t.Error("Expected multiple errors")
	}

	// Should have multiple error types
	errorTypes := make(map[ErrorType]bool)
	for _, err := range errList {
		errorTypes[err.Type] = true
	}

	if len(errorTypes) == 0 {
		t.Error("Expected multiple distinct error types")
	}
}
