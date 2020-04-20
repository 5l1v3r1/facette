// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"

	"batou.dev/httprouter"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/errors"
	"facette.io/facette/pkg/labels"
	"facette.io/facette/pkg/store"
)

func listOptionsFromRequest(r *http.Request) (*api.ListOptions, error) {
	offset, err := uintFromQuery(r, "offset")
	if err != nil {
		return nil, errors.Wrap(api.ErrInvalid, "invalid offset")
	}

	limit, err := uintFromQuery(r, "limit")
	if err != nil {
		return nil, errors.Wrap(api.ErrInvalid, "invalid limit")
	}

	opts := &api.ListOptions{
		Filter: make(map[string]interface{}),
		Sort:   httprouter.QueryParam(r, "sort"),
		Offset: offset,
		Limit:  limit,
	}

	name := httprouter.QueryParam(r, "name")
	if strings.HasPrefix(name, "~") {
		pattern := strings.TrimPrefix(name, "~")

		err = compileRegexpPattern(pattern)
		if err != nil {
			return nil, errors.Wrapf(api.ErrInvalid, "invalid name pattern: %s", err)
		}

		opts.Filter["name"] = store.RegexpPattern(pattern)
	} else if name != "" {
		opts.Filter["name"] = name
	}

	switch filepath.Base(r.URL.Path) {
	case "charts", "dashboards":
		switch httprouter.QueryParam(r, "kind") {
		case "plain":
			opts.Filter["template"] = false

		case "template":
			opts.Filter["template"] = true
		}

	case "providers":
		enabled, ok, err := boolFromQuery(r, "enabled")
		if err != nil {
			return nil, errors.Wrapf(api.ErrInvalid, "invalid enabled flag")
		}

		if ok {
			opts.Filter["enabled"] = enabled
		}
	}

	return opts, nil
}

func matcherFromRequest(r *http.Request) (labels.Matcher, error) {
	match := httprouter.QueryParam(r, "match")
	if match != "" {
		return labels.ParseMatcher(match)
	}

	return nil, nil
}

func boolFromQuery(r *http.Request, key string) (bool, bool, error) {
	value := httprouter.QueryParam(r, key)
	if value != "" {
		v, err := strconv.ParseBool(value)
		return v, true, err
	}

	return false, false, nil
}

func uintFromQuery(r *http.Request, key string) (uint, error) {
	value := httprouter.QueryParam(r, key)
	if value != "" {
		v, err := strconv.ParseUint(value, 10, 32)
		return uint(v), err
	}

	return 0, nil
}

func compileRegexpPattern(pattern string) error {
	_, err := regexp.Compile(pattern)
	if err != nil {
		var xerr *syntax.Error

		ok := errors.As(err, &xerr)
		if ok {
			return fmt.Errorf("%s: `%s`", xerr.Code, xerr.Expr)
		}

		return err
	}

	return nil
}
