// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"net/http"
	"time"

	"batou.dev/httprouter"
	"golang.org/x/net/context"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/connector"
	"facette.io/facette/pkg/errors"
	httpjson "facette.io/facette/pkg/http/json"
	"facette.io/facette/pkg/poller"
)

func (h handler) PollProvider(rw http.ResponseWriter, r *http.Request) {
	if httprouter.ContextParam(r, "id") != nil {
		provider := &api.Provider{ObjectMeta: metaFromRequest(r)}

		err := h.store.Get(provider, false, nil)
		if err != nil {
			h.WriteError(rw, err)
			return
		}

		h.poller.Notify(poller.EventPoll, provider)
	} else {
		h.poller.Notify(poller.EventPoll, nil)
	}

	rw.WriteHeader(http.StatusAccepted)
}

func (h handler) TestProvider(rw http.ResponseWriter, r *http.Request) {
	obj := &api.Provider{}

	err := httpjson.Unmarshal(r, obj)
	if err != nil {
		writeTestResponse(rw, err)
		return
	}

	err = obj.Validate()
	if err != nil {
		writeTestResponse(rw, err)
		return
	}

	conn, err := connector.New(obj.Connector.Type, "test", obj.Connector.Settings)
	if err != nil {
		writeTestResponse(rw, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = conn.Test(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = errors.New("timeout reached")
		}

		writeTestResponse(rw, err)

		return
	}

	writeTestResponse(rw, nil)
}

func writeTestResponse(rw http.ResponseWriter, err error) {
	var resp api.Response

	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Data = map[string]interface{}{"success": true}
	}

	httpjson.Write(rw, resp, http.StatusOK)
}
