package token

type TokenType string // its actually more performant using byte or int

// define all token the language will have
const (
	ILLEGAL TokenType = "ILLEGAL" // token that didn't define in the language
	EOF     TokenType = "EOF"

	IDENT  TokenType = "IDENT"  // user-defined. e.g variable name
	BILBUL TokenType = "BILBUL" // bilangan bulan / int

	// operator
	ASSIGN     TokenType = "="
	PLUS       TokenType = "PLUS"
	MINUS      TokenType = "-"
	BANG       TokenType = "!"
	ASTERISK   TokenType = "*"
	SLASH      TokenType = "/"
	LT         TokenType = "<"
	GT         TokenType = ">"
	SAMA       TokenType = "=="
	TIDAK_SAMA TokenType = "!="

	// delimiter
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	RBRACE TokenType = "{"
	LBRACE TokenType = "}"

	// keyword (reserved word specific to the programming language - non user-defined)
	FUNGSI     TokenType = "FUNGSI"
	BUAT       TokenType = "BUAT"
	BENAR      TokenType = "BENAR"
	SALAH      TokenType = "SALAH"
	JIKA       TokenType = "JIKA"
	LAINNYA    TokenType = "LAINNYA"
	KEMBALIKAN TokenType = "KEMBALIKAN"
)

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(t TokenType, lit string) Token {
	return Token{Type: t, Literal: lit}
}

var keywords = map[string]TokenType{
	"fungsi":     FUNGSI,
	"buat":       BUAT,
	"benar":      BENAR,
	"salah":      SALAH,
	"jika":       JIKA,
	"lainnya":    LAINNYA,
	"kembalikan": KEMBALIKAN,
}

func LookUpIdent(lit string) TokenType {
	tok, ok := keywords[lit]
	if ok {
		return tok
	}
	return IDENT
}
