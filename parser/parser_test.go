package parser

import (
	"testing"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
)

type test_struct struct {
	expectedIdent string
}

var input_one string = `
buat x = 1;
buat y = 2;
buat buzz = 12345;
`
var input_one_test_struct []test_struct = []test_struct{
	{expectedIdent: "x"},
	{expectedIdent: "y"},
	{expectedIdent: "buzz"},
}

func TestBuatStatement(t *testing.T) {
	input := input_one
	test := input_one_test_struct

	lex := lexer.NewLex(input)
	pars := NewPars(lex)

	code := pars.ParsCode()
	checkPeekError(t, pars) // check if there error in parsing stage

	if code == nil {
		t.Fatal("ParsCode() returned nil")
	}
	if len(code.Statements) != 3 {
		t.Fatalf("code.Statements does not contain 3 statement. got: %d", len(code.Statements))
	}

	for i, tt := range test {
		s := code.Statements[i]

		if s.TokenLiteral() != "buat" {
			t.Fatalf("s.TokenLiteral() is not 'buat'. got: %v", s.TokenLiteral())
		}

		buatStatementStruct, ok := s.(*ast.BuatStatement) // type assertion. get the underlying concrete type (BuatStatement) from  s (Statement) inteface
		if !ok {
			t.Fatalf("s is not *ast.BuatStatement. got: %T", s)
		}

		if buatStatementStruct.Name.Value != tt.expectedIdent {
			t.Fatalf("buatStatementStruct.Name.Value not '%s'", tt.expectedIdent)
		}

		if buatStatementStruct.Name.TokenLiteral() != tt.expectedIdent {
			t.Fatalf("s.Name not '%s'. got: %s", tt.expectedIdent, buatStatementStruct.Name)
		}

	}
}

func checkPeekError(t *testing.T, pars *Parser) {
	errors := pars.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("There's %d error in parsing stage.", len(errors))
	for i, msg := range errors {
		t.Errorf("%d: %s", i+1, msg)
	}
	t.FailNow()
}
