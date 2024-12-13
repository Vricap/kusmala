package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
)

type test_struct struct {
	expected string
}

// buat statement test
var buat_input string = `
buat x = 1;
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

		if buatStatementStruct.Name.Value != tt.expected {
			t.Fatalf("buatStatementStruct.Name.Value not '%s'", tt.expected)
		}

		if buatStatementStruct.Name.TokenLiteral() != tt.expected {
			t.Fatalf("s.Name not '%s'. got: %s", tt.expected, buatStatementStruct.Name)
		}

	}
}

// kembalikan statement test
// kembalikan <expresion>;
var kembalikan_input string = `
kembalikan 5;
kembalikan 10;
kembalikan add(15);
`

var kembalikan_input_test_struct []test_struct = []test_struct{
	{expected: "5"},
	{expected: "10"},
	{expected: "add(15)"},
}

func TestKembalikanStatement(t *testing.T) {
	input := kembalikan_input
	test := kembalikan_input_test_struct

	lex := lexer.NewLex(input)
	pars := NewPars(lex)
	tree := pars.ConstructTree()
	checkPeekError(t, pars) // check if there error in parsing stage

	if tree == nil {
		t.Fatal("ConstructTree() returned nil")
	}
	if len(tree.Statements) != 3 {
		t.Fatalf("tree.Statements does not contain 3 statement. got: %d", len(tree.Statements))
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

	lex := lexer.NewLex(input)
	pars := NewPars(lex)

	tree := pars.ConstructTree()
	checkPeekError(t, pars)

	if len(tree.Statements) != 1 {
		t.Fatalf("Tree has not enough statements. got: %d", len(tree.Statements))
	}

	statement, ok := tree.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("tree.Statements[0] is not *ast.ExpressionStatement. got: %T", tree.Statements[0])
	}

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
	lex := lexer.NewLex(input)
	pars := NewPars(lex)

	tree := pars.ConstructTree()
	checkPeekError(t, pars)

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

	lex := lexer.NewLex(input)
	pars := NewPars(lex)
	tree := pars.ConstructTree()
	checkPeekError(t, pars)

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
		input    string
		left     int
		operator string
		right    int
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTest {
		lex := lexer.NewLex(tt.input)
		pars := NewPars(lex)
		tree := pars.ConstructTree()
		checkPeekError(t, pars)

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
		if !checkIntegerLiteral(t, infixExpr.Left, tt.left) {
			return
		}
		if infixExpr.Operator != tt.operator {
			t.Fatalf("infixExpr.Operator is not %s. got: %s", tt.operator, infixExpr.Operator)
		}
		if !checkIntegerLiteral(t, infixExpr.Right, tt.right) {
			return
		}
	}
}

// func TestBooleanLiteral(t *testing.T) {
// 	input := []struct {
// 		in     string
// 		expect string
// 	}{
// 		{"benar;", "benar"},
// 		{"salah;", "salah"},
// 		{"salah;", "salah"},
// 	}
// }

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

func InfixTreeToString(exp *ast.InfixExpression) string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	le, ok := exp.Left.(*ast.InfixExpression)
	if ok {
		recursiveInfix(le, &buffer)
	} else {
		// TODO: what if left is *ast.PrefixExpression?
		left := exp.Left.(*ast.IntegerLiteral)
		buffer.WriteString(left.Token.Literal)
		buffer.WriteString(exp.Operator)
	}

	re, k := exp.Right.(*ast.InfixExpression)
	if k {
		recursiveInfix(re, &buffer)
	} else {
		// TODO: what if left is *ast.PrefixExpression?
		right := exp.Left.(*ast.IntegerLiteral)
		buffer.WriteString(right.Token.Literal)
		buffer.WriteString(")")
	}
	return buffer.String()
}

func recursiveInfix(exp *ast.InfixExpression, buffer *bytes.Buffer) {
	buffer.WriteString("(")
	le, ok := exp.Left.(*ast.InfixExpression)
	if ok {
		recursiveInfix(le, buffer)
	} else {
		left := exp.Left.(*ast.IntegerLiteral)
		buffer.WriteString(left.Token.Literal)
		buffer.WriteString(exp.Operator)
	}

	re, k := exp.Right.(*ast.InfixExpression)
	if k {
		recursiveInfix(re, buffer)
	} else {
		right := exp.Left.(*ast.IntegerLiteral)
		buffer.WriteString(right.Token.Literal)
		buffer.WriteString(")")
	}
}
