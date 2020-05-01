// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package connector provides upstream time series support.
package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"facette.io/facette/pkg/catalog"
)

// Connector is an upstream time series connector interface.
type Connector interface {
	Metrics(ctx context.Context, ch chan<- catalog.Metric, errCh chan<- error)
	Query(ctx context.Context, q *Query) (Result, error)
	Test(ctx context.Context) error
}

// New creates a new upstream time series connector instance.
func New(typ Type, name string, settings json.RawMessage) (Connector, error) {
	fn, ok := connectors[typ]
	if !ok {
		return nil, fmt.Errorf("unsupported connector: %s", typ)
	}

	return fn(name, settings)
}

// Type is an upstream time series connector type.
type Type string

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (t *Type) UnmarshalText(b []byte) error {
	switch v := Type(b); v {
	case TypeKairosDB, TypePrometheus, TypeRRDtool:
		*t = v
		return nil
	}

	return fmt.Errorf("unsupported connector: %s", b)
}

// Types:
const (
	TypeKairosDB   Type = "kairosdb"
	TypePrometheus Type = "prometheus"
	TypeRRDtool    Type = "rrdtool"
)

// Fn is an upstream time series connector function.
type Fn func(name string, settings json.RawMessage) (Connector, error)

// Connectors returns a list of supported time series connectors.
func Connectors() []string {
	result := []string{}
	for typ := range connectors {
		result = append(result, string(typ))
	}

	return result
}

// Register registers a new upstream time series connector support.
func Register(typ Type, fn Fn) {
	connectors[typ] = fn
}

var connectors = make(map[Type]Fn)
