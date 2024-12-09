package parser

import (
	"fmt"

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

type Parser struct {
	lex    *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	// map to define the token type assosiations with prefix or infix
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

	pars.prefixParsMap = map[token.TokenType]prefixParsFunc{}
	pars.registerPrefix(token.IDENT, pars.parsIdent)
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

func (pars *Parser) parsIdent() ast.Expression {
	return &ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}
}

func (pars *Parser) parsExpressionStatement() *ast.ExpressionStatement {
	exprStmnt := &ast.ExpressionStatement{
		Token: pars.currToken,
	}
	exprStmnt.Expression = pars.parsExpression(LOWEST)
	if pars.peekToken.Type == token.SEMICOLON {
		pars.parsNextToken()
	}
	return exprStmnt
}

func (pars *Parser) parsExpression(precedence int) ast.Expression {
	prefix := pars.prefixParsMap[pars.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (pars *Parser) expectPeek(tok token.TokenType) bool {
	return pars.peekToken.Type == tok
}

func (pars *Parser) peekError(expectTok token.TokenType) {
	msg := fmt.Sprintf("Expected next token is %s, but got %s", expectTok, pars.peekToken)
	pars.errors = append(pars.errors, msg)
}

func (pars *Parser) registerPrefix(tokType token.TokenType, f prefixParsFunc) {
	pars.prefixParsMap[tokType] = f
}
func (pars *Parser) registerInfix(tokeType token.TokenType, f infixParsFunc) {
	pars.infixParsMap[tokeType] = f
}
