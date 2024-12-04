package parser

import (
	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

type Parser struct {
	lex *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func NewPars(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}

	// call twice so currToken point to first token
	p.parsNextToken()
	p.parsNextToken()

	return p
}

func (pars *Parser) parsNextToken() {
	pars.currToken = pars.peekToken
	pars.peekToken = pars.lex.NextToken()
}

func (pars *Parser) ParsCode() *ast.Code {
	statement := []ast.Statement{}
	for pars.currToken.Type != token.EOF {

		if pars.currToken.Type == token.BUAT {
			statement = append(statement, pars.parsBuatStatement())
		}
		pars.parsNextToken()
	}
	return &ast.Code{Statements: statement}
}

func (pars *Parser) parsBuatStatement() *ast.BuatStatement {
	statement := &ast.BuatStatement{
		Token: pars.currToken,
	}
	pars.parsNextToken() // currToken now have to be point to ident name

	ident := pars.parsIdent() // parse the ident name

	// now the peekToken should be '=', but if it not, panic
	if pars.peekToken.Type != token.ASSIGN {
		pars.parsError("Tanda '=' tidak ditemukan!")
	}
	pars.parsNextToken()

	expr := pars.parsExpression()
	statement.Name = ident
	statement.Expression = expr

	return statement
}

func (pars *Parser) parsIdent() ast.Identifier {
	return ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}
}

func (pars *Parser) parsExpression() string {
	pars.parsNextToken()
	expr := ""
	for pars.currToken.Type != token.SEMICOLON {
		expr += pars.currToken.Literal
		pars.parsNextToken()
	}
	return expr
}

func (pars *Parser) parsError(msg string) {
	panic(msg)
}
