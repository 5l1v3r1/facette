// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package store

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type logger struct{}

func (l logger) Print(v ...interface{}) {
	if v[0] == "sql" {
		args := []string{}

		for _, arg := range v[4].([]interface{}) {
			format := "%v"

			switch arg.(type) {
			case string, fmt.Stringer:
				format = "%q"
			}

			args = append(args, fmt.Sprintf(format, arg))
		}

		zap.S().Named("store/sql").Debugf("%s [%s] in %s\n", v[3], strings.Join(args, ", "), v[2])
	} else {
		zap.S().Named("store/sql").Debugf("%s", v[2])
	}
}
