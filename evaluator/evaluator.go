package evaluator

import (
	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/object"
)

func Eval(node ast.Node) object.Object {
	return &object.Boolean{}
}
