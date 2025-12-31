package lexer

import "pidgin-lang/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number for error reporting
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

// readChar advances to the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL signifies EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tok.Line = l.line

	switch l.ch {
	case '=':
		tok = l.newToken(token.ASSIGN, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '!':
		tok = l.newToken(token.BANG, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '<':
		tok = l.newToken(token.LT, l.ch)
	case '>':
		tok = l.newToken(token.GT, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token
func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.line}
}

// readIdentifier reads an identifier (letters and underscores)
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads an integer
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString() string {
	position := l.position + 1 // skip the opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
		}
	}
	str := l.input[position:l.position]
	l.readChar() // consume the closing quote
	return str
}

// skipWhitespace skips spaces, tabs, newlines, carriage returns, and comments
func (l *Lexer) skipWhitespace() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
			l.readChar()
		} else if l.ch == '\n' {
			l.line++
			l.readChar()
		} else if l.ch == '/' && l.peekChar() == '/' {
			// Skip single-line comment
			l.skipComment()
		} else {
			break
		}
	}
}

// skipComment skips until end of line
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// isLetter checks if a character is a letter or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}