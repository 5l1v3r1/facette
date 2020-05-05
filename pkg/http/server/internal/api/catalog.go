// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"net/http"

	"batou.dev/httprouter"
	"facette.io/facette/pkg/api"
	httpjson "facette.io/facette/pkg/http/json"
)

func (h handler) ListLabels(rw http.ResponseWriter, r *http.Request) {
	opts, err := listOptionsFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	matcher, err := matcherFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	labels := h.catalog.Labels(matcher)
	total := uint(len(labels))

	if opts.Limit > 0 {
		applyCatalogPagination(&labels, total, opts.Offset, opts.Limit)
	}

	httpjson.Write(rw, api.Response{Data: labels, Total: total}, http.StatusOK)
}

func (h handler) ListMetrics(rw http.ResponseWriter, r *http.Request) {
	opts, err := listOptionsFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	matcher, err := matcherFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	metrics := []string{}
	for _, metric := range h.catalog.Metrics(matcher) {
		metrics = append(metrics, metric.String())
	}

	total := uint(len(metrics))

	if opts.Limit > 0 {
		applyCatalogPagination(&metrics, total, opts.Offset, opts.Limit)
	}

	httpjson.Write(rw, api.Response{Data: metrics, Total: total}, http.StatusOK)
}

func (h handler) ListValues(rw http.ResponseWriter, r *http.Request) {
	opts, err := listOptionsFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	matcher, err := matcherFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	values := h.catalog.Values(httprouter.ContextParam(r, "label").(string), matcher)
	total := uint(len(values))

	if opts.Limit > 0 {
		applyCatalogPagination(&values, total, opts.Offset, opts.Limit)
	}

	httpjson.Write(rw, api.Response{Data: values, Total: total}, http.StatusOK)
}

func applyCatalogPagination(entries *[]string, total, offset, limit uint) {
	if offset < total {
		end := offset + limit
		if limit > 0 && total > end {
			*entries = (*entries)[offset:end]
		} else if offset > 0 {
			*entries = (*entries)[offset:]
		}
	} else {
		*entries = []string{}
	}
}
