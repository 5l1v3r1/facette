// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"fmt"
	"strconv"
	"strings"

	"facette.io/facette/pkg/labels"
	"facette.io/facette/pkg/parser"
)

// Expr is a series expression.
type Expr interface {
	expr()
}

func (AggregateExpr) expr() {}
func (AliasExpr) expr()     {}
func (MatcherExpr) expr()   {}
func (SampleExpr) expr()    {}
func (ScaleExpr) expr()     {}

// ParseExpr parses input string for a series expression.
func ParseExpr(text string) (Expr, error) {
	p := parser.New(strings.NewReader(text))

	expr, err := parseExpr(p)
	if err != nil {
		return nil, err
	}

	_, err = p.Expect(parser.EOF)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func parseExpr(p *parser.Parser) (Expr, error) {
	var (
		expr Expr
		err  error
	)

	tok := p.Peek()

	switch tok.Type {
	case parser.LBRACE:
		v, err := labels.MatcherFromParser(p)
		if err != nil {
			return nil, err
		}

		expr = &MatcherExpr{Matcher: v}

	case parser.IDENT:
		var fn Func

		err = fn.UnmarshalText([]byte(tok.Text))
		if err != nil || p.PeekRune() != '(' {
			v, err := labels.MatcherFromParser(p)
			if err != nil {
				return nil, err
			}

			expr = &MatcherExpr{Matcher: v}

			break
		}

		// Skip both function ident and left parenthesis
		p.Next()
		p.Next()

		switch fn {
		case FuncAlias:
			expr, err = parseAliasExpr(p)

		case FuncAvg, FuncSum:
			expr, err = parseAggregateExpr(p, fn)

		case FuncSample:
			expr, err = parseSampleExpr(p)

		case FuncScale:
			expr, err = parseScaleExpr(p)
		}

		if err != nil {
			return nil, err
		}

		// Allow extraneous comma
		if p.Peek().Type == parser.COMMA {
			p.Next()
		}

		_, err = p.Expect(parser.RPAREN)
		if err != nil {
			return nil, err
		}

	default:
		// Expect LBRACE per default to ensure proper parsing of nested
		// expression
		_, err = p.Expect(parser.LBRACE)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

// AggregateExpr is an aggregate expression.
type AggregateExpr struct {
	Exprs []Expr
	Op    AggregateOp
}

func parseAggregateExpr(p *parser.Parser, fn Func) (*AggregateExpr, error) {
	aggr := &AggregateExpr{}

	switch fn {
	case FuncAvg:
		aggr.Op = AggregateOpAverage

	case FuncSum:
		aggr.Op = AggregateOpSum
	}

	for {
		v, err := parseExpr(p)
		if err != nil {
			return nil, err
		}

		aggr.Exprs = append(aggr.Exprs, v)

		if p.Peek().Type != parser.COMMA {
			break
		}

		p.Next()

		// Allow extraneous comma at the end
		if p.Peek().Type == parser.RPAREN {
			break
		}
	}

	return aggr, nil
}

// AggregateOp is a series aggregate expression operator.
type AggregateOp int

// Aggregation expression operators:
const (
	_ AggregateOp = iota
	AggregateOpAverage
	AggregateOpSum
)

// AliasExpr is an alias expression.
type AliasExpr struct {
	Expr  Expr
	Alias string
}

func parseAliasExpr(p *parser.Parser) (*AliasExpr, error) {
	v, err := parseExpr(p)
	if err != nil {
		return nil, err
	}

	_, err = p.Expect(parser.COMMA)
	if err != nil {
		return nil, err
	}

	tok, err := p.Expect(parser.STRING)
	if err != nil {
		return nil, err
	}

	return &AliasExpr{Expr: v, Alias: tok.String()}, nil
}

// MatcherExpr is an alias expression.
type MatcherExpr struct {
	Matcher labels.Matcher
}

// SampleExpr is a sample expression.
type SampleExpr struct {
	Expr Expr
	Mode SampleMode
}

func parseSampleExpr(p *parser.Parser) (*SampleExpr, error) {
	v, err := parseExpr(p)
	if err != nil {
		return nil, err
	}

	var mode SampleMode

	tok := p.Peek()

	if tok.Type != parser.COMMA && tok.Type != parser.RPAREN {
		tok, err = p.Expect(parser.IDENT)
		if err != nil {
			return nil, err
		} else if tok.Text != "with" {
			return nil, fmt.Errorf(`expected "with" keyword but got %q at %s`, tok.Text, tok.Pos)
		}

		tok, err = p.Expect(parser.IDENT)
		if err != nil {
			return nil, err
		}

		err = mode.UnmarshalText([]byte(tok.String()))
		if err != nil {
			return nil, fmt.Errorf("expected sample mode but got %q at %s", tok.Text, tok.Pos)
		}
	} else {
		mode = SampleModeAverage
	}

	return &SampleExpr{Expr: v, Mode: mode}, nil
}

// SampleMode is a series sample expression mode.
type SampleMode int

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (s *SampleMode) UnmarshalText(b []byte) error {
	v, ok := sampleModes[string(b)]
	if !ok {
		return fmt.Errorf("unsupported sample mode: %s", b)
	}

	*s = v

	return nil
}

// Sample expression modes:
const (
	_ SampleMode = iota
	SampleModeAverage
	SampleModeFirst
	SampleModeLast
	SampleModeMax
	SampleModeMin
	SampleModeSum
)

var sampleModes = map[string]SampleMode{
	"avg":   SampleModeAverage,
	"first": SampleModeFirst,
	"last":  SampleModeLast,
	"max":   SampleModeMax,
	"min":   SampleModeMin,
	"sum":   SampleModeSum,
}

// ScaleExpr is a scale expression.
type ScaleExpr struct {
	Expr   Expr
	Factor float64
}

func parseScaleExpr(p *parser.Parser) (*ScaleExpr, error) {
	v, err := parseExpr(p)
	if err != nil {
		return nil, err
	}

	_, err = p.Expect(parser.COMMA)
	if err != nil {
		return nil, err
	}

	tok, err := p.Expect(parser.NUMBER)
	if err != nil {
		return nil, err
	}

	f, err := strconv.ParseFloat(tok.String(), 64)
	if err != nil {
		return nil, err
	}

	return &ScaleExpr{Expr: v, Factor: f}, nil
}

// MatchersFromExprs returns labels matchers present in expressions.
func MatchersFromExprs(exprs ...Expr) []labels.Matcher {
	hashes := make(map[uint64]struct{})
	matchers := []labels.Matcher{}

	stack := exprs
	for len(stack) > 0 {
		var cur Expr
		cur, stack = stack[0], stack[1:]

		switch v := cur.(type) {
		case *AggregateExpr:
			stack = append(stack, v.Exprs...)

		case *AliasExpr:
			stack = append(stack, v.Expr)

		case *MatcherExpr:
			// Deduplicate matcher from expressions
			hash := v.Matcher.Hash()

			_, ok := hashes[hash]
			if !ok {
				matchers = append(matchers, v.Matcher)

				hashes[hash] = struct{}{}
			}

		case *SampleExpr:
			stack = append(stack, v.Expr)

		case *ScaleExpr:
			stack = append(stack, v.Expr)
		}
	}

	return matchers
}
