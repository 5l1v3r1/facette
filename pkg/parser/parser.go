// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package parser provides expressions lexical parsing.
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

const eof = rune(0)

// Parser is a lexical parser.
type Parser struct {
	scanner   io.RuneScanner
	pos, last Position
	peeked    *Token
}

// New creates a new lexical parser instance.
func New(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewReader(r),
		pos:     Position{Line: 1, Char: 1},
	}
}

// Expect returns the next token ensuring for its type.
func (p *Parser) Expect(typ Type) (*Token, error) {
	tok := p.Next()
	if tok.Type != typ {
		return tok, fmt.Errorf("expected %s but got %s at %s", typ, tok.Type, tok.Pos)
	}

	return tok, nil
}

// Next returns the next token.
func (p *Parser) Next() *Token {
	var tok *Token

	if p.peeked != nil {
		tok = p.peeked
		p.peeked = nil
	} else {
		tok = p.scan()
	}

	return tok
}

// Peek returns the next token without consuming it.
func (p *Parser) Peek() *Token {
	if p.peeked == nil {
		tok := p.scan()
		p.peeked = tok
	}

	return p.peeked
}

// PeekRune returns the next rune without consuming it.
// FIXME: find a way to ditch this as it only serves for function ident
//        detection in series expression.
func (p *Parser) PeekRune() rune {
	r, _ := p.read()
	p.unread()

	return r
}

func (p *Parser) read() (rune, Position) {
	ch, _, err := p.scanner.ReadRune()
	if err != nil {
		return eof, p.pos
	}

	p.last = p.pos

	if ch == '\n' {
		p.pos.Line++
		p.pos.Char = 1
	} else {
		p.pos.Char++
	}

	return ch, p.last
}

func (p *Parser) run(set []byte) []byte {
	b := bytes.NewBuffer(nil)

	for {
		ch, _ := p.read()
		if !bytes.ContainsRune(set, ch) {
			break
		}

		b.WriteRune(ch)
	}

	p.unread()

	return b.Bytes()
}

func (p *Parser) scan() *Token {
	c, pos := p.read()
	for isSpace(c) {
		c, pos = p.read()
	}

	switch {
	case c == eof:
		return newToken(EOF, pos, "")

	case c == '-' || isDigit(c):
		p.unread()
		return p.scanNumber()

	case isIdentChar(c):
		p.unread()
		return p.scanIdent()

	case c == '"' || c == '\'':
		p.unread()
		return p.scanString()

	case c == '=':
		next, _ := p.read()
		if next == '~' {
			return newToken(EQREGEXP, pos, "")
		}

		p.unread()

		return newToken(EQ, pos, "")

	case c == '!':
		next, _ := p.read()
		if next == '=' {
			return newToken(NEQ, pos, "")
		} else if next == '~' {
			return newToken(NEQREGEXP, pos, "")
		}

		p.unread()

	case c == '{':
		return newToken(LBRACE, pos, "")

	case c == '}':
		return newToken(RBRACE, pos, "")

	case c == '(':
		return newToken(LPAREN, pos, "")

	case c == ')':
		return newToken(RPAREN, pos, "")

	case c == ',':
		return newToken(COMMA, pos, "")
	}

	return newToken(INVALID, pos, string(c))
}

func (p *Parser) scanIdent() *Token {
	start := p.pos

	buf := bytes.NewBuffer(nil)

loop:
	for {
		ch, _ := p.read()
		switch {
		case ch == eof:
			break loop

		case isIdentChar(ch):
			buf.WriteRune(ch)

		default:
			p.unread()
			break loop
		}
	}

	return newToken(IDENT, start, buf.String())
}

func (p *Parser) scanNumber() *Token {
	digits := []byte("0123456789")
	start := p.pos

	buf := bytes.NewBuffer(nil)

	next, _ := p.read()
	if next == '-' {
		buf.WriteRune(next)
	} else {
		p.unread()
	}

	buf.Write(p.run(digits))

	next, _ = p.read()
	if next == '.' {
		buf.WriteRune(next)
		buf.Write(p.run(digits))
	} else {
		p.unread()
	}

	next, _ = p.read()
	if bytes.ContainsRune([]byte("eE"), next) {
		buf.WriteRune(next)

		next, _ = p.read()
		if bytes.ContainsRune([]byte("+-"), next) {
			buf.WriteRune(next)
		} else {
			p.unread()
		}

		buf.Write(p.run(digits))
	} else {
		p.unread()
	}

	return newToken(NUMBER, start, buf.String())
}

func (p *Parser) scanString() *Token {
	quote, start := p.read()

	buf := bytes.NewBuffer(nil)

	for {
		ch, pos := p.read()

		switch ch {
		case quote:
			return newToken(STRING, start, buf.String())

		case eof:
			return newToken(EOF, pos, buf.String())

		case '\n':
			return newToken(NEWLINE, pos, buf.String())

		case '\\':
			next, _ := p.read()
			switch next {
			case '\\', '"', '\'':
				buf.WriteRune(next)

			default:
				return newToken(BADESCAPE, pos, string(ch+next))
			}

		default:
			buf.WriteRune(ch)
		}
	}
}

func (p *Parser) unread() {
	err := p.scanner.UnreadRune()
	if err == nil {
		p.pos = p.last
	}
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isIdentChar(c rune) bool {
	return isDigit(c) || isAlpha(c) || c == '_'
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}
