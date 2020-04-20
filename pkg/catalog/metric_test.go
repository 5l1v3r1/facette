// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"facette.io/facette/pkg/labels"
)

var (
	testMetrics = []Metric{
		{
			Labels: labels.Labels{
				{Name: labels.Name, Value: "foo"},
				{Name: "abc", Value: "123"},
			},
			Attributes: map[string]interface{}{
				"foo": "abc",
				"bar": 123,
			},
			section: testSection,
		},
		{
			Labels: labels.Labels{
				{Name: labels.Name, Value: "foo"},
				{Name: "abc", Value: "456"},
			},
			section: testSection,
		},
		{
			Labels: labels.Labels{
				{Name: labels.Name, Value: "bar"},
				{Name: "abc", Value: "123"},
				{Name: "def", Value: "456"},
			},
			section: testSection,
		},
	}

	testMetricsStrings = []string{
		`foo{abc="123"}`,
		`foo{abc="456"}`,
		`bar{abc="123",def="456"}`,
	}
)

func Test_Metric_String(t *testing.T) {
	for idx, metric := range testMetrics {
		assert.Equal(t, testMetricsStrings[idx], metric.String())
	}
}
