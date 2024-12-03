package parser

import (
	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

type Parser struct {
	lex *lexer.Lexer

	currToken token.Token // current token under examination
	peekToken token.Token // next token
}

func NewPars(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}

	// call ParsNextToken() twice so currToken point to first token and peekToken to the next
	p.ParsNextToken()
	p.ParsNextToken()
	return p
}

func (pars *Parser) ParsNextToken() {
	pars.currToken = pars.peekToken
	pars.peekToken = pars.lex.NextToken()
}

func (pars *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for pars.currToken.Type != token.EOF {
		statement := pars.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		pars.ParsNextToken()
	}
	return program
}

func (pars *Parser) parseStatement() ast.Statement {
	switch pars.currToken.Type {
	case token.BUAT:
		return pars.parseBuatStatement()
	default:
		return nil
	}
}

func (pars *Parser) parseBuatStatement() *ast.BuatStatement {
	statement := &ast.BuatStatement{Token: pars.currToken}

	if !pars.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: pars.currToken, Value: pars.currToken.Literal}

	if !pars.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: we're skipping the expression until we
	// encounter a semicolon
	for !pars.currTokenIs(token.SEMICOLON) {
		pars.ParsNextToken()
	}
	return statement
}

func (pars *Parser) currTokenIs(tok token.TokenType) bool {
	return pars.currToken.Type == tok
}

func (pars *Parser) peekTokenIs(tok token.TokenType) bool {
	return pars.peekToken.Type == tok
}

func (pars *Parser) expectPeek(tok token.TokenType) bool {
	if pars.peekTokenIs(tok) {
		pars.ParsNextToken()
		return true
	} else {
		return false
	}
}
