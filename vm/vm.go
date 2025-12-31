package vm

import (
	"fmt"
	"strings"
)

const (
	STACK_MAX  = 65536 // Maximum stack depth
	FRAMES_MAX = 1024  // Maximum call depth
)

// VM represents the virtual machine that executes bytecode
type VM struct {
	// Value stack (pre-allocated, no growth)
	stack    [STACK_MAX]Value
	stackTop int // Points to next free slot

	// Call frames (pre-allocated)
	frames     [FRAMES_MAX]CallFrame
	frameCount int

	// Global variables
	globals map[string]Value

	// Current chunk being executed (for single-chunk execution)
	chunk *Chunk

	// Instruction pointer (for single-chunk execution without call frames)
	ip int
}

// CallFrame represents a single function call on the call stack
type CallFrame struct {
	function *Function // Function being executed
	ip       int       // Instruction pointer for this frame
	slots    int       // Base pointer: where this frame's locals start on stack
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	return &VM{
		globals:    make(map[string]Value),
		stackTop:   0,
		frameCount: 0,
	}
}

// Reset clears the VM state for reuse
func (vm *VM) Reset() {
	vm.stackTop = 0
	vm.frameCount = 0
	vm.ip = 0
}

// ============================================================================
// Stack Operations
// ============================================================================

// push adds a value to the top of the stack
func (vm *VM) push(value Value) {
	if vm.stackTop >= STACK_MAX {
		panic("Stack overflow")
	}
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

// pop removes and returns the value from the top of the stack
func (vm *VM) pop() Value {
	if vm.stackTop <= 0 {
		panic("Stack underflow")
	}
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

// peek returns the value at distance from the top of the stack without removing it
func (vm *VM) peek(distance int) Value {
	return vm.stack[vm.stackTop-1-distance]
}

// ============================================================================
// Execution
// ============================================================================

// Run executes bytecode in a chunk and returns the result
func (vm *VM) Run(chunk *Chunk) (Value, error) {
	vm.Reset()
	vm.chunk = chunk
	vm.ip = 0

	return vm.execute()
}

// execute is the main execution loop with direct threading dispatch
func (vm *VM) execute() (Value, error) {
	// Cache hot values in local variables (Go compiler may map these to registers)
	var (
		stackTop = vm.stackTop
		ip       = vm.ip
		code     = vm.chunk.Code
		a, b     Value // For binary operations
	)

	// Inline helper to read next byte
	readByte := func() byte {
		b := code[ip]
		ip++
		return b
	}

	readShort := func() uint16 {
		b1 := uint16(code[ip])
		b2 := uint16(code[ip+1])
		ip += 2
		return (b1 << 8) | b2
	}

	// Main dispatch loop
dispatch:
	for {
		// Fetch instruction
		instruction := Opcode(readByte())

		// Dispatch
		switch instruction {

		// ====================================================================
		// Literals and Constants
		// ====================================================================

		case OP_CONST_0:
			vm.stack[stackTop] = NewInt(0)
			stackTop++
			goto dispatch

		case OP_CONST_1:
			vm.stack[stackTop] = NewInt(1)
			stackTop++
			goto dispatch

		case OP_CONST_MINUS1:
			vm.stack[stackTop] = NewInt(-1)
			stackTop++
			goto dispatch

		case OP_NOTHING:
			vm.stack[stackTop] = NewNothing()
			stackTop++
			goto dispatch

		case OP_TRU:
			vm.stack[stackTop] = NewBool(true)
			stackTop++
			goto dispatch

		case OP_LIE:
			vm.stack[stackTop] = NewBool(false)
			stackTop++
			goto dispatch

		case OP_CONST_I8:
			val := int8(readByte())
			vm.stack[stackTop] = NewInt(int64(val))
			stackTop++
			goto dispatch

		case OP_CONST_I16:
			val := int16(readShort())
			vm.stack[stackTop] = NewInt(int64(val))
			stackTop++
			goto dispatch

		case OP_CONSTANT:
			idx := readShort()
			vm.stack[stackTop] = vm.chunk.Constants[idx]
			stackTop++
			goto dispatch

		// ====================================================================
		// Arithmetic
		// ====================================================================

		case OP_ADD:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			// Fast path for integers (most common case)
			if a.IsInt() && b.IsInt() {
				result := a.AsInt() + b.AsInt()
				vm.stack[stackTop] = NewInt(result)
				stackTop++
				goto dispatch
			}

			// Slow path for string concatenation
			if a.IsString() || b.IsString() {
				strA := vm.valueToString(a)
				strB := vm.valueToString(b)
				result := strA + strB
				interned := vm.chunk.InternString(result)
				vm.stack[stackTop] = NewString(interned)
				stackTop++
				goto dispatch
			}

			// Type error
			vm.stackTop = stackTop
			vm.ip = ip
			return NewNothing(), vm.runtimeError(
				"I no fit add %s and %s", a.TypeName(), b.TypeName(),
			)

		case OP_SUB:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			if !a.IsInt() || !b.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit subtract %s and %s", a.TypeName(), b.TypeName(),
				)
			}

			vm.stack[stackTop] = NewInt(a.AsInt() - b.AsInt())
			stackTop++
			goto dispatch

		case OP_MUL:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			if !a.IsInt() || !b.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit multiply %s and %s", a.TypeName(), b.TypeName(),
				)
			}

			vm.stack[stackTop] = NewInt(a.AsInt() * b.AsInt())
			stackTop++
			goto dispatch

		case OP_DIV:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			if !a.IsInt() || !b.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit divide %s and %s", a.TypeName(), b.TypeName(),
				)
			}

			if b.AsInt() == 0 {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError("I no fit divide by zero o!")
			}

			vm.stack[stackTop] = NewInt(a.AsInt() / b.AsInt())
			stackTop++
			goto dispatch

		case OP_NEGATE:
			a = vm.stack[stackTop-1]

			if !a.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit negate %s", a.TypeName(),
				)
			}

			vm.stack[stackTop-1] = NewInt(-a.AsInt())
			goto dispatch

		// ====================================================================
		// Comparison
		// ====================================================================

		case OP_EQUAL:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			vm.stack[stackTop] = NewBool(a.Equals(b))
			stackTop++
			goto dispatch

		case OP_NOT_EQUAL:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			vm.stack[stackTop] = NewBool(!a.Equals(b))
			stackTop++
			goto dispatch

		case OP_GREATER:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			if !a.IsInt() || !b.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit compare %s and %s", a.TypeName(), b.TypeName(),
				)
			}

			vm.stack[stackTop] = NewBool(a.AsInt() > b.AsInt())
			stackTop++
			goto dispatch

		case OP_LESS:
			b = vm.stack[stackTop-1]
			a = vm.stack[stackTop-2]
			stackTop -= 2

			if !a.IsInt() || !b.IsInt() {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no fit compare %s and %s", a.TypeName(), b.TypeName(),
				)
			}

			vm.stack[stackTop] = NewBool(a.AsInt() < b.AsInt())
			stackTop++
			goto dispatch

		// ====================================================================
		// Logical
		// ====================================================================

		case OP_NOT:
			a = vm.stack[stackTop-1]
			vm.stack[stackTop-1] = NewBool(a.IsFalsey())
			goto dispatch

		// ====================================================================
		// Stack Manipulation
		// ====================================================================

		case OP_POP:
			stackTop--
			goto dispatch

		case OP_DUP:
			vm.stack[stackTop] = vm.stack[stackTop-1]
			stackTop++
			goto dispatch

		// ====================================================================
		// Variables (simplified for Phase 1 - global only)
		// ====================================================================

		case OP_GET_GLOBAL:
			idx := readShort()
			name := vm.chunk.Constants[idx].AsString()
			val, ok := vm.globals[*name]
			if !ok {
				vm.stackTop = stackTop
				vm.ip = ip
				return NewNothing(), vm.runtimeError(
					"I no sabi dis one: %s", *name,
				)
			}
			vm.stack[stackTop] = val
			stackTop++
			goto dispatch

		case OP_SET_GLOBAL:
			idx := readShort()
			name := vm.chunk.Constants[idx].AsString()
			vm.globals[*name] = vm.stack[stackTop-1]
			goto dispatch

		// ====================================================================
		// Control Flow (simplified for Phase 1)
		// ====================================================================

		case OP_JUMP:
			offset := int16(readShort())
			ip += int(offset)
			goto dispatch

		case OP_JUMP_IF_LIE:
			offset := int16(readShort())
			condition := vm.stack[stackTop-1]
			stackTop--
			if condition.IsFalsey() {
				ip += int(offset)
			}
			goto dispatch

		case OP_JUMP_IF_TRU:
			offset := int16(readShort())
			condition := vm.stack[stackTop-1]
			stackTop--
			if condition.IsTruthy() {
				ip += int(offset)
			}
			goto dispatch

		case OP_LOOP:
			offset := readShort()
			ip -= int(offset)
			goto dispatch

		// ====================================================================
		// Special
		// ====================================================================

		// ====================================================================
		// Builtins
		// ====================================================================

		case OP_YARN:
			argCount := readByte()

			// Print each argument
			for i := byte(0); i < argCount; i++ {
				idx := stackTop - int(argCount) + int(i)
				val := vm.stack[idx]
				fmt.Print(vm.valueToString(val))
			}
			fmt.Println()

			// Pop arguments
			stackTop -= int(argCount)

			// Push nothing as result
			vm.stack[stackTop] = NewNothing()
			stackTop++
			goto dispatch

		case OP_HALT:
			vm.stackTop = stackTop
			vm.ip = ip
			// Return top of stack or nothing if stack is empty
			if stackTop > 0 {
				return vm.stack[stackTop-1], nil
			}
			return NewNothing(), nil

		default:
			vm.stackTop = stackTop
			vm.ip = ip
			return NewNothing(), vm.runtimeError(
				"Unknown opcode: %d", instruction,
			)
		}
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// valueToString converts a value to a string for printing/concatenation
func (vm *VM) valueToString(v Value) string {
	return v.String()
}

// runtimeError creates a runtime error with line information
func (vm *VM) runtimeError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)

	// Get line number from current instruction
	line := 0
	if vm.ip > 0 && vm.ip <= len(vm.chunk.Lines) {
		line = vm.chunk.Lines[vm.ip-1]
	}

	var errorMsg strings.Builder
	errorMsg.WriteString("Wahala dey o! ") // "There's trouble!"
	errorMsg.WriteString(message)
	errorMsg.WriteString(fmt.Sprintf(" [line %d]\n", line))

	return fmt.Errorf("%s", errorMsg.String())
}
