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

type HTMLRenderer interface {
	HTML() string
	String() string
}

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node struct {
	NodeType   int
	Parent     *Node
	Next       *Node
	Previous   *Node
	FirstChild *Node
	LastChild  *Node
	RawText    string
	Tokens     items

	HTMLRenderer
}

func (n *Node) Unlink() {
	if nil != n.Previous {
		n.Previous.Next = n.Next
	} else if nil != n.Parent {
		n.Parent.FirstChild = n.Next
	}
	if nil != n.Next {
		n.Next.Previous = n.Previous
	} else if nil != n.Parent {
		n.Parent.LastChild = n.Previous
	}
	n.Parent = nil
	n.Next = nil
	n.Previous = nil
}

func (n *Node) InsertAfter(sibling *Node) {
	sibling.Unlink()
	sibling.Next = n.Next
	if nil != sibling.Next {
		sibling.Next.Previous = sibling
	}
	sibling.Previous = n
	n.Next = sibling
	sibling.Parent = n
	if nil == sibling.Next {
		sibling.Parent.LastChild = sibling
	}
}

func (n *Node) InsertBefore(sibling *Node) {
	sibling.Unlink()
	sibling.Previous = n.Previous
	if nil != sibling.Previous {
		sibling.Previous.Next = sibling
	}
	sibling.Next = n
	n.Previous = sibling
	sibling.Parent = n.Parent
	if nil == sibling.Previous {
		sibling.Parent.FirstChild = sibling
	}
}

func (n *Node) Append(child *Node) {
	child.Unlink()
	child.Parent = n
	if nil != n.LastChild {
		n.LastChild.Next = child
		child.Previous = n.LastChild
		n.LastChild = child
	} else {
		n.FirstChild = child
		n.LastChild = child
	}
}

func (v *Node) Children() (ret []*Node) {
	for child := v.FirstChild; nil != child; child = child.Next {
		ret = append(ret, child)
	}

	return
}

type items []*item

func (tokens items) Tokens() items {
	return tokens
}

func (tokens items) isEOF() bool {
	return 1 == len(tokens) && (tokens)[0].isEOF()
}

func (tokens items) rawText() (ret string) {
	for i := 0; i < len(tokens); i++ {
		ret += (tokens)[i].val
	}

	return
}

const (
	NodeRoot = iota
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
	*Node
	Pos int
	items
	*Tree
}

func (n *Root) String() string {
	return fmt.Sprintf("%s", n.Children())
}

func (n *Root) HTML() string {
	content := html(n.Children())

	return fmt.Sprintf("%s", content)
}

type Table struct {
	*Node
	int
	*Tree

	Align alignType
}

type TableRow struct {
	*Node
	int
	*Tree
}

type TableCell struct {
	*Node
	int
	*Tree
}

type HTML struct {
	*Node
	int
	*Tree
	Value string
}

type Code struct {
	*Node
	int
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

type Definition struct {
	*Node
	int
	*Tree

	Association
	Resource
}

type FootnoteDefinition struct {
	*Node
	int
	*Tree

	Association
}

type Text struct {
	*Node
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

type Emphasis struct {
	*Node
	items
	*Tree
}

func (n *Emphasis) String() string {
	return fmt.Sprintf("*%v*", n.Children())
}

func (n *Emphasis) HTML() string {
	content := html(n.Children())

	return fmt.Sprintf("<em>%s</em>", content)
}

type Strong struct {
	*Node
	items
	*Tree
}

func (n *Strong) String() string {
	return fmt.Sprintf("**%v**", n.Children())
}

func (n *Strong) HTML() string {
	content := html(n.Children())

	return fmt.Sprintf("<strong>%s</strong>", content)
}

type Delete struct {
	*Node
	int
	items
	*Tree
}

func (n *Delete) String() string {
	return fmt.Sprintf("~~%v~~", n.Children())
}

func (n *Delete) HTML() string {
	content := html(n.Children())

	return fmt.Sprintf("<del>%s</del>", content)
}

type InlineCode struct {
	*Node
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

type Break struct {
	*Node
	int
	items
	*Tree
}

func (n *Break) String() string {
	return fmt.Sprint("\n")
}

func (n *Break) HTML() string {
	return fmt.Sprintf("\n")
}

type Link struct {
	*Node
	int
	*Tree

	Resource
}

type Image struct {
	*Node
	int
	*Tree

	Resource
	Alternative
}

func (n Image) String() string {
	return fmt.Sprintf("%s", n.URL)
}

type LinkReference struct {
	*Node
	int
	*Tree

	Reference
}

type ImageReference struct {
	*Node
	int
	*Tree

	Reference
	Alternative
}

type Footnote struct {
	*Node
	int
	*Tree
}

type FootnoteReference struct {
	*Node
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

func html(children []*Node) string {
	var ret string
	for _, c := range children {
		ret += c.HTML()
	}

	return ret
}
