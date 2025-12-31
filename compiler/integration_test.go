package compiler

import (
	"testing"

	"pidgin-lang/lexer"
	"pidgin-lang/parser"
	"pidgin-lang/vm"
)

// Integration tests: Source → Lexer → Parser → Compiler → VM → Result

func compileAndRun(input string) (vm.Value, error) {
	// Parse
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	// Compile
	compiler := New()
	chunk, err := compiler.Compile(program)
	if err != nil {
		return vm.NewNothing(), err
	}

	// Execute
	vmachine := vm.NewVM()
	return vmachine.Run(chunk)
}

// ============================================================================
// Arithmetic Integration Tests
// ============================================================================

func TestIntegration_Arithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5 + 3", 8},
		{"10 - 4", 6},
		{"6 * 7", 42},
		{"20 / 4", 5},
		{"(5 + 3) * 2", 16},
		{"10 + 5 * 2", 20},
		{"-5", -5},
		{"-5 + 10", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsInt() {
				t.Fatalf("expected int result, got %s", result.TypeName())
			}

			if got := result.AsInt(); got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

// ============================================================================
// Comparison Integration Tests
// ============================================================================

func TestIntegration_Comparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 be 5", true},
		{"5 be 3", false},
		{"5 big pass 3", true},
		{"3 big pass 5", false},
		{"3 no reach 5", true},
		{"5 no reach 3", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsBool() {
				t.Fatalf("expected bool result, got %s", result.TypeName())
			}

			if got := result.AsBool(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// ============================================================================
// Variable Integration Tests
// ============================================================================

func TestIntegration_GlobalVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"make x be 42\nx", 42},
		{"make x be 5\nmake y be 10\nx + y", 15},
		{"make x be 10\nmake x be 20\nx", 20}, // Variable reassignment
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsInt() {
				t.Fatalf("expected int result, got %s", result.TypeName())
			}

			if got := result.AsInt(); got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

// ============================================================================
// String Integration Tests
// ============================================================================

func TestIntegration_Strings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Wetin dey happen?"`, "Wetin dey happen?"},
		{`"Hello" + " " + "World"`, "Hello World"},
		{"make name be \"Chidi\"\nname", "Chidi"}, // Use actual newline
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsString() {
				t.Fatalf("expected string result, got %s", result.TypeName())
			}

			if got := *result.AsString(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// ============================================================================
// Control Flow Integration Tests
// ============================================================================

func TestIntegration_SupposeExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"suppose 5 big pass 3 { 1 } abi { 0 }", 1},
		{"suppose 3 big pass 5 { 1 } abi { 0 }", 0},
		{"suppose tru { 42 } abi { 0 }", 42},
		{"suppose lie { 42 } abi { 99 }", 99},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsInt() {
				t.Fatalf("expected int result, got %s", result.TypeName())
			}

			if got := result.AsInt(); got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestIntegration_WhileLoop(t *testing.T) {
	input := `
	make counter be 0
	dey do while counter no reach 5 {
		make counter be counter + 1
	}
	counter
	`

	result, err := compileAndRun(input)
	if err != nil {
		t.Fatalf("execution error: %v", err)
	}

	if !result.IsInt() {
		t.Fatalf("expected int result, got %s", result.TypeName())
	}

	expected := int64(5)
	if got := result.AsInt(); got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

// ============================================================================
// Short-Circuit Integration Tests
// ============================================================================

func TestIntegration_ShortCircuitAnd(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"tru and tru", true},
		{"tru and lie", false},
		{"lie and tru", false},
		{"lie and lie", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compileAndRun(tt.input)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if !result.IsBool() {
				t.Fatalf("expected bool result, got %s", result.TypeName())
			}

			if got := result.AsBool(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// ============================================================================
// Complex Integration Tests
// ============================================================================

func TestIntegration_ComplexExpression(t *testing.T) {
	input := `
	make x be 10
	make y be 20
	suppose x no reach y {
		x * 2 + y
	} abi {
		x + y * 2
	}
	`

	result, err := compileAndRun(input)
	if err != nil {
		t.Fatalf("execution error: %v", err)
	}

	expected := int64(40) // (10 * 2) + 20 = 40
	if got := result.AsInt(); got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestIntegration_FizzBuzzPattern(t *testing.T) {
	// Simplified FizzBuzz pattern (without modulo)
	input := `
	make i be 0
	make count be 0
	dey do while i no reach 10 {
		make i be i + 1
		make count be count + 1
	}
	count
	`

	result, err := compileAndRun(input)
	if err != nil {
		t.Fatalf("execution error: %v", err)
	}

	expected := int64(10)
	if got := result.AsInt(); got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

// ============================================================================
// Error Handling Integration Tests
// ============================================================================

func TestIntegration_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"undefined variable", "undefined_var"},
		{"division by zero", "10 / 0"},
		{"type error subtract", "5 - tru"},
		{"type error multiply", "5 * lie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := compileAndRun(tt.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
