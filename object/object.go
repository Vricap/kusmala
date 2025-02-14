package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vricap/kusmala/ast"
)

type ObjectType string

const (
	OBJECT_INTEGER    ObjectType = "INTEGER"
	OBJECT_BOOLEAN               = "BOOLEAN"
	OBJECT_NIL                   = "NIL"
	OBJECT_KEMBALIKAN            = "OBJECT_KEMBALIKAN"
	OBJECT_ERR                   = "ERROR"
	OBJECT_STRING                = "STRING"
	OBEJCT_IDENTIFIER            = "IDENTIFIER"
	OBJECT_FUNGSI                = "FUNGSI"
	OBJECT_JIKA                  = "JIKA"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Line() int
}

type Integer struct {
	Value int
	Ln    int
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return OBJECT_INTEGER
}
func (i *Integer) Line() int {
	return i.Ln
}

type Boolean struct {
	Value bool
	Ln    int
}

func (b *Boolean) Inspect() string {
	var s string
	if b.Value {
		s = "benar"
	} else {
		s = "salah"
	}
	return s
}

func (b *Boolean) Type() ObjectType {
	return OBJECT_BOOLEAN
}
func (i *Boolean) Line() int {
	return i.Ln
}

type Nil struct{}

func (n *Nil) Inspect() string {
	return "NIL"
}
func (n *Nil) Type() ObjectType {
	return OBJECT_NIL
}
func (i *Nil) Line() int {
	return 0
}

type Kembalikan struct {
	Value Object
	Ln    int
}

func (k *Kembalikan) Inspect() string {
	return k.Value.Inspect()
}
func (k *Kembalikan) Type() ObjectType {
	return OBJECT_KEMBALIKAN
}
func (i *Kembalikan) Line() int {
	return i.Ln
}

type Error struct {
	Msg string
}

func (e *Error) Inspect() string {
	return "ERROR di baris " + e.Msg
}
func (e *Error) Type() ObjectType {
	return OBJECT_ERR
}
func (i *Error) Line() int {
	return 0
}

type String struct {
	Value string
	Ln    int
}

func (s *String) Inspect() string {
	return s.Value
}
func (s *String) Type() ObjectType {
	return OBJECT_STRING
}
func (s *String) Line() int {
	return s.Ln
}

type Ident struct {
	Value Object
	Ln    int
}

func (i *Ident) Type() ObjectType {
	return OBEJCT_IDENTIFIER
}
func (i *Ident) Inspect() string {
	return i.Value.Inspect()
}
func (i *Ident) Line() int {
	return i.Ln
}

func NewEnv() *Environment {
	s := map[string]Object{}
	return &Environment{store: s}
}

type Environment struct {
	store  map[string]Object
	master *Environment // the master Environment of this Environment if any
}

func NewChildEnv(master *Environment) *Environment {
	child := NewEnv()
	child.master = master
	return child
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if e.master != nil && !ok { // if the ident if not exist in this env AND this env have master, we will return the master one.
		obj, ok := e.master.store[name]
		return obj, ok
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type JikaStatement struct {
	Env *Environment
	Ln  int
}

func (js *JikaStatement) Line() int {
	return js.Ln
}
func (js *JikaStatement) Type() ObjectType {
	return OBJECT_JIKA
}
func (js *JikaStatement) Inspect() string {
	return "jika"
}

type FungsiLiteral struct {
	Param []*ast.Identifier
	Body  *ast.BlockStatement
	Env   *Environment
	Ln    int
}

func (fl *FungsiLiteral) Line() int {
	return fl.Ln
}
func (fl *FungsiLiteral) Type() ObjectType {
	return OBJECT_FUNGSI
}

func (fl *FungsiLiteral) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Param {
		params = append(params, p.TokenLiteral())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fl.Body.TokenLiteral())
	out.WriteString("\n}")
	return out.String()
}
