// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package types

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"facette.io/facette/pkg/api"
)

// Provider is a back-end storage provider object.
type Provider struct {
	ObjectMeta
	Connector    ProviderConnector `gorm:"type:varchar(32);not null"`
	Filters      ProviderFilters   `gorm:"type:text"`
	PollInterval int               `gorm:"not null"`
	Enabled      bool              `gorm:"not null"`
}

func providerFromAPI(provider *api.Provider) *Provider {
	return &Provider{
		ObjectMeta:   objectMetaFromAPI(provider.ObjectMeta),
		Connector:    ProviderConnector(provider.Connector),
		Filters:      ProviderFilters(provider.Filters),
		PollInterval: int(provider.PollInterval.Seconds()),
		Enabled:      provider.Enabled,
	}
}

// BeforeSave handles the back-end storage ORM "BeforeSave" callback.
func (p *Provider) BeforeSave(db *gorm.DB) error {
	return p.ObjectMeta.beforeSave(db)
}

// Copy copies back-end storage provider data into an API object.
func (p Provider) Copy(dst api.Object) error {
	provider, ok := dst.(*api.Provider)
	if !ok {
		return fmt.Errorf("invalid provider object: %T", dst)
	}

	*provider = api.Provider{
		ObjectMeta:   p.ObjectMeta.ToAPI(),
		Connector:    api.ProviderConnector(p.Connector),
		Filters:      api.ProviderFilters(p.Filters),
		PollInterval: time.Duration(p.PollInterval) * time.Second,
		Enabled:      p.Enabled,
	}

	return nil
}

// ProviderList is a back-end storage list of provider objects.
type ProviderList []Provider

// Copy satisfies the ObjectList interface.
func (p ProviderList) Copy(dst api.ObjectList) error {
	list, ok := dst.(*api.ProviderList)
	if !ok {
		return fmt.Errorf("invalid providers list: %T", dst)
	}

	var err error

	*list = make(api.ProviderList, len(p))

	for idx, provider := range p {
		err = provider.Copy(&(*list)[idx])
		if err != nil {
			return err
		}
	}

	return nil
}

// New satisfies the ObjectList interface.
func (p ProviderList) New() Object {
	return &Provider{}
}

// Objects satisfies the ObjectList interface.
func (p ProviderList) Objects() []Object {
	l := make([]Object, len(p))

	for idx, provider := range p {
		x := provider
		l[idx] = &x
	}

	return l
}
