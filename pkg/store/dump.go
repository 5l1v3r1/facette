// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package store

import (
	"context"

	"gorm.io/gorm"

	"facette.io/facette/pkg/api"

	"facette.io/facette/pkg/store/internal/types"
)

// Dump dumps back-end storage objects.
func (s *Store) Dump(ch chan<- api.Object) error {
	defer close(ch)

	tx := s.db.Begin()
	defer tx.Commit()

	var err error

	for _, objects := range []types.ObjectList{
		&types.ChartList{},
		&types.DashboardList{},
		&types.ProviderList{},
	} {
		err = s.dump(tx, objects, ch)
		if err != nil {
			return err
		}
	}

	// TODO: check for dashboards parents order

	return nil
}

func (s *Store) dump(tx *gorm.DB, objects types.ObjectList, ch chan<- api.Object) error {
	q := tx.Model(objects)

	switch objects.(type) {
	case *types.ChartList:
		q = q.Order("template DESC").Order("created_at")

	case *types.DashboardList:
		q = q.Order("template DESC").Order("parent_id").Order("created_at")

	case *types.ProviderList:
		q = q.Order("created_at")
	}

	rows, err := q.Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		v := objects.New()

		err = tx.ScanRows(rows, v)
		if err != nil {
			return err
		}

		obj, err := types.ToAPI(v)
		if err != nil {
			return err
		}

		ch <- obj
	}

	return nil
}

// Restore restores back-end storage data from objects.
func (s *Store) Restore(ctx context.Context, ch <-chan api.Object) error {
	var err error

	tx := s.db.Begin()

	for _, table := range []string{
		"charts",
		"dashboards",
		"providers",
	} {
		err = tx.Exec(s.driver.TruncateStmt(table)).Error
		if err != nil {
			goto stop
		}
	}

	for {
		select {
		case obj := <-ch:
			v, err := types.FromAPI(obj)
			if err != nil {
				goto stop
			}

			err = tx.Create(v).Error
			if err != nil {
				goto stop
			}

		case <-ctx.Done():
			err = context.Canceled
			goto stop
		}
	}

stop:
	if err != nil {
		tx.Rollback()
		return s.driver.Error(err)
	}

	tx.Commit()

	return nil
}
