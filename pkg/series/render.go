// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"bytes"
	"text/template"

	"facette.io/facette/pkg/catalog"
	"facette.io/facette/pkg/labels"
)

// Render renders series given time boundaries, an expression and a map of
// data points.
func Render(q *Query, metrics map[*catalog.Metric][]Point) []Series {
	series := []Series{}
	for _, expr := range q.Exprs {
		series = append(series, renderExpr(expr, q, metrics, "")...)
	}

	// Get rid of steady gaps in points
	// TODO: reimplement
	// removeGaps(series)

	for idx := range series {
		series[idx].Summarize()
	}

	return series
}

func renderExpr(expr Expr, q *Query, metrics map[*catalog.Metric][]Point, alias string) []Series {
	var series []Series

	switch v := expr.(type) {
	case *AggregateExpr:
		// TODO: implement

	case *AliasExpr:
		series = append(series, renderExpr(v.Expr, q, metrics, v.Alias)...)

	case *MatcherExpr:
		for metric, points := range metrics {
			if metric.Labels.Match(v.Matcher) {
				var name string
				if alias != "" {
					name = renderAlias(metric, alias)
				} else {
					name = metric.String()
				}

				series = append(series, Series{Name: name, Points: points})
			}
		}

	case *SampleExpr:
		for _, s := range renderExpr(v.Expr, q, metrics, "") {
			series = append(series, sampleSeries(s, q, v.Mode))
		}

	case *ScaleExpr:
		for _, s := range renderExpr(v.Expr, q, metrics, "") {
			for idx, p := range s.Points {
				if !p.Value.IsNaN() {
					s.Points[idx].Value *= Value(v.Factor)
				}
			}

			series = append(series, s)
		}
	}

	return series
}

func renderAlias(metric *catalog.Metric, alias string) string {
	funcs := make(template.FuncMap)
	for _, label := range metric.Labels {
		funcs[label.Name] = mapLabel(metric.Labels, label.Name)
	}

	tmpl, err := template.New("alias").Funcs(funcs).Parse(alias)
	if err != nil {
		// Silently fallback to metric name
		return metric.String()
	}

	b := bytes.NewBuffer(nil)
	tmpl.Execute(b, nil) // nolint:errcheck,gosec

	return b.String()
}

func mapLabel(labels labels.Labels, name string) func() string {
	return func() string {
		return labels.Get(name)
	}
}
