package lexer

import (
	"github.com/vricap/kusmala/token"
)

type Lexer struct {
	input   string // the whole code input
	pos     int    // current position in the input - point to current char
	peekPos int    // peek the next of the current position
	char    byte   // current char under examination
	Line    int
}

func NewLex(input string) *Lexer {
	l := &Lexer{input: input, Line: 1}
	l.readChar() // so that l.char point to the actual first char in input and not just 0
	return l
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token
	lex.skipWhiteSpace()
	if lex.isComment() {
		for lex.isComment() {
			lex.skipComment()
			lex.skipWhiteSpace()
		}
	}

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
	case '[':
		tok = token.NewToken(token.LBRACKET, string(lex.char))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(lex.char))
	case ',':
		tok = token.NewToken(token.COMMA, string(lex.char))
	case ';':
		tok = token.NewToken(token.SEMICOLON, string(lex.char))
	case '"':
		tok = token.NewToken(token.STRING, lex.readString())
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

func (lex *Lexer) readIdentifier() string {
	pos := lex.pos
	for isLetter(lex.char) {
		lex.readChar()
	}
	return lex.input[pos:lex.pos]
}

func (lex *Lexer) readNumber() string {
	startPos := lex.pos
	for isDigit(lex.char) {
		lex.readChar()
	}
	return lex.input[startPos:lex.pos]
}

func (lex *Lexer) readChar() {
	if lex.peekPos >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.peekPos]
	}
	lex.pos = lex.peekPos
	lex.peekPos++
}

// TODO: quick hack. refactor
func (lex *Lexer) skipWhiteSpace() {
	if lex.char == ' ' || lex.char == '\t' {
		lex.skipSpace()
	} else if lex.char == '\n' || lex.char == '\r' {
		lex.skipNewLine()
	}
}

func (lex *Lexer) skipSpace() {
	for lex.char == ' ' || lex.char == '\t' {
		lex.readChar()
		if lex.char == '\n' || lex.char == '\r' {
			lex.skipNewLine()
		}
	}
}

func (lex *Lexer) skipNewLine() {
	for lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
		lex.Line++
		if lex.char == ' ' || lex.char == '\t' {
			lex.skipSpace()
		}
	}
}

func (lex *Lexer) isComment() bool {
	return lex.char == '/' && lex.peekChar() == '/'
}

func (lex *Lexer) skipComment() {
	for lex.char != 10 { // 10 ascii code for new line
		lex.readChar()
	}
}

func (lex *Lexer) peekChar() byte {
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

func (lex *Lexer) readString() string {
	var str string
	lex.readChar()
	for lex.char != 34 {
		str += string(lex.char)
		lex.readChar()
	}
	return str
}
