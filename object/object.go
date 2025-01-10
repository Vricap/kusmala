package object

import "fmt"

type ObjectType string

const (
	OBJECT_INTEGER ObjectType = "INTEGER"
	OBJECT_BOOLEAN            = "BOOLEAN"
	OBJECT_NULL               = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return OBJECT_INTEGER
}

type Boolean struct {
	Value bool
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

type Null struct{}

func (n *Null) Inspect() string {
	return "NULL"
}
func (n *Null) Type() ObjectType {
	return OBJECT_NULL
}
