// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package client

import (
	"errors"
	"net"
	"net/url"
)

// Error cleans up error information received from URL errors.
func Error(err error) error {
	var opErr *net.OpError

	ok := errors.As(err, &opErr)
	if ok {
		return opErr
	}

	var urlErr *url.Error

	ok = errors.As(err, &urlErr)
	if ok {
		return urlErr.Err
	}

	return err
}
