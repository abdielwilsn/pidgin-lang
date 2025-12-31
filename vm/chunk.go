package vm

import (
	"fmt"
)

// Chunk represents a sequence of bytecode instructions with associated metadata
type Chunk struct {
	Code      []byte            // Bytecode instructions
	Constants []Value           // Constant pool (NaN-boxed values)
	Lines     []int             // Line numbers for each instruction (for error reporting)
	strings   map[string]*string // Interned strings for deduplication
}

// NewChunk creates a new empty chunk
func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0, 256),
		Constants: make([]Value, 0, 32),
		Lines:     make([]int, 0, 256),
		strings:   make(map[string]*string),
	}
}

// ============================================================================
// Code Generation
// ============================================================================

// WriteByte appends a byte to the chunk's code array
func (c *Chunk) WriteByte(b byte, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

// WriteOpcode appends an opcode to the chunk
func (c *Chunk) WriteOpcode(op Opcode, line int) {
	c.WriteByte(byte(op), line)
}

// WriteBytes appends multiple bytes to the chunk
func (c *Chunk) WriteBytes(bytes []byte, line int) {
	for _, b := range bytes {
		c.WriteByte(b, line)
	}
}

// Count returns the number of bytes in the chunk
func (c *Chunk) Count() int {
	return len(c.Code)
}

// ============================================================================
// Constants Pool
// ============================================================================

// AddConstant adds a value to the constant pool and returns its index
// Deduplicates identical values to save memory
func (c *Chunk) AddConstant(value Value) int {
	// Check if we already have this constant
	for i, existing := range c.Constants {
		if existing.Equals(value) {
			return i
		}
	}

	// Add new constant
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// GetConstant retrieves a constant by index
func (c *Chunk) GetConstant(index int) Value {
	if index < 0 || index >= len(c.Constants) {
		return NewNothing()
	}
	return c.Constants[index]
}

// ============================================================================
// String Interning
// ============================================================================

// InternString returns a pointer to an interned string
// This ensures that identical strings share the same memory location
func (c *Chunk) InternString(s string) *string {
	if existing, ok := c.strings[s]; ok {
		return existing
	}

	ptr := new(string)
	*ptr = s
	c.strings[s] = ptr
	return ptr
}

// ============================================================================
// Line Number Tracking
// ============================================================================

// GetLine returns the source line number for a bytecode offset
func (c *Chunk) GetLine(offset int) int {
	if offset < 0 || offset >= len(c.Lines) {
		return 0
	}
	return c.Lines[offset]
}

// ============================================================================
// Disassembly (for debugging)
// ============================================================================

// Disassemble prints the entire chunk with human-readable instruction names
func (c *Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < len(c.Code); {
		offset = c.DisassembleInstruction(offset)
	}
}

// DisassembleInstruction prints a single instruction and returns the next offset
func (c *Chunk) DisassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)

	// Print line number (with run-length encoding for readability)
	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", c.Lines[offset])
	}

	instruction := Opcode(c.Code[offset])

	switch instruction {
	// Simple instructions (no operands)
	case OP_CONST_0, OP_CONST_1, OP_CONST_MINUS1,
		OP_NOTHING, OP_TRU, OP_LIE,
		OP_ADD, OP_SUB, OP_MUL, OP_DIV, OP_NEGATE,
		OP_EQUAL, OP_NOT_EQUAL, OP_GREATER, OP_LESS, OP_NOT,
		OP_GET_LOCAL_0, OP_GET_LOCAL_1, OP_GET_LOCAL_2, OP_GET_LOCAL_3,
		OP_SET_LOCAL_0, OP_SET_LOCAL_1,
		OP_CALL_0, OP_CALL_1, OP_CALL_2,
		OP_RETURN, OP_BRING,
		OP_POP, OP_DUP, OP_CONCAT, OP_HALT:
		return c.simpleInstruction(instruction, offset)

	// Byte operand instructions
	case OP_CONST_I8:
		return c.byteInstruction(instruction, offset)

	// Short operand instructions
	case OP_CONST_I16:
		return c.shortInstruction(instruction, offset)

	// Constant instructions (2-byte index into constants pool)
	case OP_CONSTANT:
		return c.constantInstruction(instruction, offset)

	// Local variable instructions (1-byte slot)
	case OP_GET_LOCAL, OP_SET_LOCAL:
		return c.byteInstruction(instruction, offset)

	// Global variable instructions (2-byte index)
	case OP_GET_GLOBAL, OP_SET_GLOBAL:
		return c.shortInstruction(instruction, offset)

	// Jump instructions (2-byte offset)
	case OP_JUMP, OP_JUMP_IF_LIE, OP_JUMP_IF_TRU:
		return c.jumpInstruction(instruction, 1, offset)

	// Loop instruction (2-byte backward offset)
	case OP_LOOP:
		return c.jumpInstruction(instruction, -1, offset)

	// Call instructions with argument count
	case OP_CALL, OP_YARN:
		return c.byteInstruction(instruction, offset)

	// Closure instruction
	case OP_CLOSURE:
		return c.shortInstruction(instruction, offset)

	// Builtin instruction (builtin index + arg count)
	case OP_BUILTIN:
		offset++
		builtinIdx := c.Code[offset]
		offset++
		argCount := c.Code[offset]
		fmt.Printf("%-16s %4d (args: %d)\n", instruction.String(), builtinIdx, argCount)
		return offset + 1

	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

// simpleInstruction disassembles instructions with no operands
func (c *Chunk) simpleInstruction(op Opcode, offset int) int {
	fmt.Printf("%s\n", op.String())
	return offset + 1
}

// byteInstruction disassembles instructions with a 1-byte operand
func (c *Chunk) byteInstruction(op Opcode, offset int) int {
	slot := c.Code[offset+1]
	fmt.Printf("%-16s %4d\n", op.String(), slot)
	return offset + 2
}

// shortInstruction disassembles instructions with a 2-byte operand
func (c *Chunk) shortInstruction(op Opcode, offset int) int {
	value := uint16(c.Code[offset+1])<<8 | uint16(c.Code[offset+2])
	fmt.Printf("%-16s %4d\n", op.String(), value)
	return offset + 3
}

// constantInstruction disassembles OP_CONSTANT (shows the constant value)
func (c *Chunk) constantInstruction(op Opcode, offset int) int {
	constantIdx := uint16(c.Code[offset+1])<<8 | uint16(c.Code[offset+2])
	fmt.Printf("%-16s %4d '", op.String(), constantIdx)

	if int(constantIdx) < len(c.Constants) {
		constant := c.Constants[constantIdx]
		fmt.Printf("%s", constant.String())
	}

	fmt.Printf("'\n")
	return offset + 3
}

// jumpInstruction disassembles jump instructions
func (c *Chunk) jumpInstruction(op Opcode, sign int, offset int) int {
	jump := int16(uint16(c.Code[offset+1])<<8 | uint16(c.Code[offset+2]))
	target := offset + 3 + sign*int(jump)
	fmt.Printf("%-16s %4d -> %d\n", op.String(), offset, target)
	return offset + 3
}
