package vm

import (
	"testing"
)

// ============================================================================
// Benchmarks for Phase 1: Core VM Performance
// ============================================================================

func BenchmarkVM_ArithmeticAdd(b *testing.B) {
	// Benchmark: 5 + 3
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(3, 1)
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

func BenchmarkVM_ArithmeticComplex(b *testing.B) {
	// Benchmark: (5 + 3) * (10 - 2) / 4
	// Expected: (8) * (8) / 4 = 16
	chunk := NewChunk()

	// 5 + 3
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(3, 1)
	chunk.WriteOpcode(OP_ADD, 1)

	// 10 - 2
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(10, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(2, 1)
	chunk.WriteOpcode(OP_SUB, 1)

	// Multiply results
	chunk.WriteOpcode(OP_MUL, 1)

	// Divide by 4
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(4, 1)
	chunk.WriteOpcode(OP_DIV, 1)

	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

func BenchmarkVM_Comparison(b *testing.B) {
	// Benchmark: 5 > 3
	chunk := NewChunk()
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(5, 1)
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(3, 1)
	chunk.WriteOpcode(OP_GREATER, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

func BenchmarkVM_GlobalVariableAccess(b *testing.B) {
	// Benchmark: make x be 42; x + x
	chunk := NewChunk()

	varName := chunk.InternString("x")
	nameIdx := chunk.AddConstant(NewString(varName))

	// Set global x = 42
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(42, 1)
	chunk.WriteOpcode(OP_SET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Get global x (twice)
	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Add x + x
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

func BenchmarkVM_StringConcatenation(b *testing.B) {
	// Benchmark: "Hello" + " " + "World"
	chunk := NewChunk()

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

	// Add
	chunk.WriteOpcode(OP_ADD, 1)

	// Push "World"
	chunk.WriteOpcode(OP_CONSTANT, 1)
	chunk.WriteByte(byte(idx3>>8), 1)
	chunk.WriteByte(byte(idx3&0xFF), 1)

	// Add
	chunk.WriteOpcode(OP_ADD, 1)
	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

func BenchmarkVM_LoopSimulation(b *testing.B) {
	// Simulate a simple counting loop:
	// make counter be 0
	// dey do while counter no reach 100 {
	//     counter be counter + 1
	// }

	chunk := NewChunk()

	counterName := chunk.InternString("counter")
	nameIdx := chunk.AddConstant(NewString(counterName))

	// Set counter = 0
	chunk.WriteOpcode(OP_CONST_0, 1)
	chunk.WriteOpcode(OP_SET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Loop start (offset 7)
	loopStart := chunk.Count()

	// Get counter
	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Push 100
	chunk.WriteOpcode(OP_CONST_I8, 1)
	chunk.WriteByte(100, 1)

	// Check if counter < 100
	chunk.WriteOpcode(OP_LESS, 1)

	// Jump if false (exit loop)
	exitJumpPos := chunk.Count()
	chunk.WriteOpcode(OP_JUMP_IF_LIE, 1)
	chunk.WriteByte(0, 1) // Placeholder
	chunk.WriteByte(0, 1) // Placeholder

	// Get counter
	chunk.WriteOpcode(OP_GET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Push 1
	chunk.WriteOpcode(OP_CONST_1, 1)

	// Add counter + 1
	chunk.WriteOpcode(OP_ADD, 1)

	// Set counter
	chunk.WriteOpcode(OP_SET_GLOBAL, 1)
	chunk.WriteByte(byte(nameIdx>>8), 1)
	chunk.WriteByte(byte(nameIdx&0xFF), 1)

	// Loop back
	loopOffset := chunk.Count() - loopStart + 3
	chunk.WriteOpcode(OP_LOOP, 1)
	chunk.WriteByte(byte(loopOffset>>8), 1)
	chunk.WriteByte(byte(loopOffset&0xFF), 1)

	// Patch exit jump
	exitOffset := chunk.Count() - exitJumpPos - 3
	chunk.Code[exitJumpPos+1] = byte(exitOffset >> 8)
	chunk.Code[exitJumpPos+2] = byte(exitOffset & 0xFF)

	chunk.WriteOpcode(OP_HALT, 1)

	vm := NewVM()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(chunk)
	}
}

// ============================================================================
// NaN Boxing Performance Benchmarks
// ============================================================================

func BenchmarkNaNBox_NewInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewInt(42)
	}
}

func BenchmarkNaNBox_IsInt(b *testing.B) {
	v := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.IsInt()
	}
}

func BenchmarkNaNBox_AsInt(b *testing.B) {
	v := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.AsInt()
	}
}

func BenchmarkNaNBox_Equals(b *testing.B) {
	v1 := NewInt(42)
	v2 := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Equals(v2)
	}
}

// ============================================================================
// Stack Operations Benchmarks
// ============================================================================

func BenchmarkVM_StackPushPop(b *testing.B) {
	vm := NewVM()
	value := NewInt(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.push(value)
		_ = vm.pop()
	}
}

func BenchmarkVM_StackPeek(b *testing.B) {
	vm := NewVM()
	vm.push(NewInt(42))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vm.peek(0)
	}
}
