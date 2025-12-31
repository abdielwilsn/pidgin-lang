# Quick Start Guide

## Installation

### Option 1: Build from Source (Recommended for Now)

```bash
# Clone the repository
git clone https://github.com/abdielwilsn/pidgin-lang.git
cd pidgin-lang

# Build the executable
go build -o pidgin

# Test it works
./pidgin --version
```

### Option 2: Install with Go (After Publishing)

```bash
go install github.com/abdielwilsn/pidgin-lang@latest
```

## Your First Program

Create a file called `greeting.pdg`:

```pidgin
make name be "Ada"
yarn("How far, " + name + "!")
```

Run it:

```bash
./pidgin greeting.pdg
```

Output:
```
How far, Ada!
```

## Using the REPL

Start the interactive shell:

```bash
./pidgin
```

You'll see:

```
╔═══════════════════════════════════════════════════════════════╗
║                    PIDGIN-LANG v0.1                           ║
║         Na programming language wey dey use Pidgin            ║
╠═══════════════════════════════════════════════════════════════╣
║  Try these:                                                   ║
║    make name be "Chidi"                                       ║
║    yarn("How far, " + name + "!")                             ║
║    suppose 5 big pass 3 { yarn("E plenty!") }                 ║
║                                                               ║
║  Type 'comot' to exit                                         ║
╚═══════════════════════════════════════════════════════════════╝
🚀 Using bytecode VM (Level 2 performance)

pidgin>>
```

Try some commands:

```pidgin
pidgin>> make x be 10
pidgin>> make y be 5
pidgin>> x + y
15
pidgin>> suppose x big pass y { yarn("X win!") }
X win!
pidgin>> comot
We go see later! 👋
```

## Language Basics

### 1. Variables

```pidgin
make name be "Emeka"
make age be 25
make isStudent be tru
```

### 2. Arithmetic

```pidgin
make result be 10 + 5
make product be result * 2
make final be (product - 10) / 2
```

### 3. Print to Console

```pidgin
yarn("Hello World")
yarn("The answer is: ")
yarn(42)
```

### 4. Conditionals

```pidgin
suppose age big pass 18 {
    yarn("You don mature!")
} abi {
    yarn("You still young!")
}
```

### 5. Loops

```pidgin
make i be 0
dey do while i no reach 5 {
    yarn(i)
    make i be i + 1
}
```

### 6. Comparisons

```pidgin
5 be 5           # true (equals)
5 no be 3        # true (not equals)
10 big pass 5    # true (greater than)
3 no reach 10    # true (less than)
```

### 7. Boolean Logic

```pidgin
tru and lie      # false
tru and tru      # true
```

## Example Programs

### Counter (examples/counter.pdg)

```pidgin
make count be 0

dey do while count no reach 5 {
    yarn("Count: ")
    yarn(count)
    make count be count + 1
}

yarn("Don finish!")
```

Run:
```bash
./pidgin examples/counter.pdg
```

### Conditional Logic (examples/hello.pdg)

```pidgin
make name be "Chidi"
yarn("How far, " + name + "!")

suppose name be "Chidi" {
    yarn("Na correct guy!")
}
```

## Command-Line Options

```bash
# Use the bytecode VM (default, fastest)
./pidgin program.pdg

# Use legacy tree-walking interpreter
./pidgin --vm=false program.pdg

# Show version
./pidgin --version

# Show help
./pidgin --help
```

## Performance

Pidgin-Lang uses a high-performance bytecode VM:

- **5-50 ns per operation**
- **10-20x faster than Python/Ruby**
- **Competitive with Lua 5.4**
- **Zero allocations** for arithmetic and primitives

See [PERFORMANCE.md](PERFORMANCE.md) for benchmarks.

## Next Steps

1. **Explore Examples**: Check out the [examples/](examples/) directory
2. **Read the Manual**: See [README.md](README.md) for full language reference
3. **Performance Details**: Read [PERFORMANCE.md](PERFORMANCE.md)
4. **Contributing**: See contribution guidelines in [README.md](README.md)

## Getting Help

- **Issues**: https://github.com/abdielwilsn/pidgin-lang/issues
- **Discussions**: https://github.com/abdielwilsn/pidgin-lang/discussions

## Language Cheat Sheet

| Feature | Pidgin Syntax | English Equivalent |
|---------|---------------|-------------------|
| Variable | `make x be 10` | `x = 10` |
| Print | `yarn("Hello")` | `print("Hello")` |
| If/Else | `suppose x be 5 { ... } abi { ... }` | `if (x == 5) { ... } else { ... }` |
| While Loop | `dey do while x no reach 10 { ... }` | `while (x < 10) { ... }` |
| Equals | `x be y` | `x == y` |
| Not Equals | `x no be y` | `x != y` |
| Greater | `x big pass y` | `x > y` |
| Less | `x no reach y` | `x < y` |
| True | `tru` | `true` |
| False | `lie` | `false` |
| And | `x and y` | `x && y` |
| Exit REPL | `comot` | `exit` |

---

**Na now you be Pidgin programmer! 🚀**
