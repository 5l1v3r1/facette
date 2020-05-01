// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package series provides time series parsing and computation.
package series

import (
	"encoding/json"
	"math"
	"time"
)

// Series is a time series.
type Series struct {
	Name    string           `json:"name"`
	Points  []Point          `json:"points"`
	Summary map[string]Value `json:"summary"`
}

// Summarize calculates the minimum, maximum, average, first, and last values
// for the time series.
func (s *Series) Summarize() {
	var total, valid Value

	min := Value(math.NaN())
	max := Value(math.NaN())
	last := Value(math.NaN())

	for _, p := range s.Points {
		if !p.Value.IsNaN() {
			last = p.Value

			if min.IsNaN() || p.Value < min {
				min = p.Value
			}

			if max.IsNaN() || p.Value > max {
				max = p.Value
			}

			total += p.Value
			valid++
		}
	}

	s.Summary = map[string]Value{
		"min":  min,
		"max":  max,
		"avg":  total / valid,
		"last": last,
	}
}

// Point is a time series data point.
type Point struct {
	Time  time.Time
	Value Value
}

// MarshalJSON satisfies the json.Marshaler interface.
func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		float64(p.Time.UnixNano()) * 1e-9,
		p.Value,
	})
}

// Value is a time series data point value.
type Value float64

// IsNaN returns whether or not the value is an IEEE 754 'not-a-number' value.
func (v Value) IsNaN() bool {
	return math.IsNaN(float64(v))
}

// MarshalJSON satisfies the json.Marshaler interface.
func (v Value) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) {
		return json.Marshal(nil)
	}

	return json.Marshal(f)
}
