package compiler

import (
	"fmt"
	"sort"

	"pidgin-lang/ast"
	"pidgin-lang/vm"
)

// Compiler compiles Pidgin AST into bytecode
type Compiler struct {
	chunk       *vm.Chunk      // Output bytecode chunk
	symbolTable *SymbolTable   // Symbol table for variable tracking
	scopeDepth  int            // Current scope nesting level
	constants   map[string]int // Cache for constant indices
}

// New creates a new compiler
func New() *Compiler {
	symbolTable := NewSymbolTable()

	// Define builtins
	symbolTable.DefineBuiltin(0, "yarn")
	symbolTable.DefineBuiltin(1, "len")
	symbolTable.DefineBuiltin(2, "type")

	return &Compiler{
		chunk:       vm.NewChunk(),
		symbolTable: symbolTable,
		scopeDepth:  0,
		constants:   make(map[string]int),
	}
}

// Compile compiles an AST program into bytecode
func (c *Compiler) Compile(program *ast.Program) (*vm.Chunk, error) {
	numStmts := len(program.Statements)

	for i, stmt := range program.Statements {
		isLast := (i == numStmts-1)

		if err := c.compileStatementWithContext(stmt, isLast); err != nil {
			return nil, err
		}
	}

	// Emit halt at the end
	c.emit(vm.OP_HALT)

	return c.chunk, nil
}

// compileStatementWithContext compiles a statement with knowledge of whether it's the last one
func (c *Compiler) compileStatementWithContext(stmt ast.Statement, isLast bool) error {
	// For the last statement, don't pop expression results
	if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok && isLast {
		return c.compileExpression(exprStmt.Expression)
	}

	return c.compileStatement(stmt)
}

// ============================================================================
// Statement Compilation
// ============================================================================

func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch node := stmt.(type) {

	case *ast.ExpressionStatement:
		// Compile the expression
		if err := c.compileExpression(node.Expression); err != nil {
			return err
		}
		// Pop the result (expression statements don't use their value)
		c.emit(vm.OP_POP)
		return nil

	case *ast.MakeStatement:
		// Compile the value first
		if err := c.compileExpression(node.Value); err != nil {
			return err
		}

		// Define or set the variable
		name := node.Name.Value
		symbol, exists := c.symbolTable.Resolve(name)

		if c.scopeDepth == 0 {
			// Global scope
			if !exists {
				symbol = c.symbolTable.Define(name)
			}

			// Add variable name to constants pool
			nameStr := c.chunk.InternString(name)
			idx := c.addConstant(vm.NewString(nameStr))

			c.emitShort(vm.OP_SET_GLOBAL, uint16(idx))
		} else {
			// Local scope
			if !exists {
				symbol = c.symbolTable.Define(name)
			}

			// Emit optimized SET_LOCAL for first two locals
			if symbol.Index == 0 {
				c.emit(vm.OP_SET_LOCAL_0)
			} else if symbol.Index == 1 {
				c.emit(vm.OP_SET_LOCAL_1)
			} else {
				c.emitByte(vm.OP_SET_LOCAL, byte(symbol.Index))
			}
		}

		return nil

	case *ast.BringStatement:
		// Compile the return value
		if node.ReturnValue != nil {
			if err := c.compileExpression(node.ReturnValue); err != nil {
				return err
			}
		} else {
			c.emit(vm.OP_NOTHING)
		}

		c.emit(vm.OP_BRING)
		return nil

	case *ast.BlockStatement:
		numStmts := len(node.Statements)
		for i, stmt := range node.Statements {
			isLast := (i == numStmts-1)

			// Don't pop the last expression in a block (it becomes the block's value)
			if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok && isLast {
				if err := c.compileExpression(exprStmt.Expression); err != nil {
					return err
				}
			} else {
				if err := c.compileStatement(stmt); err != nil {
					return err
				}
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

// ============================================================================
// Expression Compilation
// ============================================================================

func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch node := expr.(type) {

	case *ast.IntegerLiteral:
		return c.compileIntegerLiteral(node)

	case *ast.StringLiteral:
		return c.compileStringLiteral(node)

	case *ast.Boolean:
		return c.compileBooleanLiteral(node)

	case *ast.NothingLiteral:
		c.emit(vm.OP_NOTHING)
		return nil

	case *ast.Identifier:
		return c.compileIdentifier(node)

	case *ast.PrefixExpression:
		return c.compilePrefixExpression(node)

	case *ast.InfixExpression:
		return c.compileInfixExpression(node)

	case *ast.SupposeExpression:
		return c.compileSupposeExpression(node)

	case *ast.WhileExpression:
		return c.compileWhileExpression(node)

	case *ast.CallExpression:
		return c.compileCallExpression(node)

	default:
		return fmt.Errorf("unknown expression type: %T", expr)
	}
}

// ============================================================================
// Literal Compilation
// ============================================================================

func (c *Compiler) compileIntegerLiteral(node *ast.IntegerLiteral) error {
	val := node.Value

	// Optimize for common small integers
	if val == 0 {
		c.emit(vm.OP_CONST_0)
	} else if val == 1 {
		c.emit(vm.OP_CONST_1)
	} else if val == -1 {
		c.emit(vm.OP_CONST_MINUS1)
	} else if val >= -128 && val <= 127 {
		// 8-bit inline integer
		c.emitByte(vm.OP_CONST_I8, byte(int8(val)))
	} else if val >= -32768 && val <= 32767 {
		// 16-bit inline integer
		c.emitShort(vm.OP_CONST_I16, uint16(int16(val)))
	} else {
		// Large integer - use constant pool
		idx := c.addConstant(vm.NewInt(val))
		c.emitShort(vm.OP_CONSTANT, uint16(idx))
	}

	return nil
}

func (c *Compiler) compileStringLiteral(node *ast.StringLiteral) error {
	strPtr := c.chunk.InternString(node.Value)
	idx := c.addConstant(vm.NewString(strPtr))
	c.emitShort(vm.OP_CONSTANT, uint16(idx))
	return nil
}

func (c *Compiler) compileBooleanLiteral(node *ast.Boolean) error {
	if node.Value {
		c.emit(vm.OP_TRU)
	} else {
		c.emit(vm.OP_LIE)
	}
	return nil
}

// ============================================================================
// Identifier Compilation
// ============================================================================

func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	name := node.Value
	symbol, ok := c.symbolTable.Resolve(name)

	if !ok {
		return fmt.Errorf("I no sabi dis one: %s", name)
	}

	switch symbol.Scope {
	case SCOPE_GLOBAL:
		// Add variable name to constants pool
		nameStr := c.chunk.InternString(name)
		idx := c.addConstant(vm.NewString(nameStr))
		c.emitShort(vm.OP_GET_GLOBAL, uint16(idx))

	case SCOPE_LOCAL:
		// Emit optimized GET_LOCAL for first 4 locals
		if symbol.Index == 0 {
			c.emit(vm.OP_GET_LOCAL_0)
		} else if symbol.Index == 1 {
			c.emit(vm.OP_GET_LOCAL_1)
		} else if symbol.Index == 2 {
			c.emit(vm.OP_GET_LOCAL_2)
		} else if symbol.Index == 3 {
			c.emit(vm.OP_GET_LOCAL_3)
		} else {
			c.emitByte(vm.OP_GET_LOCAL, byte(symbol.Index))
		}

	case SCOPE_BUILTIN:
		// Push builtin function reference
		idx := c.addConstant(vm.NewBuiltin(symbol.Index))
		c.emitShort(vm.OP_CONSTANT, uint16(idx))

	default:
		return fmt.Errorf("unknown symbol scope: %v", symbol.Scope)
	}

	return nil
}

// ============================================================================
// Operator Compilation
// ============================================================================

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	// Compile the operand
	if err := c.compileExpression(node.Right); err != nil {
		return err
	}

	// Emit the operator
	switch node.Operator {
	case "-":
		c.emit(vm.OP_NEGATE)
	case "!", "no":
		c.emit(vm.OP_NOT)
	case "no be":
		c.emit(vm.OP_NOT)
	default:
		return fmt.Errorf("unknown prefix operator: %s", node.Operator)
	}

	return nil
}

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
	// Handle short-circuit operators specially
	if node.Operator == "and" {
		return c.compileShortCircuitAnd(node)
	}
	if node.Operator == "abi" || node.Operator == "or" {
		return c.compileShortCircuitOr(node)
	}

	// Compile left operand
	if err := c.compileExpression(node.Left); err != nil {
		return err
	}

	// Compile right operand
	if err := c.compileExpression(node.Right); err != nil {
		return err
	}

	// Emit the operator
	switch node.Operator {
	case "+":
		c.emit(vm.OP_ADD)
	case "-":
		c.emit(vm.OP_SUB)
	case "*":
		c.emit(vm.OP_MUL)
	case "/":
		c.emit(vm.OP_DIV)
	case "be", "na", "==":
		c.emit(vm.OP_EQUAL)
	case "no be", "!=":
		c.emit(vm.OP_NOT_EQUAL)
	case "big pass", ">":
		c.emit(vm.OP_GREATER)
	case "no reach", "<":
		c.emit(vm.OP_LESS)
	default:
		return fmt.Errorf("unknown infix operator: %s", node.Operator)
	}

	return nil
}

// compileShortCircuitAnd implements: a and b
// If a is false, skip b and return false
func (c *Compiler) compileShortCircuitAnd(node *ast.InfixExpression) error {
	// Compile left operand
	if err := c.compileExpression(node.Left); err != nil {
		return err
	}

	// Duplicate it for the jump condition
	c.emit(vm.OP_DUP)

	// Jump if left is false (short-circuit)
	jumpIfFalse := c.emitJump(vm.OP_JUMP_IF_LIE)

	// Left was true, pop it and evaluate right
	c.emit(vm.OP_POP)
	if err := c.compileExpression(node.Right); err != nil {
		return err
	}

	// Patch the jump
	c.patchJump(jumpIfFalse)

	return nil
}

// compileShortCircuitOr implements: a abi b (or a or b)
// If a is true, skip b and return true
func (c *Compiler) compileShortCircuitOr(node *ast.InfixExpression) error {
	// Compile left operand
	if err := c.compileExpression(node.Left); err != nil {
		return err
	}

	// Duplicate it for the jump condition
	c.emit(vm.OP_DUP)

	// Jump if left is true (short-circuit)
	jumpIfTrue := c.emitJump(vm.OP_JUMP_IF_TRU)

	// Left was false, pop it and evaluate right
	c.emit(vm.OP_POP)
	if err := c.compileExpression(node.Right); err != nil {
		return err
	}

	// Patch the jump
	c.patchJump(jumpIfTrue)

	return nil
}

// ============================================================================
// Control Flow Compilation
// ============================================================================

func (c *Compiler) compileSupposeExpression(node *ast.SupposeExpression) error {
	// Compile condition
	if err := c.compileExpression(node.Condition); err != nil {
		return err
	}

	// Jump if condition is false
	jumpIfFalse := c.emitJump(vm.OP_JUMP_IF_LIE)

	// Compile consequence
	if err := c.compileStatement(node.Consequence); err != nil {
		return err
	}

	// Jump over alternative
	jumpEnd := c.emitJump(vm.OP_JUMP)

	// Patch the jump to alternative
	c.patchJump(jumpIfFalse)

	// Compile alternative (or push nothing if there isn't one)
	if node.Alternative != nil {
		if err := c.compileStatement(node.Alternative); err != nil {
			return err
		}
	} else {
		c.emit(vm.OP_NOTHING)
	}

	// Patch the jump to end
	c.patchJump(jumpEnd)

	return nil
}

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) error {
	// Mark loop start
	loopStart := c.chunk.Count()

	// Compile condition
	if err := c.compileExpression(node.Condition); err != nil {
		return err
	}

	// Jump if condition is false (exit loop)
	exitJump := c.emitJump(vm.OP_JUMP_IF_LIE)

	// Compile loop body
	if err := c.compileStatement(node.Body); err != nil {
		return err
	}

	// Loop back to start
	c.emitLoop(loopStart)

	// Patch exit jump
	c.patchJump(exitJump)

	// Push nothing as the result (loops return nothing)
	c.emit(vm.OP_NOTHING)

	return nil
}

// ============================================================================
// Function Call Compilation
// ============================================================================

func (c *Compiler) compileCallExpression(node *ast.CallExpression) error {
	// Check if this is a builtin call
	if ident, ok := node.Function.(*ast.Identifier); ok {
		if symbol, found := c.symbolTable.Resolve(ident.Value); found && symbol.Scope == SCOPE_BUILTIN {
			return c.compileBuiltinCall(ident.Value, symbol.Index, node.Arguments)
		}
	}

	// Regular function call
	// Compile arguments
	for _, arg := range node.Arguments {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Compile function
	if err := c.compileExpression(node.Function); err != nil {
		return err
	}

	// Emit optimized CALL for 0-2 arguments
	argCount := len(node.Arguments)
	if argCount == 0 {
		c.emit(vm.OP_CALL_0)
	} else if argCount == 1 {
		c.emit(vm.OP_CALL_1)
	} else if argCount == 2 {
		c.emit(vm.OP_CALL_2)
	} else {
		c.emitByte(vm.OP_CALL, byte(argCount))
	}

	return nil
}

func (c *Compiler) compileBuiltinCall(name string, builtinIdx int, args []ast.Expression) error {
	// Special optimization for yarn (print)
	if name == "yarn" {
		// Compile arguments
		for _, arg := range args {
			if err := c.compileExpression(arg); err != nil {
				return err
			}
		}

		// Emit optimized YARN instruction
		c.emitByte(vm.OP_YARN, byte(len(args)))
		return nil
	}

	// Other builtins
	// Compile arguments
	for _, arg := range args {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Emit BUILTIN instruction
	c.emitBytes(vm.OP_BUILTIN, byte(builtinIdx), byte(len(args)))
	return nil
}

// ============================================================================
// Code Generation Helpers
// ============================================================================

func (c *Compiler) emit(op vm.Opcode) int {
	pos := c.chunk.Count()
	c.chunk.WriteOpcode(op, 0) // Line number 0 for now (Phase 2 doesn't track lines)
	return pos
}

func (c *Compiler) emitByte(op vm.Opcode, operand byte) int {
	pos := c.emit(op)
	c.chunk.WriteByte(operand, 0)
	return pos
}

func (c *Compiler) emitBytes(op vm.Opcode, operands ...byte) int {
	pos := c.emit(op)
	for _, b := range operands {
		c.chunk.WriteByte(b, 0)
	}
	return pos
}

func (c *Compiler) emitShort(op vm.Opcode, operand uint16) int {
	pos := c.emit(op)
	c.chunk.WriteByte(byte(operand>>8), 0)
	c.chunk.WriteByte(byte(operand&0xFF), 0)
	return pos
}

func (c *Compiler) emitJump(op vm.Opcode) int {
	c.emit(op)
	c.chunk.WriteByte(0xFF, 0) // Placeholder
	c.chunk.WriteByte(0xFF, 0) // Placeholder
	return c.chunk.Count() - 2
}

func (c *Compiler) patchJump(offset int) {
	// Calculate jump distance
	jump := c.chunk.Count() - offset - 2

	if jump > 65535 {
		panic("Too much code to jump over")
	}

	c.chunk.Code[offset] = byte(jump >> 8)
	c.chunk.Code[offset+1] = byte(jump & 0xFF)
}

func (c *Compiler) emitLoop(loopStart int) {
	c.emit(vm.OP_LOOP)

	offset := c.chunk.Count() - loopStart + 2
	if offset > 65535 {
		panic("Loop body too large")
	}

	c.chunk.WriteByte(byte(offset>>8), 0)
	c.chunk.WriteByte(byte(offset&0xFF), 0)
}

func (c *Compiler) addConstant(value vm.Value) int {
	// Check if we already have this constant
	key := value.String()
	if idx, exists := c.constants[key]; exists {
		return idx
	}

	idx := c.chunk.AddConstant(value)
	c.constants[key] = idx
	return idx
}

// ============================================================================
// Debugging
// ============================================================================

// Disassemble prints the compiled bytecode
func (c *Compiler) Disassemble(name string) {
	c.chunk.Disassemble(name)
}

// Chunk returns the compiled bytecode chunk
func (c *Compiler) Chunk() *vm.Chunk {
	return c.chunk
}

// SymbolTable returns the symbol table
func (c *Compiler) SymbolTable() *SymbolTable {
	return c.symbolTable
}

// Constants returns a sorted list of constants for debugging
func (c *Compiler) Constants() []string {
	keys := make([]string, 0, len(c.constants))
	for k := range c.constants {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
