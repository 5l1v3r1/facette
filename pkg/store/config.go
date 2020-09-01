// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package store

import "facette.io/facette/pkg/store/driver"

// Config is a back-end storage configuration.
type Config struct {
	Debug  bool           `yaml:"debug"`
	Driver *driver.Config `yaml:"driver"`
}

// DefaultConfig returns a default back-end storage configuration.
func DefaultConfig() *Config {
	return &Config{
		Debug:  false,
		Driver: driver.DefaultConfig(),
	}
}
