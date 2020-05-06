// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"facette.io/facette/pkg/api"
	httpjson "facette.io/facette/pkg/http/json"
	"go.uber.org/zap"
)

func (h handler) ExecBulk(rw http.ResponseWriter, r *http.Request) {
	bulk := []api.BulkRequest{}

	err := httpjson.Unmarshal(r, &bulk)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	results := make([]api.BulkResult, len(bulk))

	for idx, unit := range bulk {
		rec := httptest.NewRecorder()

		req, err := requestFromBulkUnit(r.Context(), unit)
		if err != nil {
			results[idx].Status = http.StatusInternalServerError

			h.log.Error("cannot create sub-request", zap.Error(err))

			continue
		}

		h.router.ServeHTTP(rec, req)

		results[idx] = api.BulkResult{
			Status:  rec.Code,
			Headers: rec.Header(),
		}

		if rec.Body.Len() > 0 {
			err = json.Unmarshal(rec.Body.Bytes(), &results[idx].Response)
			if err != nil {
				results[idx].Status = http.StatusInternalServerError

				h.log.Error("cannot unmarshal sub-request JSON", zap.Error(err))
			}
		}
	}

	httpjson.Write(rw, api.Response{Data: results}, http.StatusOK)
}

func requestFromBulkUnit(ctx context.Context, r api.BulkRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		r.Method,
		api.Prefix+"/"+strings.TrimLeft(r.Endpoint, "/"),
		bytes.NewReader(r.Data),
	)
	if err != nil {
		return nil, err
	}

	switch r.Method {
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		req.Header.Set("Content-Type", "application/json")
	}

	q := req.URL.Query()

	for key, value := range r.Params {
		q.Set(key, fmt.Sprintf("%v", value))
	}

	req.URL.RawQuery = q.Encode()

	return req, nil
}
