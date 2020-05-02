// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

var digits = []byte("0123456789")

// ParseTime parses either a RFC3339 formatted date or a relative time.
func ParseTime(text string) (time.Time, error) {
	t, err := time.ParseInLocation(time.RFC3339Nano, text, time.UTC)
	if err == nil {
		return t.Round(time.Second), nil
	}

	t = time.Now().UTC().Round(time.Second)
	orig := text

	if text == "" {
		return time.Time{}, fmt.Errorf("invalid time: %s", orig)
	} else if text == "now" {
		return t, nil
	}

	var neg bool

	c := text[0]
	if c == '-' || c == '+' {
		neg = c == '-'
		text = text[1:]
	}

	if text == "" {
		return time.Time{}, fmt.Errorf("invalid time: %s", orig)
	}

	var d, m, y int

	for text != "" {
		buf := bytes.NewBuffer(nil)

		for _, c := range text {
			if !bytes.ContainsRune(digits, c) {
				break
			}

			buf.WriteRune(c)
		}

		text = text[buf.Len():]
		if text == "" {
			return time.Time{}, fmt.Errorf("invalid time: %s", orig)
		}

		v, err := strconv.Atoi(buf.String())
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time: %s", orig)
		}

		if neg {
			v = -v
		}

		switch text[0] {
		case 's':
			t = t.Add(time.Duration(v) * time.Second)

		case 'm':
			t = t.Add(time.Duration(v) * time.Minute)

		case 'h':
			t = t.Add(time.Duration(v) * time.Hour)

		case 'd':
			d += v

		case 'M':
			m += v

		case 'y':
			y += v

		default:
			return time.Time{}, fmt.Errorf("invalid time: %s", orig)
		}

		text = text[1:]
	}

	if d != 0 || m != 0 || y != 0 {
		t = t.AddDate(y, m, d)
	}

	return t, nil
}
