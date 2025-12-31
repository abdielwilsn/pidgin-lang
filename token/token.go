package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL" // Unknown or invalid character/token
	EOF     TokenType = "EOF"     // End of file/input

	// Identifiers and literals
	IDENT  TokenType = "IDENT"  // User-defined names: variables, functions (e.g., money, calculateSalary)
	INT    TokenType = "INT"    // Integer literals (e.g., 42, 1000)
	STRING TokenType = "STRING" // String literals (e.g., "How you dey?", 'Wetin dey happen?')

	// Operators
	ASSIGN   TokenType = "="   // =  (assignment)
	PLUS     TokenType = "+"        // +  (addition or string concatenation)
	MINUS    TokenType = "-"        // -  (subtraction)
	BANG     TokenType = "!"        // !  (bang / not)
	ASTERISK TokenType = "*"        // *  (multiplication)
	SLASH    TokenType = "/"        // /  (division)

	// Comparison (single character)
	LT TokenType = "<" // <  (less than)
	GT TokenType = ">" // >  (greater than)

	// Delimiters
	COMMA     TokenType = "," // ,  (separator in arguments, lists)
	SEMICOLON TokenType = ";" // ;  (optional statement terminator)
	LPAREN    TokenType = "(" // (  (left parenthesis – grouping, function calls)
	RPAREN    TokenType = ")" // )  (right parenthesis)
	LBRACE    TokenType = "{" // {  (start of block – functions, conditionals, loops)
	RBRACE    TokenType = "}" // }  (end of block)

	// Keywords (Pidgin-flavored control flow and declarations)
	MAKE      TokenType = "MAKE"      // make – variable declaration (e.g., make name = "John")
	BE        TokenType = "BE"        // be – used in equality check or assignment context (e.g., if x be 5)
	NA        TokenType = "NA"        // na – "is" / equality (common in Pidgin: "na true")
	SUPPOSE   TokenType = "SUPPOSE"   // suppose – if statement (e.g., suppose x > 5 { ... })
	ABI       TokenType = "ABI"       // abi – else (e.g., suppose ... { ... } abi { ... })
	DEY       TokenType = "DEY"       // dey – part of loop/existence (e.g., dey while condition)
	DO        TokenType = "DO"        // do – function definition or loop action
	WHILE     TokenType = "WHILE"     // while – loop keyword (combined with dey/do for "dey do while")
	BRING     TokenType = "BRING"     // bring – return statement (e.g., bring result;)
	YARN      TokenType = "YARN"      // yarn – print/output to console (e.g., yarn "Hello world!")
	TRU       TokenType = "TRU"       // tru – boolean true
	LIE       TokenType = "LIE"       // lie – boolean false
	NOTHING   TokenType = "NOTHING"   // nothing – null / none value
	AND       TokenType = "AND"       // and – logical AND
	NO        TokenType = "NO"        // no – logical NOT or part of negation (e.g., "no be")
	BIG       TokenType = "BIG"       // big – part of comparison (e.g., "big pass" for >)
	PASS      TokenType = "PASS"      // pass – greater than in Pidgin style ("big pass") or no-op
	REACH     TokenType = "REACH"     // reach – part of comparison (e.g., "no reach" for <)
)

var keywords = map[string]TokenType{
	"make":    MAKE,
	"be":      BE,
	"na":      NA,
	"suppose": SUPPOSE,
	"abi":     ABI,
	"dey":     DEY,
	"do":      DO,
	"while":   WHILE,
	"bring":   BRING,
	"yarn":    YARN,
	"tru":     TRU,
	"lie":     LIE,
	"nothing": NOTHING,
	"and":     AND,
	"no":      NO,
	"big":     BIG,
	"pass":    PASS,
	"reach":   REACH,
}

// LookupIdent checks if an identifier is a keyword and returns the corresponding TokenType.
// If not found in the map, it returns IDENT.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}