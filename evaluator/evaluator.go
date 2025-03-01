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
		if err, ok := eval.(*object.Error); ok {
			fmt.Println("\t", err.Inspect())
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
	case *ast.ExpressionStatement:
		return evalExpression(s.Expression, env)
	case *ast.BlockStatement:
		var ret_obj object.Object
		return evalBlockStatement(s, env, ret_obj)
	case *ast.CetakStatement:
		return evalCetakStatement(s, env)
	case *ast.ReassignStatement:
		return evalReassignStatement(s, env, s.Ln)
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
	case *ast.FungsiExpression:
		return evalFungsiLiteral(e, env)
	case *ast.CallExpression:
		fn := evalExpression(e.Function, env) // return the function ident name
		if fn.Type() == object.OBJECT_ERR {
			return fn
		}
		return runFunction(fn, e, env)
	case *ast.StringLiteral:
		return &object.String{Value: e.Value, Ln: e.Ln}
	case *ast.PanjangFungsi:
		return evalPanjangFungsi(e, e.Ln, env)
	case *ast.ArrayLiteral:
		return evalArray(e, e.Ln, env)
	case *ast.IndexExpression:
		left := evalExpression(e.Left, env)
		if left.Type() == object.OBJECT_ERR {
			return left
		}
		index := evalExpression(e.Index, env)
		if index.Type() == object.OBJECT_ERR {
			return index
		}
		return evalIndexExpression(left, index, e.Ln, env)
	default:
		return newError("ekspresi tidak diketahui atau tidak ditempatnya", e.TokenLiteral(), e.Line())
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	if k, ok := right.(*object.Kembalikan); ok {
		right = k.Value
	}
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
	if k, ok := left.(*object.Kembalikan); ok {
		left = k.Value
	}
	if k, ok := right.(*object.Kembalikan); ok {
		right = k.Value
	}
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
	// newChildEnv := object.NewChildEnv(env) // TODO: this fuck recursive function
	if condIsTrue(cond) {
		return evalStatement(jk.JikaBlock, env)
	} else if jk.LainnyaBlock != nil {
		return evalStatement(jk.LainnyaBlock, env)
	} else {
		return &object.Nil{}
	}
}

func evalFungsiLiteral(fl *ast.FungsiExpression, env *object.Environment) object.Object {
	return &object.FungsiLiteral{Param: fl.Params, Body: fl.Body, Env: env, Ln: fl.Line()}
}

func evalArguments(a []ast.Expression, env *object.Environment) []object.Object {
	var obj []object.Object
	for _, e := range a {
		eval := evalExpression(e, env)
		if eval.Type() == object.OBJECT_ERR {
			return []object.Object{eval}
		}
		obj = append(obj, eval)
	}
	return obj
}

// TODO: i honestly doesn't know how all of this work
func runFunction(fn object.Object, e *ast.CallExpression, env *object.Environment) object.Object {
	args := evalArguments(e.Arguments, env)
	if len(args) == 1 && args[0].Type() == object.OBJECT_ERR { // if there's error
		return args[0]
	}
	f, ok := fn.(*object.FungsiLiteral)
	if !ok {
		return newError("bukan sebuah fungsi", fn.Inspect(), fn.Line())
	}
	if len(args) != len(f.Param) {
		s := fmt.Sprintf("fungsi membutuhkan %d parameter namun menemukan %d argumen", len(f.Param), len(args))
		return newError(s, e.TokenLiteral(), e.Line())
	}
	childEnv := extendFuncEnv(f, args)
	eval := evalStatement(f.Body, childEnv)
	return eval
}

func extendFuncEnv(f *object.FungsiLiteral, args []object.Object) *object.Environment {
	env := object.NewChildEnv(f.Env)
	for i, p := range f.Param {
		env.Set(p.Value, args[i]) // assign each params ident to arguments value
	}
	return env
}

// TODO: goodluck trying to understand all of this
// var ret_obj object.Object

/*
TODO: BUG if we have this:

	buat x = fungsi() {
		kembalikan 2 * 2;
	}

	buat f = fungsi() {
		buat a = x(); function will stop and return here because somehow kembalikan statement at x is also triggered in f
		cetak("helo");
	}
*/

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment, ret_obj object.Object) object.Object {
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
			fmt.Println("\t", err.Inspect())
			continue // so that all error from the blocks from parent to all its child is outputted. change to break to negate
		}
	}
	// only return the last statement from the block
	return obj
}

func evalKembalikanStatement(ks *ast.KembalikanStatement, env *object.Environment) object.Object {
	if ks.Expression != nil {
		// return &object.Kembalikan{Value: evalExpression(ks.Expression, env), Ln: ks.Line()}
		return evalExpression(ks.Expression, env)
	}
	// return &object.Kembalikan{Value: &object.Nil{}, Ln: ks.Line()}
	return &object.Nil{}
}

func evalCetakStatement(cs *ast.CetakStatement, env *object.Environment) object.Object {
	var obj object.Object
	for _, e := range cs.Expression {
		obj = evalExpression(e, env)
		if obj.Type() == object.OBJECT_ERR {
			return obj
		}

		// cetak statement is just calling Go fmt.Println
		fmt.Print(obj.Inspect() + " ")
	}
	fmt.Print("\n")
	// only return the last expression
	return obj
}

func evalReassignStatement(rs *ast.ReassignStatement, env *object.Environment, l int) object.Object {
	expr := evalExpression(rs.NewValue, env)
	if expr.Type() == object.OBJECT_ERR {
		return expr
	}
	_, ok := env.Get(rs.Ident.Value)
	if !ok {
		return newError("pengenal tidak diketahui", rs.Ident.TokenLiteral(), l)
	}
	traverseEnv(rs.Ident.Value, env, expr)
	env.Set(rs.Ident.Value, expr)
	return &object.Nil{}
}

// TODO: find out more about this
func traverseEnv(name string, env *object.Environment, expr object.Object) {
	if env == nil {
		return
	}
	_, ok := env.Get(name) // we traverse the env and change every value we found
	if ok {
		env.Set(name, expr)
		if env.Master != nil {
			if _, ok := env.Master.Get(name); !ok {
				return
			}

		}
	}
	traverseEnv(name, env.Master, expr)
}

func evalPanjangFungsi(e *ast.PanjangFungsi, l int, env *object.Environment) object.Object {
	arg := evalExpression(e.Argument, env)
	if arg.Type() != object.OBJECT_STRING && arg.Type() != object.OBJECT_ARRAY {
		return newError("argumen panjang hanya menerima string atau array", arg.Inspect(), l)
	}
	var val int
	switch a := arg.(type) {
	case *object.Array:
		val = len(a.El)
	case *object.String:
		val = len(a.Value)
	}
	return &object.Integer{Ln: l, Value: val}
}

func evalArray(a *ast.ArrayLiteral, l int, env *object.Environment) object.Object {
	arr := &object.Array{Ln: a.Ln}
	for _, v := range a.Elements {
		arr.El = append(arr.El, evalExpression(v, env))
	}
	return arr
}

func evalIndexExpression(left object.Object, index object.Object, l int, env *object.Environment) object.Object {
	le := evalLeftIndex(left, l)
	if le.Type() == object.OBJECT_ERR {
		return le
	}
	val := evalIndex(le, index, l)
	if val.Type() == object.OBJECT_ERR {
		return val
	}
	return val
}

func evalLeftIndex(left object.Object, l int) object.Object {
	switch t := left.(type) {
	case *object.Kembalikan:
		switch k := t.Value.(type) {
		case *object.Array:
			return k
		default:
			return newError("struktur data tidak didukung operator index", k.Inspect(), l)
		}
	case *object.Array:
		return t
	default:
		return newError("struktur data tidak didukung operator index", left.Inspect(), l)
	}
}

func evalIndex(le object.Object, index object.Object, l int) object.Object {
	i, ok := index.(*object.Integer)
	if !ok {
		return newError("argumen index harus sebuah integer", fmt.Sprintf("[%s]", index.Inspect()), l)
	}
	arr := le.(*object.Array)
	if i.Value < 0 {
		return newError("argumen index tidak boleh negatif", fmt.Sprintf("[%s]", i.Inspect()), l)
	} else if i.Value > len(arr.El)-1 {
		return newError("argumen index melebihi panjang array", fmt.Sprintf("[%s]", i.Inspect()), l)
	}
	return arr.El[i.Value]
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
