// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package driver

// Config is a back-end storage driver configuration.
type Config struct {
	Type Type   `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

// DefaultConfig returns a default logger configuration.
func DefaultConfig() *Config {
	return &Config{
		Type: Type("sqlite"),
		DSN:  "var/store.db",
	}
}
