package parser

import (
	"testing"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
)

type test_struct struct {
	expectedIdent string
}

// buat statement test
var buat_input string = `
buat x = 1;
buat y = 2;
buat buzz = 12345;
`
var buat_input_test_struct []test_struct = []test_struct{
	{expectedIdent: "x"},
	{expectedIdent: "y"},
	{expectedIdent: "buzz"},
}

func TestBuatStatement(t *testing.T) {
	input := buat_input
	test := buat_input_test_struct

	lex := lexer.NewLex(input)
	pars := NewPars(lex)

	tree := pars.ConstructTree()
	checkPeekError(t, pars) // check if there error in parsing stage

	if tree == nil {
		t.Fatal("ParsCode() returned nil")
	}
	if len(tree.Statements) != 3 {
		t.Fatalf("tree.Statements does not contain 3 statement. got: %d", len(tree.Statements))
	}

	for i, tt := range test {
		s := tree.Statements[i]

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

// kembalikan statement test
// kembalikan <expresion>;
// var kembalikan_input string = `
// kembalikan 5;
// kembalikan 10;
// kembalikan add(15);
// `

// var kembalikan_input_test_struct []test_struct = []test_struct{
// 	{expectedIdent: "5"},
// 	{expectedIdent: "10"},
// 	{expectedIdent: "add(15)"},
// }

// func TestKembalikanStatement(t *testing.T) {
// 	input := kembalikan_input
// 	test := kembalikan_input_test_struct

// 	lex := lexer.NewLex(input)
// 	pars := NewPars(lex)
// 	code := pars.ParsCode()
// }
