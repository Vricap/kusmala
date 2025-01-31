package evaluator

import (
	"fmt"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/object"
)

func Eval(tree *ast.Tree, env *object.Environment) []object.Object {
	var evals []object.Object
	for _, s := range tree.Statements {
		eval := evalStatement(s, env)
		if ks, ok := eval.(*object.Kembalikan); ok {
			eval = ks.Value
			ret_obj = nil
		}
		if err, ok := eval.(*object.Error); ok {
			fmt.Println("\t" + err.Inspect())
			break // TODO: remove this when running test case
		}
		evals = append(evals, eval)
	}
	return evals
}

func evalStatement(stmt ast.Statement, env *object.Environment) object.Object {
	switch s := stmt.(type) {
	case *ast.BuatStatement:
		val := evalExpression(s.Expression, env)
		env.Set(s.Name.Value, val)
		return val
	case *ast.JikaStatement:
		return evalJikaStatement(s, env)
	case *ast.BlockStatement: // TODO: this problably shouldn't be here
		return evalBlockStatement(s, env)
	case *ast.ExpressionStatement:
		return evalExpression(s.Expression, env)
	case *ast.CetakStatement:
		return evalCetakStatement(s, env)
	// TODO: kembalikan statement isn't allowed in global scope. only inside a block statement
	// case *ast.KembalikanStatement:
	// 	return evalKembalikanStatement(s)
	default:
		return newError("statement tidak diketahui atau tidak ditempatnya", s.TokenLiteral(), s.Line())
	}
}

func evalExpression(expr ast.Expression, env *object.Environment) object.Object {
	switch e := expr.(type) {
	case *ast.Identifier:
		return evalIdentifier(e, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: e.Value, Ln: e.Ln}
	case *ast.PrefixExpression:
		right := evalExpression(e.Right, env)
		return evalPrefixExpression(e.Operator, right)
	case *ast.InfixExpression:
		left := evalExpression(e.Left, env)
		right := evalExpression(e.Right, env)
		return evalInfixExpression(e.Operator, left, right)
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: e.Value, Ln: e.Ln}
	// case *ast.FungsiExpression:
	// case *ast.CallExpression:
	case *ast.StringLiteral:
		return &object.String{Value: e.Value, Ln: e.Ln}
	default:
		return newError("ekspresi tidak diketahui atau tidak ditempatnya", e.TokenLiteral(), e.Line())
	}
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
			return newError("operator tidak didukung", fmt.Sprintf("%s%s", op, right.Inspect()), right.Line())
		}
		i := right.(*object.Integer)
		return &object.Integer{Value: -(i.Value)}
	}
	return newError("operator tidak didukung", fmt.Sprintf("%s%s", op, right.Inspect()), right.Line())
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	// TODO: add support for inter-expression between boolean and integer e.g: 1 == benar
	if left.Type() == object.OBJECT_INTEGER && right.Type() == object.OBJECT_INTEGER {
		return evalInfixIntegerExpression(op, left, right)
	}
	if left.Type() == object.OBJECT_BOOLEAN && right.Type() == object.OBJECT_BOOLEAN {
		return evalInifxBooelanExpression(op, left, right)
	}

	if left.Type() == object.OBJECT_STRING && right.Type() == object.OBJECT_STRING {
		return evalInifxStringExpression(op, left, right)
	}
	return newError("kesalahan tipe", fmt.Sprintf("%v %v %v", left.Inspect(), op, right.Inspect()), left.Line())
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
		return newError("operator tidak didukung", fmt.Sprintf("%v %v %v", left.Inspect(), op, right.Inspect()), left.Line())
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
	default:
		return newError("operator tidak didukung", fmt.Sprintf("%v %v %v", left.Inspect(), op, right.Inspect()), left.Line())
	}
}

func evalInifxStringExpression(op string, left object.Object, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: l + r} // string concatenation
	default:
		return newError("operator tidak didukung", fmt.Sprintf("%v %v %v", left.Inspect(), op, right.Inspect()), left.Line())
	}
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(i.Value)
	if !ok {
		return newError("pengenal tidak diketahui", i.Value, i.Ln)
	}
	return val
}

func evalJikaStatement(jk *ast.JikaStatement, env *object.Environment) object.Object {
	cond := evalExpression(jk.Condition, env)
	if condIsTrue(cond) {
		return evalBlockStatement(jk.JikaBlock, env)
	} else if jk.LainnyaBlock != nil {
		return evalBlockStatement(jk.LainnyaBlock, env)
	} else {
		return &object.Nil{}
	}
}

// TODO: goodluck trying to understand all of this
var ret_obj object.Object

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object

	for _, s := range bs.Statements {
		if ks, ok := s.(*ast.KembalikanStatement); ok {
			if ret_obj == nil {
				ret_obj = evalKembalikanStatement(ks, env)
				return ret_obj
			}
		}
		if ret_obj != nil {
			return ret_obj
		}
		obj = evalStatement(s, env)
		if err, ok := obj.(*object.Error); ok {
			fmt.Println("\t" + err.Inspect())
			continue // so that all error from the blocks from parent to all its child is outputted. change to break to negate
		}
	}
	// only return the last statement from the block
	return obj
}

func evalKembalikanStatement(ks *ast.KembalikanStatement, env *object.Environment) object.Object {
	return &object.Kembalikan{Value: evalExpression(ks.Expression, env), Ln: ks.Line()}
}

func evalCetakStatement(cs *ast.CetakStatement, env *object.Environment) object.Object {
	var obj object.Object
	for _, e := range cs.Expression {
		obj = evalExpression(e, env)
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

func newError(msg string, a any, l int) *object.Error {
	return &object.Error{Msg: fmt.Sprintf("%d: %s dekat '%v'", l, msg, a)}
}
