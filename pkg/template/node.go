// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package template

// Node is a template node.
type Node interface {
	node()
}

// TextNode is a template text node.
type TextNode struct {
	Text string
}

func (TextNode) node() {}

// VariableNode is a template variable node.
type VariableNode struct {
	Name string
}

func (VariableNode) node() {}
