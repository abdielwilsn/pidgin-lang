package lexer

import (
	"testing"

	"pidgin-lang/token"
)

func TestNextToken(t *testing.T) {
	input := `make x be 5
make y na 10
make name be "Chidi"
yarn("hello")
tru
lie
nothing
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MAKE, "make"},
		{token.IDENT, "x"},
		{token.BE, "be"},
		{token.INT, "5"},
		{token.MAKE, "make"},
		{token.IDENT, "y"},
		{token.NA, "na"},
		{token.INT, "10"},
		{token.MAKE, "make"},
		{token.IDENT, "name"},
		{token.BE, "be"},
		{token.STRING, "Chidi"},
		{token.YARN, "yarn"},
		{token.LPAREN, "("},
		{token.STRING, "hello"},
		{token.RPAREN, ")"},
		{token.TRU, "tru"},
		{token.LIE, "lie"},
		{token.NOTHING, "nothing"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
