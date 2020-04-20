// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package logger

import "fmt"

// Config is a logger configuration.
type Config struct {
	Level    string   `yaml:"level"`
	Path     string   `yaml:"path"`
	Encoding Encoding `yaml:"encoding"`
}

// DefaultConfig returns a default logger configuration.
func DefaultConfig() *Config {
	return &Config{
		Level:    "info",
		Path:     "stdout",
		Encoding: "console",
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

	switch c.Encoding {
	case EncodingConsole, EncodingJSON:
		return nil

	default:
		return fmt.Errorf("unsupported log encoding: %s", c.Encoding)
	}
}

// Encoding is logger encoding.
type Encoding string

// Encodings:
const (
	EncodingConsole Encoding = "console"
	EncodingJSON    Encoding = "json"
)
