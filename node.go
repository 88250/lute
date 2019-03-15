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

import (
	"fmt"
)

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node interface {
	Type() NodeType
	Position() Pos
	String() string
	HTML() string
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Children represents the children nodes of a tree node.
type Children []Node

func (t NodeType) Type() NodeType {
	return t
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
	*Tree
	Children
}

func (n *Root) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Root) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("%s\n", content)
}

func (n *Root) append(c Node) {
	n.Children = append(n.Children, c)
}

type Paragraph struct {
	NodeType
	Pos
	*Tree
	Children

	OpenTag, CloseTag string
}

func (n *Paragraph) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Paragraph) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag, content)
}

func (n *Paragraph) append(c Node) {
	n.Children = append(n.Children, c)
}

func (n *Paragraph) trim() {
	size := len(n.Children)
	if 1 > size {
		return
	}

	initialNoneWhitespace := 0
	for i := initialNoneWhitespace; i < size/2; i++ {
		if NodeBreak == n.Children[i].Type() {
			initialNoneWhitespace++
		}
	}

	finalNoneWhitespace := size
	for i := finalNoneWhitespace - 1; size/2 <= i; i-- {
		if NodeBreak == n.Children[i].Type() {
			finalNoneWhitespace--
		}
	}

	n.Children = n.Children[initialNoneWhitespace:finalNoneWhitespace]
}

type Heading struct {
	NodeType
	Pos
	*Tree
	Children

	Depth int
}

func (n Heading) String() string {
	return fmt.Sprintf("# %s", n.Children)
}

func (n *Heading) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<h%d>%s</h%d>", n.Depth, content, n.Depth)
}

type ThematicBreak struct {
	NodeType
	Pos
}

func (n *ThematicBreak) String() string {
	return fmt.Sprintf("'***'")
}

func (n *ThematicBreak) HTML() string {
	return fmt.Sprintf("<hr>")
}

type Blockquote struct {
	NodeType
	Pos
	Children
}

func (n *Blockquote) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Blockquote) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<blockquote>%s</blockquote>", content)
}

type List struct {
	NodeType
	Pos
	*Tree
	Children

	Ordered bool
	Start   int
	Spread  bool

	Marker string
	Indent int
}

func (n *List) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *List) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<ul>\n%s</ul>", content)
}

func (n *List) append(c Node) {
	n.Children = append(n.Children, c)
}

type ListItem struct {
	NodeType
	Pos
	*Tree
	Children

	Checked bool
	Spread  bool // loose or tight

	Spaces int
}

func (n *ListItem) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *ListItem) HTML() string {
	var content string
	for _, c := range n.Children {
		if !n.Spread && NodeParagraph == c.Type() {
			p := c.(*Paragraph)
			p.OpenTag, p.CloseTag = "", ""
		}

		content += c.HTML()
	}

	return fmt.Sprintf("<li>\n%s</li>\n", content)
}

func (n *ListItem) append(c Node) {
	n.Children = append(n.Children, c)
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
	*Tree
	Value string

	Lang string
	Meta string
}

func (n *Code) String() string {
	return fmt.Sprintf("```%s```", n.Value)
}

func (n *Code) HTML() string {
	return fmt.Sprintf("<pre><code>%s</code></pre>", n.Value)
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
	*Tree
	Value string
}

func (n *Text) String() string {
	return fmt.Sprintf("'%s'", n.Value)
}

func (n *Text) HTML() string {
	return fmt.Sprintf("%s", n.Value)
}

type Emphasis struct {
	NodeType
	Pos
	*Tree
	Children
}

func (n *Emphasis) String() string {
	return fmt.Sprintf("*%v*", n.Children)
}

func (n *Emphasis) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<em>%s</em>", content)
}

type Strong struct {
	NodeType
	Pos
	*Tree
	Children
}

func (n *Strong) String() string {
	return fmt.Sprintf("**%v**", n.Children)
}

func (n *Strong) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<strong>%s</strong>", content)
}

type Delete struct {
	NodeType
	Pos
	*Tree
	Children
}

func (n *Delete) String() string {
	return fmt.Sprintf("~~%v~~", n.Children)
}

func (n *Delete) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<del>%s</del>", content)
}

type InlineCode struct {
	NodeType
	Pos
	*Tree
	Value string
}

func (n InlineCode) String() string {
	return fmt.Sprintf("`%s`", n.Value)
}

func (n InlineCode) HTML() string {
	return fmt.Sprintf("<code>%s</code>", n.Value)
}

type Break struct {
	NodeType
	Pos
	*Tree
}

func (n *Break) String() string {
	return fmt.Sprint("\n")
}

func (n *Break) HTML() string {
	return fmt.Sprintf("\n")
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
