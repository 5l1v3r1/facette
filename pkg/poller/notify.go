// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package poller

import (
	"go.uber.org/zap"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/store"
)

// Event is a metrics poller event.
type Event int

// Event types:
const (
	_ Event = iota
	EventInsert
	EventPoll
	EventRemove
)

// Notify sends a metrics poller event.
func (p *Poller) Notify(typ Event, provider *api.Provider) {
	p.evCh <- event{typ, provider}
}

// InsertFromStore insert new jobs into the metrics poller based on providers
// present and enabled in the back-end storage.
func (p *Poller) InsertFromStore(store *store.Store) {
	providers := api.ProviderList{}

	_, err := store.List(&providers, &api.ListOptions{Filter: api.ListFilter{"enabled": true}})
	if err != nil {
		zap.L().Error("cannot list providers", zap.Error(err))
		return
	}

	for _, provider := range providers {
		x := provider
		p.Notify(EventInsert, &x)
	}
}

type event struct {
	Type     Event
	Provider *api.Provider
}
