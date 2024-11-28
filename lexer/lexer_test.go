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
	input := "=+(){},;"

	test := []testStruct{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
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
