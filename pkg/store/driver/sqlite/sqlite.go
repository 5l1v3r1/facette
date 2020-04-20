// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package sqlite provides a SQLite back-end storage driver.
package sqlite

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	"facette.io/facette/pkg/store"
	"facette.io/facette/pkg/store/driver"
)

var dialect gorm.Dialect

// Driver is a SQLite back-end storage driver.
type Driver struct {
	config *driver.Config
}

// New creates a new SQLite back-end storage driver instance.
func New(config *driver.Config) (driver.Driver, error) {
	return &Driver{config}, nil
}

// Error normalizes driver-specific errors.
func (d *Driver) Error(err error) error {
	if err == nil {
		return nil
	}

	var xerr sqlite3.Error

	ok := errors.As(err, &xerr)
	if !ok {
		return err
	}

	switch xerr.ExtendedCode {
	case sqlite3.ErrConstraintPrimaryKey, sqlite3.ErrConstraintUnique:
		return errors.Wrapf(api.ErrConflict, "conflicting field: %s", fieldName(xerr))

	case sqlite3.ErrConstraintForeignKey:
		return errors.Wrapf(api.ErrInvalid, "unknown field reference")

	case sqlite3.ErrConstraintNotNull:
		return errors.Wrapf(api.ErrInvalid, "missing field: %s", fieldName(xerr))
	}

	return errors.Wrapf(api.ErrUnhandled, "database error: %s", err)
}

// Init initializes the back-end storage driver.
func (d *Driver) Init(db *gorm.DB) error {
	return db.Exec("PRAGMA foreign_keys = ON").Error
}

// Open opens a new back-end storage database connection.
func (d *Driver) Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3_custom", d.config.Path)
}

// WhereClause returns driver-specific WHERE clause.
func (d *Driver) WhereClause(column string, v interface{}) (string, interface{}) {
	operator := "="
	extra := ""

	switch tv := v.(type) {
	case store.RegexpPattern:
		operator = "REGEXP"
		v = "(?i)" + string(tv)

	case string:
		if tv == "null" {
			operator = "IS"
			v = nil
		}
	}

	return fmt.Sprintf("%s %s ?%s", column, operator, extra), v
}

// ZapFields returns SQLite-specific back-end storage zap fields.
func (d *Driver) ZapFields() []zap.Field {
	return []zap.Field{
		zap.Any("type", d.config.Type),
		zap.String("path", d.config.Path),
	}
}

var errRegexp = regexp.MustCompile(`^.+: [^.]+\.(.+?)(?:, .+)?$`)

func fieldName(err error) string {
	m := errRegexp.FindStringSubmatch(err.Error())
	if len(m) > 1 {
		return m[1]
	}

	return ""
}

func init() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("regexp", regexp.MatchString, true)
		},
	})

	var ok bool

	dialect, ok = gorm.GetDialect("sqlite3")
	if !ok {
		panic("cannot find SQLite dialect")
	}

	gorm.RegisterDialect("sqlite3_custom", dialect)

	driver.Register(driver.TypeSQLite, New)
}
