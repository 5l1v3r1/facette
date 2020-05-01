// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"math"
	"time"
)

func sampleSeries(series Series, q *Query, mode SampleMode) Series {
	n := int(q.To.Sub(q.From.Time) / q.Step.Duration)
	result := Series{Name: series.Name, Points: []Point{}}

	if series.Points == nil {
		return result
	}

	buckets := make([]bucket, n)

	// Initialize buckets
	for i := 0; i < n; i++ {
		buckets[i] = bucket{
			Points: []Point{},
			Time:   q.From.Add(time.Duration(i) * q.Step.Duration).Round(time.Second),
		}
	}

	for _, point := range series.Points {
		if point.Time.After(q.From.Time) && point.Time.Before(q.To.Time) {
			idx := int(float64(point.Time.UnixNano()-q.From.UnixNano()) / float64(q.Step.Nanoseconds()))
			if idx >= n {
				// Skip point if index goes beyond current number of points
				continue
			}

			buckets[idx].Points = append(buckets[idx].Points, point)
		}
	}

	for _, b := range buckets {
		result.Points = append(result.Points, b.Sample(mode))
	}

	return result
}

type bucket struct {
	Points []Point
	Time   time.Time
}

func (b bucket) Sample(mode SampleMode) Point {
	point := Point{Time: b.Time}

	switch len(b.Points) {
	case 0:
		point.Value = Value(math.NaN())
		return point

	case 1:
		point.Value = b.Points[0].Value
		return point
	}

	switch mode {
	case SampleModeAverage, SampleModeSum:
		point.Value = sampleAverageSum(b.Points, mode)

	case SampleModeFirst:
		point.Value = sampleFirst(b.Points)

	case SampleModeLast:
		point.Value = sampleLast(b.Points)

	case SampleModeMax:
		point.Value = sampleMax(b.Points)

	case SampleModeMin:
		point.Value = sampleMin(b.Points)
	}

	return point
}

func sampleAverageSum(points []Point, mode SampleMode) Value {
	var (
		sum   float64
		count int
	)

	for _, p := range points {
		if !p.Value.IsNaN() {
			sum += float64(p.Value)
			count++
		}
	}

	if count > 0 {
		if mode == SampleModeAverage {
			return Value(sum / float64(count))
		}

		return Value(sum)
	}

	return Value(math.NaN())
}

func sampleFirst(points []Point) Value {
	for _, p := range points {
		if !p.Value.IsNaN() {
			return p.Value
		}
	}

	return Value(math.NaN())
}

func sampleLast(points []Point) Value {
	idx := len(points) - 1
	for idx >= 0 {
		if !points[idx].Value.IsNaN() {
			return points[idx].Value
		}
		idx--
	}

	return Value(math.NaN())
}

func sampleMax(points []Point) Value {
	value := Value(math.NaN())

	for _, p := range points {
		if !p.Value.IsNaN() && p.Value > value || value.IsNaN() {
			value = p.Value
		}
	}

	return value
}

func sampleMin(points []Point) Value {
	value := Value(math.NaN())

	for _, p := range points {
		if !p.Value.IsNaN() && p.Value < value || value.IsNaN() {
			value = p.Value
		}
	}

	return value
}
