// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package client

import (
	"crypto/tls"
	"net/http"

	"facette.io/facette/pkg/internal/version"
)

// NewRoundTripper creates a new HTTP RoundTripper.
func NewRoundTripper(skipVerify bool) http.RoundTripper {
	rt := &roundTripper{http.DefaultTransport}

	if skipVerify {
		rt.rt.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // nolint:gosec
	}

	return rt
}

type roundTripper struct {
	rt http.RoundTripper
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "facette/"+version.Version)
	return rt.rt.RoundTrip(r)
}
