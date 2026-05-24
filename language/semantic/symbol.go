package semantic

import "fmt"

// ValueType represents the type of a value in the DSL
type ValueType string

const (
	IntType     ValueType = "INT"
	BoolType    ValueType = "BOOL"
	StringType  ValueType = "STRING"
	MutexType   ValueType = "MUTEX"
	UnknownType ValueType = "UNKNOWN"
)

// Symbol represents a declared symbol (variable, mutex, etc.)
type Symbol struct {
	Name          string
	Type          ValueType
	IsShared      bool   // true if shared, false if local
	IsMutex       bool   // true if this is a mutex
	Line          int    // line where declared
	Scope         int    // scope level (0 = global)
	ThreadName    string // if non-empty, this symbol is local to this thread
	UsageCount    int    // number of times referenced
	LastUsageLine int    // line of last usage
}

// String returns a string representation of the symbol
func (s Symbol) String() string {
	if s.IsMutex {
		return fmt.Sprintf("mutex %s (line %d)", s.Name, s.Line)
	}
	if s.IsShared {
		return fmt.Sprintf("shared %s: %s (line %d)", s.Name, s.Type, s.Line)
	}
	if s.ThreadName != "" {
		return fmt.Sprintf("local %s: %s in thread %q (line %d)", s.Name, s.Type, s.ThreadName, s.Line)
	}
	return fmt.Sprintf("%s: %s (line %d)", s.Name, s.Type, s.Line)
}

// Scope represents a lexical scope with its own symbol table
type Scope struct {
	Level      int
	Symbols    map[string]*Symbol
	Parent     *Scope
	ThreadName string // non-empty if this is a thread scope
}

// NewScope creates a new scope
func NewScope(level int, parent *Scope) *Scope {
	return &Scope{
		Level:   level,
		Symbols: make(map[string]*Symbol),
		Parent:  parent,
	}
}

// NewThreadScope creates a new scope for a thread
func NewThreadScope(level int, parent *Scope, threadName string) *Scope {
	s := NewScope(level, parent)
	s.ThreadName = threadName
	return s
}

// AddSymbol adds a symbol to this scope
// Returns an error if symbol already exists in this scope
func (s *Scope) AddSymbol(symbol *Symbol) error {
	if _, exists := s.Symbols[symbol.Name]; exists {
		return fmt.Errorf("symbol %q already exists in this scope", symbol.Name)
	}
	s.Symbols[symbol.Name] = symbol
	return nil
}

// LookupSymbol looks up a symbol in this scope and parent scopes
// Returns the symbol and whether it was found
func (s *Scope) LookupSymbol(name string) (*Symbol, bool) {
	if sym, exists := s.Symbols[name]; exists {
		return sym, true
	}
	if s.Parent != nil {
		return s.Parent.LookupSymbol(name)
	}
	return nil, false
}

// LookupLocalSymbol looks up a symbol only in this scope
func (s *Scope) LookupLocalSymbol(name string) (*Symbol, bool) {
	sym, exists := s.Symbols[name]
	return sym, exists
}

// AllSymbols returns all symbols in this scope (not including parent scopes)
func (s *Scope) AllSymbols() map[string]*Symbol {
	return s.Symbols
}

// SymbolTable manages all symbols and scopes for a program
type SymbolTable struct {
	GlobalScope  *Scope
	CurrentScope *Scope
	ThreadScopes map[string]*Scope // thread name -> scope
	Threads      []string          // list of thread names in order
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	global := NewScope(0, nil)
	return &SymbolTable{
		GlobalScope:  global,
		CurrentScope: global,
		ThreadScopes: make(map[string]*Scope),
		Threads:      []string{},
	}
}

// PushScope creates and enters a new scope
func (st *SymbolTable) PushScope() *Scope {
	newScope := NewScope(st.CurrentScope.Level+1, st.CurrentScope)
	st.CurrentScope = newScope
	return newScope
}

// PushThreadScope creates and enters a new thread scope
func (st *SymbolTable) PushThreadScope(threadName string) *Scope {
	newScope := NewThreadScope(st.CurrentScope.Level+1, st.CurrentScope, threadName)
	st.CurrentScope = newScope
	st.ThreadScopes[threadName] = newScope
	st.Threads = append(st.Threads, threadName)
	return newScope
}

// PopScope exits the current scope
func (st *SymbolTable) PopScope() {
	if st.CurrentScope.Parent != nil {
		st.CurrentScope = st.CurrentScope.Parent
	}
}

// AddSymbol adds a symbol to the current scope
func (st *SymbolTable) AddSymbol(symbol *Symbol) error {
	return st.CurrentScope.AddSymbol(symbol)
}

// LookupSymbol looks up a symbol in the current scope and parent scopes
func (st *SymbolTable) LookupSymbol(name string) (*Symbol, bool) {
	return st.CurrentScope.LookupSymbol(name)
}

// LookupLocalSymbol looks up a symbol only in the current scope
func (st *SymbolTable) LookupLocalSymbol(name string) (*Symbol, bool) {
	return st.CurrentScope.LookupLocalSymbol(name)
}

// RecordUsage records a usage of a symbol
func (st *SymbolTable) RecordUsage(name string, line int) {
	if sym, ok := st.CurrentScope.LookupSymbol(name); ok {
		sym.UsageCount++
		sym.LastUsageLine = line
	}
}

// GetGlobalSymbols returns all symbols defined at global scope
func (st *SymbolTable) GetGlobalSymbols() map[string]*Symbol {
	return st.GlobalScope.AllSymbols()
}

// GetThreadSymbols returns all symbols defined in a thread scope
func (st *SymbolTable) GetThreadSymbols(threadName string) map[string]*Symbol {
	if scope, exists := st.ThreadScopes[threadName]; exists {
		return scope.AllSymbols()
	}
	return make(map[string]*Symbol)
}

// CurrentScopeInfo returns information about the current scope
func (st *SymbolTable) CurrentScopeInfo() (level int, isThreadScope bool, threadName string) {
	return st.CurrentScope.Level, st.CurrentScope.ThreadName != "", st.CurrentScope.ThreadName
}

// String returns a debug representation of the symbol table
func (st *SymbolTable) String() string {
	result := "Global Scope:\n"
	for name, sym := range st.GlobalScope.AllSymbols() {
		result += fmt.Sprintf("  %s: %s\n", name, sym)
	}

	for threadName, scope := range st.ThreadScopes {
		result += fmt.Sprintf("Thread %q Scope:\n", threadName)
		for name, sym := range scope.AllSymbols() {
			result += fmt.Sprintf("  %s: %s\n", name, sym)
		}
	}

	return result
}
