package vm

import (
	"fmt"
	"unsafe"
)

// NaN Boxing Value Representation
//
// We pack type tags and values into 64-bit words using IEEE 754 NaN space:
// IEEE 754 double-precision: Sign (1) | Exponent (11) | Mantissa (52)
//
// NaN values have exponent = 0x7FF and non-zero mantissa
// Quiet NaN: 0x7FF8_0000_0000_0000 to 0x7FFF_FFFF_FFFF_FFFF
//
// Our encoding uses the mantissa (52 bits) for:
// - Type tag (3 bits) in bits 48-50
// - Payload (48 bits) in bits 0-47

const (
	// Quiet NaN base with sign bit = 0
	QNAN_BASE = 0x7FF8000000000000

	// Type tags (3 bits = 8 types)
	TAG_INT     = 0 // 000 - 48-bit signed integer
	TAG_BOOL    = 1 // 001 - boolean (0 or 1 in low bit)
	TAG_NOTHING = 2 // 010 - null/nothing
	TAG_STRING  = 3 // 011 - pointer to string (48-bit address)
	TAG_FUNC    = 4 // 100 - pointer to function object
	TAG_BUILTIN = 5 // 101 - builtin function index
	TAG_ERROR   = 6 // 110 - pointer to error object
	// TAG_RESERVED = 7 for future use

	TAG_SHIFT    = 48
	TAG_MASK     = 0x7
	PAYLOAD_MASK = 0x0000FFFFFFFFFFFF

	// 48-bit integer limits
	MAX_INT_48 = 140737488355327  // 2^47 - 1
	MIN_INT_48 = -140737488355328 // -2^47
)

// Value is a NaN-boxed 64-bit value that can represent:
// integers, booleans, nothing, or pointers to heap objects
type Value uint64

// Type names for debugging and error messages
const (
	TypeInt     = "number"
	TypeBool    = "boolean"
	TypeNothing = "nothing"
	TypeString  = "string"
	TypeFunc    = "function"
	TypeBuiltin = "builtin"
	TypeError   = "error"
)

// ============================================================================
// Constructors
// ============================================================================

// NewInt creates a NaN-boxed integer value
func NewInt(i int64) Value {
	// Sign-extend to 48 bits
	i48 := uint64(i) & PAYLOAD_MASK
	return Value(QNAN_BASE | (TAG_INT << TAG_SHIFT) | i48)
}

// NewBool creates a NaN-boxed boolean value
func NewBool(b bool) Value {
	var val uint64
	if b {
		val = 1
	}
	return Value(QNAN_BASE | (TAG_BOOL << TAG_SHIFT) | val)
}

// NewNothing creates a NaN-boxed nothing/null value
func NewNothing() Value {
	return Value(QNAN_BASE | (TAG_NOTHING << TAG_SHIFT))
}

// NewString creates a NaN-boxed string value (stores pointer)
func NewString(s *string) Value {
	ptr := uint64(uintptr(unsafe.Pointer(s))) & PAYLOAD_MASK
	return Value(QNAN_BASE | (TAG_STRING << TAG_SHIFT) | ptr)
}

// NewFunc creates a NaN-boxed function value (stores pointer)
func NewFunc(f *Function) Value {
	ptr := uint64(uintptr(unsafe.Pointer(f))) & PAYLOAD_MASK
	return Value(QNAN_BASE | (TAG_FUNC << TAG_SHIFT) | ptr)
}

// NewBuiltin creates a NaN-boxed builtin function value (stores index)
func NewBuiltin(index int) Value {
	return Value(QNAN_BASE | (TAG_BUILTIN << TAG_SHIFT) | uint64(index))
}

// NewError creates a NaN-boxed error value (stores pointer)
func NewError(e *RuntimeError) Value {
	ptr := uint64(uintptr(unsafe.Pointer(e))) & PAYLOAD_MASK
	return Value(QNAN_BASE | (TAG_ERROR << TAG_SHIFT) | ptr)
}

// ============================================================================
// Type Checking
// ============================================================================

// GetTag extracts the type tag from a NaN-boxed value
func (v Value) GetTag() uint8 {
	return uint8((v >> TAG_SHIFT) & TAG_MASK)
}

// IsInt checks if the value is an integer
func (v Value) IsInt() bool {
	return v.GetTag() == TAG_INT
}

// IsBool checks if the value is a boolean
func (v Value) IsBool() bool {
	return v.GetTag() == TAG_BOOL
}

// IsNothing checks if the value is nothing/null
func (v Value) IsNothing() bool {
	return v.GetTag() == TAG_NOTHING
}

// IsString checks if the value is a string
func (v Value) IsString() bool {
	return v.GetTag() == TAG_STRING
}

// IsFunc checks if the value is a function
func (v Value) IsFunc() bool {
	return v.GetTag() == TAG_FUNC
}

// IsBuiltin checks if the value is a builtin function
func (v Value) IsBuiltin() bool {
	return v.GetTag() == TAG_BUILTIN
}

// IsError checks if the value is an error
func (v Value) IsError() bool {
	return v.GetTag() == TAG_ERROR
}

// ============================================================================
// Value Extraction
// ============================================================================

// AsInt extracts an integer value (sign-extends from 48-bit to 64-bit)
func (v Value) AsInt() int64 {
	// Extract 48-bit value
	i48 := int64(v & PAYLOAD_MASK)

	// Sign-extend if the sign bit (bit 47) is set
	if i48&0x800000000000 != 0 {
		return i48 | ^int64(PAYLOAD_MASK) // sign extend
	}
	return i48
}

// AsBool extracts a boolean value
func (v Value) AsBool() bool {
	return (v & 1) != 0
}

// AsString extracts a string pointer
func (v Value) AsString() *string {
	ptr := uintptr(v & PAYLOAD_MASK)
	return (*string)(unsafe.Pointer(ptr))
}

// AsFunc extracts a function pointer
func (v Value) AsFunc() *Function {
	ptr := uintptr(v & PAYLOAD_MASK)
	return (*Function)(unsafe.Pointer(ptr))
}

// AsBuiltin extracts a builtin function index
func (v Value) AsBuiltin() int {
	return int(v & PAYLOAD_MASK)
}

// AsError extracts an error pointer
func (v Value) AsError() *RuntimeError {
	ptr := uintptr(v & PAYLOAD_MASK)
	return (*RuntimeError)(unsafe.Pointer(ptr))
}

// ============================================================================
// Type Name
// ============================================================================

// TypeName returns the human-readable type name
func (v Value) TypeName() string {
	switch v.GetTag() {
	case TAG_INT:
		return TypeInt
	case TAG_BOOL:
		return TypeBool
	case TAG_NOTHING:
		return TypeNothing
	case TAG_STRING:
		return TypeString
	case TAG_FUNC:
		return TypeFunc
	case TAG_BUILTIN:
		return TypeBuiltin
	case TAG_ERROR:
		return TypeError
	default:
		return "unknown"
	}
}

// ============================================================================
// String Representation
// ============================================================================

// String returns a string representation of the value for debugging
func (v Value) String() string {
	switch v.GetTag() {
	case TAG_INT:
		return fmt.Sprintf("%d", v.AsInt())
	case TAG_BOOL:
		if v.AsBool() {
			return "tru"
		}
		return "lie"
	case TAG_NOTHING:
		return "nothing"
	case TAG_STRING:
		return *v.AsString()
	case TAG_FUNC:
		fn := v.AsFunc()
		if fn.Name != "" {
			return fmt.Sprintf("<function %s>", fn.Name)
		}
		return "<function>"
	case TAG_BUILTIN:
		return fmt.Sprintf("<builtin %d>", v.AsBuiltin())
	case TAG_ERROR:
		return v.AsError().Message
	default:
		return "<unknown>"
	}
}

// ============================================================================
// Truthiness (for conditionals)
// ============================================================================

// IsFalsey returns true if the value is "falsey" in Pidgin
// Only 'lie' (false) and 'nothing' are falsey
func (v Value) IsFalsey() bool {
	if v.IsBool() {
		return !v.AsBool()
	}
	if v.IsNothing() {
		return true
	}
	return false
}

// IsTruthy returns true if the value is "truthy" in Pidgin
func (v Value) IsTruthy() bool {
	return !v.IsFalsey()
}

// ============================================================================
// Equality
// ============================================================================

// Equals compares two values for equality
func (v Value) Equals(other Value) bool {
	// Fast path: bit-identical values
	if v == other {
		return true
	}

	// Different types are never equal
	if v.GetTag() != other.GetTag() {
		return false
	}

	// String comparison (compare actual string contents)
	if v.IsString() {
		return *v.AsString() == *other.AsString()
	}

	// For other types, bit equality is sufficient
	return false
}

// ============================================================================
// Helper Types
// ============================================================================

// Function represents a compiled function
type Function struct {
	Arity      int    // Number of parameters
	Chunk      *Chunk // Bytecode chunk
	Name       string // Function name (for debugging)
	LocalCount int    // Total local variables (including parameters)
}

// RuntimeError represents a runtime error
type RuntimeError struct {
	Message string
	Line    int
}
