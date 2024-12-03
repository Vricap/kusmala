package ast

import (
	"github.com/vricap/kusmala/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node            // the token literal string
	statementNode() // the token type itself
}

type Expression interface {
	Node             // the token literal string
	expressionNode() // the token type itself
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // ident token e.g variable name
	Value string
}

// Identifier need to satisfy Expression interface. since some identifier DO produce value.
// example:
//
//	ident := 1 + 1;
//	var   := ident; // ident here IS a value
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}

type BuatStatement struct {
	Token token.Token // buat token
	Name  *Identifier // iden struct containing ident token and the name
	Value Expression
}

// to satisfy Statement interface
func (bs *BuatStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BuatStatement) statementNode()       {}
