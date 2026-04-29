package ast

import (
	"testing"

	"deadlock/language/token"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	program := &Program{
		Statements: []Statement{
			&VariableDeclaration{
				Token: token.Token{Type: token.ASSIGN, Literal: "="},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}

	assert.Equal("foo = bar", program.String())
}
