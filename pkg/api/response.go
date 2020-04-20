// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

// Response is an API response.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Total uint        `json:"total,omitempty"`
	Error string      `json:"error,omitempty"`
}
