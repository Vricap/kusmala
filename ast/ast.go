package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode() // marker to identify a Statement node
}

type BuatStatement struct {
}

func (bs *BuatStatement) TokenLiteral() string {
	return ""
}

func (bs *BuatStatement) statementNode() {}
