// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package template

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"facette.io/facette/pkg/set"
)

// Template is a representation of a parsed template.
type Template struct {
	r         *bufio.Reader
	pos, last position
	text      string
	nodes     []Node
}

// New creates a new template instance.
func New() *Template {
	return &Template{}
}

// Nodes returns the list of parsed template nodes.
func (t *Template) Nodes() []Node {
	return t.nodes
}

// Parse parses text as template data.
func (t *Template) Parse(text string) error {
	t.r = bufio.NewReader(strings.NewReader(text))
	t.pos = position{Line: 1, Char: 1}
	t.text = text
	t.nodes = []Node{}

	for {
		next, err := t.peek(2)
		if err == io.EOF {
			// Append remainder as text node
			if len(next) > 0 {
				t.nodes = append(t.nodes, &TextNode{string(next)})
			}

			break
		} else if err != nil {
			return err
		}

		var node Node

		if next[0] == '$' && (isVariableChar(next[1]) || next[1] == '{') {
			t.discard()

			node, err = t.readVariable()
			if err != nil {
				return err
			}

			t.nodes = append(t.nodes, node)

			continue
		}

		node, err = t.readText()
		if err != nil {
			return err
		}

		t.nodes = append(t.nodes, node)
	}

	return nil
}

// Variables returns the list of parsed variable names.
func (t *Template) Variables() []string {
	keys := set.New()

	for _, node := range t.nodes {
		v, ok := node.(*VariableNode)
		if ok {
			keys.Add(v.Name)
		}
	}

	result := set.StringSlice(keys)
	sort.Strings(result)

	return result
}

func (t *Template) discard() {
	_, _ = t.read()
}

func (t *Template) peek(n int) ([]byte, error) {
	return t.r.Peek(n)
}

func (t *Template) read() (byte, error) {
	b, err := t.r.ReadByte()
	if err == nil {
		t.last = t.pos

		if b == '\n' {
			t.pos.Line++
			t.pos.Char = 1
		} else {
			t.pos.Char++
		}
	}

	return b, err
}

func (t *Template) readText() (*TextNode, error) {
	var prev byte

	buf := bytes.NewBuffer(nil)

	for {
		c, err := t.read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if c == '$' && prev != '\\' {
			t.unread()

			next, err := t.peek(2)
			if err == nil && (isVariableChar(next[1]) || next[1] == '{') {
				break
			}

			t.discard()
		}

		buf.WriteByte(c)
		prev = c
	}

	return &TextNode{buf.String()}, nil
}

func (t *Template) readVariable() (*VariableNode, error) {
	var (
		lbrace, rbrace bool
		bracePos       position
	)

	next, err := t.peek(1)
	if err == nil && next[0] == '{' {
		lbrace, bracePos = true, t.pos
		t.discard()
	}

	buf := bytes.NewBuffer(nil)

loop:
	for {
		c, err := t.read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch {
		case lbrace && c == '}':
			rbrace = true
			break loop

		case isVariableChar(c):
			buf.WriteByte(c)

		default:
			t.unread()
			break loop
		}
	}

	if lbrace && !rbrace {
		return nil, fmt.Errorf(`unbalanced brace at %s`, bracePos)
	}

	return &VariableNode{buf.String()}, nil
}

func (t *Template) unread() {
	err := t.r.UnreadByte()
	if err == nil && t.last.Line != 0 {
		t.pos = t.last
	}
}

type position struct {
	Line int
	Char int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Char)
}

func isVariableChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}
