// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package series

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"facette.io/facette/pkg/errors"
)

// DefaultSample is the time series query default sample value.
const DefaultSample = 400

// Query is a time series query.
type Query struct {
	From  Moment
	To    Moment
	Step  Step
	Exprs []Expr
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (q *Query) UnmarshalJSON(b []byte) error {
	proxy := struct {
		From  *Moment  `json:"from"`
		To    *Moment  `json:"to"`
		Step  *Step    `json:"step"`
		Exprs []string `json:"exprs"`
	}{
		From: &q.From,
		To:   &q.To,
		Step: &q.Step,
	}

	err := json.Unmarshal(b, &proxy)
	if err != nil {
		return err
	}

	if q.To.IsZero() {
		q.To.Time = time.Now().UTC()
	}

	if q.From.IsZero() {
		q.From.Time = q.To.Add(-1 * time.Hour)
	}

	q.From.Time = q.From.Time.Round(1 * time.Second)
	q.To.Time = q.To.Time.Round(1 * time.Second)

	if !q.To.After(q.From.Time) {
		return errors.New(`invalid time range: "to" must come after "from"`)
	}

	if q.Step.Duration == 0 {
		q.Step.Duration = q.To.Sub(q.From.Time) / time.Duration(DefaultSample)
	}

	for idx, text := range proxy.Exprs {
		expr, err := ParseExpr(text)
		if err != nil {
			return fmt.Errorf("invalid #%d series expression: %s", idx, err.Error())
		}

		q.Exprs = append(q.Exprs, expr)
	}

	return nil
}

// Result is a time series query result.
type Result struct {
	From   Moment   `json:"from"`
	To     Moment   `json:"to"`
	Step   Step     `json:"step"`
	Series []Series `json:"series"`
}

// Moment is a time boundary.
//
// A valid moment value can be a RFC3339 string or a relative time reference
// based on the current one (e.g. -1h, -30d or now).
type Moment struct {
	time.Time
}

// MarshalJSON satisfies the json.Marshaler interface.
func (m Moment) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Time.Format(time.RFC3339Nano))
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (m *Moment) UnmarshalJSON(b []byte) error {
	v, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	if v == "" {
		v = "now"
	}

	m.Time, err = ParseTime(v)
	if err != nil {
		return err
	}

	return nil
}

// Step is a time series step duration.
type Step struct {
	time.Duration
}

// MarshalJSON satisfies the json.Marshaler interface.
func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Duration.Seconds())
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (s *Step) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.Duration)
}
