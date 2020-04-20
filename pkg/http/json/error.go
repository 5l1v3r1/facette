// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"facette.io/facette/pkg/errors"
)

func formatError(err error) error {
	if err == nil {
		return nil
	}

	var xerr *json.UnmarshalTypeError

	ok := errors.As(err, &xerr)
	if ok {
		expected := xerr.Type.Name()
		if expected == "" {
			expected = xerr.Type.Kind().String()
		}

		if xerr.Field != "" {
			return fmt.Errorf("%s: expected a %s", xerr.Field, strings.ToLower(expected))
		}

		return fmt.Errorf("expected a %s", strings.ToLower(expected))
	}

	return err
}
