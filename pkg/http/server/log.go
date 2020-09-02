// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package server

import (
	"net/http"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	r   *http.Request
	log *zap.SugaredLogger
}

func (rw responseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)

	rw.log.Debugf(
		`received "%s %s %s" from %s, returned: %d`,
		rw.r.Method,
		rw.r.URL,
		rw.r.Proto,
		rw.r.RemoteAddr,
		status,
	)
}

func debugLog(hh http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hh.ServeHTTP(responseWriter{rw, r, zap.S().Named("http/server")}, r)
	})
}
