package compiler

// SymbolScope represents the scope level of a symbol
type SymbolScope int

const (
	SCOPE_GLOBAL  SymbolScope = iota // Global variable
	SCOPE_LOCAL                       // Local variable (function parameter or local)
	SCOPE_UPVALUE                     // Captured variable from enclosing scope
	SCOPE_BUILTIN                     // Builtin function
)

// Symbol represents a variable or function in the symbol table
type Symbol struct {
	Name  string      // Variable/function name
	Scope SymbolScope // Scope level
	Index int         // Slot index for locals, or global index
}

// SymbolTable tracks variable bindings and their scopes
type SymbolTable struct {
	outer          *SymbolTable       // Enclosing scope (for nested functions)
	store          map[string]Symbol  // Symbol storage
	numDefinitions int                // Number of definitions in this scope
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

// NewEnclosedSymbolTable creates a new symbol table enclosed in another
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		outer: outer,
		store: make(map[string]Symbol),
	}
}

// Define adds a new symbol to the current scope
func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: st.numDefinitions,
	}

	// Determine scope based on whether we have an outer scope
	if st.outer == nil {
		symbol.Scope = SCOPE_GLOBAL
	} else {
		symbol.Scope = SCOPE_LOCAL
	}

	st.store[name] = symbol
	st.numDefinitions++

	return symbol
}

// DefineBuiltin adds a builtin function to the symbol table
func (st *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: SCOPE_BUILTIN,
		Index: index,
	}

	st.store[name] = symbol
	return symbol
}

// Resolve looks up a symbol by name, searching enclosing scopes if necessary
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	// Try to find in current scope
	symbol, ok := st.store[name]
	if ok {
		return symbol, true
	}

	// If not found and we have an outer scope, search there
	if st.outer != nil {
		symbol, ok = st.outer.Resolve(name)
		if !ok {
			return symbol, false
		}

		// If the symbol is a local in an outer scope, it becomes an upvalue
		if symbol.Scope == SCOPE_LOCAL {
			// For Phase 2, we'll treat upvalues as globals for simplicity
			// Phase 4 will implement proper closure support
			return symbol, true
		}

		return symbol, true
	}

	// Symbol not found
	return Symbol{}, false
}

// NumDefinitions returns the number of symbols defined in this scope
func (st *SymbolTable) NumDefinitions() int {
	return st.numDefinitions
}

// ScopeName returns a human-readable name for a scope
func (s SymbolScope) String() string {
	switch s {
	case SCOPE_GLOBAL:
		return "GLOBAL"
	case SCOPE_LOCAL:
		return "LOCAL"
	case SCOPE_UPVALUE:
		return "UPVALUE"
	case SCOPE_BUILTIN:
		return "BUILTIN"
	default:
		return "UNKNOWN"
	}
}
