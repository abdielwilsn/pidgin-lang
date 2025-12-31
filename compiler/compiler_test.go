package compiler

import (
	"testing"

	"pidgin-lang/ast"
	"pidgin-lang/lexer"
	"pidgin-lang/parser"
	"pidgin-lang/vm"
)

// Helper function to parse Pidgin source code
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// ============================================================================
// Integer Literal Tests
// ============================================================================

func TestCompileIntegerLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		// Optimized small integers
		{"0", []byte{byte(vm.OP_CONST_0), byte(vm.OP_POP), byte(vm.OP_HALT)}},
		{"1", []byte{byte(vm.OP_CONST_1), byte(vm.OP_POP), byte(vm.OP_HALT)}},
		{"42", []byte{byte(vm.OP_CONST_I8), 42, byte(vm.OP_POP), byte(vm.OP_HALT)}},
		// Note: -1 is parsed as prefix expression, not a literal
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			compiler := New()

			chunk, err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}

			if len(chunk.Code) != len(tt.expected) {
				t.Errorf("wrong bytecode length. want=%d, got=%d", len(tt.expected), len(chunk.Code))
				t.Logf("Expected: %v", tt.expected)
				t.Logf("Got:      %v", chunk.Code)
			}

			for i, expectedByte := range tt.expected {
				if i >= len(chunk.Code) {
					break
				}
				if chunk.Code[i] != expectedByte {
					t.Errorf("wrong byte at position %d. want=%d, got=%d", i, expectedByte, chunk.Code[i])
				}
			}
		})
	}
}

// ============================================================================
// Boolean Literal Tests
// ============================================================================

func TestCompileBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected vm.Opcode
	}{
		{"tru", vm.OP_TRU},
		{"lie", vm.OP_LIE},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			compiler := New()

			chunk, err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}

			if vm.Opcode(chunk.Code[0]) != tt.expected {
				t.Errorf("wrong opcode. want=%s, got=%s", tt.expected, vm.Opcode(chunk.Code[0]))
			}
		})
	}
}

// ============================================================================
// Arithmetic Tests
// ============================================================================

func TestCompileArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		hasOp    vm.Opcode
	}{
		{"5 + 3", vm.OP_ADD},
		{"10 - 4", vm.OP_SUB},
		{"6 * 7", vm.OP_MUL},
		{"20 / 4", vm.OP_DIV},
		{"-42", vm.OP_NEGATE},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			compiler := New()

			chunk, err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}

			// Check that the expected opcode appears in the bytecode
			found := false
			for _, b := range chunk.Code {
				if vm.Opcode(b) == tt.hasOp {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected opcode %s not found in bytecode", tt.hasOp)
			}
		})
	}
}

// ============================================================================
// Comparison Tests
// ============================================================================

func TestCompileComparison(t *testing.T) {
	tests := []struct {
		input    string
		hasOp    vm.Opcode
	}{
		{"5 be 5", vm.OP_EQUAL},
		// Note: "no be" is not a single operator in the parser
		{"5 big pass 3", vm.OP_GREATER},
		{"3 no reach 5", vm.OP_LESS},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			compiler := New()

			chunk, err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}

			// Check that the expected opcode appears in the bytecode
			found := false
			for _, b := range chunk.Code {
				if vm.Opcode(b) == tt.hasOp {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected opcode %s not found in bytecode", tt.hasOp)
			}
		})
	}
}

// ============================================================================
// Variable Tests
// ============================================================================

func TestCompileGlobalVariables(t *testing.T) {
	input := `
	make x be 42
	x
	`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check for SET_GLOBAL and GET_GLOBAL opcodes
	hasSet := false
	hasGet := false
	for _, b := range chunk.Code {
		if vm.Opcode(b) == vm.OP_SET_GLOBAL {
			hasSet = true
		}
		if vm.Opcode(b) == vm.OP_GET_GLOBAL {
			hasGet = true
		}
	}

	if !hasSet {
		t.Error("expected OP_SET_GLOBAL not found")
	}
	if !hasGet {
		t.Error("expected OP_GET_GLOBAL not found")
	}
}

// ============================================================================
// String Tests
// ============================================================================

func TestCompileStringLiterals(t *testing.T) {
	input := `"Wetin dey happen?"`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check that we have a CONSTANT opcode
	if vm.Opcode(chunk.Code[0]) != vm.OP_CONSTANT {
		t.Errorf("expected OP_CONSTANT, got %s", vm.Opcode(chunk.Code[0]))
	}

	// Check that the string is in the constants pool
	if len(chunk.Constants) == 0 {
		t.Error("expected constant in pool")
	}
}

// ============================================================================
// Control Flow Tests
// ============================================================================

func TestCompileSupposeExpression(t *testing.T) {
	input := `suppose 5 big pass 3 { 1 } abi { 0 }`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check for jump instructions
	hasJumpIfLie := false
	hasJump := false
	for _, b := range chunk.Code {
		if vm.Opcode(b) == vm.OP_JUMP_IF_LIE {
			hasJumpIfLie = true
		}
		if vm.Opcode(b) == vm.OP_JUMP {
			hasJump = true
		}
	}

	if !hasJumpIfLie {
		t.Error("expected OP_JUMP_IF_LIE not found")
	}
	if !hasJump {
		t.Error("expected OP_JUMP not found")
	}
}

func TestCompileWhileExpression(t *testing.T) {
	input := `dey do while tru { 42 }`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check for jump and loop instructions
	hasJumpIfLie := false
	hasLoop := false
	for _, b := range chunk.Code {
		if vm.Opcode(b) == vm.OP_JUMP_IF_LIE {
			hasJumpIfLie = true
		}
		if vm.Opcode(b) == vm.OP_LOOP {
			hasLoop = true
		}
	}

	if !hasJumpIfLie {
		t.Error("expected OP_JUMP_IF_LIE not found")
	}
	if !hasLoop {
		t.Error("expected OP_LOOP not found")
	}
}

// ============================================================================
// Short-Circuit Tests
// ============================================================================

func TestCompileShortCircuitAnd(t *testing.T) {
	input := `tru and lie`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check for DUP and JUMP_IF_LIE (short-circuit pattern)
	hasDup := false
	hasJumpIfLie := false
	for _, b := range chunk.Code {
		if vm.Opcode(b) == vm.OP_DUP {
			hasDup = true
		}
		if vm.Opcode(b) == vm.OP_JUMP_IF_LIE {
			hasJumpIfLie = true
		}
	}

	if !hasDup {
		t.Error("expected OP_DUP not found (needed for short-circuit)")
	}
	if !hasJumpIfLie {
		t.Error("expected OP_JUMP_IF_LIE not found")
	}
}

// NOTE: OR operator ("abi") is not yet supported by the parser as an infix operator
// It's only used in "suppose...abi" contexts
// This test is skipped for Phase 2

// ============================================================================
// Builtin Function Tests
// ============================================================================

func TestCompileBuiltinYarn(t *testing.T) {
	input := `yarn("Hello")`

	program := parse(input)
	compiler := New()

	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// Check for optimized YARN instruction
	hasYarn := false
	for _, b := range chunk.Code {
		if vm.Opcode(b) == vm.OP_YARN {
			hasYarn = true
			break
		}
	}

	if !hasYarn {
		t.Error("expected OP_YARN not found")
	}
}

// ============================================================================
// Error Tests
// ============================================================================

func TestCompileUndefinedVariable(t *testing.T) {
	input := `undefined_var`

	program := parse(input)
	compiler := New()

	_, err := compiler.Compile(program)
	if err == nil {
		t.Error("expected compilation error for undefined variable, got nil")
	}
}

// ============================================================================
// Disassembly Test
// ============================================================================

func TestDisassemble(t *testing.T) {
	input := `5 + 3`

	program := parse(input)
	compiler := New()

	_, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compilation error: %v", err)
	}

	// This should not panic
	compiler.Disassemble("test")
}
