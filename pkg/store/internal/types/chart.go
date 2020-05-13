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
	"facette.io/facette/pkg/template"
)

// Chart is a back-end storage chart object.
// nolint:lll
type Chart struct {
	ObjectMeta
	Options  ChartOptions   `gorm:"type:text"`
	Series   SeriesList     `gorm:"type:text"`
	Link     sql.NullString `gorm:"type:varchar(36); DEFAULT NULL REFERENCES charts (id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Template bool           `gorm:"not null"`
}

func chartFromAPI(chart *api.Chart) *Chart {
	return &Chart{
		ObjectMeta: ObjectMetaFromAPI(chart.ObjectMeta),
		Options:    ChartOptions(chart.Options),
		Series:     SeriesList(chart.Series),
		Link:       sql.NullString{String: chart.Link, Valid: chart.Link != ""},
		Template:   chart.Template,
	}
}

// BeforeSave handles the back-end storage ORM "BeforeSave" callback.
func (c *Chart) BeforeSave(scope *gorm.Scope) error {
	err := c.ObjectMeta.beforeSave(scope)
	if err != nil {
		return err
	}

	// Cleanup extraneous options
	switch {
	case c.Link.String != "": // linked: only keep variables
		c.Options = ChartOptions{Variables: c.Options.Variables}

	case c.Template: // template: discard any defined variables
		c.Options.Variables = nil
	}

	return nil
}

// Copy copies back-end storage chart data into an API object.
func (c Chart) Copy(dst api.Object) error {
	chart, ok := dst.(*api.Chart)
	if !ok {
		return fmt.Errorf("invalid chart object: %T", dst)
	}

	*chart = api.Chart{
		ObjectMeta: c.ObjectMeta.ToAPI(),
		Options:    api.ChartOptions(c.Options),
		Series:     api.SeriesList(c.Series),
		Link:       c.Link.String,
		Template:   c.Template,
	}

	return nil
}

// Resolve resolves th back-end storage chart from the linked object given data.
func (c *Chart) Resolve(data map[string]string, get func(api.Object) (Object, error)) error {
	if c.Template {
		return nil
	}

	var (
		proxy *Chart
		err   error
	)

	curData := make(map[string]string)

	if c.Link.Valid {
		tmpl := &api.Chart{ObjectMeta: api.ObjectMeta{ID: c.Link.String}}

		v, err := get(tmpl)
		if err != nil {
			return err
		}

		var ok bool

		proxy, ok = v.(*Chart)
		if !ok {
			return fmt.Errorf("expected *Chart but got %T", v)
		}

		err = mergo.Merge(&proxy.Options, c.Options, mergo.WithOverride)
		if err != nil {
			return err
		}

		for _, variable := range c.Options.Variables {
			if !variable.Dynamic {
				curData[variable.Name] = variable.Value
			}
		}
	} else {
		proxy = c
	}

	if data != nil {
		err := mergo.Merge(&curData, data, mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	for idx := range proxy.Series {
		proxy.Series[idx].Expr, err = template.Render(proxy.Series[idx].Expr, curData)
		if err != nil {
			return err
		}
	}

	if proxy.Options.Axes.Y.Left.Label != "" {
		proxy.Options.Axes.Y.Left.Label, err = template.Render(proxy.Options.Axes.Y.Left.Label, curData)
		if err != nil {
			return err
		}
	}

	if proxy.Options.Axes.Y.Right.Label != "" {
		proxy.Options.Axes.Y.Right.Label, err = template.Render(proxy.Options.Axes.Y.Right.Label, curData)
		if err != nil {
			return err
		}
	}

	proxy.Options.Title, err = template.Render(proxy.Options.Title, curData)
	if err != nil {
		return err
	}

	proxy.ObjectMeta = c.ObjectMeta
	proxy.Options.Variables = nil
	proxy.Link = c.Link
	proxy.Template = false

	*c = *proxy

	return nil
}

// ChartList is a back-end storage list of chart objects.
type ChartList []Chart

// Copy copies back-end storage charts list data into an API objects list.
func (c ChartList) Copy(dst api.ObjectList) error {
	list, ok := dst.(*api.ChartList)
	if !ok {
		return fmt.Errorf("invalid charts list: %T", dst)
	}

	var err error

	*list = make(api.ChartList, len(c))

	for idx, chart := range c {
		err = chart.Copy(&(*list)[idx])
		if err != nil {
			return err
		}
	}

	return nil
}
