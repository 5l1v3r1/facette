// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package poller provides the metrics polling system.
package poller

import (
	"context"

	"go.uber.org/zap"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/catalog"
	"facette.io/facette/pkg/store"
)

// Poller is a metrics poller.
type Poller struct {
	config  *Config
	catalog *catalog.Catalog
	jobs    map[string]*job
	evCh    chan event
	log     *zap.Logger
}

// New creates a new metrics poller instance.
func New(config *Config, catalog *catalog.Catalog) *Poller {
	return &Poller{
		config:  config,
		catalog: catalog,
		jobs:    make(map[string]*job),
		evCh:    make(chan event),
		log:     zap.L().Named("poller"),
	}
}

// Run satisfies the fantask.Task interface.
func (p *Poller) Run(ctx context.Context) error {
	p.log.Info("poller started")
	defer p.log.Info("poller stopped")

loop:
	for {
		select {
		case ev := <-p.evCh:
			p.handleEvent(ctx, ev)

		case <-ctx.Done():
			break loop
		}
	}

	return nil
}

// InsertFromStore insert jobs into the metrics poller based on providers
// present and enabled into the back-end storage.
func (p *Poller) InsertFromStore(store *store.Store) {
	for _, job := range p.jobs {
		p.Notify(EventRemove, job.provider)
	}

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

func (p *Poller) handleEvent(ctx context.Context, ev event) {
	switch ev.Type {
	case EventInsert:
		// Shutdown any preexisting job instance
		prevJob, ok := p.jobs[ev.Provider.ID]
		if ok {
			prevJob.Shutdown()
			delete(p.jobs, ev.Provider.ID)
		}

		job, err := newJob(ctx, ev.Provider, p.catalog)
		if err != nil {
			p.log.Error("cannot create provider job", zap.Error(err), zap.String("provider", ev.Provider.ID))
			return
		}

		p.jobs[ev.Provider.ID] = job

		go job.Run()

	case EventPoll:
		// Polling event received. If provider isn't nil, trigger its polling
		// else trigger polling on all registered jobs.
		if ev.Provider != nil {
			job, ok := p.jobs[ev.Provider.ID]
			if ok {
				job.Poll()
			}
		} else {
			for _, job := range p.jobs {
				go job.Poll()
			}
		}

	case EventRemove:
		job, ok := p.jobs[ev.Provider.ID]
		if ok {
			job.Shutdown()
			delete(p.jobs, ev.Provider.ID)
		}
	}
}
