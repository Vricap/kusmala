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
	INDEX       // array[]
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
	token.LPAREN:     CALL,
	token.LBRACKET:   INDEX,
}

type Parser struct {
	lex       *lexer.Lexer
	Errors    []string
	DevErrors []string

	currToken token.Token
	peekToken token.Token

	// a map to define the token type assosiations with prefix or infix function
	prefixParsMap map[token.TokenType]prefixParsFunc
	infixParsMap  map[token.TokenType]infixParsFunc
}

func NewPars(lex *lexer.Lexer) *Parser {
	pars := &Parser{
		lex:    lex,
		Errors: []string{},
	}

	// call twice so currToken point to first token
	pars.parsNextToken()
	pars.parsNextToken()

	// PREFIX EXPRESSION
	pars.prefixParsMap = map[token.TokenType]prefixParsFunc{} // initialize empty prefixParsMap map
	pars.registerPrefix(token.IDENT, pars.parsIdent)          // register token type ident with parsIdent function which match prefixParsFunc function type
	pars.registerPrefix(token.INTEGER, pars.parsIntegerLiteral)
	pars.registerPrefix(token.BANG, pars.parsPrefix)
	pars.registerPrefix(token.MINUS, pars.parsPrefix)
	pars.registerPrefix(token.BENAR, pars.parsBooleanLiteral)
	pars.registerPrefix(token.SALAH, pars.parsBooleanLiteral)
	pars.registerPrefix(token.FUNGSI, pars.parsFungsiLiteral)
	pars.registerPrefix(token.PANJANG, pars.parsPanjangFungsi)
	pars.registerPrefix(token.STRING, pars.parsStringLiteral)
	pars.registerPrefix(token.LBRACKET, pars.parsArrayLiteral)

	// i decide if-else is a statement and NOT a expression
	// pars.registerPrefix(token.JIKA, pars.parsJikaExpression)
	// kusmala does not grouped expression... i dont know how to implement it :(
	// pars.registerPrefix(token.LPAREN, pars.parseGroupedExpression)

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
	pars.registerInfix(token.LPAREN, pars.parsCallExpression)
	pars.registerInfix(token.LBRACKET, pars.parsIndexExpression)
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
	case token.JIKA:
		return pars.parsJikaStatement()
	case token.CETAK:
		return pars.parsCetakStatement()
	case token.IDENT:
		// TODO: find out more about this
		if pars.expectPeek(token.ASSIGN) { // only run if it is not function call
			return pars.parsReassignmentStatement()
		}
		return pars.parsExpressionStatement()
	default:
		// since the real statement in the language in only 2 (buat & kembalikan), then other statement must be expression statement
		return pars.parsExpressionStatement()
	}
}

func (pars *Parser) parsBuatStatement() *ast.BuatStatement {
	statement := &ast.BuatStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}
	if !pars.expectPeek(token.IDENT) {
		// pars.Errors("Sebuah buat statement membutuhkan nama!")
		pars.peekError(token.IDENT, pars.lex.Line)
	}

	pars.parsNextToken() // currToken now have to be point to ident name
	statement.Name = &ast.Identifier{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
	}

	if !pars.expectPeek(token.ASSIGN) {
		// pars.Errors("Tanda '=' tidak ditemukan!")
		pars.peekError(token.ASSIGN, pars.lex.Line)
	}
	pars.parsNextToken()

	_, ok := pars.prefixParsMap[pars.peekToken.Type]
	if !ok {
		pars.Errors = append(pars.Errors, fmt.Sprintf("ERROR di baris %d: \n\tMengharapkan Nilai atau Ekspresi, tetapi mendapatkan '%s'.", pars.lex.Line, pars.peekToken.Literal))
	}
	pars.parsNextToken()
	statement.Expression = pars.parsExpression(LOWEST)

	// TODO: find out more about this
	if pars.peekToken.Type == token.SEMICOLON {
		pars.parsNextToken()
	}
	return statement
}

func (pars *Parser) parsKembalikanStatement() *ast.KembalikanStatement {
	statement := &ast.KembalikanStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}
	if pars.expectPeek(token.SEMICOLON) {
		pars.parsNextToken()
		return statement
	}
	_, ok := pars.prefixParsMap[pars.peekToken.Type]
	if !ok {
		pars.Errors = append(pars.Errors, fmt.Sprintf("ERROR di baris %d: \n\tMengharapkan Nilai atau Ekspresi, tetapi mendapatkan %s.", pars.lex.Line, pars.peekToken.Literal))
	}
	pars.parsNextToken()
	statement.Expression = pars.parsExpression(LOWEST)
	if pars.expectPeek(token.SEMICOLON) {
		pars.parsNextToken()
	}
	return statement
}

func (pars *Parser) parsJikaStatement() *ast.JikaStatement {
	jika := &ast.JikaStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}

	if !pars.expectPeek(token.LPAREN) {
		pars.peekError(token.LPAREN, pars.lex.Line)
	}
	pars.parsNextToken()
	if pars.expectPeek(token.RPAREN) {
		pars.Errors = append(pars.Errors, "Kondisi tidak boleh kosong!")
	}
	pars.parsNextToken()
	jika.Condition = pars.parsExpression(LOWEST)

	if !pars.expectPeek(token.RPAREN) {
		pars.peekError(token.RPAREN, pars.lex.Line)
	}
	pars.parsNextToken()
	if !pars.expectPeek(token.LBRACE) {
		pars.peekError(token.LBRACE, pars.lex.Line)
	}
	pars.parsNextToken()
	pars.parsNextToken()
	jika.JikaBlock = pars.parsBlockStatement()

	// TODO: this is stupid
	if pars.expectPeek(token.LAINNYA) {
		pars.parsNextToken()
		if !pars.expectPeek(token.LBRACE) {
			pars.peekError(token.LBRACE, pars.lex.Line)
		}
		pars.parsNextToken()
		pars.parsNextToken()
		jika.LainnyaBlock = pars.parsBlockStatement()
	}
	return jika
}

func (pars *Parser) parsBlockStatement() *ast.BlockStatement {
	stmnt := &ast.BlockStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}
	for pars.currToken.Type != token.RBRACE {
		if pars.currToken.Type == token.EOF {
			break
		}
		stmnt.Statements = append(stmnt.Statements, pars.parsStatement())
		pars.parsNextToken()
	}
	// TODO: quick hack
	if !pars.expectCurr(token.RBRACE) {
		pars.currError(token.RBRACE, pars.lex.Line)
	}
	return stmnt
}

func (pars *Parser) parsCetakStatement() *ast.CetakStatement {
	cetak := &ast.CetakStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}
	if !pars.expectPeek(token.LPAREN) {
		pars.peekError(token.LPAREN, pars.lex.Line)
	}
	pars.parsNextToken()
	if pars.expectPeek(token.RPAREN) {
		return cetak
	}
	pars.parsNextToken()
	cetak.Expression = pars.parsArguments()

	if pars.peekToken.Type == token.SEMICOLON {
		pars.parsNextToken()
	}
	return cetak
}

func (pars *Parser) parsReassignmentStatement() *ast.ReassignStatement {
	rs := &ast.ReassignStatement{Token: pars.currToken, Ln: pars.lex.Line}
	rs.Ident = &ast.Identifier{Token: pars.currToken, Ln: pars.lex.Line, Value: pars.currToken.Literal}
	if !pars.expectPeek(token.ASSIGN) {
		pars.peekError(token.ASSIGN, pars.lex.Line)
	}
	pars.parsNextToken()
	pars.parsNextToken()
	rs.NewValue = pars.parsExpression(LOWEST)
	if pars.expectPeek(token.SEMICOLON) {
		pars.parsNextToken()
	}
	return rs
}

/*******************************************
*			EXPRESSION PARSING			   *
*******************************************/

func (pars *Parser) parsExpressionStatement() *ast.ExpressionStatement {
	exprStmnt := &ast.ExpressionStatement{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
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
		pars.Errors = append(pars.Errors, fmt.Sprintf("ERROR di baris %d: \n\tToken tidak diharapkan ditemukan '%s'", pars.lex.Line, pars.currToken.Literal))
		pars.DevErrors = append(pars.DevErrors, fmt.Sprintf("There's not function assosiated with %v, literal: %s", pars.currToken.Type, pars.currToken.Literal))
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
		Ln:    pars.lex.Line,
	}
}

func (pars *Parser) parsIntegerLiteral() ast.Expression {
	literal, err := strconv.Atoi(pars.currToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse literal: %s to integer", pars.currToken.Literal)
		pars.DevErrors = append(pars.DevErrors, msg)
		return nil
	}
	int := &ast.IntegerLiteral{
		Token: pars.currToken,
		Value: literal,
		Ln:    pars.lex.Line,
	}
	return int
}

func (pars *Parser) parsBooleanLiteral() ast.Expression {
	bool := &ast.BooleanLiteral{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}
	if pars.currToken.Literal == "benar" {
		bool.Value = true
	} else if pars.currToken.Literal == "salah" {
		bool.Value = false
	}
	return bool
}

func (pars *Parser) parsPrefix() ast.Expression {
	prefix := &ast.PrefixExpression{
		Token:    pars.currToken,
		Operator: pars.currToken.Literal,
		Ln:       pars.lex.Line,
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
		Ln:       pars.lex.Line,
	}
	precedence := pars.currPrecedence()
	pars.parsNextToken()
	exp.Right = pars.parsExpression(precedence)

	return exp
}

// func (pars *Parser) parseGroupedExpression() ast.Expression {
// 	pars.parsNextToken()
// 	exp := pars.parsExpression(LOWEST)
// 	if !pars.expectPeek(token.RPAREN) {
// 		return nil
// 	}
// 	return exp
// }

func (pars *Parser) parsFungsiLiteral() ast.Expression {
	fung := &ast.FungsiExpression{
		Token: pars.currToken,
		Ln:    pars.lex.Line,
	}

	if !pars.expectPeek(token.LPAREN) {
		pars.peekError(token.LPAREN, pars.lex.Line)
	}
	pars.parsNextToken()
	if pars.expectPeek(token.RPAREN) {
		pars.parsNextToken()
	} else {
		pars.parsNextToken()
		fung.Params = pars.parsParams()
	}
	if !pars.expectPeek(token.LBRACE) {
		pars.peekError(token.LBRACE, pars.lex.Line)
	}
	pars.parsNextToken()
	pars.parsNextToken()
	fung.Body = pars.parsBlockStatement()
	return fung
}

func (pars *Parser) parsParams() []*ast.Identifier {
	iden := []*ast.Identifier{}

	for pars.currToken.Type != token.RPAREN {
		if pars.currToken.Type == token.COMMA {
			pars.parsNextToken()
			continue
		}
		if pars.currToken.Type == token.EOF {
			break
		}
		iden = append(iden, &ast.Identifier{Token: pars.currToken, Value: pars.currToken.Literal, Ln: pars.lex.Line})
		pars.parsNextToken()
	}
	// TODO: quick hack
	if !pars.expectCurr(token.RPAREN) {
		pars.currError(token.RPAREN, pars.lex.Line)
	}
	return iden
}

func (pars *Parser) parsPanjangFungsi() ast.Expression {
	pf := &ast.PanjangFungsi{Token: pars.currToken, Ln: pars.lex.Line}
	if !pars.expectPeek(token.LPAREN) {
		pars.peekError(token.LPAREN, pars.lex.Line)
	}
	pars.parsNextToken()
	pars.parsNextToken()
	pf.Argument = pars.parsExpression(LOWEST)
	if !pars.expectPeek(token.RPAREN) {
		pars.peekError(token.RPAREN, pars.lex.Line)
	}
	pars.parsNextToken()

	return pf
}

// add(1, 2 * 3, 1 - 2)`
func (pars *Parser) parsCallExpression(ident ast.Expression) ast.Expression {
	ce := &ast.CallExpression{Token: pars.currToken, Function: ident, Ln: pars.lex.Line}
	pars.parsNextToken()
	ce.Arguments = pars.parsArguments()
	return ce
}

func (pars *Parser) parsArguments() []ast.Expression {
	expr := []ast.Expression{}

	for pars.currToken.Type != token.RPAREN {
		if pars.currToken.Type == token.COMMA {
			pars.parsNextToken()
			continue
		}
		if pars.currToken.Type == token.EOF {
			break
		}
		expr = append(expr, pars.parsExpression(LOWEST))
		pars.parsNextToken()
	}
	// TODO: quick hack
	if !pars.expectCurr(token.RPAREN) {
		pars.currError(token.RPAREN, pars.lex.Line)
	}
	return expr
}

func (pars *Parser) parsStringLiteral() ast.Expression {
	str := &ast.StringLiteral{
		Token: pars.currToken,
		Value: pars.currToken.Literal,
		Ln:    pars.lex.Line,
	}
	return str
}

func (pars *Parser) parsArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: pars.currToken, Ln: pars.lex.Line}
	pars.parsNextToken()
	if pars.expectCurr(token.RBRACKET) {
		return arr
	}
	arr.Elements = pars.parsArrElements()
	return arr
}

func (pars *Parser) parsIndexExpression(left ast.Expression) ast.Expression {
	index := &ast.IndexExpression{Token: pars.currToken, Ln: pars.lex.Line, Left: left}
	pars.parsNextToken()
	if pars.expectCurr(token.RBRACKET) {
		return index
	}
	index.Index = pars.parsExpression(LOWEST)
	if !pars.expectPeek(token.RBRACKET) {
		pars.peekError(token.RBRACKET, pars.lex.Line)
	}
	pars.parsNextToken()
	return index
}

func (pars *Parser) parsArrElements() []ast.Expression {
	el := []ast.Expression{}

	for pars.currToken.Type != token.RBRACKET {
		if pars.currToken.Type == token.COMMA {
			pars.parsNextToken()
			continue
		}
		if pars.currToken.Type == token.EOF {
			break
		}
		el = append(el, pars.parsExpression(LOWEST))
		pars.parsNextToken()
	}
	// TODO: quick hack
	if !pars.expectCurr(token.RBRACKET) {
		pars.currError(token.RBRACKET, pars.lex.Line)
	}
	return el
}

/*******************************************
*			HELPER METHOD   			   *
*******************************************/

func (pars *Parser) expectPeek(tok token.TokenType) bool {
	return pars.peekToken.Type == tok
}

func (pars *Parser) expectCurr(tok token.TokenType) bool {
	return pars.currToken.Type == tok
}

func (pars *Parser) peekError(expectTok token.TokenType, l int) {
	msg := fmt.Sprintf("ERROR di baris %d: \n\tToken selanjutnya mengharapkan %s, tetapi menemukan '%s'", l, expectTok, pars.peekToken.Literal)
	pars.Errors = append(pars.Errors, msg)
}

func (pars *Parser) currError(expectTok token.TokenType, l int) {
	msg := fmt.Sprintf("ERROR di baris %d: \n\tToken sekarang mengharapkan %s, tetapi menemukan '%s'", l, expectTok, pars.currToken.Literal)
	pars.Errors = append(pars.Errors, msg)
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

func (pars *Parser) skipTokenTo(tok token.TokenType) {
	for pars.currToken.Type != tok {
		if pars.currToken.Type == token.EOF {
			return
			// EOF_FOUND()
		}
		pars.parsNextToken()
	}
}

func EOF_FOUND() {
	panic("EOF ditemukan! Program keluar!")
}

/*
This infix expression: 1-2+3 translate to this:
type InfinixExpression struct {
	Token: token.PLUS
	Right: type IntegerExpression struct {
		Token: token.INTEGER
		Value: 3
	}
	Operator: "+"
	Left: type InfinixExpression struct {
		Token: token.MINUS
		Right: type IntegerExpression struct {
			Token: token.INTEGER
			Value: 2
		}
		Operator: "-"
		Left: type IntegerExpression struct {
			Token: token.INTEGER
			Value: 1
		}
	}
}
*/
