package object

import "fmt"

type ObjectType string

const (
	OBJECT_INTEGER    ObjectType = "INTEGER"
	OBJECT_BOOLEAN               = "BOOLEAN"
	OBJECT_NIL                   = "NIL"
	OBJECT_KEMBALIKAN            = "OBJECT_KEMBALIKAN"
	OBJECT_ERR                   = "ERROR"
	OBJECT_STRING                = "STRING"
	OBEJCT_IDENTIFIER            = "IDENTIFIER"
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

// TODO: the environment doesn't have idea about global scope or block scope.
type Environment struct {
	store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
