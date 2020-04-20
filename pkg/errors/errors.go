// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package errors provides error wrapping.
package errors

import (
	"errors"
	"fmt"
)

// Stdlib functions:
var (
	As     = errors.As
	Is     = errors.Is
	New    = errors.New
	Unwrap = errors.Unwrap
)

// Wrap wraps an error given a text.
func Wrap(err error, text string) error {
	return &wrappedErr{err, text}
}

// Wrapf wraps an error and formats its text given format specifier and
// arguments.
func Wrapf(err error, format string, args ...interface{}) error {
	var text string
	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	} else {
		text = format
	}

	return &wrappedErr{err, text}
}

type wrappedErr struct {
	err  error
	text string
}

func (w *wrappedErr) Error() string {
	return w.text
}

func (w *wrappedErr) Unwrap() error {
	return w.err
}
