package lexer

import (
	"testing"

	"github.com/vricap/kusmala/token"
)

type testStruct struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `buat five = 5;
buat ten = 10;
buat add = fungsi(x, y) {
x + y;
};
buat result = add(five, ten);
`

	test := []testStruct{
		{token.BUAT, "buat"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.BILBUL, "5"},
		{token.SEMICOLON, ";"},
		{token.BUAT, "buat"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.BILBUL, "10"},
		{token.SEMICOLON, ";"},
		{token.BUAT, "buat"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNGSI, "fungsi"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.BUAT, "buat"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := NewLex(input)

	for i, tokTest := range test {
		tok := lex.NextToken()

		if tok.Type != tokTest.expectedType {
			t.Fatalf("tokenType wrong at [%d] - expected (%s), got (%s)", i, tokTest.expectedType, tok.Type)
		}

		if tok.Literal != tokTest.expectedLiteral {
			t.Fatalf("tokenLiteral wrong at [%d] - expected (%s), got (%s)", i, tokTest.expectedLiteral, tok.Literal)
		}
	}
}
