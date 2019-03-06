// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
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

import "fmt"

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node interface {
	Type() NodeType
	Position() Pos
	String() string
}

// NodeType identifies the type of a parse tree node.
type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeParent NodeType = iota
	NodeRoot
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

type Parent struct {
	NodeType
	Pos
	Children []Node // element nodes in lexical order
}

type Literal struct {
	NodeType
	Pos
	Value string
}

func (n *Literal) String() string {
	return fmt.Sprintf("%s", n.Value)
}

type Root struct {
	Parent
}

func (n *Root) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Root) append(c Node) {
	n.Children = append(n.Children, c)
}

type Paragraph struct {
	Parent
	Children []Node
}

func (n *Paragraph) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Paragraph) append(c Node) {
	n.Children = append(n.Children, c)
}

type Heading struct {
	Parent
	Depth    int8
	Children []Node
}

type ThematicBreak struct {
	NodeType
	Pos
}

type Blockquote struct {
	Parent
	Children []Node
}

type List struct {
	Parent
	Ordered  bool
	Start    int8
	Spread   bool
	Children []Node
}

type ListItem struct {
	Parent
	Checked  bool
	Spread   bool
	Children []Node
}

type Table struct {
	Parent
	Align    alignType
	Children []Node
}

type TableRow struct {
	Parent
	Children []Node
}

type TableCell struct {
	Parent
	Children []Node
}

type HTML struct {
	Literal
}

type Code struct {
	Literal
	Lang string
	Meta string
}

type YAML struct {
	Literal
}

type Definition struct {
	NodeType
	Pos
	Association
	Resource
}

type FootnoteDefinition struct {
	Parent
	Association
	Children []Node
}

type Text struct {
	Literal
}

func (n Text) String() string {
	return fmt.Sprintf("'%s'", n.Value)
}

type Emphasis struct {
	Parent
	Children []Node
}

func (n Emphasis) String() string {
	return fmt.Sprintf("*%v*", n.Children)
}

type Strong struct {
	Parent
	Children []Node
}

func (n Strong) String() string {
	return fmt.Sprintf("**%v**", n.Children)
}

type Delete struct {
	Parent
	Children []Node
}

func (n Delete) String() string {
	return fmt.Sprintf("~~%v~~", n.Children)
}

type InlineCode struct {
	Literal
}

func (n InlineCode) String() string {
	return fmt.Sprintf("`%s`", n.Value)
}

type Break struct {
	NodeType
	Pos
}

type Link struct {
	Parent
	Resource
	Children []Node
}

type Image struct {
	NodeType
	Pos
	Resource
	Alternative
}

type LinkReference struct {
	Parent
	Reference
	Children []Node
}

type ImageReference struct {
	NodeType
	Pos
	Reference
	Alternative
}

type Footnote struct {
	Parent
	Children []Node
}

type FootnoteReference struct {
	NodeType
	Pos
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
