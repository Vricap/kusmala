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
	}
	return nil
}
