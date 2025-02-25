package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vricap/kusmala/ast"
)

/*******************************************
* util function to print the tree nicely   *
*******************************************/
func PrintTree(tree []ast.Statement) {
	fmt.Println("AST_TREE:")
	space := 1
	var b bytes.Buffer
	for _, t := range tree {
		printStatement(t, &b, space)
	}
	fmt.Println(b.String())
}

func printStatement(s ast.Statement, b *bytes.Buffer, space int) {
	switch s.(type) {
	case *ast.BuatStatement:
		bu := s.(*ast.BuatStatement)
		printBuatStatement(bu, b, space)
	case *ast.KembalikanStatement:
		k := s.(*ast.KembalikanStatement)
		printKembalikanStatement(k, b, space)
	case *ast.JikaStatement:
		j := s.(*ast.JikaStatement)
		printJikaStatement(j, b, space)
	case *ast.ExpressionStatement:
		e := s.(*ast.ExpressionStatement)
		printExpression(e.Expression, b, space)
	case *ast.CetakStatement:
		c := s.(*ast.CetakStatement)
		printCetakStatement(c, b, space)
	}
	space = 1
	b.WriteString("\n")
}

// TODO: panjang builtin function and array is not yet implemented here
func printExpression(expr ast.Expression, b *bytes.Buffer, space int) {
	switch expr.(type) {
	case *ast.Identifier:
		i := expr.(*ast.Identifier)
		printIdent(i, b, space)
	case *ast.IntegerLiteral:
		i := expr.(*ast.IntegerLiteral)
		printIntegerLiteral(i, b, space)
	case *ast.PrefixExpression:
		p := expr.(*ast.PrefixExpression)
		printPrefixExpression(p, b, space)
	case *ast.InfixExpression:
		in := expr.(*ast.InfixExpression)
		printInfixExpression(in, b, space)
	case *ast.BooleanLiteral:
		bl := expr.(*ast.BooleanLiteral)
		printBooleanLiteral(bl, b, space)
	case *ast.FungsiExpression:
		f := expr.(*ast.FungsiExpression)
		printFungsiExpression(f, b, space)
	case *ast.CallExpression:
		c := expr.(*ast.CallExpression)
		printCallExpression(c, b, space)
	case *ast.StringLiteral:
		s := expr.(*ast.StringLiteral)
		printStringLiteral(s, b, space)
	}
}

func printBuatStatement(bu *ast.BuatStatement, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "BUAT_STATEMENT:\n")
	space++
	printIdent(bu.Name, b, space)
	printExpression(bu.Expression, b, space)
}

func printKembalikanStatement(k *ast.KembalikanStatement, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "KEMBALIKAN_STATEMENT:\n")
	space++
	printExpression(k.Expression, b, space)
}

func printJikaStatement(j *ast.JikaStatement, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "JIKA_STATEMENT:\n")
	space++
	b.WriteString(addSpace(space) + "CONDITION:\n")
	space++
	printExpression(j.Condition, b, space)
	b.WriteString(addSpace(space) + "JIKA_BLOCK: \n")
	space++
	printBlockStatement(j.JikaBlock, b, space)
	if j.LainnyaBlock != nil {
		// remove the last "\n" from buffer
		rmBuffNl(b)
		space--
		b.WriteString(addSpace(space) + "LAINNYA_BLOCK: \n")
		space++
		printBlockStatement(j.LainnyaBlock, b, space)
	}
}

func printIdent(ident *ast.Identifier, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "IDENT: " + ident.Value + "\n")
}

func printInfixExpression(i *ast.InfixExpression, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "INFIX_EXPRESSION:\n")
	space++
	printExpression(i.Left, b, space)
	b.WriteString(addSpace(space) + "OEPERATOR: " + i.Operator + "\n")
	printExpression(i.Right, b, space)
}

func printPrefixExpression(p *ast.PrefixExpression, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "PREFIX_EXPRESSION:\n")
	space++
	b.WriteString(addSpace(space) + "OPERATOR: " + p.Operator + "\n")
	printExpression(p.Right, b, space)
}

func printIntegerLiteral(i *ast.IntegerLiteral, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "INTEGER_LITERAL: " + i.Token.Literal + "\n")
}

func printStringLiteral(s *ast.StringLiteral, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "STRING_LITERAL: " + s.Token.Literal + "\n")
}

func printBooleanLiteral(bl *ast.BooleanLiteral, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "BOOLEAN_LITERAL: " + bl.Token.Literal + "\n")
}

func printFungsiExpression(f *ast.FungsiExpression, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "FUNGSI_EXPRESSION: \n")
	space++
	if f.Params != nil {
		b.WriteString(addSpace(space) + "PARAMS: \n")
		space++
		printParams(f.Params, b, space)
		space--
	}
	b.WriteString(addSpace(space) + "FUNGSI_BODY: \n")
	space++
	printBlockStatement(f.Body, b, space)

}

func printCallExpression(c *ast.CallExpression, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "CALL_EXPRESSION: \n")
	space++
	printExpression(c.Function, b, space)
	b.WriteString(addSpace(space) + "ARGUMENTS: \n")
	space++
	printArguments(c.Arguments, b, space)
}

func printBlockStatement(be *ast.BlockStatement, b *bytes.Buffer, space int) {
	for _, s := range be.Statements {
		printStatement(s, b, space)
	}
}

func printParams(i []*ast.Identifier, b *bytes.Buffer, space int) {
	for _, n := range i {
		printIdent(n, b, space)
	}
}

func printArguments(a []ast.Expression, b *bytes.Buffer, space int) {
	for _, e := range a {
		printExpression(e, b, space)
	}
}

func printCetakStatement(c *ast.CetakStatement, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "CETAK_STATEMENT: \n")
	space++
	b.WriteString(addSpace(space) + "EXPRESSION: \n")
	space++
	printArguments(c.Expression, b, space)
}

func addSpace(r int) string {
	s := strings.Repeat("  ", r)
	return s
}

func rmBuffNl(b *bytes.Buffer) {
	c := b.Bytes()
	b.Truncate(len(c) - 1)
}
