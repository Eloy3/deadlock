package semantic

import "fmt"

// ErrorType categorizes semantic errors
type ErrorType string

const (
	UndeclaredVariable ErrorType = "UNDECLARED_VARIABLE"
	DuplicateVariable  ErrorType = "DUPLICATE_VARIABLE"
	UndeclaredMutex    ErrorType = "UNDECLARED_MUTEX"
	DuplicateMutex     ErrorType = "DUPLICATE_MUTEX"
	InvalidLockOp      ErrorType = "INVALID_LOCK_OP"
	InvalidUnlockOp    ErrorType = "INVALID_UNLOCK_OP"
	TypeMismatch       ErrorType = "TYPE_MISMATCH"
	InvalidOperation   ErrorType = "INVALID_OPERATION"
	InvalidAssignment  ErrorType = "INVALID_ASSIGNMENT"
	ScopeViolation     ErrorType = "SCOPE_VIOLATION"
	InvalidThreadBody  ErrorType = "INVALID_THREAD_BODY"
	InvalidFunctionArg ErrorType = "INVALID_FUNCTION_ARG"
)

// Severity levels for errors
type Severity string

const (
	Error   Severity = "ERROR"
	Warning Severity = "WARNING"
	Info    Severity = "INFO"
)

// SemanticError represents a semantic analysis error
type SemanticError struct {
	Type     ErrorType
	Severity Severity
	Message  string
	Line     int
	Token    string // The token causing the error
}

// Error implements the error interface
func (se SemanticError) Error() string {
	return fmt.Sprintf("[%s] Line %d: %s (token: %q)", se.Type, se.Line, se.Message, se.Token)
}

// String returns a human-readable error message
func (se SemanticError) String() string {
	return fmt.Sprintf("Line %d: %s", se.Line, se.Message)
}

// ErrorList is a collection of semantic errors
type ErrorList []SemanticError

// Error implements the error interface for ErrorList
func (el ErrorList) Error() string {
	if len(el) == 0 {
		return "no semantic errors"
	}
	msg := fmt.Sprintf("%d semantic error(s):\n", len(el))
	for i, err := range el {
		msg += fmt.Sprintf("  %d. %s\n", i+1, err)
	}
	return msg
}

// HasErrors returns true if there are any errors
func (el ErrorList) HasErrors() bool {
	for _, err := range el {
		if err.Severity == Error {
			return true
		}
	}
	return false
}

// Helper functions to create errors

// NewUndeclaredVariableError creates an undeclared variable error
func NewUndeclaredVariableError(name string, line int) SemanticError {
	return SemanticError{
		Type:     UndeclaredVariable,
		Severity: Error,
		Message:  fmt.Sprintf("undeclared variable %q", name),
		Line:     line,
		Token:    name,
	}
}

// NewDuplicateVariableError creates a duplicate variable error
func NewDuplicateVariableError(name string, line int, prevLine int) SemanticError {
	return SemanticError{
		Type:     DuplicateVariable,
		Severity: Error,
		Message:  fmt.Sprintf("variable %q already declared at line %d", name, prevLine),
		Line:     line,
		Token:    name,
	}
}

// NewUndeclaredMutexError creates an undeclared mutex error
func NewUndeclaredMutexError(name string, line int) SemanticError {
	return SemanticError{
		Type:     UndeclaredMutex,
		Severity: Error,
		Message:  fmt.Sprintf("undeclared mutex %q", name),
		Line:     line,
		Token:    name,
	}
}

// NewDuplicateMutexError creates a duplicate mutex error
func NewDuplicateMutexError(name string, line int, prevLine int) SemanticError {
	return SemanticError{
		Type:     DuplicateMutex,
		Severity: Error,
		Message:  fmt.Sprintf("mutex %q already declared at line %d", name, prevLine),
		Line:     line,
		Token:    name,
	}
}

// NewInvalidLockOpError creates an invalid lock operation error
func NewInvalidLockOpError(name string, line int) SemanticError {
	return SemanticError{
		Type:     InvalidLockOp,
		Severity: Error,
		Message:  fmt.Sprintf("cannot lock non-mutex %q", name),
		Line:     line,
		Token:    name,
	}
}

// NewInvalidUnlockOpError creates an invalid unlock operation error
func NewInvalidUnlockOpError(name string, line int) SemanticError {
	return SemanticError{
		Type:     InvalidUnlockOp,
		Severity: Error,
		Message:  fmt.Sprintf("cannot unlock non-mutex %q", name),
		Line:     line,
		Token:    name,
	}
}

// NewTypeMismatchError creates a type mismatch error
func NewTypeMismatchError(expected string, got string, line int, token string) SemanticError {
	return SemanticError{
		Type:     TypeMismatch,
		Severity: Error,
		Message:  fmt.Sprintf("type mismatch: expected %s, got %s", expected, got),
		Line:     line,
		Token:    token,
	}
}

// NewInvalidOperationError creates an invalid operation error
func NewInvalidOperationError(op string, operandType string, line int) SemanticError {
	return SemanticError{
		Type:     InvalidOperation,
		Severity: Error,
		Message:  fmt.Sprintf("invalid operation %q on type %s", op, operandType),
		Line:     line,
		Token:    op,
	}
}

// NewInvalidAssignmentError creates an invalid assignment error
func NewInvalidAssignmentError(varName string, line int) SemanticError {
	return SemanticError{
		Type:     InvalidAssignment,
		Severity: Error,
		Message:  fmt.Sprintf("invalid assignment to %q", varName),
		Line:     line,
		Token:    varName,
	}
}

// NewScopeViolationError creates a scope violation error
func NewScopeViolationError(name string, scope string, line int) SemanticError {
	return SemanticError{
		Type:     ScopeViolation,
		Severity: Error,
		Message:  fmt.Sprintf("cannot access %s variable %q outside %s scope", scope, name, scope),
		Line:     line,
		Token:    name,
	}
}
