// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// +build ignore

package main

import (
	"path/filepath"
	"runtime"

	"facette.io/facette/pkg/internal/generate"
)

const template = `// Copyright (c) 2020, The Facette Authors
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

{{- range $type := .Types }}

type {{ $type }} api.{{ $type }}

func (t *{{ $type }}) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return json.Unmarshal(x, t)

	case string:
		return json.Unmarshal([]byte(x), t)

	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
}

func (t {{ $type }}) Value() (driver.Value, error) {
	return json.Marshal(t)
}
{{- end }}
`

func main() {
	_, path, _, _ := runtime.Caller(0)
	dst := filepath.Dir(path)

	err := generate.Generate(dst, filepath.Join(dst, "/../../../api"), "type", template)
	if err != nil {
		panic(err)
	}
}
