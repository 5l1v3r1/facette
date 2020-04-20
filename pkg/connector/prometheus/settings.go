// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package prometheus

import (
	"encoding/json"
	"net/url"

	"facette.io/facette/pkg/errors"
)

// Settings are Prometheus time series connector settings.
type Settings struct {
	URL        string `json:"url"`
	Filter     string `json:"filter"`
	SkipVerify bool   `json:"skipVerify"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (s *Settings) UnmarshalJSON(b []byte) error {
	type Proxy Settings

	var proxy *Proxy = (*Proxy)(s)

	err := json.Unmarshal(b, proxy)
	if err != nil {
		return err
	}

	if s.URL == "" {
		return errors.New("missing connector setting: url")
	}

	url, err := url.Parse(s.URL)
	if err != nil || url.Scheme == "" || url.Host == "" {
		return errors.New("invalid connector setting: url")
	}

	return nil
}
