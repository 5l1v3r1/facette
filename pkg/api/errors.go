// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import "facette.io/facette/pkg/errors"

// Errors:
var (
	ErrConflict        = errors.New("conflicting data")
	ErrInvalid         = errors.New("invalid data")
	ErrNotFound        = errors.New("not found")
	ErrUnhandled       = errors.New("unhandled error")
	ErrUnsupportedType = errors.New("unsupported type")
)
