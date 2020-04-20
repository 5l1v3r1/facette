// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package catalog

import (
	"fmt"
	"regexp"
	"regexp/syntax"

	"facette.io/facette/pkg/errors"
)

// FilterRule is a catalog metrics filter rule.
type FilterRule struct {
	Action  FilterAction      `json:"action"`
	Label   string            `json:"label"`
	Pattern FilterPattern     `json:"pattern"`
	Into    string            `json:"into,omitempty"`
	Targets map[string]string `json:"targets,omitempty"`
}

// FilterAction is a catalog metrics filter action.
type FilterAction string

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (f *FilterAction) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return errors.New("invalid filter action")
	}

	switch v := FilterAction(b); v {
	case FilterActionDiscard, FilterActionRelabel, FilterActionRewrite, FilterActionSieve:
		*f = v
		return nil
	}

	return fmt.Errorf("invalid filter action: %s", b)
}

// Filter actions:
const (
	FilterActionDiscard FilterAction = "discard"
	FilterActionRelabel FilterAction = "relabel"
	FilterActionRewrite FilterAction = "rewrite"
	FilterActionSieve   FilterAction = "sieve"
)

// FilterPattern is a catalog metrics filter pattern.
type FilterPattern struct {
	s  string
	re *regexp.Regexp
}

// MarshalText satisfies the encoding.TextMarshaler interface.
func (f FilterPattern) MarshalText() ([]byte, error) {
	return []byte(f.s), nil
}

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (f *FilterPattern) UnmarshalText(b []byte) error {
	f.s = string(b)

	var err error

	f.re, err = regexp.Compile(f.s)
	if err != nil {
		var xerr *syntax.Error

		ok := errors.As(err, &xerr)
		if ok {
			return fmt.Errorf("invalid filter pattern: %s: `%s`", xerr.Code, xerr.Expr)
		}

		return err
	}

	return nil
}
