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
	offset, err := int64FromQuery(r, "offset")
	if err != nil {
		return nil, errors.Wrap(api.ErrInvalid, "invalid offset")
	}

	limit, err := int64FromQuery(r, "limit")
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

	section := filepath.Base(r.URL.Path)

	switch section {
	case "charts", "dashboards":
		switch httprouter.QueryParam(r, "kind") {
		case "plain":
			opts.Filter["template"] = false

		case "template":
			opts.Filter["template"] = true
		}

		if section == "dashboards" {
			parent := httprouter.QueryParam(r, "parent")
			if parent != "" {
				opts.Filter["parent_id"] = parent
			}
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

func int64FromQuery(r *http.Request, key string) (int64, error) {
	value := httprouter.QueryParam(r, key)
	if value != "" {
		return strconv.ParseInt(value, 10, 64)
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
