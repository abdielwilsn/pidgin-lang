package vm

import (
	"testing"
)

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNewInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{"zero", 0, 0},
		{"positive small", 42, 42},
		{"negative small", -42, -42},
		{"positive large", 1000000, 1000000},
		{"negative large", -1000000, -1000000},
		{"max 48-bit", MAX_INT_48, MAX_INT_48},
		{"min 48-bit", MIN_INT_48, MIN_INT_48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewInt(tt.input)
			if !v.IsInt() {
				t.Errorf("NewInt(%d).IsInt() = false, want true", tt.input)
			}
			if got := v.AsInt(); got != tt.expected {
				t.Errorf("NewInt(%d).AsInt() = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNewBool(t *testing.T) {
	trueVal := NewBool(true)
	if !trueVal.IsBool() {
		t.Error("NewBool(true).IsBool() = false, want true")
	}
	if !trueVal.AsBool() {
		t.Error("NewBool(true).AsBool() = false, want true")
	}

	falseVal := NewBool(false)
	if !falseVal.IsBool() {
		t.Error("NewBool(false).IsBool() = false, want true")
	}
	if falseVal.AsBool() {
		t.Error("NewBool(false).AsBool() = true, want false")
	}
}

func TestNewNothing(t *testing.T) {
	v := NewNothing()
	if !v.IsNothing() {
		t.Error("NewNothing().IsNothing() = false, want true")
	}
}

func TestNewString(t *testing.T) {
	str := "Wetin dey happen?"
	v := NewString(&str)

	if !v.IsString() {
		t.Error("NewString().IsString() = false, want true")
	}

	retrieved := v.AsString()
	if retrieved == nil {
		t.Fatal("AsString() returned nil")
	}
	if *retrieved != str {
		t.Errorf("AsString() = %q, want %q", *retrieved, str)
	}
}

// ============================================================================
// Type Tag Tests
// ============================================================================

func TestGetTag(t *testing.T) {
	tests := []struct {
		name        string
		value       Value
		expectedTag uint8
	}{
		{"int", NewInt(42), TAG_INT},
		{"bool true", NewBool(true), TAG_BOOL},
		{"bool false", NewBool(false), TAG_BOOL},
		{"nothing", NewNothing(), TAG_NOTHING},
		{"string", NewString(stringPtr("test")), TAG_STRING},
		{"builtin", NewBuiltin(0), TAG_BUILTIN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.GetTag(); got != tt.expectedTag {
				t.Errorf("GetTag() = %d, want %d", got, tt.expectedTag)
			}
		})
	}
}

// ============================================================================
// Type Checking Tests
// ============================================================================

func TestTypeChecking(t *testing.T) {
	intVal := NewInt(42)
	boolVal := NewBool(true)
	nothingVal := NewNothing()
	strVal := NewString(stringPtr("test"))
	builtinVal := NewBuiltin(0)

	// Test IsInt
	if !intVal.IsInt() {
		t.Error("int value: IsInt() = false, want true")
	}
	if boolVal.IsInt() {
		t.Error("bool value: IsInt() = true, want false")
	}

	// Test IsBool
	if !boolVal.IsBool() {
		t.Error("bool value: IsBool() = false, want true")
	}
	if intVal.IsBool() {
		t.Error("int value: IsBool() = true, want false")
	}

	// Test IsNothing
	if !nothingVal.IsNothing() {
		t.Error("nothing value: IsNothing() = false, want true")
	}
	if intVal.IsNothing() {
		t.Error("int value: IsNothing() = true, want false")
	}

	// Test IsString
	if !strVal.IsString() {
		t.Error("string value: IsString() = false, want true")
	}
	if intVal.IsString() {
		t.Error("int value: IsString() = true, want false")
	}

	// Test IsBuiltin
	if !builtinVal.IsBuiltin() {
		t.Error("builtin value: IsBuiltin() = false, want true")
	}
	if intVal.IsBuiltin() {
		t.Error("int value: IsBuiltin() = true, want false")
	}
}

// ============================================================================
// Value Extraction Tests
// ============================================================================

func TestAsInt(t *testing.T) {
	tests := []int64{0, 1, -1, 42, -42, 12345, -12345, MAX_INT_48, MIN_INT_48}

	for _, expected := range tests {
		v := NewInt(expected)
		if got := v.AsInt(); got != expected {
			t.Errorf("NewInt(%d).AsInt() = %d, want %d", expected, got, expected)
		}
	}
}

func TestSignExtension(t *testing.T) {
	// Test that negative numbers are properly sign-extended
	negatives := []int64{-1, -42, -1000, -1000000, MIN_INT_48}

	for _, n := range negatives {
		v := NewInt(n)
		got := v.AsInt()
		if got != n {
			t.Errorf("Sign extension failed for %d: got %d", n, got)
		}
		if got >= 0 {
			t.Errorf("Negative number %d became positive: %d", n, got)
		}
	}
}

func TestAsBool(t *testing.T) {
	trueVal := NewBool(true)
	falseVal := NewBool(false)

	if !trueVal.AsBool() {
		t.Error("true value: AsBool() = false, want true")
	}
	if falseVal.AsBool() {
		t.Error("false value: AsBool() = true, want false")
	}
}

// ============================================================================
// Type Name Tests
// ============================================================================

func TestTypeName(t *testing.T) {
	tests := []struct {
		value    Value
		expected string
	}{
		{NewInt(42), TypeInt},
		{NewBool(true), TypeBool},
		{NewNothing(), TypeNothing},
		{NewString(stringPtr("test")), TypeString},
		{NewBuiltin(0), TypeBuiltin},
	}

	for _, tt := range tests {
		if got := tt.value.TypeName(); got != tt.expected {
			t.Errorf("TypeName() = %q, want %q", got, tt.expected)
		}
	}
}

// ============================================================================
// String Representation Tests
// ============================================================================

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"int zero", NewInt(0), "0"},
		{"int positive", NewInt(42), "42"},
		{"int negative", NewInt(-42), "-42"},
		{"bool true", NewBool(true), "tru"},
		{"bool false", NewBool(false), "lie"},
		{"nothing", NewNothing(), "nothing"},
		{"string", NewString(stringPtr("Abeg, wetin dey happen?")), "Abeg, wetin dey happen?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// Truthiness Tests
// ============================================================================

func TestIsFalsey(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected bool
	}{
		{"bool false", NewBool(false), true},
		{"nothing", NewNothing(), true},
		{"bool true", NewBool(true), false},
		{"int zero", NewInt(0), false},
		{"int positive", NewInt(42), false},
		{"string empty", NewString(stringPtr("")), false},
		{"string non-empty", NewString(stringPtr("test")), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsFalsey(); got != tt.expected {
				t.Errorf("IsFalsey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected bool
	}{
		{"bool true", NewBool(true), true},
		{"int zero", NewInt(0), true},
		{"int positive", NewInt(42), true},
		{"string", NewString(stringPtr("test")), true},
		{"bool false", NewBool(false), false},
		{"nothing", NewNothing(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsTruthy(); got != tt.expected {
				t.Errorf("IsTruthy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// Equality Tests
// ============================================================================

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		a        Value
		b        Value
		expected bool
	}{
		{"int equal", NewInt(42), NewInt(42), true},
		{"int not equal", NewInt(42), NewInt(43), false},
		{"bool equal true", NewBool(true), NewBool(true), true},
		{"bool equal false", NewBool(false), NewBool(false), true},
		{"bool not equal", NewBool(true), NewBool(false), false},
		{"nothing equal", NewNothing(), NewNothing(), true},
		{"different types int/bool", NewInt(1), NewBool(true), false},
		{"different types int/nothing", NewInt(0), NewNothing(), false},
		{"string equal", NewString(stringPtr("test")), NewString(stringPtr("test")), true},
		{"string not equal", NewString(stringPtr("test")), NewString(stringPtr("other")), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.expected {
				t.Errorf("Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// Edge Cases and Boundary Tests
// ============================================================================

func TestInt48BitLimits(t *testing.T) {
	// Test maximum 48-bit value
	maxVal := NewInt(MAX_INT_48)
	if got := maxVal.AsInt(); got != MAX_INT_48 {
		t.Errorf("MAX_INT_48: got %d, want %d", got, MAX_INT_48)
	}

	// Test minimum 48-bit value
	minVal := NewInt(MIN_INT_48)
	if got := minVal.AsInt(); got != MIN_INT_48 {
		t.Errorf("MIN_INT_48: got %d, want %d", got, MIN_INT_48)
	}
}

func TestMultipleStringPointers(t *testing.T) {
	// Test that different string pointers are handled correctly
	s1 := "test1"
	s2 := "test2"

	v1 := NewString(&s1)
	v2 := NewString(&s2)

	if v1.Equals(v2) {
		t.Error("Different string pointers should not be equal")
	}

	// Same string content but different pointers
	s3 := "test1"
	v3 := NewString(&s3)

	if !v1.Equals(v3) {
		t.Error("Same string content should be equal")
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func stringPtr(s string) *string {
	return &s
}
