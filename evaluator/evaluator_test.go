package evaluator

import (
	"testing"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/object"
	"github.com/vricap/kusmala/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	test := []struct {
		in     string
		expect int
	}{
		{"4", 4},
		{"12", 12},
	}

	for _, tt := range test {
		eval := testVal(tt.in)
		testIntegerObject(t, eval, tt.expect)
	}
}

func testVal(in string) object.Object {
	lex := lexer.NewLex(in)
	pars := parser.NewPars(lex)
	tree := pars.ConstructTree()

	return Eval(tree)
}

func testIntegerObject(t *testing.T, eval object.Object, expect int) {
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Fatalf("object is not *object.Integer. got: %T", i)
	}
	if i.Value != expect {
		t.Fatalf("i.Value is not %d. got: %d", expect, i.Value)
	}
}
