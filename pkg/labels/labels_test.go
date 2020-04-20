// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package labels

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLabels Labels

func Test_Labels_New(t *testing.T) {
	testLabels = New(
		Label{"foo", "abc"},
		Label{"bar", "123"},
	)

	assert.Equal(t, Labels{
		Label{"foo", "abc"},
		Label{"bar", "123"},
	}, testLabels)
}

func Test_Labels_Append(t *testing.T) {
	ls := testLabels.Copy()
	ls.Append(Label{"baz", "value"})

	assert.Equal(t, Labels{
		Label{"foo", "abc"},
		Label{"bar", "123"},
		Label{"baz", "value"},
	}, ls)
}

func Test_Labels_Copy(t *testing.T) {
	ls := testLabels.Copy()
	ls.Append(Label{"baz", "value"})

	assert.Equal(t, Labels{
		Label{"foo", "abc"},
		Label{"bar", "123"},
	}, testLabels)
}

func Test_Labels_Get(t *testing.T) {
	v := testLabels.Get("foo")
	assert.Equal(t, "abc", v)

	v = testLabels.Get("unknown")
	assert.Equal(t, "", v)
}

func Test_Labels_Match(t *testing.T) {
	ls := testLabels.Copy()
	ls.Append(Label{"baz", "value"})

	for _, test := range []struct {
		matcher  Matcher
		expected bool
	}{
		{
			matcher:  nil,
			expected: true,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpEq, "foo", "abc"),
			},
			expected: true,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEq, "foo", "abc"),
			},
			expected: false,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpEqRegexp, "foo", "^a"),
			},
			expected: true,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEqRegexp, "foo", "^a"),
			},
			expected: false,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpEq, "foo", "abc"),
				mustMatchCond(t, OpEqRegexp, "baz", "^v"),
			},
			expected: true,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEq, "foo", "abc"),
				mustMatchCond(t, OpEqRegexp, "baz", "^v"),
			},
			expected: false,
		},
		{
			matcher: Matcher{
				mustMatchCond(t, OpNotEq, "unknown", "abc"),
			},
			expected: false,
		},
	} {
		assert.Equal(t, test.expected, ls.Match(test.matcher), "matcher: %s", test.matcher)
	}
}

func Test_Labels_Pop(t *testing.T) {
	ls := testLabels.Copy()

	v := ls.Pop("foo")
	assert.Equal(t, "abc", v)
	assert.Equal(t, Labels{Label{"bar", "123"}}, ls)

	v = ls.Pop("unknown")
	assert.Equal(t, "", v)
	assert.Equal(t, Labels{Label{"bar", "123"}}, ls)
}

func Test_Labels_Sort(t *testing.T) {
	ls := testLabels.Copy()
	sort.Sort(ls)

	assert.Equal(t, Labels{
		Label{"bar", "123"},
		Label{"foo", "abc"},
	}, ls)
}

func Test_Labels_String(t *testing.T) {
	assert.Equal(t, `{foo="abc",bar="123"}`, testLabels.String())
}

func Test_Labels_Validate(t *testing.T) {
	err := testLabels.Validate()
	assert.Nil(t, err)

	ls := testLabels.Copy()
	ls.Append(Label{"invalid", ""})
	err = ls.Validate()
	assert.EqualError(t, err, "empty label value")

	ls = testLabels.Copy()
	ls.Append(Label{"invalid", "\xff"})
	err = ls.Validate()
	assert.EqualError(t, err, "invalid label value: \xff")

	ls = testLabels.Copy()
	ls.Append(Label{"", "value"})
	err = ls.Validate()
	assert.EqualError(t, err, "empty label name")

	ls = testLabels.Copy()
	ls.Append(Label{"invalid!", "value"})
	err = ls.Validate()
	assert.EqualError(t, err, "invalid label name: invalid!")

	ls = testLabels.Copy()
	ls.Append(Label{"foo", "duplicate"})
	sort.Sort(ls)
	err = ls.Validate()
	assert.EqualError(t, err, "duplicate label name: foo")
}
