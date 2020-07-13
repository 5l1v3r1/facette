// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package poller

import (
	"facette.io/facette/pkg/api"
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

type event struct {
	Type     Event
	Provider *api.Provider
}
