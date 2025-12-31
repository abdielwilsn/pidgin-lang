package main

import (
	"testing"

	"pidgin-lang/compiler"
	"pidgin-lang/evaluator"
	"pidgin-lang/lexer"
	"pidgin-lang/object"
	"pidgin-lang/parser"
	"pidgin-lang/vm"
)

// Benchmark helpers

func parseSource(input string) *parser.Parser {
	l := lexer.New(input)
	return parser.New(l)
}

func runInterpreter(input string) (object.Object, error) {
	p := parseSource(input)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return evaluator.Eval(program, env), nil
}

func runVM(input string) (vm.Value, error) {
	p := parseSource(input)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, err := c.Compile(program)
	if err != nil {
		return vm.NewNothing(), err
	}

	vmachine := vm.NewVM()
	return vmachine.Run(chunk)
}

// ============================================================================
// Arithmetic Benchmarks
// ============================================================================

func BenchmarkArithmetic_Interpreter(b *testing.B) {
	input := "5 + 3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkArithmetic_VM(b *testing.B) {
	input := "5 + 3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

func BenchmarkComplexArithmetic_Interpreter(b *testing.B) {
	input := "(5 + 3) * (10 - 2) / 4"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkComplexArithmetic_VM(b *testing.B) {
	input := "(5 + 3) * (10 - 2) / 4"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Variable Access Benchmarks
// ============================================================================

func BenchmarkVariableAccess_Interpreter(b *testing.B) {
	input := `
	make x be 42
	make y be 10
	x + y
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkVariableAccess_VM(b *testing.B) {
	input := `
	make x be 42
	make y be 10
	x + y
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Control Flow Benchmarks
// ============================================================================

func BenchmarkSupposeExpression_Interpreter(b *testing.B) {
	input := `suppose 5 big pass 3 { 42 } abi { 0 }`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkSupposeExpression_VM(b *testing.B) {
	input := `suppose 5 big pass 3 { 42 } abi { 0 }`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Loop Benchmarks
// ============================================================================

func BenchmarkLoop10_Interpreter(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 10 {
		make counter be counter + 1
	}
	counter
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkLoop10_VM(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 10 {
		make counter be counter + 1
	}
	counter
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

func BenchmarkLoop100_Interpreter(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 100 {
		make counter be counter + 1
	}
	counter
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkLoop100_VM(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 100 {
		make counter be counter + 1
	}
	counter
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// String Benchmarks
// ============================================================================

func BenchmarkStringConcat_Interpreter(b *testing.B) {
	input := `"Hello" + " " + "World"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkStringConcat_VM(b *testing.B) {
	input := `"Hello" + " " + "World"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Nested Operations Benchmark
// ============================================================================

func BenchmarkNestedOperations_Interpreter(b *testing.B) {
	input := `
	make x be 10
	make y be 20
	suppose x no reach y {
		x * 2 + y
	} abi {
		x + y * 2
	}
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkNestedOperations_VM(b *testing.B) {
	input := `
	make x be 10
	make y be 20
	suppose x no reach y {
		x * 2 + y
	} abi {
		x + y * 2
	}
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Comparison Benchmarks
// ============================================================================

func BenchmarkComparison_Interpreter(b *testing.B) {
	input := `5 big pass 3`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkComparison_VM(b *testing.B) {
	input := `5 big pass 3`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}

// ============================================================================
// Mixed Operations Benchmark
// ============================================================================

func BenchmarkMixedOperations_Interpreter(b *testing.B) {
	input := `
	make a be 5
	make b be 10
	make c be 15
	suppose a no reach b and b no reach c {
		(a + b) * c
	} abi {
		a * b + c
	}
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runInterpreter(input)
	}
}

func BenchmarkMixedOperations_VM(b *testing.B) {
	input := `
	make a be 5
	make b be 10
	make c be 15
	suppose a no reach b and b no reach c {
		(a + b) * c
	} abi {
		a * b + c
	}
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = runVM(input)
	}
}
