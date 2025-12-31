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

// These benchmarks measure EXECUTION time only (not parsing/compilation)

// ============================================================================
// Arithmetic Execution Benchmarks
// ============================================================================

func BenchmarkExecution_Arithmetic_Interpreter(b *testing.B) {
	input := "5 + 3"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_Arithmetic_VM(b *testing.B) {
	input := "5 + 3"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)
	// VM reused

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vmachine.Run(chunk)
	}
}

func BenchmarkExecution_ComplexArithmetic_Interpreter(b *testing.B) {
	input := "(5 + 3) * (10 - 2) / 4"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_ComplexArithmetic_VM(b *testing.B) {
	input := "(5 + 3) * (10 - 2) / 4"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)
	// VM reused

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vmachine.Run(chunk)
	}
}

// ============================================================================
// Loop Execution Benchmarks
// ============================================================================

func BenchmarkExecution_Loop10_Interpreter(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 10 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := object.NewEnvironment()
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_Loop10_VM(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 10 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VM reused
		_, _ = vmachine.Run(chunk)
	}
}

func BenchmarkExecution_Loop100_Interpreter(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 100 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := object.NewEnvironment()
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_Loop100_VM(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 100 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VM reused
		_, _ = vmachine.Run(chunk)
	}
}

func BenchmarkExecution_Loop1000_Interpreter(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 1000 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := object.NewEnvironment()
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_Loop1000_VM(b *testing.B) {
	input := `
	make counter be 0
	dey do while counter no reach 1000 {
		make counter be counter + 1
	}
	counter
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VM reused
		_, _ = vmachine.Run(chunk)
	}
}

// ============================================================================
// Variable Access Benchmarks
// ============================================================================

func BenchmarkExecution_VariableAccess_Interpreter(b *testing.B) {
	input := `
	make x be 42
	make y be 10
	x + y
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := object.NewEnvironment()
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_VariableAccess_VM(b *testing.B) {
	input := `
	make x be 42
	make y be 10
	x + y
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VM reused
		_, _ = vmachine.Run(chunk)
	}
}

// ============================================================================
// Nested Operations Benchmark
// ============================================================================

func BenchmarkExecution_NestedOperations_Interpreter(b *testing.B) {
	input := `
	make x be 10
	make y be 20
	suppose x no reach y {
		x * 2 + y
	} abi {
		x + y * 2
	}
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := object.NewEnvironment()
		_ = evaluator.Eval(program, env)
	}
}

func BenchmarkExecution_NestedOperations_VM(b *testing.B) {
	input := `
	make x be 10
	make y be 20
	suppose x no reach y {
		x * 2 + y
	} abi {
		x + y * 2
	}
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.New()
	chunk, _ := c.Compile(program)

	vmachine := vm.NewVM()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VM reused
		_, _ = vmachine.Run(chunk)
	}
}
