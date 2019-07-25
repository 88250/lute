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
	Is(NodeType) bool
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
	AppendRawText(string)
	Tokens() items
	AddTokens(items)

	IsOpen() bool
	IsClosed() bool
	Close()
	Finalize()
	Continue(*Context) int
	AcceptLines() bool
	CanContain(NodeType) bool
	LastLineBlank() bool
	SetLastLineBlank(lastLineBlank bool)
	LastLineChecked() bool
	SetLastLineChecked(bool)
}

type BaseNode struct {
	typ             NodeType
	parent          Node
	next            Node
	previous        Node
	firstChild      Node
	lastChild       Node
	rawText         string
	tokens          items
	close           bool
	lastLineBlank   bool
	lastLineChecked bool
}

func (n *BaseNode) Type() NodeType {
	return n.typ
}

func (n *BaseNode) Is(nodeType NodeType) bool {
	return nodeType == n.typ
}

func (n *BaseNode) IsOpen() bool {
	return !n.close
}

func (n *BaseNode) IsClosed() bool {
	return n.close
}

func (n *BaseNode) Close() {
	n.close = true
}

func (n *BaseNode) Finalize() {
}

func (n *BaseNode) Continue(context *Context) int {
	return 0
}

func (n *BaseNode) AcceptLines() bool {
	return false
}

func (n *BaseNode) CanContain(nodeType NodeType) bool {
	return NodeListItem != nodeType
}

func (n *BaseNode) LastLineBlank() bool {
	return n.lastLineBlank
}

func (n *BaseNode) SetLastLineBlank(lastLineBlank bool) {
	n.lastLineBlank = lastLineBlank
}

func (n *BaseNode) LastLineChecked() bool {
	return n.lastLineChecked
}

func (n *BaseNode) SetLastLineChecked(lastLineChecked bool) {
	n.lastLineChecked = lastLineChecked
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

func (n *BaseNode) AppendRawText(rawText string) {
	n.rawText += rawText
}

func (n *BaseNode) Tokens() items {
	return n.tokens
}

func (n *BaseNode) AddTokens(tokens items) {
	n.tokens = append(n.tokens, tokens...)
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

type NodeType int

const (
	NodeRoot NodeType = iota
	NodeBlankLine
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
	NodeInlineHTML
	NodeCodeBlock
	NodeText
	NodeEmphasis
	NodeStrong
	NodeDelete
	NodeInlineCode
	NodeHardBreak
	NodeSoftBreak
	NodeLink
	NodeImage
)

// Nodes.

type Table struct {
	*BaseNode
	Align string
}

type TableRow struct {
	*BaseNode
}

type TableCell struct {
	*BaseNode
	int
	*Tree
}

type InlineHTML struct {
	*BaseNode
	Value string
}

type Text struct {
	*BaseNode
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
}

type InlineCode struct {
	*BaseNode
	Value string
}

type HardBreak struct {
	*BaseNode
}

type SoftBreak struct {
	*BaseNode
}

type Link struct {
	*BaseNode
	Destination string
	Title       string
}

type Image struct {
	*BaseNode
	URL   string
	Title string
}
