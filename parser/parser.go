package parser

import (
	"fmt"
	"strconv"

	"pidgin-lang/ast"
	"pidgin-lang/lexer"
	"pidgin-lang/token"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	OR          // abi (when used as or)
	AND         // and
	EQUALS      // be, na
	LESSGREATER // big pass, no reach
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -x, !x, no be x
	CALL        // function(x)
)

// Precedence table mapping tokens to their precedence
var precedences = map[token.TokenType]int{
	token.BE:       EQUALS,
	token.NA:       EQUALS,
	token.NO:       LESSGREATER, // for "no reach"
	token.BIG:      LESSGREATER, // for "big pass"
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.AND:      AND,
	token.ABI:      OR,
	token.LPAREN:   CALL,
}

// Parser holds the state for parsing tokens into an AST
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// New creates a new Parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Register prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRU, p.parseBoolean)
	p.registerPrefix(token.LIE, p.parseBoolean)
	p.registerPrefix(token.NOTHING, p.parseNothingLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NO, p.parsePrefixExpression) // for "no be"
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.SUPPOSE, p.parseSupposeExpression)
	p.registerPrefix(token.DEY, p.parseWhileExpression)
	p.registerPrefix(token.DO, p.parseDoExpression)
	p.registerPrefix(token.YARN, p.parseYarnExpression)

	// Register infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.BE, p.parseInfixExpression)
	p.registerInfix(token.NA, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.BIG, p.parseCompoundComparison)   // big pass
	p.registerInfix(token.NO, p.parseCompoundComparison)    // no reach
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Errors returns the list of parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: no prefix parse function for %s found",
		p.curToken.Line, t)
	p.errors = append(p.errors, msg)
}

// ParseProgram parses the entire program and returns the AST root
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.MAKE:
		return p.parseMakeStatement()
	case token.BRING:
		return p.parseBringStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseMakeStatement parses: make x be 5
func (p *Parser) parseMakeStatement() *ast.MakeStatement {
	stmt := &ast.MakeStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect 'be' or 'na' after identifier
	if !p.peekTokenIs(token.BE) && !p.peekTokenIs(token.NA) {
		p.errors = append(p.errors, fmt.Sprintf("line %d: expected 'be' or 'na' after variable name, got %s",
			p.peekToken.Line, p.peekToken.Type))
		return nil
	}
	p.nextToken()

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBringStatement parses: bring x
func (p *Parser) parseBringStatement() *ast.BringStatement {
	stmt := &ast.BringStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("line %d: could not parse %q as integer",
			p.curToken.Line, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRU)}
}

func (p *Parser) parseNothingLiteral() ast.Expression {
	return &ast.NothingLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// Handle "no be" as a compound prefix operator
	if p.curTokenIs(token.NO) && p.peekTokenIs(token.BE) {
		p.nextToken() // consume 'be'
		expression.Operator = "no be"
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseCompoundComparison handles "big pass" and "no reach"
func (p *Parser) parseCompoundComparison(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left,
	}

	if p.curTokenIs(token.BIG) {
		if !p.expectPeek(token.PASS) {
			return nil
		}
		expression.Operator = "big pass"
	} else if p.curTokenIs(token.NO) {
		if !p.expectPeek(token.REACH) {
			return nil
		}
		expression.Operator = "no reach"
	}

	p.nextToken()
	expression.Right = p.parseExpression(LESSGREATER)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseSupposeExpression parses: suppose condition { ... } abi { ... }
func (p *Parser) parseSupposeExpression() ast.Expression {
	expression := &ast.SupposeExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ABI) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parseWhileExpression parses: dey do while condition { ... }
func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	// Expect 'do' after 'dey'
	if !p.expectPeek(token.DO) {
		return nil
	}

	// Expect 'while' after 'do'
	if !p.expectPeek(token.WHILE) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

// parseDoExpression parses function definition: do name(params) { body }
func (p *Parser) parseDoExpression() ast.Expression {
	expression := &ast.DoExpression{Token: p.curToken}

	// Check if next token is an identifier (named function)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		expression.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expression.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip comma
		p.nextToken() // move to next identifier
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseYarnExpression parses: yarn("message") or yarn(expression)
func (p *Parser) parseYarnExpression() ast.Expression {
	// Treat yarn as an identifier so it can be called
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.peekTokenIs(token.LPAREN) {
		return ident
	}

	// Parse as a call expression
	p.nextToken()
	return p.parseCallExpression(ident)
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip comma
		p.nextToken() // move to next expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}