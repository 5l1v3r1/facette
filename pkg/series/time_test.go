// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ParseTime(t *testing.T) {
	now := time.Now().UTC()

	for _, test := range []struct {
		input    string
		expected time.Time
		fail     bool
	}{
		{
			input:    "2016-05-04T03:02:01Z",
			expected: time.Date(2016, 5, 4, 3, 2, 1, 0, time.UTC),
		},
		{
			input:    "now",
			expected: now,
		},
		{
			input:    "-1s",
			expected: now.Add(-1 * time.Second),
		},
		{
			input:    "+1s",
			expected: now.Add(1 * time.Second),
		},
		{
			input:    "-12m",
			expected: now.Add(-12 * time.Minute),
		},
		{
			input:    "-123h",
			expected: now.Add(-123 * time.Hour),
		},
		{
			input:    "-1d",
			expected: now.AddDate(0, 0, -1),
		},
		{
			input:    "-2M",
			expected: now.AddDate(0, -2, 0),
		},
		{
			input:    "-3y",
			expected: now.AddDate(-3, 0, 0),
		},
		{
			input:    "-1h30m15s",
			expected: now.Add(-(1*time.Hour + 30*time.Minute + 15*time.Second)),
		},
		{
			input:    "-1y2M3d1h30m15s",
			expected: now.AddDate(-1, -2, -3).Add(-(1*time.Hour + 30*time.Minute + 15*time.Second)),
		},
		{
			input:    "-1s2s",
			expected: now.Add(-3 * time.Second),
		},
		{
			input:    "-1d2d1M",
			expected: now.AddDate(0, -1, -3),
		},
		{
			input: "",
			fail:  true,
		},
		{
			input: "-",
			fail:  true,
		},
		{
			input: "+",
			fail:  true,
		},
		{
			input: "0",
			fail:  true,
		},
		{
			input: "1",
			fail:  true,
		},
		{
			input: "+1",
			fail:  true,
		},
		{
			input: "-1",
			fail:  true,
		},
		{
			input: "-1x",
			fail:  true,
		},
	} {
		tm, err := ParseTime(test.input)

		if test.fail {
			assert.EqualError(t, err, "invalid time: "+test.input)
			assert.Equal(t, time.Time{}, tm, test.input)
		} else {
			assert.Nil(t, err, test.input)
			assert.Equal(t, test.expected.Add(time.Since(now)).Round(time.Second), tm, test.input)
		}
	}
}
