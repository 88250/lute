// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"fmt"
)

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node interface {
	Type() NodeType
	Raw() RawText
	Children() Children
	Tokens() items
	Append(child Node)
	String() string
	HTML() string
}

// NodeType identifies the type of a parse tree node.
type NodeType int

func (nt NodeType) Type() NodeType {
	return nt
}

// Children represents the children nodes of a tree node.
type Children []Node

type RawText string

func (r RawText) Raw() RawText {
	return r
}

type items []*item

func (tokens items) Tokens() items {
	return tokens
}

func (tokens items) isEOF() bool {
	return 1 == len(tokens) && (tokens)[0].isEOF()
}

func (tokens items) rawText() (ret RawText) {
	for i := 0; i < len(tokens); i++ {
		ret += RawText((tokens)[i].val)
	}

	return
}

const (
	NodeRoot NodeType = iota
	NodeParagraph
	NodeHeading
	NodeThematicBreak
	NodeBlockquote
	NodeList
	NodeListItem
	NodeTable
	NodeTableRow
	NodeTableCell
	NodeHTML
	NodeCode
	NodeYAML
	NodeDefinition
	NodeFootnoteDefinition
	NodeText
	NodeEmphasis
	NodeStrong
	NodeDelete
	NodeInlineCode
	NodeBreak
	NodeLink
	NodeImage
	NodeLinkReference
	NodeImageReference
	NodeFootnote
	NodeFootnoteReference
)

// Nodes.

type Root struct {
	NodeType
	Pos int
	RawText
	items
	*Tree
	Subnodes Children
}

func (n *Root) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Root) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("%s", content)
}

func (n *Root) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Root) Children() Children {
	return n.Subnodes
}

type Table struct {
	NodeType
	int
	*Tree
	Children

	Align alignType
}

type TableRow struct {
	NodeType
	int
	*Tree
	Children
}

type TableCell struct {
	NodeType
	int
	*Tree
	Children
}

type HTML struct {
	NodeType
	int
	*Tree
	Value string
}

type Code struct {
	NodeType
	int
	RawText
	items
	*Tree
	Value string

	Lang string
	Meta string
}

func (n *Code) String() string {
	return fmt.Sprintf("```%s```", n.Value)
}

func (n *Code) HTML() string {
	return fmt.Sprintf("<pre><code>%s</code></pre>\n", n.Value)
}

func (n *Code) Append(c Node) {
}

func (n *Code) Children() Children {
	return nil
}

type YAML struct {
	NodeType
	int
	*Tree
	Value string
}

type Definition struct {
	NodeType
	int
	*Tree

	Association
	Resource
}

type FootnoteDefinition struct {
	NodeType
	int
	*Tree
	Children

	Association
}

type Text struct {
	NodeType
	RawText
	items
	*Tree
	Value string
}

func (n *Text) String() string {
	return fmt.Sprintf("'%s'", n.Value)
}

func (n *Text) HTML() string {
	return fmt.Sprintf("%s", n.Value)
}

func (n *Text) Append(child Node) {
}

func (n *Text) Children() Children {
	return nil
}

type Emphasis struct {
	NodeType
	RawText
	items
	*Tree
	Subnodes Children
}

func (n *Emphasis) String() string {
	return fmt.Sprintf("*%v*", n.Subnodes)
}

func (n *Emphasis) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<em>%s</em>", content)
}

func (n *Emphasis) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Emphasis) Children() Children {
	return n.Subnodes
}

type Strong struct {
	NodeType
	RawText
	items
	*Tree
	Subnodes Children
}

func (n *Strong) String() string {
	return fmt.Sprintf("**%v**", n.Subnodes)
}

func (n *Strong) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<strong>%s</strong>", content)
}

func (n *Strong) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Strong) Children() Children {
	return n.Subnodes
}

type Delete struct {
	NodeType
	int
	RawText
	items
	*Tree
	Subnodes Children
}

func (n *Delete) String() string {
	return fmt.Sprintf("~~%v~~", n.Subnodes)
}

func (n *Delete) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<del>%s</del>", content)
}

func (n *Delete) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Delete) Children() Children {
	return n.Subnodes
}

type InlineCode struct {
	NodeType
	RawText
	items
	*Tree
	Value string
}

func (n *InlineCode) String() string {
	return fmt.Sprintf("`%s`", n.Value)
}

func (n *InlineCode) HTML() string {
	return fmt.Sprintf("<code>%s</code>", n.Value)
}

func (n *InlineCode) Append(c Node) {}

func (n *InlineCode) Children() Children {
	return nil
}

type Break struct {
	NodeType
	int
	RawText
	items
	*Tree
}

func (n *Break) String() string {
	return fmt.Sprint("\n")
}

func (n *Break) HTML() string {
	return fmt.Sprintf("\n")
}

func (n *Break) Append(c Node) {}

func (n *Break) Children() Children {
	return nil
}

type Link struct {
	NodeType
	int
	*Tree
	Children

	Resource
}

type Image struct {
	NodeType
	int
	*Tree

	Resource
	Alternative
}

func (n Image) String() string {
	return fmt.Sprintf("%s", n.URL)
}

type LinkReference struct {
	NodeType
	int
	*Tree
	Children

	Reference
}

type ImageReference struct {
	NodeType
	int
	*Tree

	Reference
	Alternative
}

type Footnote struct {
	NodeType
	int
	*Tree
	Children
}

type FootnoteReference struct {
	NodeType
	int
	*Tree

	Association
}

// Mixins.

type Resource struct {
	URL   string
	Title string
}

type Association struct {
	Identifier string
	Label      string
}

type Reference struct {
	ReferenceType referenceType
	Association
}

type Alternative struct {
	Alt string
}

// Enumerations.

type alignType string
type referenceType string

func html(children Children) string {
	var ret string
	for _, c := range children {
		ret += c.HTML()
	}

	return ret
}
