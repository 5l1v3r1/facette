// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"net/http"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/connector"
	httpjson "facette.io/facette/pkg/http/json"
)

func (handler) getOptions(rw http.ResponseWriter, r *http.Request) {
	httpjson.Write(rw, api.Response{Data: struct {
		Connectors []string `json:"connectors"`
	}{
		Connectors: connector.Connectors(),
	}}, http.StatusOK)
}
