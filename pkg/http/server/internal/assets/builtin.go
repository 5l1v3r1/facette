// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

//go:generate statik -src=../../../../../dist/assets -p generated

// +build !dev

package assets

import (
	"net/http"
	"strings"

	"batou.dev/httprouter"
	statikfs "github.com/rakyll/statik/fs"

	_ "facette.io/facette/pkg/http/server/internal/assets/generated" // Import static data
)

var fs http.FileSystem

// Register registers assets serving.
func Register(router *httprouter.Router) {
	var err error

	fs, err = statikfs.New()
	if err != nil {
		panic(err)
	}

	router.Endpoint("/*").Get(handle)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	// Fallback to index if requested path doesn't contain an extension
	if r.URL.Path != "/" && !strings.Contains(r.URL.Path, ".") {
		r.URL.Path = "/"
	}

	http.FileServer(fs).ServeHTTP(rw, r)
}
