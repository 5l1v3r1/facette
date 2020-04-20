// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package set provides values set type.
package set

import (
	"fmt"
	"sync"
)

// Set is a values set.
type Set struct {
	entries map[interface{}]struct{}
	l       sync.RWMutex
}

// New creates a new values set instance.
func New(values ...interface{}) *Set {
	s := &Set{
		entries: make(map[interface{}]struct{}),
	}

	if len(values) > 0 {
		s.Add(values...)
	}

	return s
}

// Add inserts new values into the set.
func (s *Set) Add(values ...interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, v := range values {
		s.entries[v] = struct{}{}
	}
}

// Has returns whether values are present in the set or not.
func (s *Set) Has(values ...interface{}) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	ok := true
	for _, v := range values {
		_, ok = s.entries[v]
		if !ok {
			break
		}
	}

	return ok
}

// Len returns the number of values present in the set.
func (s *Set) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.entries)
}

// Remove removes values from the set.
func (s *Set) Remove(values ...interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, v := range values {
		delete(s.entries, v)
	}
}

// Slice returns a slice of all set values.
func (s *Set) Slice() []interface{} {
	s.l.RLock()
	defer s.l.RUnlock()

	result := []interface{}{}
	for v := range s.entries {
		result = append(result, v)
	}

	return result
}

// StringSlice returns a strings slice representation of all set values.
func StringSlice(s *Set) []string {
	result := []string{}
	for _, v := range s.Slice() {
		result = append(result, fmt.Sprintf("%v", v))
	}

	return result
}
