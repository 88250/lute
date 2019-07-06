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

type Node interface {
	Type() NodeType
	Unlink()
	Parent() Node
	SetParent(Node)
	Next() Node
	SetNext(Node)
	Previous() Node
	SetPrevious(Node)
	FirstChild() Node
	SetFirstChild(Node)
	LastChild() Node
	SetLastChild(Node)
	Children() []Node
	AppendChild(this, child Node)
	InsertAfter(this Node, sibling Node)
	RawText() string
	SetRawText(string)
	Tokens() items
}

type BaseNode struct {
	typ NodeType
	parent     Node
	next       Node
	previous   Node
	firstChild Node
	lastChild  Node
	rawText    string
	tokens     items
}

func (n*BaseNode) Type() NodeType {
	return n.typ
}

func (n *BaseNode) Unlink() {
	if nil != n.previous {
		n.previous.SetNext(n.next)
	} else if nil != n.parent {
		n.parent.SetFirstChild(n.next)
	}
	if nil != n.next {
		n.next.SetPrevious(n.previous)
	} else if nil != n.parent {
		n.parent.SetLastChild(n.previous)
	}
	n.parent = nil
	n.next = nil
	n.previous = nil
}

func (n *BaseNode) Parent() Node {
	return n.parent
}

func (n *BaseNode) SetParent(parent Node) {
	n.parent = parent
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) SetNext(next Node) {
	n.next = next
}

func (n *BaseNode) Previous() Node {
	return n.previous
}

func (n *BaseNode) SetPrevious(previous Node) {
	n.previous = previous
}

func (n *BaseNode) FirstChild() Node {
	return n.firstChild
}

func (n *BaseNode) SetFirstChild(firstChild Node) {
	n.firstChild = firstChild
}

func (n *BaseNode) LastChild() Node {
	return n.lastChild
}

func (n *BaseNode) SetLastChild(lastChild Node) {
	n.lastChild = lastChild
}

func (n *BaseNode) Children() (ret []Node) {
	for child := n.firstChild; nil != child; child = child.Next() {
		ret = append(ret, child)
	}

	return
}

func (n *BaseNode) RawText() string {
	return n.rawText
}

func (n *BaseNode) SetRawText(rawText string) {
	n.rawText = rawText
}

func (n *BaseNode) Tokens() items {
	return n.tokens
}

func (n *BaseNode) InsertAfter(this Node, sibling Node) {
	sibling.Unlink()
	sibling.SetNext(n.next)
	if nil != sibling.Next() {
		sibling.Next().SetPrevious(sibling)
	}
	sibling.SetPrevious(this)
	n.next = sibling
	sibling.SetParent(this)
	if nil == sibling.Next() {
		sibling.Parent().SetLastChild(sibling)
	}
}

func (n *BaseNode) InsertBefore(this Node, sibling Node) {
	sibling.Unlink()
	sibling.SetPrevious(n.previous)
	if nil != sibling.Previous() {
		sibling.Previous().SetNext(sibling)
	}
	sibling.SetNext(this)
	n.previous = sibling
	sibling.SetParent(n.parent)
	if nil == sibling.Previous() {
		sibling.Parent().SetFirstChild(sibling)
	}
}

func (n *BaseNode) AppendChild(this, child Node) {
	child.Unlink()
	child.SetParent(this)
	if nil != n.lastChild {
		n.lastChild.SetNext(child)
		child.SetPrevious(n.lastChild)
		n.lastChild = child
	} else {
		n.firstChild = child
		n.lastChild = child
	}
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

type NodeType int

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
	*BaseNode
}

type Table struct {
	*BaseNode
	int
	*Tree

	Align alignType
}

type TableRow struct {
	*BaseNode
	int
	*Tree
}

type TableCell struct {
	*BaseNode
	int
	*Tree
}

type HTML struct {
	*BaseNode
	int
	*Tree
	Value string
}

type Code struct {
	*BaseNode
	int
	*Tree
	Value string

	Lang string
	Meta string
}

type Definition struct {
	*BaseNode
	int
	*Tree

	Association
	Resource
}

type FootnoteDefinition struct {
	*BaseNode
	int
	*Tree

	Association
}

type Text struct {
	*BaseNode
	*Tree
	Value string
}

type Emphasis struct {
	*BaseNode
}

type Strong struct {
	*BaseNode
}

type Delete struct {
	*BaseNode
	int
	items
	*Tree
}

type InlineCode struct {
	*BaseNode
	*Tree
	Value string
}

type Break struct {
	*BaseNode
	int
	items
	*Tree
}

type Link struct {
	*BaseNode
	int
	*Tree

	Resource
}

type Image struct {
	*BaseNode
	int
	*Tree

	Resource
	Alternative
}

type LinkReference struct {
	*BaseNode
	int
	*Tree

	Reference
}

type ImageReference struct {
	*BaseNode
	int
	*Tree

	Reference
	Alternative
}

type Footnote struct {
	*BaseNode
	int
	*Tree
}

type FootnoteReference struct {
	*BaseNode
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
