package ast

import (
	"bytes"

	"github.com/vricap/kusmala/token"
)

// node in the tree
type Node interface {
	TokenLiteral() string
	String() string
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
func (t *Tree) String() string {
	var out bytes.Buffer

	for _, s := range t.Statements {
		out.WriteString(s.String()) // append statement string to the buffer
	}
	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // the name of the variable, function name etc...
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

// example of buat statement: buat x = 1 + 1;
type BuatStatement struct {
	Token      token.Token // token.BUAT
	Name       *Identifier // the ident name (x)
	Expression Expression  // the value (1 + 1)
}

func (bs *BuatStatement) TokenLiteral() string {
	// return fmt.Sprintf("Token: %v | Name: %v | Expression: %s", bs.Token, bs.Name, bs.Expression)
	return bs.Token.Literal
}
func (bs *BuatStatement) statementNode() {}
func (bs *BuatStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral() + " ")
	out.WriteString(bs.Name.String())
	out.WriteString(" = ")

	if bs.Expression != nil {
		out.WriteString(bs.Expression.String())
	}
	out.WriteString(";")
	return out.String()
}

type KembalikanStatement struct {
	Token      token.Token
	Expression Expression // the value expression that will be returned
}

func (ks *KembalikanStatement) TokenLiteral() string {
	return ks.Token.Literal
}
func (ks *KembalikanStatement) statementNode() {}
func (ks *KembalikanStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ks.TokenLiteral() + " ")
	if ks.Expression != nil {
		out.WriteString(ks.Expression.String())
	}
	return out.String()
}

// ExpressionStatement is statement that consist solely of one expression. it's a wrapper so that we could insert this in Tree Statements slice
type ExpressionStatement struct {
	Token      token.Token // the first token in the ExpressionStatement
	Expression Expression
}

func (ex *ExpressionStatement) TokenLiteral() string {
	return ex.Token.Literal
}
func (ex *ExpressionStatement) statementNode() {}
func (ex *ExpressionStatement) String() string {
	if ex.Expression != nil {
		return ex.Expression.String()
	}
	return ""
}
