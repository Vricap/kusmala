package evaluator

import (
	"fmt"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/object"
)

func Eval(tree *ast.Tree) []object.Object {
	var evals []object.Object
	for _, s := range tree.Statements {
		evals = append(evals, evalStatement(s))
	}
	return evals
}

func evalStatement(stmt ast.Statement) object.Object {
	switch s := stmt.(type) {
	case *ast.BuatStatement:
	// TODO: kembalikan statement isn't allowed in global scope. only inside a block statement
	// case *ast.KembalikanStatement:
	// 	return evalKembalikanStatement(s)
	case *ast.JikaStatement:
		return evalJikaStatement(s)
	case *ast.BlockStatement:
		return evalBlockStatement(s)
	case *ast.ExpressionStatement:
		return evalExpression(s.Expression)
	case *ast.CetakStatement:
		return evalCetakStatement(s)
	}
	return &object.Nil{}
}

func evalExpression(expr ast.Expression) object.Object {
	switch e := expr.(type) {
	case *ast.Identifier:
	case *ast.IntegerLiteral:
		return &object.Integer{Value: e.Value}
	case *ast.PrefixExpression:
		right := evalExpression(e.Right)
		return evalPrefixExpression(e.Operator, right)
	case *ast.InfixExpression:
		left := evalExpression(e.Left)
		right := evalExpression(e.Right)
		return evalInfixExpression(e.Operator, left, right)
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: e.Value}
	case *ast.FungsiExpression:
	case *ast.CallExpression:
	}
	return &object.Nil{}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		if right.Inspect() == "benar" {
			return &object.Boolean{Value: false}
		} else if right.Inspect() == "salah" {
			return &object.Boolean{Value: true}
		} else {
			return &object.Boolean{Value: false}
		}
	case "-":
		if right.Type() != object.OBJECT_INTEGER {
			return &object.Nil{}
		}
		i := right.(*object.Integer)
		return &object.Integer{Value: -(i.Value)}
	}
	return &object.Nil{}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	// TODO: add support for inter-expression between boolean and integer e.g: 1 == benar
	if left.Type() == object.OBJECT_INTEGER && right.Type() == object.OBJECT_INTEGER {
		return evalInfixIntegerExpression(op, left, right)
	}
	if left.Type() == object.OBJECT_BOOLEAN && right.Type() == object.OBJECT_BOOLEAN {
		return evalInifxBooelanExpression(op, left, right)
	}
	return &object.Nil{}
}

func evalInifxBooelanExpression(op string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Boolean).Value
	r := right.(*object.Boolean).Value
	switch op {
	case "==":
		return &object.Boolean{Value: l == r}
	case "!=":
		return &object.Boolean{Value: l != r}
	default:
		return &object.Nil{}
	}
}

func evalInfixIntegerExpression(op string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case "<":
		return &object.Boolean{Value: l < r}
	case ">":
		return &object.Boolean{Value: l > r}
	case "==":
		return &object.Boolean{Value: l == r}
	case "!=":
		return &object.Boolean{Value: l != r}
	}
	return &object.Nil{}
}

func evalJikaStatement(jk *ast.JikaStatement) object.Object {
	cond := evalExpression(jk.Condition)
	if condIsTrue(cond) {
		return evalBlockStatement(jk.JikaBlock)
	} else if jk.LainnyaBlock != nil {
		return evalBlockStatement(jk.LainnyaBlock)
	} else {
		return &object.Nil{}
	}
}

// TODO: goodluck trying to understand all of this
// var KEM_ARR []object.Object = []object.Object{} // which is better, this or below
var ret bool = false // switch to check if already encounter a return stataement

func evalBlockStatement(bs *ast.BlockStatement) object.Object {
	var obj object.Object

	for _, s := range bs.Statements {
		// if len(KEM_ARR) != 0 {
		// 	v := KEM_ARR[0]
		// 	return KEM_ARR[0]
		// }

		if ks, ok := s.(*ast.KembalikanStatement); ok && !ret {
			// KEM_ARR = append(KEM_ARR, evalExpression(ks.Expression))
			ret = true
			obj = evalExpression(ks.Expression)
		}
		if ret {
			// if already encounter return statement, return the last obj which is return statement obj
			return obj
		}
		obj = evalStatement(s)
	}
	ret = false // switch back
	// only return the last statement from the block
	return obj
}

func evalCetakStatement(cs *ast.CetakStatement) object.Object {
	var obj object.Object
	for _, e := range cs.Expression {
		obj = evalExpression(e)
		// cetak statement is just calling Go fmt.Println
		fmt.Print(obj.Inspect() + " ")
	}
	fmt.Print("\n")
	// only return the last expression
	return obj
}

func printEval(eval object.Object) {
	fmt.Printf("%s\n", eval.Inspect())
}

func condIsTrue(cond object.Object) bool {
	switch c := cond.(type) {
	case *object.Boolean:
		return c.Value
	case *object.Nil:
		return false
	default:
		// truthy
		return true
	}
}
