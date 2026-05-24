package semantic

import (
	"deadlock/language/ast"
)

// Analyzer performs semantic analysis on an AST
type Analyzer struct {
	program           *ast.Program
	symTable          *SymbolTable
	errors            []SemanticError
	inThreadScope     bool
	currentThreadName string
}

// NewAnalyzer creates a new semantic analyzer
func NewAnalyzer(program *ast.Program) *Analyzer {
	return &Analyzer{
		program:           program,
		symTable:          NewSymbolTable(),
		errors:            []SemanticError{},
		inThreadScope:     false,
		currentThreadName: "",
	}
}

// Analyze performs semantic analysis on the program
// Returns the symbol table and list of errors
func (a *Analyzer) Analyze() (*SymbolTable, ErrorList) {
	if a.program == nil {
		return a.symTable, ErrorList(a.errors)
	}

	// Phase 1: Collect symbols
	a.collectSymbols()

	// Phase 2: Validate scopes
	a.validateScopes()

	// Phase 3: Validate mutexes and locks
	a.validateMutexes()

	// Phase 4: Validate types
	a.validateTypes()

	return a.symTable, ErrorList(a.errors)
}

// collectSymbols walks the AST and populates the symbol table
func (a *Analyzer) collectSymbols() {
	for _, stmt := range a.program.Statements {
		a.collectSymbolsFromStatement(stmt)
	}
}

// collectSymbolsFromStatement collects symbols from a statement
func (a *Analyzer) collectSymbolsFromStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		a.collectFromVariableDeclaration(s)
	case *ast.MutexStatement:
		a.collectFromMutexStatement(s)
	case *ast.BlockStatement:
		if s != nil {
			for _, stmt := range s.Statements {
				a.collectSymbolsFromStatement(stmt)
			}
		}
	case *ast.ExpressionStatement:
		if s != nil && s.Expression != nil {
			a.collectSymbolsFromExpression(s.Expression)
		}
	}
}

// collectSymbolsFromExpression collects symbols from an expression
func (a *Analyzer) collectSymbolsFromExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.SharedBlock:
		a.collectFromSharedBlock(e)
	case *ast.ThreadExpression:
		a.collectFromThreadExpression(e)
	}
}

// collectFromVariableDeclaration collects a variable declaration
func (a *Analyzer) collectFromVariableDeclaration(vd *ast.VariableDeclaration) {
	if vd == nil {
		return
	}

	sym := &Symbol{
		Name:     vd.Name.Value,
		Type:     a.inferExpressionType(vd.Value),
		IsShared: false,
		IsMutex:  false,
		Line:     vd.Name.Token.Line,
	}

	// Check for duplicate
	if existing, ok := a.symTable.LookupLocalSymbol(vd.Name.Value); ok {
		a.errors = append(a.errors, NewDuplicateVariableError(vd.Name.Value, vd.Name.Token.Line, existing.Line))
		return
	}

	err := a.symTable.AddSymbol(sym)
	if err != nil {
		a.errors = append(a.errors, SemanticError{
			Type:     DuplicateVariable,
			Severity: Error,
			Message:  err.Error(),
			Line:     vd.Name.Token.Line,
			Token:    vd.Name.Value,
		})
	}
}

// collectFromMutexStatement collects a mutex declaration
func (a *Analyzer) collectFromMutexStatement(ms *ast.MutexStatement) {
	if ms == nil || ms.Name == nil {
		return
	}

	sym := &Symbol{
		Name:     ms.Name.Value,
		Type:     MutexType,
		IsMutex:  true,
		IsShared: true, // mutexes are always accessible globally
		Line:     ms.Name.Token.Line,
	}

	// Check for duplicate
	if existing, ok := a.symTable.LookupLocalSymbol(ms.Name.Value); ok {
		a.errors = append(a.errors, NewDuplicateMutexError(ms.Name.Value, ms.Name.Token.Line, existing.Line))
		return
	}

	err := a.symTable.AddSymbol(sym)
	if err != nil {
		a.errors = append(a.errors, SemanticError{
			Type:     DuplicateMutex,
			Severity: Error,
			Message:  err.Error(),
			Line:     ms.Name.Token.Line,
			Token:    ms.Name.Value,
		})
	}
}

// collectFromSharedBlock collects symbols from a shared block
func (a *Analyzer) collectFromSharedBlock(sb *ast.SharedBlock) {
	if sb == nil {
		return
	}

	for _, decl := range sb.Declarations {
		if decl == nil {
			continue
		}
		sym := &Symbol{
			Name:     decl.Name.Value,
			Type:     a.inferExpressionType(decl.Value),
			IsShared: true,
			IsMutex:  false,
			Line:     decl.Name.Token.Line,
		}

		// Check for duplicate
		if existing, ok := a.symTable.LookupLocalSymbol(decl.Name.Value); ok {
			a.errors = append(a.errors, NewDuplicateVariableError(decl.Name.Value, decl.Name.Token.Line, existing.Line))
			continue
		}

		err := a.symTable.AddSymbol(sym)
		if err != nil {
			a.errors = append(a.errors, SemanticError{
				Type:     DuplicateVariable,
				Severity: Error,
				Message:  err.Error(),
				Line:     decl.Name.Token.Line,
				Token:    decl.Name.Value,
			})
		}
	}
}

// collectFromThreadExpression collects symbols from a thread expression
func (a *Analyzer) collectFromThreadExpression(te *ast.ThreadExpression) {
	if te == nil || te.Name == nil {
		return
	}

	threadName := te.Name.Value
	a.symTable.PushThreadScope(threadName)

	// Process thread body to collect local variables
	if te.Body != nil {
		for _, stmt := range te.Body.Statements {
			a.collectSymbolsFromStatement(stmt)
		}
	}

	a.symTable.PopScope()
}

// validateScopes validates that all identifiers are declared and in scope
func (a *Analyzer) validateScopes() {
	a.validateScopesInProgram()
}

// validateScopesInProgram validates scopes in the program
func (a *Analyzer) validateScopesInProgram() {
	for _, stmt := range a.program.Statements {
		a.validateScopesInStatement(stmt, false, "")
	}
}

// validateScopesInStatement validates scopes in a statement
func (a *Analyzer) validateScopesInStatement(stmt ast.Statement, inThread bool, threadName string) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		if s != nil && s.Expression != nil {
			a.validateScopesInExpression(s.Expression, inThread, threadName)
		}

	case *ast.BlockStatement:
		if s != nil {
			for _, stmt := range s.Statements {
				a.validateScopesInStatement(stmt, inThread, threadName)
			}
		}

	case *ast.VariableDeclaration:
		// Variable declarations are already collected; just validate the value expression
		if s != nil && s.Value != nil {
			a.validateScopesInExpression(s.Value, inThread, threadName)
		}

	case *ast.LockStatement:
		if s != nil && s.Argument != nil {
			// Lock target must be an identifier that exists in scope
			a.validateIdentifierInScope(s.Argument, inThread, threadName)
		}

	case *ast.UnlockStatement:
		if s != nil && s.Argument != nil {
			a.validateIdentifierInScope(s.Argument, inThread, threadName)
		}

	case *ast.PrintStmt:
		if s != nil && s.Value != nil {
			a.validateScopesInExpression(s.Value, inThread, threadName)
		}
	}
}

// validateScopesInExpression validates scopes in an expression
func (a *Analyzer) validateScopesInExpression(expr ast.Expression, inThread bool, threadName string) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.Identifier:
		a.validateIdentifierInScope(e, inThread, threadName)

	case *ast.AssignmentExpression:
		if e != nil {
			if e.Left != nil {
				a.validateScopesInExpression(e.Left, inThread, threadName)
			}
			if e.Value != nil {
				a.validateScopesInExpression(e.Value, inThread, threadName)
			}
		}

	case *ast.InfixExpression:
		if e != nil {
			if e.Left != nil {
				a.validateScopesInExpression(e.Left, inThread, threadName)
			}
			if e.Right != nil {
				a.validateScopesInExpression(e.Right, inThread, threadName)
			}
		}

	case *ast.PrefixExpression:
		if e != nil && e.Right != nil {
			a.validateScopesInExpression(e.Right, inThread, threadName)
		}

	case *ast.IfExpression:
		if e != nil {
			if e.Condition != nil {
				a.validateScopesInExpression(e.Condition, inThread, threadName)
			}
			if e.Consequence != nil {
				for _, stmt := range e.Consequence.Statements {
					a.validateScopesInStatement(stmt, inThread, threadName)
				}
			}
			if e.Alternative != nil {
				for _, stmt := range e.Alternative.Statements {
					a.validateScopesInStatement(stmt, inThread, threadName)
				}
			}
		}

	case *ast.IndexExpression:
		if e != nil {
			if e.Left != nil {
				a.validateScopesInExpression(e.Left, inThread, threadName)
			}
			if e.Index != nil {
				a.validateScopesInExpression(e.Index, inThread, threadName)
			}
		}

	case *ast.SharedBlock:
		if e != nil {
			for _, decl := range e.Declarations {
				if decl != nil && decl.Value != nil {
					a.validateScopesInExpression(decl.Value, false, "")
				}
			}
		}

	case *ast.ThreadExpression:
		if e != nil && e.Body != nil {
			for _, stmt := range e.Body.Statements {
				a.validateScopesInStatement(stmt, true, e.Name.Value)
			}
		}
	}
}

// validateIdentifierInScope checks if an identifier is declared and in scope
func (a *Analyzer) validateIdentifierInScope(ident *ast.Identifier, inThread bool, threadName string) {
	if ident == nil {
		return
	}

	// Need to reconstruct the symbol table state for this validation
	// This is a simplified approach; in a real implementation, you'd track scope during AST walk
	sym, ok := a.symTable.GlobalScope.LookupSymbol(ident.Value)
	if !ok {
		// Check thread scope if in a thread
		if inThread && threadName != "" {
			if threadScope, exists := a.symTable.ThreadScopes[threadName]; exists {
				sym, ok = threadScope.LookupSymbol(ident.Value)
			}
		}
	}

	if !ok {
		a.errors = append(a.errors, NewUndeclaredVariableError(ident.Value, ident.Token.Line))
	} else if sym != nil && sym.IsShared && inThread {
		// Record usage of shared variable in thread (this is valid, but we track it)
		a.symTable.RecordUsage(ident.Value, ident.Token.Line)
	}
}

// validateMutexes validates mutex and lock/unlock operations
func (a *Analyzer) validateMutexes() {
	a.validateMutexesInProgram()
}

// validateMutexesInProgram validates mutexes in the program
func (a *Analyzer) validateMutexesInProgram() {
	for _, stmt := range a.program.Statements {
		a.validateMutexesInStatement(stmt, false, "")
	}
}

// validateMutexesInStatement validates mutexes in a statement
func (a *Analyzer) validateMutexesInStatement(stmt ast.Statement, inThread bool, threadName string) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.LockStatement:
		if s != nil && s.Argument != nil {
			a.validateLockTarget(s.Argument, s.Argument.Token.Line)
		}

	case *ast.UnlockStatement:
		if s != nil && s.Argument != nil {
			a.validateUnlockTarget(s.Argument, s.Argument.Token.Line)
		}

	case *ast.BlockStatement:
		if s != nil {
			for _, stmt := range s.Statements {
				a.validateMutexesInStatement(stmt, inThread, threadName)
			}
		}

	case *ast.ExpressionStatement:
		if s != nil && s.Expression != nil {
			a.validateMutexesInExpression(s.Expression, inThread, threadName)
		}
	}
}

// validateMutexesInExpression validates mutexes in an expression
func (a *Analyzer) validateMutexesInExpression(expr ast.Expression, inThread bool, threadName string) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.IfExpression:
		if e != nil {
			if e.Consequence != nil {
				for _, stmt := range e.Consequence.Statements {
					a.validateMutexesInStatement(stmt, inThread, threadName)
				}
			}
			if e.Alternative != nil {
				for _, stmt := range e.Alternative.Statements {
					a.validateMutexesInStatement(stmt, inThread, threadName)
				}
			}
		}

	case *ast.SharedBlock:
		if e != nil {
			for _, decl := range e.Declarations {
				if decl != nil && decl.Value != nil {
					a.validateMutexesInExpression(decl.Value, false, "")
				}
			}
		}

	case *ast.ThreadExpression:
		if e != nil && e.Body != nil {
			for _, stmt := range e.Body.Statements {
				a.validateMutexesInStatement(stmt, true, e.Name.Value)
			}
		}
	}
}

// validateLockTarget checks if a lock target is a valid mutex
func (a *Analyzer) validateLockTarget(target *ast.Identifier, line int) {
	if target == nil {
		return
	}

	sym, ok := a.symTable.GlobalScope.LookupSymbol(target.Value)
	if !ok {
		a.errors = append(a.errors, NewUndeclaredMutexError(target.Value, line))
		return
	}

	if !sym.IsMutex {
		a.errors = append(a.errors, NewInvalidLockOpError(target.Value, line))
	}
}

// validateUnlockTarget checks if an unlock target is a valid mutex
func (a *Analyzer) validateUnlockTarget(target *ast.Identifier, line int) {
	if target == nil {
		return
	}

	sym, ok := a.symTable.GlobalScope.LookupSymbol(target.Value)
	if !ok {
		a.errors = append(a.errors, NewUndeclaredMutexError(target.Value, line))
		return
	}

	if !sym.IsMutex {
		a.errors = append(a.errors, NewInvalidUnlockOpError(target.Value, line))
	}
}

// validateTypes validates type consistency
func (a *Analyzer) validateTypes() {
	// Type validation is performed as we traverse the AST
	// This will be implemented in Phase 4
	// For now, we'll do basic checking
	a.validateTypesInProgram()
}

// validateTypesInProgram validates types in the program
func (a *Analyzer) validateTypesInProgram() {
	for _, stmt := range a.program.Statements {
		a.validateTypesInStatement(stmt)
	}
}

// validateTypesInStatement validates types in a statement
func (a *Analyzer) validateTypesInStatement(stmt ast.Statement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		if s != nil && s.Expression != nil {
			a.validateTypesInExpression(s.Expression)
		}

	case *ast.VariableDeclaration:
		// Types are inferred during collection
		if s != nil && s.Value != nil {
			a.validateTypesInExpression(s.Value)
		}

	case *ast.BlockStatement:
		if s != nil {
			for _, stmt := range s.Statements {
				a.validateTypesInStatement(stmt)
			}
		}

	case *ast.PrintStmt:
		if s != nil && s.Value != nil {
			a.validateTypesInExpression(s.Value)
		}
	}
}

// validateTypesInExpression validates types in an expression
func (a *Analyzer) validateTypesInExpression(expr ast.Expression) ValueType {
	if expr == nil {
		return UnknownType
	}

	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return IntType

	case *ast.StringLiteral:
		return StringType

	case *ast.Boolean:
		return BoolType

	case *ast.Identifier:
		if sym, ok := a.symTable.GlobalScope.LookupSymbol(e.Value); ok {
			return sym.Type
		}
		return UnknownType

	case *ast.PrefixExpression:
		if e != nil {
			rightType := a.validateTypesInExpression(e.Right)
			switch e.Operator {
			case "-":
				if rightType != IntType && rightType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(rightType), e.Token.Line))
				}
				return IntType
			case "!":
				if rightType != BoolType && rightType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(rightType), e.Token.Line))
				}
				return BoolType
			}
		}
		return UnknownType

	case *ast.InfixExpression:
		if e != nil {
			leftType := a.validateTypesInExpression(e.Left)
			rightType := a.validateTypesInExpression(e.Right)

			switch e.Operator {
			case "+", "-", "*", "/", "%":
				if leftType != IntType && leftType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(leftType), e.Token.Line))
				}
				if rightType != IntType && rightType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(rightType), e.Token.Line))
				}
				return IntType

			case ">", "<", ">=", "<=":
				if leftType != IntType && leftType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(leftType), e.Token.Line))
				}
				if rightType != IntType && rightType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(rightType), e.Token.Line))
				}
				return BoolType

			case "==", "!=":
				// Equality can work on multiple types, but both sides should match
				if leftType != UnknownType && rightType != UnknownType && leftType != rightType {
					a.errors = append(a.errors, NewTypeMismatchError(string(leftType), string(rightType), e.Token.Line, e.Operator))
				}
				return BoolType

			case "&&", "||":
				if leftType != BoolType && leftType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(leftType), e.Token.Line))
				}
				if rightType != BoolType && rightType != UnknownType {
					a.errors = append(a.errors, NewInvalidOperationError(e.Operator, string(rightType), e.Token.Line))
				}
				return BoolType
			}
		}
		return UnknownType

	case *ast.AssignmentExpression:
		if e != nil && e.Left != nil {
			leftType := a.validateTypesInExpression(e.Left)
			rightType := a.validateTypesInExpression(e.Value)

			if leftType != UnknownType && rightType != UnknownType && leftType != rightType {
				a.errors = append(a.errors, NewTypeMismatchError(string(leftType), string(rightType), e.Token.Line, "="))
			}
			return leftType
		}
		return UnknownType

	case *ast.IfExpression:
		if e != nil && e.Condition != nil {
			condType := a.validateTypesInExpression(e.Condition)
			if condType != BoolType && condType != UnknownType {
				a.errors = append(a.errors, NewInvalidOperationError("if", string(condType), e.Token.Line))
			}

			if e.Consequence != nil {
				for _, stmt := range e.Consequence.Statements {
					a.validateTypesInStatement(stmt)
				}
			}

			if e.Alternative != nil {
				for _, stmt := range e.Alternative.Statements {
					a.validateTypesInStatement(stmt)
				}
			}
		}
		return UnknownType
	}

	return UnknownType
}

// inferExpressionType infers the type of an expression
func (a *Analyzer) inferExpressionType(expr ast.Expression) ValueType {
	if expr == nil {
		return UnknownType
	}

	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return IntType
	case *ast.StringLiteral:
		return StringType
	case *ast.Boolean:
		return BoolType
	case *ast.Identifier:
		if sym, ok := a.symTable.GlobalScope.LookupSymbol(e.Value); ok {
			return sym.Type
		}
		return UnknownType
	case *ast.InfixExpression:
		if e != nil {
			switch e.Operator {
			case "+", "-", "*", "/", "%":
				return IntType
			case ">", "<", ">=", "<=", "==", "!=", "&&", "||":
				return BoolType
			}
		}
		return UnknownType
	case *ast.PrefixExpression:
		if e != nil {
			switch e.Operator {
			case "-":
				return IntType
			case "!":
				return BoolType
			}
		}
		return UnknownType
	default:
		return UnknownType
	}
}

// AnalyzeProgram is a convenience function to create an analyzer and run analysis
func AnalyzeProgram(program *ast.Program) (*SymbolTable, ErrorList) {
	analyzer := NewAnalyzer(program)
	return analyzer.Analyze()
}
