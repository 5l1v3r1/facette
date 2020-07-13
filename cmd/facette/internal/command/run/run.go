// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package run provides a run sub-command.
package run

import (
	"context"
	"syscall"

	"batou.dev/fantask"
	"github.com/urfave/cli/v2"

	"facette.io/facette/pkg/catalog"
	httpserver "facette.io/facette/pkg/http/server"
	"facette.io/facette/pkg/logger"
	"facette.io/facette/pkg/poller"
	"facette.io/facette/pkg/store"

	"facette.io/facette/cmd/facette/internal/config"
)

// Command is a run command.
var Command = &cli.Command{
	Name:   "run",
	Usage:  "Run service",
	Action: run,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "configuration file path",
			Value:   config.DefaultPath,
		},
	},
}

func run(ctx *cli.Context) error {
	cfg, err := config.Load(ctx.String("config"))
	if err != nil {
		return err
	}

	logger, err := logger.New(cfg.Log)
	if err != nil {
		return err
	}

	defer logger.Sync() // nolint:errcheck

	// Initialize metrics catalog and back-end storage. Database schema will be
	// created or automatically migrated upon startup.
	catalog := catalog.New()

	store, err := store.New(cfg.Store)
	if err != nil {
		return err
	}
	defer store.Close()

	// Run tasks concurrently. Tasks will be stopped as soon as one of them
	// returns or if any termination signal is received.
	poller := poller.New(cfg.Poller, catalog)
	httpserver := httpserver.New(cfg.HTTP, catalog, store, poller)

	tasks := fantask.NewGroup().
		Add(poller).
		Add(httpserver)

	go tasks.CancelWithSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go poller.InsertFromStore(store)

	return fantask.IgnoreCanceled(tasks.Run(context.Background()))
}
