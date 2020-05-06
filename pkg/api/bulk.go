// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"encoding/json"
	"net/http"
)

// BulkRequest in an API bulk request.
type BulkRequest struct {
	Method   string                 `json:"method"`
	Endpoint string                 `json:"endpoint"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Data     json.RawMessage        `json:"data,omitempty"`
}

// BulkResult in an API bulk request result.
type BulkResult struct {
	Status   int         `json:"status"`
	Headers  http.Header `json:"headers"`
	Response interface{} `json:"response,omitempty"`
}
