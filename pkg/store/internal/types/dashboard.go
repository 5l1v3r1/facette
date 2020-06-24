// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package types

import (
	"database/sql"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/set"
)

// Dashboard is a back-end storage dashboard object.
// nolint:lll
type Dashboard struct {
	ObjectMeta
	Options    *DashboardOptions `gorm:"type:text"`
	Layout     GridLayout        `gorm:"type:text"`
	Items      DashboardItems    `gorm:"type:text"`
	Parent     sql.NullString    `gorm:"type:varchar(36); DEFAULT NULL REFERENCES dashboards (id) ON DELETE SET NULL ON UPDATE SET NULL"`
	Link       sql.NullString    `gorm:"type:varchar(36); DEFAULT NULL REFERENCES dashboards (id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Template   bool              `gorm:"not null"`
	References []api.Reference   `gorm:"-"`
}

func dashboardFromAPI(dashboard *api.Dashboard) *Dashboard {
	return &Dashboard{
		ObjectMeta: ObjectMetaFromAPI(dashboard.ObjectMeta),
		Options:    (*DashboardOptions)(dashboard.Options),
		Items:      DashboardItems(dashboard.Items),
		Layout:     GridLayout(dashboard.Layout),
		Parent:     sql.NullString{String: dashboard.Parent, Valid: dashboard.Parent != ""},
		Link:       sql.NullString{String: dashboard.Link, Valid: dashboard.Link != ""},
		Template:   dashboard.Template,
		References: dashboard.References,
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
		d.Options = &DashboardOptions{Variables: d.Options.Variables}

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
		Options:    (*api.DashboardOptions)(d.Options),
		Items:      api.DashboardItems(d.Items),
		Layout:     api.GridLayout(d.Layout),
		Parent:     d.Parent.String,
		Link:       d.Link.String,
		Template:   d.Template,
		References: d.References,
	}

	return nil
}

// Resolve resolves the back-end storage dashboard from the linked object given
// data.
func (d *Dashboard) Resolve(store StoreFuncs) error {
	if d.Template {
		return nil
	}

	var (
		proxy *Dashboard
		err   error
	)

	if d.Link.Valid {
		tmpl := &api.Dashboard{ObjectMeta: api.ObjectMeta{ID: d.Link.String}}

		v, err := store.Get(tmpl)
		if err != nil {
			return err
		}

		var ok bool

		proxy, ok = v.(*Dashboard)
		if !ok {
			return fmt.Errorf("expected *Dashboard but got %T", v)
		}

		if d.Options != nil {
			err = mergo.Merge(proxy.Options, d.Options, mergo.WithOverride)
			if err != nil {
				return err
			}
		}
	} else {
		proxy = d
	}

	// Loop through dashboard items are fetch each known referenced element. If
	// those elements are resolvable, resolve them too.
	refs := map[api.DashboardItemType]*set.Set{}

	for _, item := range proxy.Items {
		_, ok := refs[item.Type]
		if !ok {
			refs[item.Type] = set.New()
		}

		switch item.Type {
		case api.DashboardItemChart:
			id, ok := item.Options["id"]
			if ok {
				refs[item.Type].Add(id)
			}
		}
	}

	if refs["chart"] != nil {
		err = resolveChartReferences(proxy, store, set.StringSlice(refs["chart"])...)
		if err != nil {
			return err
		}
	}

	proxy.ObjectMeta = d.ObjectMeta
	proxy.Link = d.Link
	proxy.Template = false

	*d = *proxy

	return nil
}

func resolveChartReferences(dashboard *Dashboard, store StoreFuncs, ids ...string) error {
	var charts api.ChartList

	list, _, err := store.List(&charts, &api.ListOptions{Filter: api.ListFilter{"id": ids}})
	if err != nil {
		return err
	}

	for _, obj := range list.Objects() {
		err = obj.(Resolver).Resolve(store)
		if err != nil {
			return err
		}

		var chart api.Chart

		err = obj.Copy(&chart)
		if err != nil {
			return err
		}

		dashboard.References = append(dashboard.References, api.Reference{
			Type:  string(api.DashboardItemChart),
			Value: chart,
		})
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

// Objects satisfies the ObjectList interface.
func (d DashboardList) Objects() []Object {
	l := make([]Object, len(d))

	for idx, dashboard := range d {
		x := dashboard
		l[idx] = &x
	}

	return l
}
