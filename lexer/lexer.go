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
	l := &lexer{input: input}
	l.readChar() // so that l.char point to the actual first char in input and not just 0
	return l
}

func (lex *lexer) NextToken() token.Token {
	var tok token.Token
	lex.skipWhiteSpace()

	switch lex.char {
	case '=':
		if lex.peekChar() == '=' { // if equal ==
			tok = token.NewToken(token.SAMA, "==")
			lex.pos++
			lex.peekPos++
		} else {
			tok = token.NewToken(token.ASSIGN, string(lex.char))
		}
	case '+':
		tok = token.NewToken(token.PLUS, string(lex.char))
	case '-':
		tok = token.NewToken(token.MINUS, string(lex.char))
	case '!':
		if lex.peekChar() == '=' { // if not equal !=
			tok = token.NewToken(token.TIDAK_SAMA, "!=")
			lex.pos++
			lex.peekPos++
		} else {
			tok = token.NewToken(token.BANG, string(lex.char))
		}
	case '/':
		tok = token.NewToken(token.SLASH, string(lex.char))
	case '*':
		tok = token.NewToken(token.ASTERISK, string(lex.char))
	case '<':
		tok = token.NewToken(token.LT, string(lex.char))
	case '>':
		tok = token.NewToken(token.GT, string(lex.char))
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
	default:
		if isLetter(lex.char) {
			tok.Literal = lex.readIdentifier()
			tokType := token.LookUpIdent(tok.Literal) // check wether the word is keyword or just identifier
			tok = token.NewToken(tokType, tok.Literal)
			return tok // return early so that readChar at the bottom didn't run again. the pos and peekPos is move up since we already readChar repeatedly inside lex.readIdentifier()
		} else if isDigit(lex.char) {
			tok.Literal = lex.readNumber()
			// fmt.Println(string(lex.char))
			// fmt.Println(tok.Literal)
			tok = token.NewToken(token.BILBUL, tok.Literal)
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, string(lex.char))
		}
	}
	lex.readChar()
	return tok
}

func (lex *lexer) readIdentifier() string {
	pos := lex.pos
	for isLetter(lex.char) {
		lex.readChar()
	}
	return lex.input[pos:lex.pos]
}

func (lex *lexer) readNumber() string {
	startPos := lex.pos
	for isDigit(lex.char) {
		lex.readChar()
	}
	return lex.input[startPos:lex.pos]
}

func (lex *lexer) readChar() {
	if lex.peekPos >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.peekPos]
	}
	lex.pos = lex.peekPos
	lex.peekPos++
}

func (lex *lexer) skipWhiteSpace() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
	}
}

func (lex *lexer) peekChar() byte {
	if lex.peekPos >= len(lex.input) {
		return 0
	} else {
		return lex.input[lex.peekPos]
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9' // '0' corresponds to ASCII value 48. '9' corresponds to ASCII value 57.
}
