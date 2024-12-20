package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vricap/kusmala/ast"
)

// util function to print the tree nicely
func PrintTree(tree []ast.Statement) string {
	fmt.Println("TREE:")
	space := 1
	var b bytes.Buffer
	for _, t := range tree {
		PrintNode(t, &b, space)
	}
	return b.String()
}

func PrintNode(node ast.Node, b *bytes.Buffer, space int) {
	switch node.(type) {
	case *ast.BuatStatement:
		buat := node.(*ast.BuatStatement)
		printBuatStatement(buat, b, space)
	}
	space = 1
	b.WriteString("\n")
}

func printExpression(expr ast.Expression, b *bytes.Buffer, space int) {
	switch expr.(type) {
	case *ast.IntegerLiteral:
		i := expr.(*ast.IntegerLiteral)
		printIntegerLiteral(i, b, space)
	case *ast.PrefixExpression:
		p := expr.(*ast.PrefixExpression)
		printPrefixExpression(p, b, space)
	case *ast.InfixExpression:
		in := expr.(*ast.InfixExpression)
		printInfixExpression(in, b, space)
	}
}

func printBuatStatement(buat *ast.BuatStatement, b *bytes.Buffer, space int) {
	b.WriteString(addSpace(space) + "BUAT_STATEMENT:\n")
	space++
	printIdent(buat.Name, b, space)
	printExpression(buat.Expression, b, space)
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

func addSpace(r int) string {
	s := strings.Repeat(" ", r)
	return s
}

// buat x = 2 - 3 * 4;
// 2 - 3 * 4 == (2 - (3 * 4))
// BuatStatement:
// 		Ident: x
// 		Infix:
// 			Left:
//
