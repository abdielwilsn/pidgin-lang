# Pidgin-Lang

**A programming language that uses Nigerian Pidgin English!**

A dynamically-typed programming language with Nigerian Pidgin English syntax, powered by a high-performance bytecode VM achieving **Level 2 performance** (competitive with Lua, 10-20x faster than Python).

## Features

- 🚀 **Fast bytecode VM** with NaN boxing and direct threading
- 🇳🇬 **Nigerian Pidgin syntax** - code like you talk!
- 🎯 **Zero allocations** during execution for primitives
- 📦 **Easy to embed** - clean Go API
- 🔧 **REPL included** for interactive development

## Performance

Pidgin-Lang achieves **Level 2 VM performance**:

- **5-50 ns/op** (competitive with Lua 5.4)
- **10-20x faster than Python/Ruby**
- **Zero heap allocations** for arithmetic and variables
- See [PERFORMANCE.md](PERFORMANCE.md) for detailed benchmarks

## Installation

### Option 1: Install with Go (Recommended)

```bash
go install github.com/abdielwilsn/pidgin-lang@latest
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/abdielwilsn/pidgin-lang.git
cd pidgin-lang

# Build the executable
go build -o pidgin

# (Optional) Move to your PATH
sudo mv pidgin /usr/local/bin/
```

### Option 3: Download Binary

Download pre-built binaries from the [Releases page](https://github.com/abdielwilsn/pidgin-lang/releases).

## Quick Start

### Interactive REPL

```bash
pidgin
```

```
╔═══════════════════════════════════════════════════════════════╗
║                    PIDGIN-LANG v0.1                           ║
║         A programming language that uses Pidgin               ║
╠═══════════════════════════════════════════════════════════════╣
║  Try these:                                                   ║
║    make name be "Chidi"                                       ║
║    yarn("How far, " + name + "!")                             ║
║    suppose 5 big pass 3 { yarn("E plenty!") }                 ║
║                                                               ║
║  Type 'comot' to exit                                         ║
╚═══════════════════════════════════════════════════════════════╝

pidgin>> make greeting be "How far!"
pidgin>> yarn(greeting)
How far!
```

### Run a File

Create `hello.pdg`:

```pidgin
make name be "Chidi"
yarn("How far, " + name + "!")

suppose name be "Chidi" {
    yarn("Na correct guy!")
}
```

Run it:

```bash
pidgin hello.pdg
```

## Language Syntax

### Variables

```pidgin
make name be "Ada"
make age be 25
make isStudent be tru
```

### Arithmetic

```pidgin
make x be 10 + 5
make y be x * 2
make z be (x + y) / 3
```

### Conditionals

```pidgin
suppose age big pass 18 {
    yarn("You don mature!")
} abi {
    yarn("You never reach!")
}
```

### Loops

```pidgin
make counter be 0
dey do while counter no reach 10 {
    yarn(counter)
    make counter be counter + 1
}
```

### Comparisons

- `be` - equals (==)
- `no be` - not equals (!=)
- `big pass` - greater than (>)
- `no reach` - less than (<)
- `and` - logical AND
- `abi` - logical OR (in conditionals)

### Booleans

- `tru` - true
- `lie` - false

### Built-in Functions

```pidgin
yarn("Hello World")  # Print to console
```

## Examples

### FizzBuzz

```pidgin
make i be 1
dey do while i no reach 16 {
    suppose i be 15 {
        yarn("FizzBuzz")
    } abi {
        suppose i be 3 {
            yarn("Fizz")
        } abi {
            suppose i be 5 {
                yarn("Buzz")
            } abi {
                yarn(i)
            }
        }
    }
    make i be i + 1
}
```

### Simple Counter

```pidgin
make count be 0
dey do while count no reach 5 {
    yarn("Count: " + count)
    make count be count + 1
}
yarn("Don finish!")
```

## Architecture

Pidgin-Lang uses a two-phase execution model:

```
Source → Lexer → Parser → AST → Compiler → Bytecode → VM → Result
```

### Key Optimizations

1. **NaN Boxing** - Pack type tags and values into 64-bit words
2. **Direct Threading** - Optimized instruction dispatch with goto
3. **Register Caching** - Keep hot values in CPU registers
4. **Zero Allocations** - Primitives stay on stack
5. **Specialized Opcodes** - Fast paths for common operations

See [PERFORMANCE.md](PERFORMANCE.md) for technical details.

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./vm
go test ./compiler
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmarks
go test -bench=BenchmarkExecution ./...

# With memory allocation stats
go test -bench=. -benchmem ./...
```

### Project Structure

```
pidgin-lang/
├── lexer/          # Tokenization
├── parser/         # AST generation
├── ast/            # AST node definitions
├── compiler/       # AST → Bytecode compiler
├── vm/             # Bytecode virtual machine
├── object/         # Object system (legacy interpreter)
├── evaluator/      # Tree-walking interpreter (legacy)
├── examples/       # Example programs
└── main.go         # CLI entry point
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/abdielwilsn/pidgin-lang.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes and add tests
5. Run tests: `go test ./...`
6. Commit: `git commit -m 'Add amazing feature'`
7. Push: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Roadmap

- [x] Phase 1: Core VM with NaN boxing
- [x] Phase 2: Compiler with bytecode generation
- [ ] Phase 3: Advanced control flow (break, continue)
- [ ] Phase 4: Functions and closures
- [ ] Phase 5: Inline caching for variables
- [ ] Phase 6: Standard library (strings, arrays, files)
- [ ] Phase 7: Final optimizations (peephole, constant folding)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Inspired by:

- [Crafting Interpreters](https://craftinginterpreters.com/) by Bob Nystrom
- Lua VM design
- Nigerian Pidgin English

## Contact

Questions? Open an issue on GitHub!

---

**What are you waiting for? Let's code!** 🚀