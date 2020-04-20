// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"net/http"

	"facette.io/facette/pkg/api"
	httpjson "facette.io/facette/pkg/http/json"
	"facette.io/facette/pkg/internal/version"
)

func (*handler) GetVersion(rw http.ResponseWriter, r *http.Request) {
	httpjson.Write(rw, api.Response{
		Data: struct {
			Version   string `json:"version,omitempty"`
			Branch    string `json:"branch,omitempty"`
			Revision  string `json:"revision,omitempty"`
			Compiler  string `json:"compiler,omitempty"`
			BuildDate string `json:"buildDate,omitempty"`
		}{
			Version:   version.Version,
			Branch:    version.Branch,
			Revision:  version.Revision,
			Compiler:  version.Compiler,
			BuildDate: version.BuildDate,
		},
	}, http.StatusOK)
}
