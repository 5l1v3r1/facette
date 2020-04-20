// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Code automatically generated; DO NOT EDIT

package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"facette.io/facette/pkg/api"
)

type ChartOptions api.ChartOptions

func (t *ChartOptions) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t ChartOptions) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type DashboardItems api.DashboardItems

func (t *DashboardItems) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t DashboardItems) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type DashboardOptions api.DashboardOptions

func (t *DashboardOptions) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t DashboardOptions) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type ProviderConnector api.ProviderConnector

func (t *ProviderConnector) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t ProviderConnector) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type ProviderFilters api.ProviderFilters

func (t *ProviderFilters) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t ProviderFilters) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type SeriesList api.SeriesList

func (t *SeriesList) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t SeriesList) Value() (driver.Value, error) {
	return json.Marshal(t)
}
