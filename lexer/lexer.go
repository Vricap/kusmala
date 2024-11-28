package lexer

import (
	"github.com/vricap/kusmala/token"
)

type lexer struct {
	input   string
	pos     int  // current position in the input - point to current char
	peekPos int  // peek the next of the current position
	char    byte // current char under examination
}

func NewLex(input string) *lexer {
	return &lexer{
		input: input,
	}
}

func (lex *lexer) NextToken() token.Token {
	var tok token.Token
	lex.readChar()

	switch lex.char {
	case '=':
		tok = token.NewToken(token.ASSIGN, string(lex.char))
	case '+':
		tok = token.NewToken(token.PLUS, string(lex.char))
	case '(':
		tok = token.NewToken(token.LPAREN, string(lex.char))
	case ')':
		tok = token.NewToken(token.RPAREN, string(lex.char))
	case '{':
		tok = token.NewToken(token.LBRACE, string(lex.char))
	case '}':
		tok = token.NewToken(token.RBRACE, string(lex.char))
	case ',':
		tok = token.NewToken(token.COMMA, string(lex.char))
	case ';':
		tok = token.NewToken(token.SEMICOLON, string(lex.char))
	case 0:
		tok = token.NewToken(token.EOF, "")
	}
	return tok
}

func (lex *lexer) readChar() {
	if lex.peekPos >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.peekPos]
		lex.pos = lex.peekPos
		lex.peekPos++
	}
}
