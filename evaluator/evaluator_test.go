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
	env := object.NewEnv()

	// TODO: this fuck my test case
	eval := Eval(tree, env)
	return eval[len(eval)-1]
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

func TestJikaStatement(t *testing.T) {
	test := []struct {
		in     string
		expect any
	}{
		{"jika (benar) { 10 }", 10},
		{"jika (salah) { 10 }", nil},
		{"jika (1) { 10 }", 10},
		{"jika (1 < 2) { 10 }", 10},
		{"jika (1 > 2) { 10 }", nil},
		{"jika (1 > 2) { 10 } lainnya { 20 }", 20},
		{"jika (1 < 2) { 10 } lainnya { 20 }", 10},
	}
	for _, tt := range test {
		eval := testVal(tt.in)
		int, ok := tt.expect.(int)
		if ok {
			testIntegerObject(t, eval, int)
		} else {
			testNilObject(t, eval)
		}
	}
}

func TestKembalikanStatement(t *testing.T) {
	test := []struct {
		in     string
		expect int
	}{
		{`
		jika (10 > 1) {
			jika (10 > 1) {
				jika (3 < 2) {
					kembalikan 10;
				} lainnya {
					kembalikan 4;
				}
			}
			jika (1 < 10) {
				2;
				kembalikan 3;
			}
			129;
			kembalikan 1;
			2;
		}`, 4},
		{
			`
		jika (1 > 3) {
			kembalikan 1;
		} lainnya {
			jika (2 < 3) {
				kembalikan 2;
			}
			kembalikan 3;
		}`, 2,
		},
	}

	for _, tt := range test {
		eval := testVal(tt.in)
		testIntegerObject(t, eval, tt.expect)
	}
}

func TestErrorHandling(t *testing.T) {
	test := []struct {
		in     string
		expect string
	}{
		{"1 + benar;", "ERROR di baris 1: kesalahan tipe dekat '1 + benar'"},
		{"1 + benar;", "ERROR di baris 1: kesalahan tipe dekat '1 + benar'"},
		{"salah + benar;", "ERROR di baris 1: operator tidak didukung dekat 'salah + benar'"},
		{"-benar;", "ERROR di baris 1: operator tidak didukung dekat '-benar'"},
		{"salah + benar;", "ERROR di baris 1: operator tidak didukung dekat 'salah + benar'"},
		{"foo;", "ERROR di baris 1: pengenal tidak diketahui dekat 'foo'"},
		{`
		jika(1 < 2) {
			jika (1 < 3) {
				kembalikan benar - salah;
			}
			kembalikan 1;
		}
		`, "ERROR di baris 4: operator tidak didukung dekat 'benar - salah'"},
	}
	for i, tt := range test {
		eval := testVal(tt.in)
		e, ok := eval.(*object.Error)
		if !ok {
			t.Errorf("ERROR in (%d): eval is not *object.Error. got: %T", i, eval)
		}
		if e.Inspect() != tt.expect {
			t.Fatalf("e.Msg is not: '%s'. got: %s", tt.expect, e.Inspect())
		}
	}
}

func TestBuatStatement(t *testing.T) {
	test := []struct {
		in     string
		expect int
	}{
		{"buat a = 5; a;", 5},
		{"buat a = 5 * 5; a;", 25},
		{"buat a = 5; buat b = a; b;", 5},
		{`buat a = 5; buat b = a; buat c = a + b + 5;`, 15},
	}
	for _, tt := range test {
		eval := testVal(tt.in)
		i, ok := eval.(*object.Integer)
		if !ok {
			t.Fatalf("eval is not *object.Integer. got: %T", eval)
		}
		testIntegerObject(t, i, tt.expect)
	}
}

func TestFungsiLiteral(t *testing.T) {
	in := `fungsi(x) { x + 2; };`
	eval := testVal(in)
	fl, ok := eval.(*object.FungsiLiteral)
	if !ok {
		t.Fatalf("eval is not *object.FungsiLiteral. got: %T", eval)
	}
	if len(fl.Param) != 1 {
		t.Fatalf("len(fl.Param) is not 1. got: %d", len(fl.Param))
	}
	if fl.Param[0].Value != "x" {
		t.Fatalf("fl.Param[0].Value is not 'x'. got: %s", fl.Param[0].Value)
	}
	// expectBody := `(x + 2)`
	// if fl.Body.TokenLiteral() != expectBody {
	// 	t.Fatalf("fl.Body.TokenLiteral() is not %s. got: %s", expectBody, fl.Body.TokenLiteral())
	// }

}

func TestFungsiCall(t *testing.T) {
	test := []struct {
		in     string
		expect int
	}{
		{"buat identity = fungsi(x) { x; }; identity(5);", 5},
		{"buat identity = fungsi(x) { kembalikan x; }; identity(5);", 5},
		{"buat double = fungsi(x) { x * 2; }; double(5);", 10},
		{"buat add = fungsi(x, y) { x + y; }; add(5, 5);", 10},
		{"buat add = fungsi(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fungsi(x) { x; }(5)", 5},
	}
	for _, tt := range test {
		eval := testVal(tt.in)
		testIntegerObject(t, eval, tt.expect)
	}
}

func TestClosures(t *testing.T) {
	input := `
buat newAdder = fungsi(x) {
fungsi(y) { x + y };
};
buat addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testVal(input), 4)
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

func testNilObject(t *testing.T, eval object.Object) {
	_, ok := eval.(*object.Nil)
	if !ok {
		t.Fatalf("eval is not NIL. got %T", eval)
	}
}
