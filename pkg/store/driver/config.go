// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package driver

import "fmt"

// Config is a back-end storage driver configuration.
type Config struct {
	Type Type `yaml:"type"`

	// SQLite:
	Path string `yaml:"path"`

	// MySQL and PostgreSQL:
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// DefaultConfig returns a default logger configuration.
func DefaultConfig() *Config {
	return &Config{
		Type: TypeSQLite,

		// SQLite:
		Path: "var/store.db",

		// MySQL and PostgreSQL:
		Host:     "localhost",
		Port:     0, // Use driver-specific default port
		User:     "facette",
		Password: "",
		Database: "facette",
	}
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Proxy Config

	var proxy *Proxy = (*Proxy)(c)

	err := unmarshal(proxy)
	if err != nil {
		return err
	}

	_, ok := drivers[c.Type]
	if !ok {
		return fmt.Errorf("unsupported driver: %s", c.Type)
	}

	return nil
}
