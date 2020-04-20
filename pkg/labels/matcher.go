// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package labels

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/cespare/xxhash/v2"

	"facette.io/facette/pkg/parser"
)

// Matcher is a labels matcher.
type Matcher []MatchCond

// MatcherFromParser gets a labels matcher from a parser.
func MatcherFromParser(p *parser.Parser) (Matcher, error) {
	matcher := Matcher{}

	tok := p.Peek()
	isName := tok.Type == parser.IDENT

	if isName {
		matcher = append(matcher, MatchCond{
			Op:    OpEq,
			Name:  Name,
			Value: tok.String(),
		})

		p.Next()
	}

	// Stop if end of input and name was at least provided
	tok = p.Peek()
	if tok.Type != parser.LBRACE && isName {
		return matcher, nil
	}

	var err error

	_, err = p.Expect(parser.LBRACE)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != parser.RBRACE {
		for {
			cond, err := parseMatchCond(p)
			if err != nil {
				return nil, err
			}

			matcher = append(matcher, cond)

			// Allow extraneous comma
			if p.Peek().Type != parser.COMMA {
				break
			}

			p.Next()

			if p.Peek().Type == parser.RBRACE {
				break
			}
		}
	}

	_, err = p.Expect(parser.RBRACE)
	if err != nil {
		return nil, err
	}

	return matcher, nil
}

// ParseMatcher parses input string for a labels matcher.
func ParseMatcher(text string) (Matcher, error) {
	p := parser.New(strings.NewReader(text))

	matcher, err := MatcherFromParser(p)
	if err != nil {
		return nil, err
	}

	_, err = p.Expect(parser.EOF)
	if err != nil {
		return nil, err
	}

	return matcher, nil
}

// Hash returns a hashed representation of the labels matcher.
func (m Matcher) Hash() uint64 {
	hash := xxhash.New()

	for _, cond := range m {
		hash.Write([]byte(cond.Name + "\x1e" + cond.Value + "\x1e")) // nolint:errcheck,gosec
	}

	return hash.Sum64()
}

// String satisfies the fmt.Stringer interface.
// nolint:gosec
func (m Matcher) String() string {
	b := bytes.NewBuffer(nil)
	b.WriteByte('{')

	for idx, cond := range m {
		if idx > 0 {
			b.WriteByte(',')
		}

		b.WriteString(cond.String())
	}

	b.WriteByte('}')

	return b.String()
}

// MatchCond is a labels match condition.
type MatchCond struct {
	Op    Op
	Name  string
	Value string
	re    *regexp.Regexp
}

// NewMatchCond creates a new labels match condition instance.
func NewMatchCond(op Op, name, value string) (MatchCond, error) {
	var (
		re  *regexp.Regexp
		err error
	)

	if op == OpEqRegexp || op == OpNotEqRegexp {
		re, err = regexp.Compile(value)
		if err != nil {
			return MatchCond{}, err
		}
	}

	return MatchCond{Op: op, Name: name, Value: value, re: re}, nil
}

// Match returns whether or not the condition matches with the given name and
// value pair.
func (m MatchCond) Match(label Label) bool {
	if label.Name != m.Name {
		return false
	}

	switch m.Op {
	case OpEq:
		return label.Value == m.Value

	case OpNotEq:
		return label.Value != m.Value

	case OpEqRegexp:
		return m.re.MatchString(label.Value)

	case OpNotEqRegexp:
		return !m.re.MatchString(label.Value)
	}

	panic("invalid match operator")
}

// String satisfies the fmt.Stringer interface.
func (m MatchCond) String() string {
	var op string

	switch m.Op {
	case OpEq:
		op = "="

	case OpNotEq:
		op = "!="

	case OpEqRegexp:
		op = "=~"

	case OpNotEqRegexp:
		op = "!~"

	default:
		panic("invalid match operator")
	}

	return fmt.Sprintf("%s%s%q", m.Name, op, m.Value)
}

func parseMatchCond(p *parser.Parser) (MatchCond, error) {
	tok, err := p.Expect(parser.IDENT)
	if err != nil {
		return MatchCond{}, err
	}

	name := tok.String()

	tok = p.Next()
	if !tok.IsOperator() {
		return MatchCond{}, fmt.Errorf("expected operator but got %s at %s", tok.Type, tok.Pos)
	}

	var op Op

	switch tok.Type {
	case parser.EQ:
		op = OpEq

	case parser.NEQ:
		op = OpNotEq

	case parser.EQREGEXP:
		op = OpEqRegexp

	case parser.NEQREGEXP:
		op = OpNotEqRegexp
	}

	tok, err = p.Expect(parser.STRING)
	if err != nil {
		return MatchCond{}, err
	}

	value := tok.String()

	return NewMatchCond(op, name, value)
}

// Op is a labels match operator.
type Op int

// Match operators:
const (
	_ Op = iota
	OpEq
	OpNotEq
	OpEqRegexp
	OpNotEqRegexp
)
