// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"batou.dev/httprouter"
	"github.com/google/uuid"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	httpjson "facette.io/facette/pkg/http/json"
	"facette.io/facette/pkg/poller"
)

const (
	sectionCharts     = "charts"
	sectionDashboards = "dashboards"
	sectionProviders  = "providers"
)

func (h handler) DeleteObject(rw http.ResponseWriter, r *http.Request) {
	section := httprouter.ContextParam(r, "section").(string)
	meta := metaFromRequest(r)

	var obj api.Object

	switch section {
	case sectionCharts:
		obj = &api.Chart{ObjectMeta: meta}

	case sectionDashboards:
		obj = &api.Dashboard{ObjectMeta: meta}

	case sectionProviders:
		obj = &api.Provider{ObjectMeta: meta}

	default:
		panic("unknown section")
	}

	err := h.store.Delete(obj)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	// Notify poller about provider deletion
	if section == sectionProviders {
		h.poller.Notify(poller.EventRemove, obj.(*api.Provider))
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h handler) GetObject(rw http.ResponseWriter, r *http.Request) {
	meta := metaFromRequest(r)

	var obj api.Object

	switch httprouter.ContextParam(r, "section").(string) {
	case sectionCharts:
		obj = &api.Chart{ObjectMeta: meta}

	case sectionDashboards:
		obj = &api.Dashboard{ObjectMeta: meta}

	case sectionProviders:
		obj = &api.Provider{ObjectMeta: meta}

	default:
		panic("unknown section")
	}

	err := h.store.Get(obj, false, nil)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	excerpt, _ := strconv.ParseBool(httprouter.QueryParam(r, "excerpt"))
	if excerpt {
		httpjson.Write(rw, api.Response{Data: obj.Excerpt()}, http.StatusOK)
	} else {
		httpjson.Write(rw, api.Response{Data: obj}, http.StatusOK)
	}
}

func (h handler) GetObjectVars(rw http.ResponseWriter, r *http.Request) {
	meta := metaFromRequest(r)

	var tmpl api.Template

	switch httprouter.ContextParam(r, "section").(string) {
	case sectionCharts:
		tmpl = &api.Chart{ObjectMeta: meta}

	case sectionDashboards:
		tmpl = &api.Dashboard{ObjectMeta: meta}

	default:
		panic("unknown section")
	}

	err := h.store.Get(tmpl, false, nil)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	var variables []string

	variables, err = tmpl.Variables()
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	httpjson.Write(rw, api.Response{Data: variables}, http.StatusOK)
}

func (h handler) ListObjects(rw http.ResponseWriter, r *http.Request) {
	opts, err := listOptionsFromRequest(r)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	var list api.ObjectList

	switch httprouter.ContextParam(r, "section").(string) {
	case sectionCharts:
		list = &api.ChartList{}

	case sectionDashboards:
		list = &api.DashboardList{}

	case sectionProviders:
		list = &api.ProviderList{}

	default:
		panic("unknown section")
	}

	total, err := h.store.List(list, opts)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	result := make([]interface{}, list.Len())
	for idx, object := range list.Objects() {
		result[idx] = object.Excerpt()
	}

	httpjson.Write(rw, api.Response{Data: result, Total: total}, http.StatusOK)
}

func (h handler) ResolveObject(rw http.ResponseWriter, r *http.Request) {
	meta := metaFromRequest(r)

	var tmpl api.Template

	switch httprouter.ContextParam(r, "section").(string) {
	case sectionCharts:
		tmpl = &api.Chart{ObjectMeta: meta}

	case sectionDashboards:
		tmpl = &api.Dashboard{ObjectMeta: meta}

	default:
		panic("unknown section")
	}

	var (
		data map[string]string
		err  error
	)

	if r.ContentLength > 0 {
		err = httpjson.Unmarshal(r, &data)
		if err != nil {
			h.WriteError(rw, err)
			return
		}
	}

	err = h.store.Get(tmpl, true, data)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	httpjson.Write(rw, api.Response{Data: tmpl}, http.StatusOK)
}

func (h handler) SaveObject(rw http.ResponseWriter, r *http.Request) {
	section := httprouter.ContextParam(r, "section").(string)

	var obj api.Object

	switch section {
	case sectionCharts:
		obj = &api.Chart{}

	case sectionDashboards:
		obj = &api.Dashboard{}

	case sectionProviders:
		obj = &api.Provider{}

	default:
		panic("unknown section")
	}

	var err error

	if r.Method == http.MethodPatch {
		// Object is being patched, thus retrieve its current state from the
		// back-end storage.
		obj.SetMeta(metaFromRequest(r))

		err = h.store.Get(obj, false, nil)
		if err != nil {
			h.WriteError(rw, err)
			return
		}
	} else {
		copy := httprouter.QueryParam(r, "copy")
		if copy != "" {
			// Object is being copied from an existing one, thus retrieve its
			// current state and reset metadata to create a whole new object.
			obj.SetMeta(api.ObjectMeta{ID: copy})

			err = h.store.Get(obj, false, nil)
			if err != nil {
				h.WriteError(rw, err)
				return
			}

			obj.SetMeta(api.ObjectMeta{})
		}
	}

	err = httpjson.Unmarshal(r, obj)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	err = obj.Validate()
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	if r.Method == http.MethodPut {
		// Object is being overwritten, thus inherit identifier from request if
		// none provided or ensure both values are matching.
		meta := obj.GetMeta()
		reqMeta := metaFromRequest(r)

		if meta.ID == "" {
			meta.ID = reqMeta.ID
			obj.SetMeta(meta)
		} else if meta.ID != reqMeta.ID {
			h.WriteError(rw, errors.Wrap(api.ErrInvalid, "mismatching identifiers"))
			return
		}
	}

	err = h.store.Save(obj, r.Method != http.MethodPost)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	// Notify poller for insertion/removal depending on provider enable state
	if section == sectionProviders {
		provider := obj.(*api.Provider)
		if provider.Enabled {
			h.poller.Notify(poller.EventInsert, provider)
		} else {
			h.poller.Notify(poller.EventRemove, provider)
		}
	}

	if r.Method == http.MethodPost {
		rw.Header().Set("Location", fmt.Sprintf("%s/providers/%s", api.Prefix, obj.GetMeta().ID))
		rw.WriteHeader(http.StatusCreated)
	} else {
		rw.WriteHeader(http.StatusNoContent)
	}
}

func metaFromRequest(r *http.Request) api.ObjectMeta {
	meta := api.ObjectMeta{}
	id := httprouter.ContextParam(r, "id").(string)

	_, err := uuid.Parse(id)
	if err == nil {
		meta.ID = id
	} else {
		meta.Name = id
	}

	return meta
}
