// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

//go:generate go run generate.go

// Package types provides back-end storage internal types.
package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
)

// ObjectMeta are back-end storage object metadata.
type ObjectMeta struct {
	ID       string    `gorm:"type:varchar(36);not null;primary_key"`
	Name     string    `gorm:"type:varchar(128);not null;unique_index"`
	Created  time.Time `gorm:"not null;default:current_timestamp"`
	Modified time.Time `gorm:"not null;default:current_timestamp"`
}

// ObjectMetaFromAPI returns a back-end storage representation of an API object
// metadata.
func ObjectMetaFromAPI(meta api.ObjectMeta) ObjectMeta {
	return ObjectMeta{
		ID:       meta.ID,
		Name:     meta.Name,
		Created:  meta.Created,
		Modified: meta.Modified,
	}
}

// ToAPI returns an API representation of the back-end storage object metadata.
func (m ObjectMeta) ToAPI() api.ObjectMeta {
	return api.ObjectMeta{
		ID:       m.ID,
		Name:     m.Name,
		Created:  m.Created,
		Modified: m.Modified,
	}
}

func (m *ObjectMeta) beforeSave(scope *gorm.Scope) error {
	if m.Name == "" {
		return errors.New("invalid name")
	}

	if m.ID == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		err = scope.SetColumn("ID", uuid.String())
		if err != nil {
			return err
		}
	} else {
		_, err := uuid.Parse(m.ID)
		if err != nil {
			return errors.New("invalid identifier")
		}
	}

	now := time.Now().UTC()

	err := scope.SetColumn("Modified", now)
	if err != nil {
		return err
	}

	if m.Created.IsZero() {
		err = scope.SetColumn("Created", now)
		if err != nil {
			return err
		}
	}

	return nil
}

// Object is a back-end storage object interface.
type Object interface {
	Copy(dst api.Object) error
}

// FromAPI returns a back-end storage representation of an API object.
func FromAPI(obj api.Object) (Object, error) {
	switch x := obj.(type) {
	case *api.Chart:
		return chartFromAPI(x), nil

	case *api.Dashboard:
		return dashboardFromAPI(x), nil

	case *api.Provider:
		return providerFromAPI(x), nil
	}

	return nil, fmt.Errorf("unhandled object: %T", obj)
}

// ObjectList is a back-end storage objects list interface.
type ObjectList interface {
	Copy(dst api.ObjectList) error
	Objects() []Object
}

// ListFromAPI returns a back-end storage representation of an API objects list.
func ListFromAPI(v interface{}) (ObjectList, error) {
	switch v.(type) {
	case *api.ChartList:
		return &ChartList{}, nil

	case *api.DashboardList:
		return &DashboardList{}, nil

	case *api.ProviderList:
		return &ProviderList{}, nil
	}

	return nil, fmt.Errorf("unhandled objects list: %T", v)
}

// Resolver is a back-end storage resolver interface.
type Resolver interface {
	Resolve(store StoreFuncs) error
}

// StoreFuncs are back-end storage functions.
type StoreFuncs struct {
	Get  func(obj api.Object) (Object, error)
	List func(objects api.ObjectList, opts *api.ListOptions) (ObjectList, uint, error)
}
