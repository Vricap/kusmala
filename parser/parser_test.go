package parser

import (
	"testing"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
)

type testStruct struct {
	expectedIdentifier string
}

var input_one string = `
buat x = 5;
buat y = 10;
buat foobar = 838383;
`

var input_one_test_struct = []testStruct{
	{"x"},
	{"y"},
	{"foobar"},
}

func TestBuatStatement(t *testing.T) {
	input := input_one
	test := input_one_test_struct

	lex := lexer.NewLex(input)
	pars := NewPars(lex)

	program := pars.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nill")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statement. got %d", len(program.Statements))
	}

	for i, tt := range test {
		statement := program.Statements[i]
		if !testBuatStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testBuatStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "buat" {
		t.Errorf("s.TokenLiteral not 'buat'. got %q", s.TokenLiteral())
		return false
	}

	buatStatement, ok := s.(*ast.BuatStatement) // gobbledygook
	if !ok {
		t.Errorf("s not *ast.BuatStatement. got %T", s)
		return false
	}

	if buatStatement.Name.Value != name {
		t.Errorf("buatStatement.Name.Value not %s. got %s", name, buatStatement.Name.Value)
		return false
	}

	if buatStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %s. got %s", name, buatStatement.Name)
		return false
	}
	return true
}
