package token

type TokenType string // its actually more performant using byte or int

// define all token the language will have
const (
	ILLEGAL TokenType = "ILLEGAL" // token that didn't define in the language
	EOF     TokenType = "EOF"

	IDENT  TokenType = "IDENT"  // e.g variable name
	BILBUL TokenType = "BILBUL" // bilangan bulan / int

	// operator
	ASSIGN TokenType = "="
	PLUS   TokenType = "PLUS"

	// delimiter
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	// keyword (reserved word specific to the programming language - non user-defined)
	FUNGSI TokenType = "FUNGSI"
	BUAT   TokenType = "BUAT"
)

type Token struct {
	Type    TokenType
	Literal string
}
