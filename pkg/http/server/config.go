// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package server

import "time"

// Config is a HTTP server configuration.
type Config struct {
	Address         string        `yaml:"address"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DefaultConfig returns a default HTTP server configuration.
func DefaultConfig() *Config {
	return &Config{
		Address:         ":12003",
		ShutdownTimeout: 10 * time.Second,
	}
}
