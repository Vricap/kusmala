package parser

import (
	"fmt"
	"strconv"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

// the idea behind pratt parser is 'token type associations'. which mean, when we encounter certain token type, we call the parser function related to it.
// in our  operator precedence parser, there's two kind of parsing. depend on the token position, is it in the prefix or infix
type prefixParsFunc func() ast.Expression
type infixParsFunc func(ast.Expression) ast.Expression

// operator precedence, the lower the number, the less the priority
const (
	_ int = iota // since its zero, so we dont need that
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedence map[token.TokenType]int = map[token.TokenType]int{
	token.SAMA:       EQUALS,
	token.TIDAK_SAMA: EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.SLASH:      PRODUCT,
	token.ASTERISK:   PRODUCT,
}

type Parser struct {
	lex    *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	// a map to define the token type assosiations with prefix or infix function
	prefixParsMap map[token.TokenType]prefixParsFunc
	infixParsMap  map[token.TokenType]infixParsFunc
}

func NewPars(lex *lexer.Lexer) *Parser {
	pars := &Parser{
		lex:    lex,
		errors: []string{},
	}

	// call twice so currToken point to first token
	pars.parsNextToken()
	pars.parsNextToken()

	// PREFIX EXPRESSION
	pars.prefixParsMap = map[token.TokenType]prefixParsFunc{} // initialize empty prefixParsMap map
	pars.registerPrefix(token.IDENT, pars.parsIdent)          // register token type ident with parsIdent function which match prefixParsFunc function type
	pars.registerPrefix(token.BILBUL, pars.parsIntegerLiteral)
	pars.registerPrefix(token.BANG, pars.parsPrefix)
	pars.registerPrefix(token.MINUS, pars.parsPrefix)
	pars.registerPrefix(token.BENAR, pars.parsBooleanLiteral)
	pars.registerPrefix(token.SALAH, pars.parsBooleanLiteral)
	pars.registerPrefix(token.LPAREN, pars.parseGroupedExpression)

	// INFIX EXPRESSION
	pars.infixParsMap = map[token.TokenType]infixParsFunc{}
	pars.registerInfix(token.PLUS, pars.parsInfix)
	pars.registerInfix(token.MINUS, pars.parsInfix)
	pars.registerInfix(token.ASTERISK, pars.parsInfix)
	pars.registerInfix(token.SLASH, pars.parsInfix)
	pars.registerInfix(token.LT, pars.parsInfix)
	pars.registerInfix(token.GT, pars.parsInfix)
	pars.registerInfix(token.SAMA, pars.parsInfix)
	pars.registerInfix(token.TIDAK_SAMA, pars.parsInfix)
	return pars
}

func (pars *Parser) parsNextToken() {
	pars.currToken = pars.peekToken
	pars.peekToken = pars.lex.NextToken()
}

func (pars *Parser) ConstructTree() *ast.Tree {
	statement := []ast.Statement{}
	for pars.currToken.Type != token.EOF {
		statement = append(statement, pars.parsStatement())
		pars.parsNextToken()
	}
	return &ast.Tree{Statements: statement}
}

/*******************************************
*			STATEMENT PARSING			   *
*******************************************/

func (pars *Parser) parsStatement() ast.Statement {
	switch pars.currToken.Type {
	case token.BUAT:
		return pars.parsBuatStatement()
	case token.KEMBALIKAN:
		return pars.parsKembalikanStatement()
	default:
		// since the real statement in the language in only 2 (buat & kembalikan), then other statement must be expression statement
		return pars.parsExpressionStatement()
	}
}

func (pars *Parser) parsBuatStatement() *ast.BuatStatement {
	statement := &ast.BuatStatement{
		Token: pars.currToken,
	}
	if !pars.expectPeek(token.IDENT) {
		// pars.Errors("Sebuah buat statement membutuhkan nama!")
		pars.peekError(token.IDENT)
	}

	pars.parsNextToken() // currToken now have to be point to ident name
	statement.Name = &ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}

	if !pars.expectPeek(token.ASSIGN) {
		// pars.Errors("Tanda '=' tidak ditemukan!")
		pars.peekError(token.ASSIGN)
	}

	// we skip the expression for now
	for pars.currToken.Type != token.SEMICOLON {
		pars.parsNextToken()
	}

	// statement.Expression = pars.parsExpression()
	return statement
}

func (pars *Parser) parsKembalikanStatement() *ast.KembalikanStatement {
	statemtent := &ast.KembalikanStatement{
		Token: pars.currToken,
	}

	// we skip the expression for now
	for pars.currToken.Type != token.SEMICOLON {
		pars.parsNextToken()
	}
	// statemtent.Expression = pars.parsExpression()
	return statemtent
}

/*******************************************
*			EXPRESSION PARSING			   *
*******************************************/

func (pars *Parser) parsExpressionStatement() *ast.ExpressionStatement {
	exprStmnt := &ast.ExpressionStatement{
		Token: pars.currToken,
	}
	exprStmnt.Expression = pars.parsExpression(LOWEST)
	if pars.peekToken.Type == token.SEMICOLON {
		pars.parsNextToken() // so currToken point to ; since we don't want to do parsStatement() again - we are in eof
	}
	return exprStmnt
}

func (pars *Parser) parsExpression(precedence int) ast.Expression {
	prefix := pars.prefixParsMap[pars.currToken.Type] // check if currToken have function assosiated with that
	if prefix == nil {
		pars.errors = append(pars.errors, fmt.Sprintf("There's not function assosiated with %v", pars.currToken.Type))
		return nil
	}
	leftExp := prefix() // if so, call it
	for pars.peekToken.Type != token.SEMICOLON && precedence < pars.peekPrecedence() {
		infix := pars.infixParsMap[pars.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		pars.parsNextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (pars *Parser) parsIdent() ast.Expression { // this signature match the prefixParsFunc function type
	return &ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}
}

func (pars *Parser) parsIntegerLiteral() ast.Expression {
	literal, err := strconv.Atoi(pars.currToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse literal: %s to integer", pars.currToken.Literal)
		pars.errors = append(pars.errors, msg)
		return nil
	}
	int := &ast.IntegerLiteral{
		Token: pars.currToken,
		Value: literal,
	}
	return int
}

func (pars *Parser) parsBooleanLiteral() ast.Expression {
	bool := &ast.BooleanLiteral{
		Token: pars.currToken,
	}
	if pars.currToken.Literal == "benar" {
		bool.Value = true
	} else if pars.currToken.Literal == "salah" {
		bool.Value = false
	}
	return bool
}

func (pars *Parser) parsPrefix() ast.Expression {
	// if !pars.expectPeek(token.BILBUL) {
	// 	pars.peekError(token.BILBUL)
	// 	return nil
	// }
	prefix := &ast.PrefixExpression{
		Token:    pars.currToken,
		Operator: pars.currToken.Literal,
	}
	pars.parsNextToken() // currToken now point to the integer

	right := pars.parsExpression(PREFIX) // essentially same as pars.parsIntegerLiteral()
	prefix.Right = right
	return prefix
}

func (pars *Parser) parsInfix(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    pars.currToken,
		Operator: pars.currToken.Literal,
		Left:     left,
	}
	precedence := pars.currPrecedence()
	pars.parsNextToken()
	exp.Right = pars.parsExpression(precedence)

	return exp
}

func (pars *Parser) parseGroupedExpression() ast.Expression {
	pars.parsNextToken()
	exp := pars.parsExpression(LOWEST)
	if !pars.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (pars *Parser) expectPeek(tok token.TokenType) bool {
	return pars.peekToken.Type == tok
}

func (pars *Parser) peekError(expectTok token.TokenType) {
	msg := fmt.Sprintf("Expected next token is %s, but got %s", expectTok, pars.peekToken)
	pars.errors = append(pars.errors, msg)
}

// register the token type to the eiter prefixParsFunc or infixParsFunc function type
func (pars *Parser) registerPrefix(tokType token.TokenType, f prefixParsFunc) {
	pars.prefixParsMap[tokType] = f
}
func (pars *Parser) registerInfix(tokeType token.TokenType, f infixParsFunc) {
	pars.infixParsMap[tokeType] = f
}

func (pars *Parser) peekPrecedence() int {
	p, ok := precedence[pars.peekToken.Type]
	if !ok {
		return LOWEST
	}
	return p
}

func (pars *Parser) currPrecedence() int {
	p, ok := precedence[pars.currToken.Type]
	if !ok {
		return LOWEST
	}
	return p
}

/*
This infix expression: 1-2+3 translate to this:
type InfinixExpression struct {
	Token: token.PLUS
	Right: type IntegerExpression struct {
		Token: token.BILBUL
		Value: 3
	}
	Operator: "+"
	Left: type InfinixExpression struct {
		Token: token.MINUS
		Right: type IntegerExpression struct {
			Token: token.BILBUL
			Value: 2
		}
		Operator: "-"
		Left: type IntegerExpression struct {
			Token: token.BILBUL
			Value: 1
		}
	}
}
*/
