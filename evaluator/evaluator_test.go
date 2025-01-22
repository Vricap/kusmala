package evaluator

import (
	"testing"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/object"
	"github.com/vricap/kusmala/parser"
)

func testVal(in string) object.Object {
	lex := lexer.NewLex(in)
	pars := parser.NewPars(lex)
	tree := pars.ConstructTree()

	return Eval(tree)
}

func TestIntegerExpression(t *testing.T) {
	test := []struct {
		in     string
		expect int
	}{
		{"4", 4},
		{"12", 12},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"3 * 3 * 3 + 10", 37},
	}

	for _, tt := range test {
		eval := testVal(tt.in)
		testIntegerObject(t, eval, tt.expect)
	}
}

func TestBoolean(t *testing.T) {
	test := []struct {
		in     string
		expect bool
	}{
		{"benar", true},
		{"salah", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"benar == salah", false},
		{"salah == salah", true},
		{"benar == salah", false},
		{"benar != salah", true},
		{"salah != benar", true},
		// {"1 == benar", true},
		// {"1 != benar", false},
	}

	for _, tt := range test {
		eval := testVal(tt.in)
		testBooleanObject(t, eval, tt.expect)
	}
}

func TestBangOperator(t *testing.T) {
	test := []struct {
		in     string
		expect bool
	}{
		{"!benar", false},
		{"!salah", true},
		{"!1", false},
		{"!!benar", true},
		{"!!salah", false},
		{"!!1", true},
	}

	for _, tt := range test {
		eval := testVal(tt.in)
		testBooleanObject(t, eval, tt.expect)
	}
}

func testIntegerObject(t *testing.T, eval object.Object, expect int) {
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Fatalf("object is not *object.Integer. got: %T", eval)
	}
	if i.Value != expect {
		t.Fatalf("i.Value is not %d. got: %d", expect, i.Value)
	}
}

func testBooleanObject(t *testing.T, eval object.Object, expect bool) {
	b, ok := eval.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not *object.Boolean.got: %T", eval)
	}
	if b.Value != expect {
		t.Fatalf("i.Value is not %t. got: %t", expect, b.Value)
	}
}
