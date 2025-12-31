package vm

// Opcode represents a single bytecode instruction
type Opcode byte

const (
	// ========================================================================
	// Literals and Constants (0-9)
	// ========================================================================

	OP_CONST_0      Opcode = 0 // Push constant 0
	OP_CONST_1      Opcode = 1 // Push constant 1
	OP_CONST_MINUS1 Opcode = 2 // Push constant -1
	OP_NOTHING      Opcode = 3 // Push nothing
	OP_TRU          Opcode = 4 // Push true
	OP_LIE          Opcode = 5 // Push false
	OP_CONSTANT     Opcode = 6 // Push constant from pool: [u16 index]
	OP_CONST_I8     Opcode = 7 // Push 8-bit int: [i8]
	OP_CONST_I16    Opcode = 8 // Push 16-bit int: [i16]

	// ========================================================================
	// Arithmetic (10-19)
	// ========================================================================

	OP_ADD    Opcode = 10 // a + b (also string concatenation)
	OP_SUB    Opcode = 11 // a - b
	OP_MUL    Opcode = 12 // a * b
	OP_DIV    Opcode = 13 // a / b
	OP_NEGATE Opcode = 14 // -a

	// ========================================================================
	// Comparison (20-29)
	// ========================================================================

	OP_EQUAL     Opcode = 20 // a be b (equality)
	OP_NOT_EQUAL Opcode = 21 // a no be b (inequality)
	OP_GREATER   Opcode = 22 // a big pass b (greater than)
	OP_LESS      Opcode = 23 // a no reach b (less than)

	// ========================================================================
	// Logical (30-34)
	// ========================================================================

	OP_NOT Opcode = 30 // !a or no be a (logical not)

	// Note: AND and OR are implemented via short-circuit jumps in the compiler
	// OP_AND and OP_OR are reserved but not used in the initial implementation

	// ========================================================================
	// Variables (35-44)
	// ========================================================================

	OP_GET_LOCAL_0 Opcode = 35 // Get local var at slot 0
	OP_GET_LOCAL_1 Opcode = 36 // Get local var at slot 1
	OP_GET_LOCAL_2 Opcode = 37 // Get local var at slot 2
	OP_GET_LOCAL_3 Opcode = 38 // Get local var at slot 3
	OP_GET_LOCAL   Opcode = 39 // Get local var: [u8 slot]
	OP_SET_LOCAL_0 Opcode = 40 // Set local var at slot 0
	OP_SET_LOCAL_1 Opcode = 41 // Set local var at slot 1
	OP_SET_LOCAL   Opcode = 42 // Set local var: [u8 slot]
	OP_GET_GLOBAL  Opcode = 43 // Get global var: [u16 index]
	OP_SET_GLOBAL  Opcode = 44 // Set global var: [u16 index]

	// ========================================================================
	// Control Flow (45-54)
	// ========================================================================

	OP_JUMP        Opcode = 45 // Unconditional jump: [i16 offset]
	OP_JUMP_IF_LIE Opcode = 46 // Jump if false: [i16 offset]
	OP_JUMP_IF_TRU Opcode = 47 // Jump if true: [i16 offset]
	OP_LOOP        Opcode = 48 // Jump backward: [u16 offset]

	// ========================================================================
	// Functions (55-64)
	// ========================================================================

	OP_CALL_0  Opcode = 55 // Call function with 0 args
	OP_CALL_1  Opcode = 56 // Call function with 1 arg
	OP_CALL_2  Opcode = 57 // Call function with 2 args
	OP_CALL    Opcode = 58 // Call function: [u8 argCount]
	OP_CLOSURE Opcode = 59 // Create closure: [u16 funcIndex]
	OP_RETURN  Opcode = 60 // Return from function
	OP_BRING   Opcode = 61 // Return value (Pidgin's 'bring')

	// ========================================================================
	// Builtins (65-74)
	// ========================================================================

	OP_YARN    Opcode = 65 // Optimized print: [u8 argCount]
	OP_BUILTIN Opcode = 66 // Call builtin: [u8 builtinIndex, u8 argCount]

	// ========================================================================
	// Stack Manipulation (75-79)
	// ========================================================================

	OP_POP Opcode = 75 // Pop and discard top value
	OP_DUP Opcode = 76 // Duplicate top value

	// ========================================================================
	// String Operations (80-84)
	// ========================================================================

	OP_CONCAT Opcode = 80 // String concatenation (polymorphic)

	// ========================================================================
	// Special (85-89)
	// ========================================================================

	OP_HALT Opcode = 85 // Stop execution
)

// OpcodeNames maps opcodes to their string names for debugging
var OpcodeNames = map[Opcode]string{
	// Literals and Constants
	OP_CONST_0:      "OP_CONST_0",
	OP_CONST_1:      "OP_CONST_1",
	OP_CONST_MINUS1: "OP_CONST_MINUS1",
	OP_NOTHING:      "OP_NOTHING",
	OP_TRU:          "OP_TRU",
	OP_LIE:          "OP_LIE",
	OP_CONSTANT:     "OP_CONSTANT",
	OP_CONST_I8:     "OP_CONST_I8",
	OP_CONST_I16:    "OP_CONST_I16",

	// Arithmetic
	OP_ADD:    "OP_ADD",
	OP_SUB:    "OP_SUB",
	OP_MUL:    "OP_MUL",
	OP_DIV:    "OP_DIV",
	OP_NEGATE: "OP_NEGATE",

	// Comparison
	OP_EQUAL:     "OP_EQUAL",
	OP_NOT_EQUAL: "OP_NOT_EQUAL",
	OP_GREATER:   "OP_GREATER",
	OP_LESS:      "OP_LESS",

	// Logical
	OP_NOT: "OP_NOT",

	// Variables
	OP_GET_LOCAL_0: "OP_GET_LOCAL_0",
	OP_GET_LOCAL_1: "OP_GET_LOCAL_1",
	OP_GET_LOCAL_2: "OP_GET_LOCAL_2",
	OP_GET_LOCAL_3: "OP_GET_LOCAL_3",
	OP_GET_LOCAL:   "OP_GET_LOCAL",
	OP_SET_LOCAL_0: "OP_SET_LOCAL_0",
	OP_SET_LOCAL_1: "OP_SET_LOCAL_1",
	OP_SET_LOCAL:   "OP_SET_LOCAL",
	OP_GET_GLOBAL:  "OP_GET_GLOBAL",
	OP_SET_GLOBAL:  "OP_SET_GLOBAL",

	// Control Flow
	OP_JUMP:        "OP_JUMP",
	OP_JUMP_IF_LIE: "OP_JUMP_IF_LIE",
	OP_JUMP_IF_TRU: "OP_JUMP_IF_TRU",
	OP_LOOP:        "OP_LOOP",

	// Functions
	OP_CALL_0:  "OP_CALL_0",
	OP_CALL_1:  "OP_CALL_1",
	OP_CALL_2:  "OP_CALL_2",
	OP_CALL:    "OP_CALL",
	OP_CLOSURE: "OP_CLOSURE",
	OP_RETURN:  "OP_RETURN",
	OP_BRING:   "OP_BRING",

	// Builtins
	OP_YARN:    "OP_YARN",
	OP_BUILTIN: "OP_BUILTIN",

	// Stack Manipulation
	OP_POP: "OP_POP",
	OP_DUP: "OP_DUP",

	// String Operations
	OP_CONCAT: "OP_CONCAT",

	// Special
	OP_HALT: "OP_HALT",
}

// String returns the name of the opcode
func (op Opcode) String() string {
	if name, ok := OpcodeNames[op]; ok {
		return name
	}
	return "UNKNOWN"
}

// OpcodeOperandCount returns the number of operand bytes for an opcode
var OpcodeOperandCount = map[Opcode]int{
	// 0 operands
	OP_CONST_0:      0,
	OP_CONST_1:      0,
	OP_CONST_MINUS1: 0,
	OP_NOTHING:      0,
	OP_TRU:          0,
	OP_LIE:          0,
	OP_ADD:          0,
	OP_SUB:          0,
	OP_MUL:          0,
	OP_DIV:          0,
	OP_NEGATE:       0,
	OP_EQUAL:        0,
	OP_NOT_EQUAL:    0,
	OP_GREATER:      0,
	OP_LESS:         0,
	OP_NOT:          0,
	OP_GET_LOCAL_0:  0,
	OP_GET_LOCAL_1:  0,
	OP_GET_LOCAL_2:  0,
	OP_GET_LOCAL_3:  0,
	OP_SET_LOCAL_0:  0,
	OP_SET_LOCAL_1:  0,
	OP_CALL_0:       0,
	OP_CALL_1:       0,
	OP_CALL_2:       0,
	OP_RETURN:       0,
	OP_BRING:        0,
	OP_POP:          0,
	OP_DUP:          0,
	OP_CONCAT:       0,
	OP_HALT:         0,

	// 1 byte operand
	OP_CONST_I8:    1,
	OP_GET_LOCAL:   1,
	OP_SET_LOCAL:   1,
	OP_CALL:        1,
	OP_YARN:        1,

	// 2 byte operands
	OP_CONST_I16:   2,
	OP_CONSTANT:    2,
	OP_GET_GLOBAL:  2,
	OP_SET_GLOBAL:  2,
	OP_JUMP:        2,
	OP_JUMP_IF_LIE: 2,
	OP_JUMP_IF_TRU: 2,
	OP_LOOP:        2,
	OP_CLOSURE:     2,
	OP_BUILTIN:     2, // builtinIndex + argCount
}

// GetOperandCount returns the number of operand bytes for an opcode
func (op Opcode) GetOperandCount() int {
	if count, ok := OpcodeOperandCount[op]; ok {
		return count
	}
	return 0
}
