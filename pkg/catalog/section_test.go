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
	testConnector = struct{ string }{"test"}
	testSection   = NewSection(testConnector)
)

func Test_Section_Insert(t *testing.T) {
	// Test invalid insert
	err := testSection.Insert(Metric{
		Labels: labels.Labels{
			{Name: labels.Name, Value: "\xff"},
		},
	})

	assert.EqualError(t, err, "invalid label value: \xff")
	assert.Len(t, testSection.metrics, len(testMetrics))
}
