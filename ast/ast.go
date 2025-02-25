package ast

import (
	"bytes"
	"strings"

	"github.com/vricap/kusmala/token"
)

// node in the tree
type Node interface {
	TokenLiteral() string
	Line() int
}

// each node type could be either a statement...
type Statement interface {
	Node
	Line() int
	statementNode() // marker to identify a Statement node
}

// or a expression
type Expression interface {
	Node
	Line() int
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
func (t *Tree) Line() int { return 0 }

/*******************************************
*			STATEMENT STRUCT			   *
*******************************************/

// example of buat statement: buat x = 1 + 1;
type BuatStatement struct {
	Token      token.Token // token.BUAT
	Name       *Identifier // the ident name (x)
	Expression Expression  // the value (1 + 1)
	Ln         int
}

func (bs *BuatStatement) TokenLiteral() string {
	// return fmt.Sprintf("Token: %v | Name: %v | Expression: %s", bs.Token, bs.Name, bs.Expression)
	return bs.Token.Literal
}
func (bs *BuatStatement) statementNode() {}
func (bs *BuatStatement) Line() int      { return bs.Ln }

type KembalikanStatement struct {
	Token      token.Token
	Expression Expression // the value expression that will be returned
	Ln         int
}

func (ks *KembalikanStatement) TokenLiteral() string {
	return ks.Token.Literal
}
func (ks *KembalikanStatement) statementNode() {}
func (bs *KembalikanStatement) Line() int      { return bs.Ln }

// ExpressionStatement is statement that consist solely of one expression. it's a wrapper so that we could insert this in Tree Statements slice
type ExpressionStatement struct {
	Token      token.Token // the first token in the ExpressionStatement
	Expression Expression  // the struct that implement Expression interfae. e.g. Identifier, IntegerLiteral, etc...
	Ln         int
}

func (ex *ExpressionStatement) TokenLiteral() string {
	return ex.Token.Literal
}
func (ex *ExpressionStatement) statementNode() {}
func (bs *ExpressionStatement) Line() int      { return bs.Ln }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
	Ln         int
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) Line() int      { return bs.Ln }

type JikaStatement struct {
	Token        token.Token
	Condition    Expression
	JikaBlock    *BlockStatement
	LainnyaBlock *BlockStatement
	Ln           int
}

func (ie *JikaStatement) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *JikaStatement) statementNode() {}
func (bs *JikaStatement) Line() int      { return bs.Ln }

type CetakStatement struct {
	Token      token.Token
	Expression []Expression
	Ln         int
}

func (cs *CetakStatement) TokenLiteral() string {
	return cs.Token.Literal
}
func (cs *CetakStatement) statementNode() {}
func (bs *CetakStatement) Line() int      { return bs.Ln }

/*******************************************
*			EXPRESSION STRUCT			   *
*******************************************/

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // the name of the variable, function name etc...
	Ln    int
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}
func (bs *Identifier) Line() int      { return bs.Ln }

type IntegerLiteral struct {
	Token token.Token
	Value int // platform dependent. 32 size in 32 bits machine, 64 size in 64 bits machine
	Ln    int
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) expressionNode() {}
func (bs *IntegerLiteral) Line() int       { return bs.Ln }

type PrefixExpression struct {
	Token    token.Token // the prefix token. e.g - or !
	Operator string
	Right    Expression // the expression struct that implement Expression interface
	Ln       int
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) expressionNode() {}
func (bs *PrefixExpression) Line() int       { return bs.Ln }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
	Ln       int
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) expressionNode() {}
func (bs *InfixExpression) Line() int       { return bs.Ln }

type BooleanLiteral struct {
	Token token.Token
	Value bool
	Ln    int
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}
func (bl *BooleanLiteral) expressionNode() {}
func (bs *BooleanLiteral) Line() int       { return bs.Ln }

type FungsiExpression struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStatement
	Ln     int
}

func (fe *FungsiExpression) TokenLiteral() string {
	return fe.Token.Literal
}
func (fe *FungsiExpression) expressionNode() {}
func (bs *FungsiExpression) Line() int       { return bs.Ln }

type CallExpression struct {
	Token     token.Token // the '('
	Function  Expression  // the ident to the function or FungsiExpression (literal)
	Arguments []Expression
	Ln        int
}

func (ce *CallExpression) TokenLiteral() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range ce.Arguments {
		params = append(params, p.TokenLiteral())
	}
	out.WriteString(ce.Function.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}
func (ce *CallExpression) expressionNode() {}
func (bs *CallExpression) Line() int       { return bs.Ln }

type StringLiteral struct {
	Token token.Token
	Value string
	Ln    int
}

func (s *StringLiteral) TokenLiteral() string {
	return s.Value
}
func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) Line() int {
	return s.Ln
}

type PanjangFungsi struct {
	Token    token.Token
	Argument Expression
	Ln       int
}

func (pf *PanjangFungsi) TokenLiteral() string {
	return pf.Token.Literal
}
func (pf *PanjangFungsi) expressionNode() {}
func (pf *PanjangFungsi) Line() int {
	return pf.Ln
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
	// Index    Expression
	Ln int
}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}
func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) Line() int {
	return a.Ln
}

type IndexExpression struct {
	Token token.Token
	Left  Expression // could be array literal, or identifier of array literal. maybe later i will add string too
	Index Expression
	Ln    int
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) Line() int {
	return ie.Ln
}
func (ie *IndexExpression) expressionNode() {}
