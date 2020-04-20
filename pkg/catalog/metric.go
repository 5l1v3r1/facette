// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package catalog

import "facette.io/facette/pkg/labels"

// Metric is a catalog metric.
type Metric struct {
	Labels     labels.Labels
	Attributes Attributes

	section *Section
}

// Connector returns the time series connector associated with the catalog
// metric.
func (m Metric) Connector() interface{} {
	if m.section == nil {
		return nil
	}

	return m.section.connector
}

// String satisfies the fmt.Stringer interface.
func (m Metric) String() string {
	ls := m.Labels.Copy()
	return ls.Pop(labels.Name) + ls.String()
}

// Attributes are catalog metric attributes.
type Attributes map[string]interface{}
