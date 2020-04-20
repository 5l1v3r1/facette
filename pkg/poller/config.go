// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package poller

// Config is a metrics poller configuration.
type Config struct {
	CachePath string `yaml:"cache_path"`
}

// DefaultConfig returns a default metrics poller configuration.
func DefaultConfig() *Config {
	return &Config{
		CachePath: "var/cache/poller",
	}
}
