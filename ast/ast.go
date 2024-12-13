package ast

import (
	"bytes"

	"github.com/vricap/kusmala/token"
)

// node in the tree
type Node interface {
	TokenLiteral() string
}

// each node type could be either a statement...
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

/*******************************************
*			STATEMENT STRUCT			   *
*******************************************/

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

// ExpressionStatement is statement that consist solely of one expression. it's a wrapper so that we could insert this in Tree Statements slice
type ExpressionStatement struct {
	Token      token.Token // the first token in the ExpressionStatement
	Expression Expression  // the struct that implement Expression interfae. e.g. Identifier, IntegerLiteral, etc...
}

func (ex *ExpressionStatement) TokenLiteral() string {
	return ex.Token.Literal
}
func (ex *ExpressionStatement) statementNode() {}

/*******************************************
*			EXPRESSION STRUCT			   *
*******************************************/

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

type IntegerLiteral struct {
	Token token.Token
	Value int // platform dependent. 32 size in 32 bits machine, 64 size in 64 bits machine
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // the prefix token. e.g - or !
	Operator string
	Right    Expression // the expression struct that implement Expression interface
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) expressionNode() {}

// TODO: see the book on this part
func (pe *PrefixExpression) String() string {
	return pe.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) expressionNode() {}

// TODO: see the book on this part
func (ie *InfixExpression) String() string {
	return ie.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}
func (bl *BooleanLiteral) ExpressionNode() {}
func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}
