// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package parser

import "fmt"

// Tokens:
const (
	INVALID Type = iota
	EOF

	IDENT     // ident
	NUMBER    // 123.45
	STRING    // "abc" or 'abc'
	NEWLINE   // \n
	BADESCAPE // \b

	beginOperators
	EQ        // =
	NEQ       // !=
	EQREGEXP  // =~
	NEQREGEXP // !~
	endOperators

	LBRACE // {
	RBRACE // }
	LPAREN // (
	RPAREN // )
	COMMA  // ,
)

var tokenTypes = map[Type]string{
	INVALID: "invalid",
	EOF:     "end of input",

	IDENT:     "identifier",
	NUMBER:    "number",
	STRING:    "string",
	NEWLINE:   "new line",
	BADESCAPE: "bad escape",

	EQ:        "equal",
	NEQ:       "not equal",
	EQREGEXP:  "equal pattern",
	NEQREGEXP: "not equal pattern",

	LBRACE: "left brace",
	RBRACE: "right brace",
	LPAREN: "left parenthesis",
	RPAREN: "right parenthesis",
	COMMA:  "comma",
}

// Type is a lexical token type.
type Type int

// String satisfies the fmt.Stringer interface.
func (t Type) String() string {
	s, ok := tokenTypes[t]
	if !ok {
		return "token"
	}

	return s
}

// Token is a lexical token.
type Token struct {
	Type Type
	Pos  Position
	Text string
}

func newToken(typ Type, pos Position, text string) *Token {
	return &Token{
		Type: typ,
		Pos:  pos,
		Text: text,
	}
}

// String satisfies the fmt.Stringer interface.
func (t Token) String() string {
	return t.Text
}

// IsOperator returns whether or not the token is an operator.
func (t Token) IsOperator() bool {
	return t.Type >= beginOperators && t.Type <= endOperators
}

// Position is a lexical token position.
type Position struct {
	Line int
	Char int
}

// String satisfies the fmt.Stringer interface.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Char)
}
