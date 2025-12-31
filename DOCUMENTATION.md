# Pidgin Programming Language Documentation

**Version:** 0.1.0
**Official Website:** [pidgin-lang.org](https://pidgin-lang.org)

Pidgin is a programming language with Nigerian Pidgin English syntax, making coding accessible and fun using familiar West African expressions.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Data Types](#data-types)
3. [Variables](#variables)
4. [Operators](#operators)
5. [Control Flow](#control-flow)
6. [Functions](#functions)
7. [Built-in Functions](#built-in-functions)
8. [Comments](#comments)
9. [Keywords Reference](#keywords-reference)
10. [Examples](#examples)

---

## Getting Started

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/pidgin-lang
cd pidgin-lang

# Build the interpreter
go build -o pidgin

# Run a Pidgin file
./pidgin demo.pdg

# Start the REPL
./pidgin
```

### Your First Program

Create a file called `hello.pdg`:

```pidgin
make name be "Chidi"
yarn("How far, " + name + "!")
```

Run it:
```bash
./pidgin hello.pdg
```

Output:
```
How far, Chidi!
```

---

## Data Types

### Integers

64-bit signed integers for whole numbers.

```pidgin
make age be 25
make year be 2024
make negative be -42
```

### Strings

Text values enclosed in single or double quotes.

```pidgin
make greeting be "How far!"
make name be 'Chidi'
make message be "Wetin dey happen?"
```

**String Concatenation:**
```pidgin
make fullName be "Ade" + " " + "Johnson"
make info be "Age: " + 25  // Automatic type conversion
```

### Booleans

Truth values: `tru` (true) or `lie` (false).

```pidgin
make isAdult be tru
make hasPermission be lie
```

**Truthy and Falsy:**
- **Falsy values:** `lie`, `nothing`
- **Truthy values:** Everything else (numbers, strings, functions, `tru`)

### Nothing

Represents absence of value (similar to `null` or `None` in other languages).

```pidgin
make result be nothing
```

### Functions

First-class callable objects that can be passed around and stored.

```pidgin
make adder be do(a, b) {
    bring a + b
}
```

---

## Variables

### Declaration with `make`

Use the `make` keyword to declare and assign variables.

```pidgin
make name be "Amaka"
make age be 30
make isStudent be tru
```

You can also use `na` instead of `be`:

```pidgin
make price na 5000
```

### Variable Scoping

Variables are block-scoped and support closures:

```pidgin
make x be 10

do outer() {
    make y be 20

    do inner() {
        bring x + y  // Can access both x and y
    }

    bring inner()
}

yarn(outer())  // Prints 30
```

---

## Operators

### Arithmetic Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `5 + 3` → `8` |
| `-` | Subtraction | `10 - 2` → `8` |
| `*` | Multiplication | `4 * 5` → `20` |
| `/` | Division | `20 / 4` → `5` |
| `-x` | Negation | `-5` → `-5` |

**Note:** Division by zero throws an error: "Omo! You no fit divide by zero o!"

```pidgin
make sum be 5 + 3
make product be 4 * 5
make result be (10 + 5) / 3
```

### Comparison Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `be` or `na` | Equals | `5 be 5` → `tru` |
| `no be` | Not equals | `5 no be 3` → `tru` |
| `big pass` | Greater than | `10 big pass 5` → `tru` |
| `no reach` | Less than | `3 no reach 10` → `tru` |
| `>` | Greater than (alt) | `10 > 5` → `tru` |
| `<` | Less than (alt) | `3 < 10` → `tru` |

```pidgin
suppose age big pass 18 {
    yarn("You don old!")
}

suppose name be "Chidi" {
    yarn("Na Chidi!")
}
```

### Logical Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `!` | Logical NOT | `!tru` → `lie` |
| `no be` | Logical NOT | `no be tru` → `lie` |
| `and` | Logical AND | `tru and tru` → `tru` |

```pidgin
suppose age big pass 18 and hasID be tru {
    yarn("You fit enter!")
}
```

### Operator Precedence

From highest to lowest:

1. Prefix operators: `-`, `!`, `no be`
2. Function calls: `()`
3. Multiplication/Division: `*`, `/`
4. Addition/Subtraction: `+`, `-`
5. Comparison: `big pass`, `no reach`, `<`, `>`
6. Equality: `be`, `na`, `no be`
7. Logical AND: `and`

Use parentheses to control evaluation order:

```pidgin
make result be (5 + 3) * 2  // 16, not 11
```

---

## Control Flow

### Conditionals: `suppose` / `abi`

The `suppose` statement is like `if`, and `abi` is like `else`.

**Basic syntax:**
```pidgin
suppose condition {
    // code if condition is true
}
```

**With else clause:**
```pidgin
suppose condition {
    // code if condition is true
} abi {
    // code if condition is false
}
```

**Examples:**

```pidgin
make age be 20

suppose age big pass 18 {
    yarn("You be adult!")
} abi {
    yarn("You still be pikin")
}
```

**Nested conditions:**
```pidgin
make score be 85

suppose score big pass 90 {
    yarn("Grade A!")
} abi {
    suppose score big pass 80 {
        yarn("Grade B!")
    } abi {
        suppose score big pass 70 {
            yarn("Grade C!")
        } abi {
            yarn("Grade D")
        }
    }
}
```

### Loops: `dey do while`

The `dey do while` construct creates a while loop.

**Basic syntax:**
```pidgin
dey do while condition {
    // code that repeats
}
```

**Examples:**

**Counting:**
```pidgin
make i be 0

dey do while i no reach 5 {
    yarn(i)
    make i be i + 1
}

// Output: 0, 1, 2, 3, 4
```

**Countdown:**
```pidgin
make count be 10

dey do while count big pass 0 {
    yarn(count)
    make count be count - 1
}

yarn("Blast off!")
```

**Complex conditions:**
```pidgin
make x be 15

dey do while (x / 3) * 3 be x {
    yarn(x)
    make x be x + 3
}
```

---

## Functions

### Function Definition

Functions are defined with the `do` keyword.

**Named function:**
```pidgin
do functionName(param1, param2) {
    // function body
    bring returnValue
}
```

**Anonymous function:**
```pidgin
do(param1, param2) {
    bring returnValue
}
```

### Function Parameters

Functions can accept zero or more parameters:

```pidgin
// No parameters
do greet() {
    yarn("How far!")
}

// One parameter
do square(n) {
    bring n * n
}

// Multiple parameters
do add(a, b) {
    bring a + b
}
```

### Return Values

Use the `bring` keyword to return a value from a function.

```pidgin
do multiply(a, b) {
    bring a * b
}

make result be multiply(5, 3)
yarn(result)  // 15
```

### Function Calls

Call functions using parentheses:

```pidgin
greet()
square(5)
add(10, 20)
```

### Closures

Functions can capture variables from their enclosing scope:

```pidgin
do makeCounter() {
    make count be 0

    do increment() {
        make count be count + 1
        bring count
    }

    bring increment
}

make counter be makeCounter()
yarn(counter())  // 1
yarn(counter())  // 2
yarn(counter())  // 3
```

### Recursion

Functions can call themselves:

```pidgin
do factorial(n) {
    suppose n no reach 2 {
        bring 1
    } abi {
        bring n * factorial(n - 1)
    }
}

yarn(factorial(5))  // 120
```

### First-Class Functions

Functions can be stored in variables and passed as arguments:

```pidgin
do apply(func, value) {
    bring func(value)
}

do double(x) {
    bring x * 2
}

make result be apply(double, 5)
yarn(result)  // 10
```

---

## Built-in Functions

### `yarn` - Print Output

Prints values to the console.

```pidgin
yarn("Hello World")
yarn(42)
yarn(tru)
yarn("Age:", 25, "Name:", "Chidi")
```

**Features:**
- Accepts multiple arguments
- Automatically converts values to strings
- Returns `nothing`

### `len` - String Length

Returns the length of a string.

```pidgin
make message be "How far"
yarn(len(message))  // 7

yarn(len("Pidgin"))  // 6
```

**Note:** Only works with strings. Using with other types causes an error.

### `type` - Type Checking

Returns the type of a value as a string.

```pidgin
yarn(type(42))          // "INTEGER"
yarn(type("text"))      // "STRING"
yarn(type(tru))         // "BOOLEAN"
yarn(type(nothing))     // "NOTHING"
yarn(type(do(x){bring x}))  // "FUNCTION"
```

**Use cases:**
- Debugging
- Type validation
- Conditional logic based on types

---

## Comments

Single-line comments start with `//`:

```pidgin
// This is a comment
make x be 5  // This is also a comment

// Comments are ignored by the interpreter
// Use them to explain your code
```

---

## Keywords Reference

| Keyword | Purpose | Example |
|---------|---------|---------|
| `make` | Variable declaration | `make x be 5` |
| `be` | Assignment/equality | `x be 5` or `5 be 5` |
| `na` | Alternative to `be` | `make x na 5` |
| `suppose` | If statement | `suppose x big pass 5 { ... }` |
| `abi` | Else statement | `abi { ... }` |
| `dey` | Loop start | `dey do while ...` |
| `do` | Function/loop keyword | `do func(x) { ... }` |
| `while` | While loop part | `dey do while condition` |
| `bring` | Return statement | `bring value` |
| `yarn` | Print function | `yarn("text")` |
| `tru` | Boolean true | `tru` |
| `lie` | Boolean false | `lie` |
| `nothing` | Null/None value | `nothing` |
| `and` | Logical AND | `a and b` |
| `no` | Negation prefix | `no be x` |
| `big` | Greater than (part 1) | `a big pass b` |
| `pass` | Greater than (part 2) | `a big pass b` |
| `reach` | Less than (part) | `a no reach b` |

---

## Examples

### Example 1: Greeting Program

```pidgin
do greetPerson(name, age) {
    yarn("How far, " + name + "!")

    suppose age big pass 18 {
        yarn("You don old o!")
    } abi {
        yarn("You still young!")
    }
}

greetPerson("Amaka", 25)
greetPerson("Tunde", 16)
```

### Example 2: Sum of Numbers

```pidgin
do sumUpTo(n) {
    make sum be 0
    make i be 1

    dey do while i no reach n + 1 {
        make sum be sum + i
        make i be i + 1
    }

    bring sum
}

yarn("Sum of 1 to 10:", sumUpTo(10))  // 55
yarn("Sum of 1 to 100:", sumUpTo(100))  // 5050
```

### Example 3: Fibonacci Sequence

```pidgin
do fibonacci(n) {
    suppose n no reach 2 {
        bring n
    } abi {
        bring fibonacci(n - 1) + fibonacci(n - 2)
    }
}

make i be 0
dey do while i no reach 10 {
    yarn("Fib(" + i + ") = " + fibonacci(i))
    make i be i + 1
}
```

### Example 4: Even or Odd

```pidgin
do isEven(n) {
    bring (n / 2) * 2 be n
}

do checkNumber(num) {
    suppose isEven(num) {
        yarn(num + " na even number")
    } abi {
        yarn(num + " na odd number")
    }
}

checkNumber(7)   // 7 na odd number
checkNumber(12)  // 12 na even number
```

### Example 5: String Manipulation

```pidgin
do repeatString(str, times) {
    make result be ""
    make i be 0

    dey do while i no reach times {
        make result be result + str
        make i be i + 1
    }

    bring result
}

yarn(repeatString("Ha", 5))  // HaHaHaHaHa
yarn(repeatString("Oya! ", 3))  // Oya! Oya! Oya!
```

### Example 6: Higher-Order Functions

```pidgin
do map(func, start, end) {
    make i be start

    dey do while i no reach end + 1 {
        yarn(func(i))
        make i be i + 1
    }
}

do square(x) {
    bring x * x
}

do cube(x) {
    bring x * x * x
}

yarn("Squares:")
map(square, 1, 5)

yarn("Cubes:")
map(cube, 1, 5)
```

### Example 7: FizzBuzz

```pidgin
do fizzBuzz(n) {
    make i be 1

    dey do while i no reach n + 1 {
        make by3 be (i / 3) * 3 be i
        make by5 be (i / 5) * 5 be i

        suppose by3 and by5 {
            yarn("FizzBuzz")
        } abi {
            suppose by3 {
                yarn("Fizz")
            } abi {
                suppose by5 {
                    yarn("Buzz")
                } abi {
                    yarn(i)
                }
            }
        }

        make i be i + 1
    }
}

fizzBuzz(15)
```

---

## Error Messages

Pidgin provides helpful error messages in Pidgin English:

- **Unknown identifier:** `"I no sabi dis one: variableName"`
- **Type mismatch:** `"I no fit do OPERATION wit TYPE and TYPE"`
- **Division by zero:** `"Omo! You no fit divide by zero o!"`
- **Wrong argument type:** `"argument to 'len' must be STRING, got TYPE"`
- **Comparison error:** `"I no fit compare TYPE wit TYPE"`

---

## Running Pidgin Programs

### Execute a File

```bash
./pidgin yourfile.pdg
```

### Interactive REPL

Start the REPL:
```bash
./pidgin
```

Type expressions and statements interactively:
```
>>> make x be 10
>>> yarn(x * 2)
20
>>> comot
```

Exit with: `comot`, `exit`, or `quit`

### Execution Modes

Pidgin has two execution engines:

1. **Bytecode VM** (default) - Fast, optimized
2. **Tree-walking Interpreter** - Legacy mode

To use the interpreter instead of VM:
```bash
./pidgin --vm=false yourfile.pdg
```

---

## Language Philosophy

Pidgin is designed to:
- Make programming accessible using Nigerian Pidgin English
- Provide familiar syntax for West African developers
- Support functional programming with closures and first-class functions
- Be simple, expressive, and fun to use

---

## Future Features

Planned features (not yet implemented):
- Arrays/Lists
- Hash tables/Dictionaries
- Break and continue statements
- File I/O
- More string manipulation functions
- Extended standard library

---

## Contributing

Contributions are welcome! Visit the GitHub repository to:
- Report bugs
- Suggest features
- Submit pull requests
- Improve documentation

---

## License

See the LICENSE file in the repository.

---

## Learn More

- **Repository:** [https://github.com/yourusername/pidgin-lang](https://github.com/yourusername/pidgin-lang)
- **Website:** [https://pidgin-lang.org](https://pidgin-lang.org)
- **Examples:** See [demo.pdg](demo.pdg) for more code samples

---

**Enjoy coding in Pidgin! How far? Make we dey code!**
