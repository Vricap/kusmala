package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

type test_struct struct {
	expected string
}

// buat statement test
var buat_input string = `
buat x = -1;
buat y = 2;
buat buzz = 12345;
`
var buat_input_test_struct []test_struct = []test_struct{
	{expected: "x"},
	{expected: "y"},
	{expected: "buzz"},
}

func TestBuatStatement(t *testing.T) {
	input := buat_input
	test := buat_input_test_struct
	tree := constructTree(t, input)

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

		if buatStatementStruct.Name.Value != tt.expected {
			t.Fatalf("buatStatementStruct.Name.Value not '%s'", tt.expected)
		}

		if buatStatementStruct.Name.TokenLiteral() != tt.expected {
			t.Fatalf("s.Name not '%s'. got: %s", tt.expected, buatStatementStruct.Name.Value)
		}

	}
}

// kembalikan statement test
// kembalikan <expresion>;
var kembalikan_input string = `
kembalikan 5;
kembalikan 10;
kembalikan add(15);
kembalikan -1;
kembalikan !benar;
`

var kembalikan_input_test_struct []test_struct = []test_struct{
	{expected: "5"},
	{expected: "10"},
	{expected: "add(15)"},
	{expected: "-1"},
	{expected: "!benar"},
}

func TestKembalikanStatement(t *testing.T) {
	input := kembalikan_input
	test := kembalikan_input_test_struct
	tree := constructTree(t, input)

	if tree == nil {
		t.Fatal("ConstructTree() returned nil")
	}
	if len(tree.Statements) != 5 {
		t.Fatalf("tree.Statements does not contain 5 statement. got: %d", len(tree.Statements))
	}

	for i, _ := range test {
		each := tree.Statements[i]
		if each.TokenLiteral() != "kembalikan" {
			t.Fatalf("each.TokenLiteral() is not 'kembalikan'. got: %v", each.TokenLiteral())
		}

		_, ok := each.(*ast.KembalikanStatement)
		if !ok {
			t.Fatalf("each is not *ast.KembalikanStatement. got: %T", each)
		}
		// if tt.expected != x.Expression {
		// 	t.Fatalf("Expected: %s. got: %s", tt.expected, x.Expression)
		// }
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"
	tree := constructTree(t, input)

	if len(tree.Statements) != 1 {
		t.Fatalf("Tree has not enough statements. got: %d", len(tree.Statements))
	}

	statement, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}

	// TODO: use the helper function below
	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("statement.Expression is not *ast.Identifier. got: %T", statement.Expression)
	}
	if ident.Value != "foo" {
		t.Fatalf("ident.Value is not 'foo'. got: %s", ident.Value)
	}
	if ident.TokenLiteral() != "foo" {
		t.Fatalf("ident.TokenLiteral() is not 'foo'. got: %s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "1;"
	tree := constructTree(t, input)

	if len(tree.Statements) != 1 {
		t.Fatalf("Tree has not enough statement. got: %d", len(tree.Statements))
	}

	statement, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression not *ast.IntegerLiteral. got: %T", statement.Expression)
	}
	if literal.Value != 1 {
		t.Errorf("litereal.Value is not 1. got: %d", literal.Value)
	}
	if literal.TokenLiteral() != "1" {
		t.Errorf("literal.TokenLiteral() is not '1'. got: %s", literal.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	input := `
!5;
-15;
`
	testStruct := []struct {
		operator string
		integer  int
	}{
		{operator: "!", integer: 5},
		{operator: "-", integer: 15},
	}
	tree := constructTree(t, input)

	fmt.Printf("%T \n", tree.Statements[0])
	if len(tree.Statements) != 2 {
		t.Fatalf("len(tree.Statements) not 2. got: %d", len(tree.Statements))
	}

	// TODO: refactor
	statement1, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}
	statement2, ok := tree.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[1] is not *ast.ExpressionStatement. got: %T", tree.Statements[1])
	}

	expr1, ok := statement1.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("statement1.Expression is not *ast.PrefixExpression. got: %T", statement1.Expression)
	}
	expr2, ok := statement2.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("statement2.Expression is not *ast.PrefixExpression. got: %T", statement2.Expression)
	}

	if expr1.Operator != testStruct[0].operator {
		t.Fatalf("expr1.Operator is not %s. got: %s", testStruct[0].operator, expr1.Operator)
	}
	if expr2.Operator != testStruct[1].operator {
		t.Fatalf("expr2.Operator is not %s. got: %s", testStruct[1].operator, expr2.Operator)
	}

	if !checkIntegerLiteral(t, expr1.Right, testStruct[0].integer) {
		return
	}
	if !checkIntegerLiteral(t, expr2.Right, testStruct[1].integer) {
		return
	}
}

func TestInfinixExpression(t *testing.T) {
	infixTest := []struct {
		in     string
		expect string
	}{
		{"1 + 2", "(1 + 2)"},
		{"1 + 2 - 1", "((1 + 2) - 1)"},
		{"1 + 2 * 1", "(1 + (2 * 1))"},
		{"1 + 2 * 1 + 3", "((1 + (2 * 1)) + 3)"},
		{"9 > 2 == salah;", "((9 > 2) == salah)"},
		// {"-(5 + 5)", "(-(5 + 5))"},
	}

	for _, tt := range infixTest {
		tree := constructTree(t, tt.in)

		if len(tree.Statements) != 1 {
			t.Fatalf("len(tree.Statements) is not 1. got: %d", len(tree.Statements))
		}

		statement, ok := tree.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
		}

		infixExpr, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement.Expression is not *ast.InfixExpression. got: %T", statement.Expression)
		}
		checkInfix(infixExpr, tt.expect)
	}
}

func TestBooleanLiteral(t *testing.T) {
	input := []struct {
		in     string
		expect bool
	}{
		{"benar;", true},
		{"salah;", false},
		// {"9 > 2 == salah;", "((9 > 2) == salah)"},
		// {"1 < 2 == benar;", "((1 < 2) == benar)"},
	}

	for _, tt := range input {
		tree := constructTree(t, tt.in)

		stmnt, ok := tree.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
		}
		bool, ok := stmnt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("stmnt.Expression is not *ast.BooleanLiteral. got: %T", stmnt.Expression)
		}

		if lookUpBool(bool.Token.Literal) != tt.expect {
			t.Errorf("bool.Token.Literal is not %v. got: %v", tt.expect, lookUpBool(bool.Token.Literal))
		}
	}
}

// func TestOperatorPrecedenceParsing(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{
// 			"1 + (2 + 3) + 4",
// 			"((1 + (2 + 3)) + 4)",
// 		},
// 		{
// 			"(5 + 5) * 2",
// 			"((5 + 5) * 2)",
// 		},
// 		{
// 			"2 / (5 + 5)",
// 			"(2 / (5 + 5))",
// 		},
// 		{
// 			"-(5 + 5)",
// 			"(-(5 + 5))",
// 		},
// 		{
// 			"!(true == true)",
// 			"(!(true == true))",
// 		},
// 	}
// 	for _, tt := range tests {
// 		tree := constructTree(t, tt.input)

// 		stmnt := tree.Statements[0].(*ast.ExpressionStatement)
// 		exp := stmnt.Expression.(*ast.InfixExpression)
// 		var buffer bytes.Buffer
// 		infixTreeToString(exp, &buffer)
// 		if buffer.String() != tt.expected {
// 			t.Fatalf("buffer.String() is not %s. got: %s", tt.expected, buffer.String())
// 		}
// 	}
// }

func TestJikaStatement(t *testing.T) {
	input := `jika(x > y) {y} lainnya {buat x = 1 + 2 * 2;}`
	tree := constructTree(t, input)

	if len(tree.Statements) != 1 {
		t.Fatalf("len(tree.Statements) not 1. got: %d", len(tree.Statements))
	}
	stmnt, ok := tree.Statements[0].(*ast.JikaStatement)
	if stmnt.Token.Type != token.JIKA {
		t.Fatalf("stmnt.Token.Type is not JIKA. got: %s", stmnt.Token.Literal)
	}
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.JikaStatement. got: %T", tree.Statements[0])
	}
	cond, ok := stmnt.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmnt.Condition is not *ast.InfixExpression. got: %T", stmnt.Condition)
	}

	checkInfix(cond, "(x > y)")

	jikaBlock, ok := stmnt.JikaBlock.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmnt.JikaBlock.Statements[0] is not *ast.ExpressionStatement. got: %T", stmnt.JikaBlock.Statements[0])
	}
	checkIdent(t, jikaBlock.Expression, "y")

	lainnyaBlock, ok := stmnt.LainnyaBlock.Statements[0].(*ast.BuatStatement)
	if !ok {
		t.Fatalf("stmnt.LainnyaBlock.Statements[0] is not *ast.BuatStatementt. got: %T", stmnt.LainnyaBlock.Statements[0])
	}
	if lainnyaBlock.Token.Type != token.BUAT {
		t.Fatalf("lainnyaBlock.Token.Type is not token.BUAT. got: %v", lainnyaBlock.Token.Type)
	}
	if lainnyaBlock.Name.Value != "x" {
		t.Fatalf("lainnyaBlock.Name.Value is not 'x'. got: %s", lainnyaBlock.Name.Value)
	}
	e, ok := lainnyaBlock.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("lainnyaBlock.Expression is not *ast.InfixExpression, got: %T", lainnyaBlock.Expression)
	}

	checkInfix(e, "(1 + (2 * 2))")

	// if !checkIntegerLiteral(t, lainnyaBlock.Expression, 1) {
	// 	return
	// }

	// checkIdent(t, lainnyaBlock.Expression, "x")

}

func TestFungsiLiteral(t *testing.T) {
	input := `fungsi(x, y) { x + y; }`
	tree := constructTree(t, input)

	if len(tree.Statements) != 1 {
		t.Fatalf("len(tree.Statements) is not 1. got: %d", len(tree.Statements))
	}
	stmnt, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}
	expr, ok := stmnt.Expression.(*ast.FungsiExpression)
	if !ok {
		t.Fatalf("stmnt.Expression is not *ast.FungsiExpression. got: %T", stmnt.Expression)
	}

	checkIdent(t, expr.Params[0], "x")
	checkIdent(t, expr.Params[1], "y")

	if len(expr.Body.Statements) != 1 {
		t.Fatalf("len(expr.Body.Statements) is not 1. got: %d", len(expr.Body.Statements))
	}

	body, ok := expr.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.Body.Statements[0] is not *ast.ExpressionStatement. got: %T", expr.Body.Statements[0])
	}

	checkInfix(body.Expression, "(x + y)")
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 1 - 2)`
	tree := constructTree(t, input)

	if len(tree.Statements) != 1 {
		t.Fatalf("len(tree.Statements) is not 1. got: %d", len(tree.Statements))
	}

	stmnt, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}
	expr, ok := stmnt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmnt.Expression is not *ast.CallExpression. got: %T", stmnt.Expression)
	}
	checkIdent(t, expr.Function, "add")
	if !checkIntegerLiteral(t, expr.Arguments[0], 1) {
		return
	}
	checkInfix(expr.Arguments[1], "(2 * 3)")
	checkInfix(expr.Arguments[2], "(1 - 2)")
}

// TODO: too lazy to write the test...
func TestCetakStatement(t *testing.T) {

}

/*******************************************
*			HELPER FUNCTION 			   *
*******************************************/

func checkPeekError(t *testing.T, pars *Parser) {
	errors := pars.Errors

	if len(errors) == 0 {
		return
	}
	t.Errorf("There's %d error in parsing stage.", len(errors))
	for i, msg := range errors {
		t.Errorf("%d: %s", i+1, msg)
	}
	t.FailNow()
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, expect int) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got: %T", il)
		return false
	}
	if integer.Value != expect {
		t.Errorf("integer.Value is not %d. got: %d", expect, integer.Value)
		return false
	}
	if integer.TokenLiteral() != strconv.Itoa(expect) {
		t.Errorf("integer.TokenLiteral() is not %s. got: %s", strconv.Itoa(expect), integer.TokenLiteral())
		return false
	}
	return true
}

func checkIdent(t *testing.T, exp ast.Expression, lit string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not *ast.Identifie. got: %T", exp)
	}
	if ident.Token.Literal != lit {
		t.Fatalf("ident.Token.Literal is not %s. got: %s", lit, ident.Token.Literal)
	}
}

func checkInfix(e ast.Expression, expect string) bool {
	var buffer bytes.Buffer
	infixTreeToString(e, &buffer)
	return buffer.String() == expect
}

// helper function to turn infix tree into readable string
func infixTreeToString(e ast.Expression, buffer *bytes.Buffer) {
	exp := e.(*ast.InfixExpression)
	buffer.WriteString("(")
	le, ok := exp.Left.(*ast.InfixExpression)
	if ok {
		infixTreeToString(le, buffer)
	} else {
		buffer.WriteString(exp.Left.TokenLiteral())
	}

	buffer.WriteString(" " + exp.Operator + " ")

	re, k := exp.Right.(*ast.InfixExpression)
	if k {
		infixTreeToString(re, buffer)
	} else {
		buffer.WriteString(exp.Right.TokenLiteral())
	}
	buffer.WriteString(")")
}

func lookUpBool(lit string) bool {
	if lit == "benar" {
		return true
	} else {
		return false
	}
}

func constructTree(t *testing.T, input string) *ast.Tree {
	lex := lexer.NewLex(input)
	pars := NewPars(lex)
	tree := pars.ConstructTree()
	checkPeekError(t, pars) // check if there error in parsing stage
	return tree
}
