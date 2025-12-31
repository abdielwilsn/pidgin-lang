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
	INT    TokenType = "INT"    // Integer literals 
	STRING TokenType = "STRING" // String literals

	// Operators
	ASSIGN   TokenType = "="   
	PLUS     TokenType = "+"         
	MINUS    TokenType = "-"         
	BANG     TokenType = "!"         
	ASTERISK TokenType = "*"         
	SLASH    TokenType = "/"       

	// Comparison (single character)
	LT TokenType = "<" // <  (less than)
	GT TokenType = ">" // >  (greater than)

	// Delimiters
	COMMA     TokenType = "," 
	SEMICOLON TokenType = ";" 
	LPAREN    TokenType = "(" 
	RPAREN    TokenType = ")" 
	LBRACE    TokenType = "{" 
	RBRACE    TokenType = "}" 

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

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}