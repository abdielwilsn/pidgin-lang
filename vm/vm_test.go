package vm

import (
	"testing"
)

// ============================================================================
// Manual Bytecode Tests
// ============================================================================

func TestManualBytecode_5Plus3(t *testing.T) {
	// Manually create bytecode for: 5 + 3
	// Expected bytecode:
	//   OP_CONST_I8 5
	//   OP_CONST_I8 3
	//   OP_ADD
	//   OP_HALT

	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(3, 1)
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	result, err := vm.Run(chunk)

	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if !result.IsInt() {
		t.Fatalf("Expected int result, got %s", result.TypeName())
	}

	if got := result.AsInt(); got != 8 {
		t.Errorf("Expected 8, got %d", got)
	}
}

func TestManualBytecode_Arithmetic(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		b        int64
		op       Opcode
		expected int64
	}{
		{"5 + 3", 5, 3, OP_ADD, 8},
		{"10 - 4", 10, 4, OP_SUB, 6},
		{"6 * 7", 6, 7, OP_MUL, 42},
		{"20 / 4", 20, 4, OP_DIV, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewChunk()
			chunk.WriteOpcode(OP_CONST_I8, 1)
			chunk.WriteByte(byte(tt.a), 1)
			chunk.WriteOpcode(OP_CONST_I8, 1)
			chunk.WriteByte(byte(tt.b), 1)
			chunk.WriteOpcode(tt.op, 1)
			chunk.WriteOpcode(OP_HALT, 1)

			vm := NewVM()
			result, err := vm.Run(chunk)

			if err != nil {
				t.Fatalf("Execution error: %v", err)
			}

			if !result.IsInt() {
				t.Fatalf("Expected int result, got %s", result.TypeName())
			}

			if got := result.AsInt(); got != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestManualBytecode_Negation(t *testing.T) {
	// Bytecode for: -42
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(42, 1)
	chunk.WriteOpcode(OP_NEGATE, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	result, err := vm.Run(chunk)

	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if got := result.AsInt(); got != -42 {
		t.Errorf("Expected -42, got %d", got)
	}
}

func TestManualBytecode_Comparison(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		b        int64
		op       Opcode
		expected bool
	}{
		{"5 == 5", 5, 5, OP_EQUAL, true},
		{"5 == 3", 5, 3, OP_EQUAL, false},
		{"5 != 3", 5, 3, OP_NOT_EQUAL, true},
		{"5 > 3", 5, 3, OP_GREATER, true},
		{"3 > 5", 3, 5, OP_GREATER, false},
		{"3 < 5", 3, 5, OP_LESS, true},
		{"5 < 3", 5, 3, OP_LESS, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewChunk()
			chunk.WriteOpcode(OP_CONST_I8, 1)
			chunk.WriteByte(byte(tt.a), 1)
			chunk.WriteOpcode(OP_CONST_I8, 1)
			chunk.WriteByte(byte(tt.b), 1)
			chunk.WriteOpcode(tt.op, 1)
			chunk.WriteOpcode(OP_HALT, 1)

			vm := NewVM()
			result, err := vm.Run(chunk)

			if err != nil {
				t.Fatalf("Execution error: %v", err)
			}

			if !result.IsBool() {
				t.Fatalf("Expected bool result, got %s", result.TypeName())
			}

			if got := result.AsBool(); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestManualBytecode_Literals(t *testing.T) {
	tests := []struct {
		name     string
		op       Opcode
		expected Value
	}{
		{"const 0", OP_CONST_0, NewInt(0)},
		{"const 1", OP_CONST_1, NewInt(1)},
		{"const -1", OP_CONST_MINUS1, NewInt(-1)},
		{"true", OP_TRU, NewBool(true)},
		{"false", OP_LIE, NewBool(false)},
		{"nothing", OP_NOTHING, NewNothing()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewChunk()
			chunk.WriteOpcode(tt.op, 1)
			chunk.WriteOpcode(OP_HALT, 1)

			vm := NewVM()
			result, err := vm.Run(chunk)

			if err != nil {
				t.Fatalf("Execution error: %v", err)
			}

			if !result.Equals(tt.expected) {
				t.Errorf("Expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func TestManualBytecode_StackOperations(t *testing.T) {
	// Test DUP: Push 5, duplicate it, add them -> 10
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_DUP, 1)
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	result, err := vm.Run(chunk)

	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if got := result.AsInt(); got != 10 {
		t.Errorf("Expected 10, got %d", got)
	}
}

func TestManualBytecode_LogicalNot(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected bool
	}{
		{"!true", NewBool(true), false},
		{"!false", NewBool(false), true},
		{"!nothing", NewNothing(), true},
		{"!42", NewInt(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewChunk()

			// Push the value onto the constant pool
			idx := chunk.AddConstant(tt.value)
			chunk.WriteOpcode(OP_CONSTANT, 1)
			chunk.WriteByte(byte(idx>>8), 1)
			chunk.WriteByte(byte(idx&0xFF), 1)

			chunk.WriteOpcode(OP_NOT, 1)
			chunk.WriteOpcode(OP_HALT, 1)

			vm := NewVM()
			result, err := vm.Run(chunk)

			if err != nil {
				t.Fatalf("Execution error: %v", err)
			}

			if got := result.AsBool(); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestManualBytecode_StringConcatenation(t *testing.T) {
	// Test: "Hello" + " " + "World"
	chunk := NewChunk()

	// Add strings to constant pool
	hello := chunk.InternString("Hello")
	space := chunk.InternString(" ")
	world := chunk.InternString("World")

	idx1 := chunk.AddConstant(NewString(hello))
	idx2 := chunk.AddConstant(NewString(space))
	idx3 := chunk.AddConstant(NewString(world))

	// Push "Hello"
	chunk.WriteOpcode(OP_CONSTANT, 1)
	chunk.WriteByte(byte(idx1>>8), 1)
	chunk.WriteByte(byte(idx1&0xFF), 1)

	// Push " "
	chunk.WriteOpcode(OP_CONSTANT, 1)
	chunk.WriteByte(byte(idx2>>8), 1)
	chunk.WriteByte(byte(idx2&0xFF), 1)

	// Add "Hello" + " "
	chunk.WriteOpcode(OP_ADD, 1)

	// Push "World"
	chunk.WriteOpcode(OP_CONSTANT, 1)
	chunk.WriteByte(byte(idx3>>8), 1)
	chunk.WriteByte(byte(idx3&0xFF), 1)

	// Add "Hello " + "World"
	chunk.WriteOpcode(OP_ADD, 1)

	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	result, err := vm.Run(chunk)

	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if !result.IsString() {
		t.Fatalf("Expected string result, got %s", result.TypeName())
	}

	expected := "Hello World"
	if got := *result.AsString(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestManualBytecode_GlobalVariables(t *testing.T) {
	// Test: make x be 42; x + 8
	chunk := NewChunk()

	// Add variable name to constant pool
	varName := chunk.InternString("x")
	nameIdx := chunk.AddConstant(NewString(varName))

	// Push 42
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(42, 1)

	// Set global x = 42
	chunk.WriteOpcode(OP_SET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Get global x
	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Push 8
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(8, 1)

	// Add x + 8
	chunk.WriteOpcode(OP_ADD, 1)

	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	result, err := vm.Run(chunk)

	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if got := result.AsInt(); got != 50 {
		t.Errorf("Expected 50, got %d", got)
	}
}

// ============================================================================
// Error Tests
// ============================================================================

func TestVM_DivisionByZero(t *testing.T) {
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(10, 1)
	chunk.WriteOpcode(OP_CONST_0, 1)
	chunk.WriteOpcode(OP_DIV, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	_, err := vm.Run(chunk)

	if err == nil {
		t.Fatal("Expected division by zero error, got nil")
	}
}

func TestVM_TypeErrorArithmetic(t *testing.T) {
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_TRU, 1)
	chunk.WriteOpcode(OP_SUB, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	_, err := vm.Run(chunk)

	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestVM_UndefinedVariable(t *testing.T) {
	chunk := NewChunk()

	varName := chunk.InternString("undefined_var")
	nameIdx := chunk.AddConstant(NewString(varName))

	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()
	_, err := vm.Run(chunk)

	if err == nil {
		t.Fatal("Expected undefined variable error, got nil")
	}
}

// ============================================================================
// Disassembly Tests
// ============================================================================

func TestChunk_Disassemble(t *testing.T) {
	// Create bytecode for: 5 + 3
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(3, 1)
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	// This should not panic
	chunk.Disassemble("test")
}
