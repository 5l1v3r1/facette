// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package json provides HTTP helpers for handling JSON.
package json

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
)

// Unmarshal unmarshals JSON data received from a HTTP request or response
// into a given interface.
func Unmarshal(v, out interface{}) error {
	var (
		body io.ReadCloser
		ct   string
	)

	switch r := v.(type) {
	case *http.Request:
		body = r.Body
		ct = r.Header.Get("Content-Type")

	case *http.Response:
		body = r.Body
		ct = r.Header.Get("Content-Type")

	default:
		return fmt.Errorf("unsupported unmarshal source: %T", r)
	}

	mt, _, err := mime.ParseMediaType(ct)
	if err != nil || mt != "application/json" {
		return api.ErrUnsupportedType
	}

	err = formatError(json.NewDecoder(body).Decode(out))
	if err != nil {
		return errors.Wrap(api.ErrInvalid, err.Error())
	}

	return nil
}

// Write marshals a given interface to JSON and writes its data onto a HTTP
// response writer.
// nolint:errcheck,gosec
func Write(rw http.ResponseWriter, v api.Response, status int) {
	b, err := json.Marshal(v)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(b)
	rw.Write([]byte("\n"))
}
