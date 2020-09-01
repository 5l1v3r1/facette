// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package store provides the back-end storage.
package store

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	"facette.io/facette/pkg/store/driver"

	"facette.io/facette/pkg/store/internal/types"
)

// Store is a back-end storage.
type Store struct {
	driver        driver.Driver
	db            *gorm.DB
	log           *zap.Logger
	restoreCancel context.CancelFunc
}

// New creates a new back-end storage instance.
func New(config *Config) (*Store, error) {
	log := zap.L().Named("store")

	driver, err := driver.New(config.Driver)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	db, err := driver.Open(&gorm.Config{
		Logger: newLogger(config.Debug),
	})
	if err != nil {
		return nil, err
	}

	log.Info("back-end storage opened", zap.Any("type", config.Driver.Type))

	// Handle driver-specific initialization steps and automatically
	// create/migrate back-end storage database schema.
	log.Debug("initializing back-end storage")

	err = driver.Init(db)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		(*types.Chart)(nil),
		(*types.Dashboard)(nil),
		(*types.Provider)(nil),
	)
	if err != nil {
		return nil, err
	}

	s := &Store{
		driver: driver,
		db:     db,
		log:    log,
	}

	log.Debug("back-end storage initialized")

	return s, nil
}

// Close closes the back-end storage database connection.
func (s *Store) Close() {
	if s.db != nil {
		db, err := s.db.DB()
		if err == nil {
			err = db.Close()
		}

		if err != nil {
			s.log.Error("cannot close back-end storage", zap.Error(err))
			return
		}

		s.log.Info("back-end storage closed")
	}
}

// DriverInfo returns back-end storage driver information.
func (s *Store) DriverInfo() driver.Info {
	return s.driver.Info()
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
func (s *Store) Get(obj api.Object, resolve bool) error {
	v, err := s.get(obj)
	if err != nil {
		return s.driver.Error(err)
	}

	// Check for resolution flag. If not set, directly return the object
	// retrieved from the back-end storage; otherwise resolve it.
	if !resolve {
		return v.Copy(obj)
	}

	resolver, ok := v.(types.Resolver)
	if !ok {
		return fmt.Errorf("expected types.Resolver but got %T", v)
	}

	err = resolver.Resolve(types.StoreFuncs{Get: s.get, List: s.list})
	if err != nil {
		return err
	}

	return v.Copy(obj)
}

func (s *Store) get(obj api.Object) (types.Object, error) {
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
		return nil, api.ErrNotFound
	}

	v, err := types.FromAPI(obj)
	if err != nil {
		return nil, err
	}

	err = s.db.Where(fmt.Sprintf("%v = ?", s.db.Statement.Quote(column)), value).First(v).Error
	if err == gorm.ErrRecordNotFound {
		return nil, api.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

// List returns a list of objects from the back-end storage.
func (s *Store) List(objects api.ObjectList, opts *api.ListOptions) (int64, error) {
	v, total, err := s.list(objects, opts)
	if err != nil {
		return 0, s.driver.Error(err)
	}

	err = v.Copy(objects)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Store) list(objects api.ObjectList, opts *api.ListOptions) (types.ObjectList, int64, error) {
	v, err := types.ListFromAPI(objects)
	if err != nil {
		return nil, 0, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, 0, err
	}

	defer tx.Commit()

	tx = tx.Model(v)

	if opts != nil {
		for k, v := range opts.Filter {
			k, v = s.driver.WhereClause(k, v)
			tx = tx.Where(k, v)
		}
	}

	var total int64

	err = tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if opts != nil {
		for _, field := range opts.SortFields() {
			fieldName := tx.NamingStrategy.ColumnName("", field.Name)

			if !tx.Migrator().HasColumn(v, fieldName) {
				return nil, 0, errors.Wrapf(api.ErrInvalid, "unknown field: %s", fieldName)
			}

			if field.Desc {
				tx = tx.Order(fieldName + " DESC")
			} else {
				tx = tx.Order(fieldName)
			}
		}

		if opts.Limit > 0 {
			tx = tx.Offset(int(opts.Offset)).Limit(int(opts.Limit))
		}
	}

	err = tx.Find(v).Error
	if err != nil {
		return nil, 0, err
	}

	return v, total, nil
}

// Save stores an object into the back-end storage.
func (s *Store) Save(obj api.Object, update bool) error {
	v, err := types.FromAPI(obj)
	if err != nil {
		return err
	}

	tx := s.db.Begin()

	if update {
		tx = tx.Save(v)
	} else {
		tx = tx.Create(v)
	}

	err = tx.Error
	if err != nil {
		goto stop
	}

	if err != nil {
		tx.Rollback()
		goto stop
	}

	err = tx.Commit().Error

stop:
	if err != nil {
		return s.driver.Error(tx.Error)
	}

	return v.Copy(obj)
}
