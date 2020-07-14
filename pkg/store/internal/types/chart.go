// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package types

import (
	"fmt"

	"github.com/imdario/mergo"
	"gorm.io/gorm"

	"facette.io/facette/pkg/api"
)

// Chart is a back-end storage chart object.
type Chart struct {
	ObjectMeta
	Options  *ChartOptions `gorm:"type:text"`
	Series   SeriesList    `gorm:"type:text"`
	LinkID   *string       `gorm:""`
	Link     *Chart        `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Template bool          `gorm:"not null"`
}

func chartFromAPI(chart *api.Chart) *Chart {
	c := &Chart{
		ObjectMeta: objectMetaFromAPI(chart.ObjectMeta),
		Options:    (*ChartOptions)(chart.Options),
		Series:     SeriesList(chart.Series),
		Template:   chart.Template,
	}

	if chart.Link != "" {
		c.LinkID = &chart.Link
	}

	return c
}

// BeforeSave handles the back-end storage ORM "BeforeSave" callback.
func (c *Chart) BeforeSave(db *gorm.DB) error {
	err := c.ObjectMeta.beforeSave(db)
	if err != nil {
		return err
	}

	// Cleanup extraneous options
	switch {
	case c.LinkID != nil: // linked: only keep variables
		c.Options = &ChartOptions{Variables: c.Options.Variables}

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
		Options:    (*api.ChartOptions)(c.Options),
		Series:     api.SeriesList(c.Series),
		Template:   c.Template,
	}

	if c.LinkID != nil {
		chart.Link = *c.LinkID
	}

	return nil
}

// Resolve resolves the back-end storage chart from the linked object given
// data.
func (c *Chart) Resolve(store StoreFuncs) error {
	if c.Template {
		return nil
	}

	var proxy *Chart

	if c.LinkID != nil {
		tmpl := &api.Chart{ObjectMeta: api.ObjectMeta{ID: *c.LinkID}}

		v, err := store.Get(tmpl)
		if err != nil {
			return err
		}

		var ok bool

		proxy, ok = v.(*Chart)
		if !ok {
			return fmt.Errorf("expected *Chart but got %T", v)
		}

		if c.Options != nil {
			err = mergo.Merge(proxy.Options, c.Options, mergo.WithOverride)
			if err != nil {
				return err
			}
		}
	} else {
		proxy = c
	}

	proxy.ObjectMeta = c.ObjectMeta
	proxy.Link = c.Link
	proxy.Template = false

	*c = *proxy

	return nil
}

// ChartList is a back-end storage list of chart objects.
type ChartList []Chart

// Copy satisfies the ObjectList interface.
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

// New satisfies the ObjectList interface.
func (c ChartList) New() Object {
	return &Chart{}
}

// Objects satisfies the ObjectList interface.
func (c ChartList) Objects() []Object {
	l := make([]Object, len(c))

	for idx, chart := range c {
		x := chart
		l[idx] = &x
	}

	return l
}
