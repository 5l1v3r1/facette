// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package version provides the version information.
package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"facette.io/facette/pkg/connector"
	"facette.io/facette/pkg/store/driver"
)

// Version information:
var (
	Version   string
	Branch    string
	Revision  string
	BuildDate string
	Compiler  = fmt.Sprintf("%s (%s)", runtime.Version(), runtime.Compiler)
)

// Print writes the version information to the standard output.
func Print() {
	fmt.Printf(`%s
   Version:     %s
   Branch:      %s
   Revision:    %s
   Compiler:    %s
   Build date:  %s
   Drivers:     %s
   Connectors:  %s
`,
		filepath.Base(os.Args[0]),
		Version,
		Branch,
		Revision,
		Compiler,
		BuildDate,
		strings.Join(driver.Drivers(), ", "),
		strings.Join(connector.Connectors(), ", "),
	)
}
