// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package driver provides back-end storage drivers.
package driver

import (
	"fmt"

	"gorm.io/gorm"
)

// Driver is a back-end storage driver interface.
type Driver interface {
	Error(err error) error
	Info() Info
	Init(db *gorm.DB) error
	Open(config *gorm.Config) (*gorm.DB, error)
	TruncateStmt(table string) string
	WhereClause(column string, v interface{}) (string, interface{})
}

// New creates a new back-end storage driver instance.
func New(config *Config) (Driver, error) {
	fn, ok := drivers[config.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported driver: %s", config.Type)
	}

	return fn(config)
}

// Info are back-end storage driver information.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Type is a back-end storage driver type.
type Type string

// Fn is a back-end storage driver function.
type Fn func(config *Config) (Driver, error)

// Drivers returns a list of supported back-end storage drivers.
func Drivers() []string {
	result := []string{}
	for typ := range drivers {
		result = append(result, string(typ))
	}

	return result
}

// Register registers a back-end storage driver support.
func Register(name Type, fn Fn) {
	drivers[name] = fn
}

var drivers = map[Type]Fn{}
