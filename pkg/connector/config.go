// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package connector

import "encoding/json"

// Config is an upstream time series connector configuration.
type Config struct {
	Type     Type            `json:"type"`
	Settings json.RawMessage `json:"settings"`
}
