// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import "strings"

// ListOptions are API listing options.
type ListOptions struct {
	Filter ListFilter
	Sort   string
	Offset uint
	Limit  uint
}

// SortFields returns options sorting fields.
func (o *ListOptions) SortFields() []SortField {
	fields := []SortField{}

	for _, part := range strings.Split(o.Sort, ",") {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "-") {
			fields = append(fields, SortField{strings.TrimPrefix(part, "-"), true})
		} else {
			fields = append(fields, SortField{part, false})
		}
	}

	return fields
}

// ListFilter is an options list filter.
type ListFilter map[string]interface{}

// SortField is an option sorting field.
type SortField struct {
	Name string
	Desc bool
}
