// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

// Package generate provides basic internal-use code generator.
package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// Generate generates code based on template text into dstPath parsing files
// from srcPath for the given generation type.
func Generate(dstPath, srcPath, typ, text string) error {
	pkgs, err := parser.ParseDir(token.NewFileSet(), srcPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse dir: %w", err)
	}

	types := []string{}

	for _, pkg := range pkgs {
		var generate bool

		ast.Inspect(pkg, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.Comment:
				t, ok := generateType(x)
				if !generate && ok && t == typ {
					generate = true
				}

			case *ast.TypeSpec:
				if generate {
					types = append(types, x.Name.Name)
				}

			default:
				if x != nil {
					generate = false
				}
			}

			return true
		})
	}

	sort.Strings(types)

	f, err := os.OpenFile( // nolint:gosec
		filepath.Join(dstPath, "zz_generated.go"),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("cannot open destination: %w", err)
	}
	defer f.Close() // nolint:errcheck

	tmpl, err := template.New("types").Funcs(template.FuncMap{"lower": strings.ToLower}).Parse(text)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	err = tmpl.Execute(f, struct{ Types []string }{types})
	if err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}

	return nil
}

func generateType(comment *ast.Comment) (string, bool) {
	text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
	if !strings.HasPrefix(text, "+store:generate=") {
		return "", false
	}

	return strings.TrimPrefix(text, "+store:generate="), true
}
