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
	case *ast.KembalikanStatement:
	case *ast.JikaStatement:
		return evalJikaStatement(s)
	case *ast.BlockStatement:
		// TODO: only work with 1 statement
		return evalStatement(s.Statements[0])
	case *ast.ExpressionStatement:
		return evalExpression(s.Expression)
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
		return evalStatement(jk.JikaBlock)
	} else if jk.LainnyaBlock != nil {
		return evalStatement(jk.LainnyaBlock)
	} else {
		return &object.Nil{}
	}
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
