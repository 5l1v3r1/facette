// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package driver provides back-end storage drivers.
package driver

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// Driver is a back-end storage driver interface.
type Driver interface {
	Error(err error) error
	Init(db *gorm.DB) error
	Open() (*gorm.DB, error)
	WhereClause(column string, v interface{}) (string, interface{})
	ZapFields() []zap.Field
}

// New creates a new back-end storage driver instance.
func New(config *Config) (Driver, error) {
	fn, ok := drivers[config.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported driver: %s", config.Type)
	}

	return fn(config)
}

// Type is a back-end storage driver type.
type Type string

// Types:
const (
	TypeMySQL      Type = "mysql"
	TypePostgreSQL Type = "pgsql"
	TypeSQLite     Type = "sqlite"
)

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
