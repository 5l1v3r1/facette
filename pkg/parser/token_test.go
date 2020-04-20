// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TokenType_String(t *testing.T) {
	for _, typ := range []Type{
		INVALID,
		EOF,
		IDENT,
		NUMBER,
		STRING,
		NEWLINE,
		BADESCAPE,
		EQ,
		NEQ,
		EQREGEXP,
		NEQREGEXP,
		LBRACE,
		RBRACE,
		LPAREN,
		RPAREN,
		COMMA,
	} {
		assert.Equal(t, tokenTypes[typ], typ.String())
	}
}
