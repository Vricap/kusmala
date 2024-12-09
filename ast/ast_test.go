package ast

import (
	"testing"

	"github.com/vricap/kusmala/token"
)

func TestString(t *testing.T) {
	tree := &Tree{
		Statements: []Statement{
			&BuatStatement{
				Token: token.Token{Type: token.BUAT, Literal: "buat"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Expression: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "anotherVar",
					},
				},
			},
		},
	}

	if tree.String() != "buat myVar = anotherVar;" {
		t.Errorf("tree.String() is wrong. got: %s", tree.String())
	}
}
