// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import "fmt"

// Func is a series expression function.
type Func string

// Series expression functions:
const (
	FuncAlias  Func = "alias"
	FuncAvg    Func = "avg"
	FuncSample Func = "sample"
	FuncScale  Func = "scale"
	FuncSum    Func = "sum"
)

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (f *Func) UnmarshalText(b []byte) error {
	switch v := Func(b); v {
	case FuncAlias, FuncAvg, FuncSample, FuncScale, FuncSum:
		*f = v
		return nil
	}

	return fmt.Errorf("unsupported function: %s", b)
}
