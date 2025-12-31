package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"pidgin-lang/compiler"
	"pidgin-lang/evaluator"
	"pidgin-lang/lexer"
	"pidgin-lang/object"
	"pidgin-lang/parser"
	"pidgin-lang/vm"
)

const VERSION = "0.1.0"
const PROMPT = "pidgin>> "

const WELCOME = `╔═══════════════════════════════════════════════════════════════╗
║                    PIDGIN-LANG v0.1                           ║
║         Na programming language wey dey use Pidgin            ║
╠═══════════════════════════════════════════════════════════════╣
║  Try these:                                                   ║
║    make name be "Chidi"                                       ║
║    yarn("How far, " + name + "!")                             ║
║    suppose 5 big pass 3 { yarn("E plenty!") }                 ║
║                                                               ║
║  Type 'comot' to exit                                         ║
╚═══════════════════════════════════════════════════════════════╝`

var (
	useVM       = flag.Bool("vm", true, "Use bytecode VM (default: true)")
	showVersion = flag.Bool("version", false, "Show version and exit")
	showHelp    = flag.Bool("help", false, "Show help and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Pidgin-Lang v%s\n", VERSION)
		fmt.Println("High-performance bytecode VM for Nigerian Pidgin English")
		return
	}

	if *showHelp {
		printHelp()
		return
	}

	args := flag.Args()
	if len(args) > 0 {
		// Run file mode
		runFile(args[0])
	} else {
		// REPL mode
		fmt.Println(WELCOME)
		if *useVM {
			fmt.Println("🚀 Using bytecode VM (Level 2 performance)")
		} else {
			fmt.Println("⚠️  Using legacy tree-walking interpreter")
		}
		fmt.Println()
		startREPL(os.Stdin, os.Stdout)
	}
}

func printHelp() {
	fmt.Println("Pidgin-Lang - Programming language wey dey use Nigerian Pidgin")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  pidgin [OPTIONS] [FILE]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --vm          Use bytecode VM (default: true)")
	fmt.Println("  --version     Show version and exit")
	fmt.Println("  --help        Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  pidgin                  # Start REPL")
	fmt.Println("  pidgin program.pdg      # Run file with VM")
	fmt.Println("  pidgin --vm=false file  # Use legacy interpreter")
	fmt.Println()
	fmt.Println("For more info, visit: https://github.com/abdielwilsn/pidgin-lang")
}

func startREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// Initialize based on VM flag
	var env *object.Environment
	var vmachine *vm.VM

	if *useVM {
		vmachine = vm.NewVM()
	} else {
		env = object.NewEnvironment()
	}

	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		if line == "comot" || line == "exit" || line == "quit" {
			fmt.Fprintln(out, "We go see later! 👋")
			return
		}

		if line == "" {
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if *useVM {
			// Use bytecode VM
			comp := compiler.New()
			chunk, err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Compile wahala: %s\n", err)
				continue
			}

			result, err := vmachine.Run(chunk)
			if err != nil {
				fmt.Fprintf(out, "Runtime wahala: %s\n", err)
				continue
			}

			// Don't print "nothing" values
			if !result.IsNothing() {
				io.WriteString(out, result.String())
				io.WriteString(out, "\n")
			}
		} else {
			// Use legacy tree-walking interpreter
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				// Don't print "nothing" for statements that don't return meaningful values
				if evaluated.Type() != object.NOTHING_OBJ {
					io.WriteString(out, evaluated.Inspect())
					io.WriteString(out, "\n")
				}
			}
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Wahala! I no fit read file: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintln(os.Stderr, "Wahala:", msg)
		}
		os.Exit(1)
	}

	if *useVM {
		// Use bytecode VM
		comp := compiler.New()
		chunk, err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Compile wahala: %s\n", err)
			os.Exit(1)
		}

		vmachine := vm.NewVM()
		_, err = vmachine.Run(chunk)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime wahala: %s\n", err)
			os.Exit(1)
		}
	} else {
		// Use legacy tree-walking interpreter
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
			fmt.Fprintln(os.Stderr, evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Wahala! Parser don confuse:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}