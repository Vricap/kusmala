package main

import (
	"bytes"
	"fmt"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/parser"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Halo %s! Ini adalah bahasa pemrograman KUSMALA!\n", user.Username)
	// fmt.Println("Silahkan untuk mengetik perintah.")
	// repl.Start(os.Stdin)

	// input := `buat x = 1 + 1 - 1;
	// buat x = 1 + 1;`
	// lex := lexer.NewLex(input)
	// pars := parser.NewPars(lex)
	// ast := pars.ParsCode()

	// for i := 0; i < len(ast.Statements); i++ {
	// 	fmt.Println(ast.Statements[i].TokenLiteral())
	// }

	// infix := ast.InfixExpression{
	// 	Left: &ast.IntegerLiteral{
	// 		Token: token.Token{
	// 			Type:    token.BILBUL,
	// 			Literal: "6",
	// 		},
	// 		Value: 6,
	// 	},
	// 	Operator: "+",
	// 	Right: &ast.InfixExpression{
	// 		Left: &ast.IntegerLiteral{
	// 			Token: token.Token{
	// 				Type:    token.BILBUL,
	// 				Literal: "3",
	// 			},
	// 			Value: 3,
	// 		},
	// 		Operator: "-",
	// 		Right: &ast.IntegerLiteral{
	// 			Token: token.Token{
	// 				Type:    token.BILBUL,
	// 				Literal: "1",
	// 			},
	// 			Value: 1,
	// 		},
	// 	},
	// }

	input := "1 > 3 == 2"
	lex := lexer.NewLex(input)
	pars := parser.NewPars(lex)
	tree := pars.ConstructTree()
	s := tree.Statements[0].(*ast.ExpressionStatement)
	infix := s.Expression.(*ast.InfixExpression)

	i := infixTreeToString(infix)
	fmt.Println(i)
}

func infixTreeToString(exp *ast.InfixExpression) string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	le, ok := exp.Left.(*ast.InfixExpression)
	if ok {
		recursiveInfix(le, &buffer)
	} else {
		left := exp.Left.(*ast.IntegerLiteral)
		buffer.WriteString(left.Token.Literal)
	}

	buffer.WriteString(exp.Operator)

	re, k := exp.Right.(*ast.InfixExpression)
	if k {
		recursiveInfix(re, &buffer)
	} else {
		right := exp.Right.(*ast.IntegerLiteral)
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
	}

	buffer.WriteString(exp.Operator)

	re, k := exp.Right.(*ast.InfixExpression)
	if k {
		recursiveInfix(re, buffer)
	} else {
		right := exp.Right.(*ast.IntegerLiteral)
		buffer.WriteString(right.Token.Literal)
		buffer.WriteString(")")
	}
}
