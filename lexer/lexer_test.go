package lexer

import (
	"testing"

	"github.com/vricap/kusmala/token"
)

type testStruct struct {
	expectedType    token.TokenType
	expectedLiteral string
}

var input_one string = `buat five = 5;
buat ten = 10;
buat add = fungsi(x, y) {
x + y;
};
buat result = add(five, ten);
`

var input_one_test_struct []testStruct = []testStruct{
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

var input_two string = `buat five = 5;
buat ten = 10;
buat add = fungsi(x, y) {
x + y;
};
buat result = add(five, ten);
!-/*5;
5 < 10 > 5;
`

var input_two_test_struct []testStruct = []testStruct{
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

	{token.BANG, "!"},
	{token.MINUS, "-"},
	{token.SLASH, "/"},
	{token.ASTERISK, "*"},
	{token.BILBUL, "5"},
	{token.SEMICOLON, ";"},

	{token.BILBUL, "5"},
	{token.LT, "<"},
	{token.BILBUL, "10"},
	{token.GT, ">"},
	{token.BILBUL, "5"},
	{token.SEMICOLON, ";"},
}

var input_three string = `buat five = 5;
buat ten = 10;
buat add = fungsi(x, y) {
x + y;
};
buat result = add(five, ten);
!-/*5;
5 < 10 > 5;
jika (5 < 10) {
kembalikan benar;
} lainnya {
kembalikan salah;
}
10 == 11;
10 != 9;
cetak
`

var input_three_test_struct []testStruct = []testStruct{
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

	{token.BANG, "!"},
	{token.MINUS, "-"},
	{token.SLASH, "/"},
	{token.ASTERISK, "*"},
	{token.BILBUL, "5"},
	{token.SEMICOLON, ";"},

	{token.BILBUL, "5"},
	{token.LT, "<"},
	{token.BILBUL, "10"},
	{token.GT, ">"},
	{token.BILBUL, "5"},
	{token.SEMICOLON, ";"},

	{token.JIKA, "jika"},
	{token.LPAREN, "("},
	{token.BILBUL, "5"},
	{token.LT, "<"},
	{token.BILBUL, "10"},
	{token.RPAREN, ")"},
	{token.LBRACE, "{"},

	{token.KEMBALIKAN, "kembalikan"},
	{token.BENAR, "benar"},
	{token.SEMICOLON, ";"},

	{token.RBRACE, "}"},
	{token.LAINNYA, "lainnya"},
	{token.LBRACE, "{"},

	{token.KEMBALIKAN, "kembalikan"},
	{token.SALAH, "salah"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},

	{token.BILBUL, "10"},
	{token.SAMA, "=="},
	{token.BILBUL, "11"},
	{token.SEMICOLON, ";"},
	{token.BILBUL, "10"},
	{token.TIDAK_SAMA, "!="},
	{token.BILBUL, "9"},
	{token.SEMICOLON, ";"},

	{token.CETAK, "cetak"},
}

func TestNextToken(t *testing.T) {
	input := input_three // the code input

	test := input_three_test_struct // the expected token

	lex := NewLex(input)

	for i, tokTest := range test {
		tok := lex.NextToken()

		if tok.Type != tokTest.expectedType {
			t.Fatalf("tokenType wrong at [%d] - expected (%s), got (%s) val(%s)", i, tokTest.expectedType, tok.Type, tokTest.expectedLiteral)
		}

		if tok.Literal != tokTest.expectedLiteral {
			t.Fatalf("tokenLiteral wrong at [%d] - expected (%s), got (%s)", i, tokTest.expectedLiteral, tok.Literal)
		}
	}
}
