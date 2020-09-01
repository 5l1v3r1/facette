// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package sqlite provides a SQLite back-end storage driver.
package sqlite

import (
	"database/sql"
	"reflect"
	"regexp"
	"strings"

	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	"facette.io/facette/pkg/store"
	"facette.io/facette/pkg/store/driver"
)

// Type is the SQLite back-end storage driver type.
const Type driver.Type = "sqlite"

// Driver is a SQLite back-end storage driver.
type Driver struct {
	config    *driver.Config
	dialector gorm.Dialector
}

// New creates a new SQLite back-end storage driver instance.
func New(config *driver.Config) (driver.Driver, error) {
	return &Driver{
		config: config,
		dialector: &dialector{
			Dialector: sqlite.Dialector{DSN: config.DSN},
		},
	}, nil
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
		return errors.Wrapf(api.ErrInvalid, "unknown reference")

	case sqlite3.ErrConstraintNotNull:
		return errors.Wrapf(api.ErrInvalid, "missing field: %s", fieldName(xerr))
	}

	return errors.Wrapf(api.ErrUnhandled, "database error: %s", err)
}

// Info returns driver information.
func (d *Driver) Info() driver.Info {
	version, _, _ := sqlite3.Version()

	return driver.Info{
		Name:    "SQLite",
		Version: version,
	}
}

// Init initializes the back-end storage driver.
func (d *Driver) Init(db *gorm.DB) error {
	return db.Exec("PRAGMA foreign_keys = ON").Error
}

// Open opens a new back-end storage database connection.
func (d *Driver) Open(config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(d.dialector, config)
}

// TruncateStmt returns a driver-specific TRUNCATE statement.
func (d *Driver) TruncateStmt(table string) string {
	var b strings.Builder

	b.WriteString("DELETE FROM ")
	d.dialector.QuoteTo(&b, table)

	return b.String()
}

// WhereClause returns driver-specific WHERE clause.
func (d *Driver) WhereClause(column string, v interface{}) (string, interface{}) {
	operator := "="
	placeholder := "?"

	if reflect.TypeOf(v).Kind() == reflect.Slice {
		operator = "IN"
		placeholder = "(?)"
	} else {
		switch x := v.(type) {
		case store.RegexpPattern:
			operator = "REGEXP"
			v = "(?i)" + string(x)

		case string:
			if x == "null" {
				operator = "IS"
				v = nil
			}
		}
	}

	return strings.Join([]string{column, operator, placeholder}, " "), v
}

var errRegexp = regexp.MustCompile(`^.+: [^.]+\.(.+?)(?:, .+)?$`)

func fieldName(err error) string {
	m := errRegexp.FindStringSubmatch(err.Error())
	if len(m) > 1 {
		return m[1]
	}

	return ""
}

type dialector struct {
	sqlite.Dialector
}

func (d dialector) Initialize(db *gorm.DB) error {
	err := d.Dialector.Initialize(db)
	if err != nil {
		return err
	}

	db.ConnPool, err = sql.Open("sqlite3_custom", d.DSN)

	return err
}

func init() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("regexp", regexp.MatchString, true)
		},
	})

	driver.Register(Type, New)
}
