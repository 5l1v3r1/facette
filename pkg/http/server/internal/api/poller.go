// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"net/http"

	"batou.dev/httprouter"
	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/poller"
)

func (h *handler) PollProvider(rw http.ResponseWriter, r *http.Request) {
	if httprouter.ContextParam(r, "id") != nil {
		provider := &api.Provider{ObjectMeta: metaFromRequest(r)}

		err := h.store.Get(provider)
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
