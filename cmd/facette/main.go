// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"facette.io/facette/cmd/facette/internal/command/run"
	"facette.io/facette/cmd/facette/internal/command/store"
	_ "facette.io/facette/cmd/facette/internal/connector" // Import connectors
	_ "facette.io/facette/cmd/facette/internal/driver"    // Import drivers
)

func main() {
	app := &cli.App{
		Name:  "facette",
		Usage: "Time series data visualization software",
		Commands: []*cli.Command{
			run.Command,
			store.Command,
		},
		HideHelpCommand: true,
		Action:          cli.ShowAppHelp,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
