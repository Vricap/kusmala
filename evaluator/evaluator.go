package evaluator

import (
	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/object"
)

func Eval(node ast.Node) object.Object {
	switch t := node.(type) {
	case *ast.Tree:
		return Eval(t.Statements[0])
	case *ast.ExpressionStatement:
		return Eval(t.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: t.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: t.Value}
	case *ast.PrefixExpression:
		right := Eval(t.Right)
		return evalPrefixExpression(t.Operator, right)
	case *ast.InfixExpression:
		left := Eval(t.Left)
		right := Eval(t.Right)
		return evalInfixExpression(t.Operator, left, right)
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
