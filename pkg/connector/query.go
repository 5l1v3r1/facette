// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package connector

import (
	"time"

	"facette.io/facette/pkg/catalog"
	"facette.io/facette/pkg/series"
)

// Query is a time series connector query.
type Query struct {
	From    time.Time
	To      time.Time
	Step    time.Duration
	Metrics []catalog.Metric
}

// Result is a time series connector result.
type Result []Sample

// Sample is a time series connector sample.
type Sample struct {
	Metric catalog.Metric
	Points []series.Point
}
