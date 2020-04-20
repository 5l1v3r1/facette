// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package labels

import (
	"regexp/syntax"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseMatcher(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected Matcher
		err      string
	}{
		{
			input: `foo`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
			},
		},
		{
			input: `foo{}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
			},
		},
		{
			input: `{__name__="foo"}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
			},
		},
		{
			input: `foo{abc="123"}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
				mustMatchCond(t, OpEq, "abc", "123"),
			},
		},
		{
			input: `{abc="123"}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, "abc", "123"),
			},
		},
		{
			input: `foo{abc="123",def!="456",ghi=~"789",jkl!~"042"}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
				mustMatchCond(t, OpEq, "abc", "123"),
				mustMatchCond(t, OpNotEq, "def", "456"),
				mustMatchCond(t, OpEqRegexp, "ghi", "789"),
				mustMatchCond(t, OpNotEqRegexp, "jkl", "042"),
			},
		},
		{
			input: `
			foo{
				abc="123",
				def="456",
			}
			`,
			expected: Matcher{
				mustMatchCond(t, OpEq, Name, "foo"),
				mustMatchCond(t, OpEq, "abc", "123"),
				mustMatchCond(t, OpEq, "def", "456"),
			},
		},
		{
			input: `{abc="1\"2\""}`,
			expected: Matcher{
				mustMatchCond(t, OpEq, "abc", `1"2"`),
			},
		},
		{
			input: ``,
			err:   "expected left brace but got end of input at 1:1",
		},
		{
			input: "foo{abc=\"1\n2\"}",
			err:   "expected string but got new line at 1:11",
		},
		{
			input: `foo{abc="\12"}`,
			err:   "expected string but got bad escape at 1:10",
		},
		{
			input: `foo{abc}`,
			err:   "expected operator but got right brace at 1:8",
		},
		{
			input: `foo{abc!"123"}`,
			err:   "expected operator but got invalid at 1:8",
		},
		{
			input: `foo{abc=123}`,
			err:   "expected string but got number at 1:9",
		},
		{
			input: `foo{"123"}`,
			err:   "expected identifier but got string at 1:5",
		},
		{
			input: `foo{abc="123"}bar`,
			err:   "expected end of input but got identifier at 1:15",
		},
		{
			input: `foo}`,
			err:   "expected end of input but got right brace at 1:4",
		},
		{
			input: `foo{abc="12`,
			err:   "expected string but got end of input at 1:12",
		},
		{
			input: `foo{abc="123"`,
			err:   "expected right brace but got end of input at 1:14",
		},
		{
			input: `
			foo{
				abc:"123"
			}
			`,
			err: "expected operator but got invalid at 3:8",
		},
		{
			input: `foo{abc=~"^1[2"}`,
			err:   (&syntax.Error{Code: syntax.ErrMissingBracket, Expr: "[2"}).Error(),
		},
	} {
		matcher, err := ParseMatcher(test.input)
		if test.err != "" {
			assert.NotNil(t, err, test.input)

			if err != nil {
				assert.Equal(t, test.err, err.Error(), test.input)
			}

			assert.Nil(t, matcher, test.input)
		} else {
			assert.Nil(t, err, test.input)
			assert.Equal(t, test.expected, matcher, test.input)
		}
	}
}

func Test_Matcher_String(t *testing.T) {
	for _, test := range []struct {
		matcher  Matcher
		expected string
	}{
		{
			matcher: Matcher{
				mustMatchCond(t, OpEq, "a", "bc"),
			},
			expected: `{a="bc"}`,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEq, "a", "bc"),
			},
			expected: `{a!="bc"}`,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpEqRegexp, "a", "^b"),
			},
			expected: `{a=~"^b"}`,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEqRegexp, "a", "^b"),
			},
			expected: `{a!~"^b"}`,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpEq, "a", "bc"),
				mustMatchCond(t, OpEqRegexp, "g", "^h"),
			},
			expected: `{a="bc",g=~"^h"}`,
		},
	} {
		assert.Equal(t, test.expected, test.matcher.String(), "matcher: %v", test.matcher)
	}

	// Test panicing condition
	assert.Panics(t, func() {
		cond, _ := NewMatchCond(0, "a", "bc")
		_ = cond.String()
	})
}

func Test_MatchCond_Match(t *testing.T) {
	for _, test := range []struct {
		cond     MatchCond
		input    string
		expected bool
	}{
		{
			cond:     mustMatchCond(t, OpEq, "a", "bc"),
			input:    "bc",
			expected: true,
		},
		{
			cond:     mustMatchCond(t, OpEq, "a", "bc"),
			input:    "bcd",
			expected: false,
		},
		{
			cond:     mustMatchCond(t, OpNotEq, "a", "bc"),
			input:    "bc",
			expected: false,
		},
		{
			cond:     mustMatchCond(t, OpNotEq, "a", "bc"),
			input:    "bcd",
			expected: true,
		},
		{
			cond:     mustMatchCond(t, OpEqRegexp, "a", "^b"),
			input:    "bc",
			expected: true,
		},
		{
			cond:     mustMatchCond(t, OpEqRegexp, "a", "^b"),
			input:    "cd",
			expected: false,
		},
		{
			cond:     mustMatchCond(t, OpNotEqRegexp, "a", "^b"),
			input:    "bcd",
			expected: false,
		},
		{
			cond:     mustMatchCond(t, OpNotEqRegexp, "a", "^b"),
			input:    "cd",
			expected: true,
		},
	} {
		assert.Equal(
			t,
			test.expected,
			test.cond.Match(Label{"a", test.input}),
			"cond: %s, input: %s", test.cond, test.input,
		)
	}

	// Test failing conditions
	_, err := NewMatchCond(OpEqRegexp, "a", "^b[c")
	_, ok := err.(*syntax.Error)
	assert.True(t, ok)

	_, err = NewMatchCond(OpNotEqRegexp, "a", "^b[c")
	_, ok = err.(*syntax.Error)
	assert.True(t, ok)

	assert.Panics(t, func() {
		cond, _ := NewMatchCond(0, "a", "bc")
		cond.Match(Label{"a", "bc"})
	})
}

func mustMatchCond(t *testing.T, op Op, name, value string) MatchCond {
	cond, err := NewMatchCond(op, name, value)
	assert.Nil(t, err)

	return cond
}
