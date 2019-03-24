// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

import (
	"fmt"
)

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node interface {
	Type() NodeType
	Position() Pos
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

type items []item

func (tokens items) Tokens() items {
	return tokens
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
	Pos
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

type Blockquote struct {
	NodeType
	Pos
	RawText
	items
	Subnodes Children
}

func (n *Blockquote) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Blockquote) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<blockquote>\n%s</blockquote>\n", content)
}

func (n *Blockquote) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Blockquote) Children() Children {
	return n.Subnodes
}

type Table struct {
	NodeType
	Pos
	*Tree
	Children

	Align alignType
}

type TableRow struct {
	NodeType
	Pos
	*Tree
	Children
}

type TableCell struct {
	NodeType
	Pos
	*Tree
	Children
}

type HTML struct {
	NodeType
	Pos
	*Tree
	Value string
}

type Code struct {
	NodeType
	Pos
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
	Pos
	*Tree
	Value string
}

type Definition struct {
	NodeType
	Pos
	*Tree

	Association
	Resource
}

type FootnoteDefinition struct {
	NodeType
	Pos
	*Tree
	Children

	Association
}

type Text struct {
	NodeType
	Pos
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
	Pos
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
	Pos
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
	Pos
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
	Pos
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
	Pos
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
	Pos
	*Tree
	Children

	Resource
}

type Image struct {
	NodeType
	Pos
	*Tree

	Resource
	Alternative
}

func (n Image) String() string {
	return fmt.Sprintf("%s", n.URL)
}

type LinkReference struct {
	NodeType
	Pos
	*Tree
	Children

	Reference
}

type ImageReference struct {
	NodeType
	Pos
	*Tree

	Reference
	Alternative
}

type Footnote struct {
	NodeType
	Pos
	*Tree
	Children
}

type FootnoteReference struct {
	NodeType
	Pos
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
