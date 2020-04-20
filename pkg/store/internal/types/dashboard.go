// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package types

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"facette.io/facette/pkg/api"
)

// Dashboard is a back-end storage dashboard object.
// nolint:lll
type Dashboard struct {
	ObjectMeta
	Options  DashboardOptions `gorm:"type:text"`
	Items    DashboardItems   `gorm:"type:text"`
	Layout   json.RawMessage  `gorm:"type:text"`
	Parent   sql.NullString   `gorm:"type:varchar(36); DEFAULT NULL REFERENCES dashboards (id) ON DELETE SET NULL ON UPDATE SET NULL"`
	Link     sql.NullString   `gorm:"type:varchar(36); DEFAULT NULL REFERENCES dashboards (id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Template bool             `gorm:"not null"`
}

func dashboardFromAPI(dashboard *api.Dashboard) *Dashboard {
	return &Dashboard{
		ObjectMeta: ObjectMetaFromAPI(dashboard.ObjectMeta),
		Options:    DashboardOptions(dashboard.Options),
		Items:      DashboardItems(dashboard.Items),
		Layout:     dashboard.Layout,
		Parent:     sql.NullString{String: dashboard.Parent, Valid: dashboard.Parent != ""},
		Link:       sql.NullString{String: dashboard.Link, Valid: dashboard.Link != ""},
		Template:   dashboard.Template,
	}
}

// BeforeSave handles the back-end storage ORM "BeforeSave" callback.
func (d *Dashboard) BeforeSave(scope *gorm.Scope) error {
	err := d.ObjectMeta.beforeSave(scope)
	if err != nil {
		return err
	}

	// Cleanup extraneous options
	switch {
	case d.Link.String != "": // linked: only keep variables
		d.Options = DashboardOptions{Variables: d.Options.Variables}

	case d.Template: // template: discard any defined variables
		d.Options.Variables = nil
	}

	return nil
}

// Copy copies back-end storage dashboard data into an API object.
func (d Dashboard) Copy(dst api.Object) error {
	dashboard, ok := dst.(*api.Dashboard)
	if !ok {
		return fmt.Errorf("invalid dashboard object: %T", dst)
	}

	*dashboard = api.Dashboard{
		ObjectMeta: d.ObjectMeta.ToAPI(),
		Options:    api.DashboardOptions(d.Options),
		Items:      api.DashboardItems(d.Items),
		Layout:     d.Layout,
		Template:   d.Template,
	}

	return nil
}

// DashboardList is a back-end storage list of dashboard objects.
type DashboardList []Dashboard

// Copy copies back-end storage dashboards list data into an API objects list.
func (d DashboardList) Copy(dst api.ObjectList) error {
	list, ok := dst.(*api.DashboardList)
	if !ok {
		return fmt.Errorf("invalid dashboards list: %T", dst)
	}

	var err error

	*list = make(api.DashboardList, len(d))

	for idx, dashboard := range d {
		err = dashboard.Copy(&(*list)[idx])
		if err != nil {
			return err
		}
	}

	return nil
}
