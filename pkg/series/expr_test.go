// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"facette.io/facette/pkg/labels"
)

func Test_ParseExpr(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected Expr
		err      string
	}{
		// matcher:
		{
			input: "foo",
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
			},
		},
		{
			input: `{__name__="foo"}`,
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
			},
		},
		{
			input: `foo{key="value"}`,
			expected: &MatcherExpr{
				Matcher: labels.Matcher{
					mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
					mustMatchCond(t, labels.OpEq, "key", "value"),
				},
			},
		},
		{
			input: `
			{
				__name__="foo",
			}`,
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
			},
		},
		{
			input: `
			foo{
				key="value",
			}`,
			expected: &MatcherExpr{
				Matcher: labels.Matcher{
					mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
					mustMatchCond(t, labels.OpEq, "key", "value"),
				},
			},
		},

		// alias:
		{
			input: `alias(foo, "bar")`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Alias: "bar",
			},
		},
		{
			input: `alias({__name__="foo"}, "bar")`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Alias: "bar",
			},
		},
		{
			input: `alias(foo{key="value"}, "bar")`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{
						mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
						mustMatchCond(t, labels.OpEq, "key", "value"),
					},
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				foo,
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				{__name__="foo"},
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				foo{key="value"},
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{
						mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
						mustMatchCond(t, labels.OpEq, "key", "value"),
					},
				},
				Alias: "bar",
			},
		},
		{
			input: "alias",
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "alias")},
			},
		},

		// avg:
		{
			input: `avg(foo, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: `avg({__name__="foo"}, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: `avg(foo{key="value"}, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: `
			avg(
				foo,
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: `
			avg(
				{__name__="foo"},
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: `
			avg(
				foo{key="value"},
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpAverage,
			},
		},
		{
			input: "avg",
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "avg")},
			},
		},

		// sum:
		{
			input: `sum(foo, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `sum({__name__="foo"}, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `sum(foo{key="value"}, bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				foo,
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				{__name__="foo"},
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				foo{key="value"},
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: "sum",
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "sum")},
			},
		},

		// sample:
		{
			input: `sample(foo)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeAverage,
			},
		},
		{
			input: `sample(foo with avg)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeAverage,
			},
		},
		{
			input: `sample(foo with first)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeFirst,
			},
		},
		{
			input: `sample(foo with last)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeLast,
			},
		},
		{
			input: `sample(foo with max)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeMax,
			},
		},
		{
			input: `sample(foo with min)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeMin,
			},
		},
		{
			input: `sample(foo with sum)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeSum,
			},
		},
		{
			input: `sample({__name__="foo"})`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeAverage,
			},
		},
		{
			input: `sample({__name__="foo"} with sum)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeSum,
			},
		},
		{
			input: `
			sample(
				foo
			)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeAverage,
			},
		},
		{
			input: `
			sample(
				{__name__="foo"} with sum
			)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeSum,
			},
		},
		{
			input: `
			sample(
				foo,
			)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeAverage,
			},
		},
		{
			input: `
			sample(
				{__name__="foo"} with sum,
			)`,
			expected: &SampleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Mode: SampleModeSum,
			},
		},

		// scale:
		{
			input: `scale(foo, 1.23)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23,
			},
		},
		{
			input: `scale({__name__="foo"}, 1.23)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23,
			},
		},
		{
			input: `scale(foo{key="value"}, 1.23)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{
						mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
						mustMatchCond(t, labels.OpEq, "key", "value"),
					},
				},
				Factor: 1.23,
			},
		},
		{
			input: `scale(foo, -1.23)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: -1.23,
			},
		},
		{
			input: `scale(foo, 123)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 123,
			},
		},
		{
			input: `scale(foo, 1.23e4)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23e+4,
			},
		},
		{
			input: `scale(foo, 1.23e+4)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23e+4,
			},
		},
		{
			input: `scale(foo, 1.23e-4)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23e-4,
			},
		},
		{
			input: `
			scale(
				foo,
				1.23,
			)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23,
			},
		},
		{
			input: `
			scale(
				{__name__="foo"},
				1.23,
			)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
				},
				Factor: 1.23,
			},
		},
		{
			input: `
			scale(
				foo{key="value"},
				1.23,
			)`,
			expected: &ScaleExpr{
				Expr: &MatcherExpr{
					Matcher: labels.Matcher{
						mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
						mustMatchCond(t, labels.OpEq, "key", "value"),
					},
				},
				Factor: 1.23,
			},
		},
		{
			input: "scale",
			expected: &MatcherExpr{
				Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "scale")},
			},
		},

		// multiple:
		{
			input: `alias(scale(foo, 1.23), "bar")`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `alias(scale({__name__=~"foo"}, 1.23), "bar")`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEqRegexp, labels.Name, "foo")},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `alias(scale(foo{key="value"}, 1.23), "bar")`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				scale(
					foo,
					1.23,
				),
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				scale(
					{__name__="foo"},
					1.23,
				),
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `
			alias(
				scale(
					foo{key="value"},
					1.23,
				),
				"bar",
			)`,
			expected: &AliasExpr{
				Expr: &ScaleExpr{
					Expr: &MatcherExpr{
						Matcher: labels.Matcher{
							mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
							mustMatchCond(t, labels.OpEq, "key", "value"),
						},
					},
					Factor: 1.23,
				},
				Alias: "bar",
			},
		},
		{
			input: `sum(scale(foo, 1.23), bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `sum(scale({__name__="foo"}, 1.23), bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `sum(scale(foo{key="value"}, 1.23), bar)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{
								mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
								mustMatchCond(t, labels.OpEq, "key", "value"),
							},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				scale(
					foo,
					1.23,
				),
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				scale(
					{__name__="foo"},
					1.23,
				),
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "foo")},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},
		{
			input: `
			sum(
				scale(
					foo{key="value"},
					1.23,
				),
				bar,
			)`,
			expected: &AggregateExpr{
				Exprs: []Expr{
					&ScaleExpr{
						Expr: &MatcherExpr{
							Matcher: labels.Matcher{
								mustMatchCond(t, labels.OpEq, labels.Name, "foo"),
								mustMatchCond(t, labels.OpEq, "key", "value"),
							},
						},
						Factor: 1.23,
					},
					&MatcherExpr{
						Matcher: labels.Matcher{mustMatchCond(t, labels.OpEq, labels.Name, "bar")},
					},
				},
				Op: AggregateOpSum,
			},
		},

		// errors:
		{
			input: `{__name__="foo"`,
			err:   "expected right brace but got end of input at 1:16",
		},
		{
			input: `foo{"bc"}`,
			err:   "expected identifier but got string at 1:5",
		},
		{
			input: `foo()`,
			err:   "expected end of input but got left parenthesis at 1:4",
		},
		{
			input: `alias(foo)`,
			err:   "expected comma but got right parenthesis at 1:10",
		},
		{
			input: `alias("bar", foo)`,
			err:   "expected left brace but got string at 1:7",
		},
		{
			input: `alias(foo, "bar"`,
			err:   "expected right parenthesis but got end of input at 1:17",
		},
		{
			input: `alias(foo, 1.23)`,
			err:   "expected string but got number at 1:12",
		},
		{
			input: `avg(foo, "bar")`,
			err:   "expected left brace but got string at 1:10",
		},
		{
			input: `sum(foo, 1.23)`,
			err:   "expected left brace but got number at 1:10",
		},
		{
			input: `sample(1.23)`,
			err:   "expected left brace but got number at 1:8",
		},
		{
			input: `sample(foo, "bar")`,
			err:   "expected right parenthesis but got string at 1:13",
		},
		{
			input: `sample(foo "bar")`,
			err:   "expected identifier but got string at 1:12",
		},
		{
			input: `sample(foo with)`,
			err:   "expected identifier but got right parenthesis at 1:16",
		},
		{
			input: `sample(foo with 1.23)`,
			err:   "expected identifier but got number at 1:17",
		},
		{
			input: `sample(foo invalid avg)`,
			err:   `expected "with" keyword but got "invalid" at 1:12`,
		},
		{
			input: `sample(foo with invalid)`,
			err:   `expected sample mode but got "invalid" at 1:17`,
		},
		{
			input: `scale(foo)`,
			err:   "expected comma but got right parenthesis at 1:10",
		},
		{
			input: `scale(1.23, foo)`,
			err:   "expected left brace but got number at 1:7",
		},
		{
			input: `scale(foo, "bar")`,
			err:   "expected number but got string at 1:12",
		},
	} {
		expr, err := ParseExpr(test.input)
		if test.err != "" {
			assert.NotNil(t, err, test.input)

			if err != nil {
				assert.Equal(t, test.err, err.Error(), test.input)
			}

			assert.Nil(t, expr, test.input)
		} else {
			assert.Nil(t, err, test.input)
			assert.Equal(t, test.expected, expr, test.input)
		}
	}
}

func mustMatchCond(t *testing.T, op labels.Op, name, value string) labels.MatchCond {
	cond, err := labels.NewMatchCond(op, name, value)
	assert.Nil(t, err)

	return cond
}
