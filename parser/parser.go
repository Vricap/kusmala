package parser

import (
	"fmt"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

type Parser struct {
	lex       *lexer.Lexer
	errors    []string
	currToken token.Token
	peekToken token.Token
}

func NewPars(lex *lexer.Lexer) *Parser {
	p := &Parser{
		lex:    lex,
		errors: []string{},
	}

	// call twice so currToken point to first token
	p.parsNextToken()
	p.parsNextToken()

	return p
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
		return nil
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

	pars.parsNextToken()              // currToken now have to be point to ident name
	statement.Name = pars.parsIdent() // parse the ident name

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

// TODO: this is simple enough, better not separate function
func (pars *Parser) parsIdent() *ast.Identifier {
	return &ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}
}

func (pars *Parser) parsExpression() string {
	expr := ""
	for pars.currToken.Type != token.SEMICOLON {
		// TODO: we just return the expression string for now
		expr += pars.currToken.Literal
		pars.parsNextToken()
	}
	return expr
}

func (pars *Parser) expectPeek(tok token.TokenType) bool {
	return pars.peekToken.Type == tok
}

// func (pars *Parser) AnyErrors(msg string) []string {
// 	return pars.errors
// }

func (pars *Parser) peekError(expectTok token.TokenType) {
	msg := fmt.Sprintf("Expected next token is %s, but got %s", expectTok, pars.peekToken)
	pars.errors = append(pars.errors, msg)
}
