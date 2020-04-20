// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package store provides the back-end storage.
package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	"facette.io/facette/pkg/store/driver"

	"facette.io/facette/pkg/store/internal/types"
)

// Store is a back-end storage.
type Store struct {
	driver driver.Driver
	db     *gorm.DB
}

// New creates a new back-end storage instance.
func New(config *Config) (*Store, error) {
	log := zap.L().Named("store")

	driver, err := driver.New(config.Driver)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	db, err := driver.Open()
	if err != nil {
		log.Error(err.Error(), driver.ZapFields()...)
		return nil, err
	}

	db.SetLogger(logger{})

	log.Info("back-end storage opened", zap.Any("type", config.Driver.Type))

	// Handle driver-specific initialization steps and automatically create
	// or migrate back-end storage database schema.
	log.Debug("initializing back-end storage")

	err = driver.Init(db)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		(*types.Chart)(nil),
		(*types.Dashboard)(nil),
		(*types.Provider)(nil),
	).Error
	if err != nil {
		return nil, err
	}

	log.Debug("back-end storage initialized")

	return &Store{
		driver: driver,
		db:     db,
	}, nil
}

// Close closes the back-end storage database connection.
func (s *Store) Close() {
	if s.db != nil {
		log := zap.L().Named("store")

		err := s.db.Close()
		if err != nil {
			log.Error("cannot close back-end storage", zap.Error(err))
			return
		}

		log.Info("back-end storage closed")
	}
}

// Delete deletes an object from the back-end storage.
func (s *Store) Delete(obj api.Object) error {
	v, err := types.FromAPI(obj)
	if err != nil {
		return err
	}

	count := s.db.Delete(v).RowsAffected
	if count == 0 {
		return api.ErrNotFound
	}

	return nil
}

// Get returns an object from the back-end storage.
func (s *Store) Get(obj api.Object) error {
	var column, value string

	meta := obj.GetMeta()

	switch {
	case meta.ID != "":
		column = "id"
		value = meta.ID

	case meta.Name != "":
		column = "name"
		value = meta.Name

	default:
		return api.ErrNotFound
	}

	v, err := types.FromAPI(obj)
	if err != nil {
		return err
	}

	err = s.db.Where(fmt.Sprintf("%v = ?", s.db.Dialect().Quote(column)), value).First(v).Error
	if err == gorm.ErrRecordNotFound {
		return api.ErrNotFound
	} else if err != nil {
		return s.driver.Error(err)
	}

	return v.Copy(obj)
}

// List returns a list of objects from the back-end storage.
func (s *Store) List(objects api.ObjectList, opts *api.ListOptions) (uint, error) {
	v, err := types.ListFromAPI(objects)
	if err != nil {
		return 0, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, s.driver.Error(tx.Error)
	}

	defer tx.Commit()

	tx = tx.Model(v)

	if opts != nil {
		for k, v := range opts.Filter {
			k, v = s.driver.WhereClause(k, v)
			tx = tx.Where(k, v)
		}
	}

	var total uint

	err = tx.Count(&total).Error
	if err != nil {
		return 0, s.driver.Error(err)
	}

	if opts != nil {
		for _, field := range opts.SortFields() {
			if !tx.Dialect().HasColumn(tx.NewScope(v).TableName(), field.Name) {
				return 0, errors.Wrapf(api.ErrInvalid, "unknown field: %s", field.Name)
			}

			if field.Desc {
				tx = tx.Order(field.Name + " DESC")
			} else {
				tx = tx.Order(field.Name)
			}
		}

		if opts.Limit > 0 {
			tx = tx.Offset(opts.Offset).Limit(opts.Limit)
		}
	}

	err = tx.Find(v).Error
	if err != nil {
		return 0, s.driver.Error(err)
	}

	err = v.Copy(objects)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Save stores an object into the back-end storage.
func (s *Store) Save(obj api.Object, update bool) error {
	v, err := types.FromAPI(obj)
	if err != nil {
		return err
	}

	var db *gorm.DB

	if update {
		db = s.db.Save(v)
	} else {
		db = s.db.Create(v)
	}

	err = s.driver.Error(db.Error)
	if err != nil {
		return err
	}

	return v.Copy(obj)
}
