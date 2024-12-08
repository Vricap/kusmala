package ast

import (
	"github.com/vricap/kusmala/token"
)

// node in the tree
type Node interface {
	TokenLiteral() string
}

// each node type could be either a statement
type Statement interface {
	Node
	statementNode() // marker to identify a Statement node
}

// or a expression
type Expression interface {
	Node
	expressionNode() // marker to identify a Expression node
}

// Tree is the root of every tree the parser will construct
type Tree struct {
	Statements []Statement // consisting of struct that implement Statement interface
}

func (t *Tree) TokenLiteral() string {
	if len(t.Statements) > 0 {
		return t.Statements[0].TokenLiteral() // return the first statement node -- the root
	} else {
		return ""
	}
}

// example of buat statement: buat x = 1 + 1;
type BuatStatement struct {
	Token      token.Token // token.BUAT
	Name       Identifier  // the ident name (x)
	Expression string      // the value (1 + 1)
}

func (bs *BuatStatement) TokenLiteral() string {
	// return fmt.Sprintf("Token: %v | Name: %v | Expression: %s", bs.Token, bs.Name, bs.Expression)
	return bs.Token.Literal
}
func (bs *BuatStatement) statementNode() {}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // the name of the variable, function name etc...
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) identifierNode() {}

type KembalikanStatement struct {
	Token token.Token
	Value string // the value expression that will be returned
}

func (ks *KembalikanStatement) TokenLitral() string {
	return ks.Token.Literal
}

func (ks *KembalikanStatement) statementNode() {}
