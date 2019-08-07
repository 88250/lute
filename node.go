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

// Node 描述了节点操作。
type Node interface {
	// Type 返回节点类型。
	Type() int

	// Unlink 用于将节点从树上移除，后一个兄弟节点会接替该节点。
	Unlink()

	// Parent 返回父节点。
	Parent() Node

	// SetParent 设置父节点。
	SetParent(Node)

	// Next 返回后一个兄弟节点。
	Next() Node

	// SetNext 设置后一个兄弟节点。
	SetNext(Node)

	// Previous 返回前一个兄弟节点。
	Previous() Node

	// SetPrevious 设置前一个兄弟节点。
	SetPrevious(Node)

	// FirstChild 返回第一个子节点。
	FirstChild() Node

	// SetFirstChild 设置第一个子节点。
	SetFirstChild(Node)

	// 返回最后一个子节点。
	LastChild() Node

	// 设置最后一个子节点。
	SetLastChild(Node)

	// Children 返回子节点列表。
	Children() []Node

	// AppendChild 添加一个子节点。
	AppendChild(this, child Node)

	// InsertAfter 添加一个兄弟节点。
	InsertAfter(this Node, sibling Node)

	// RawText 返回原始内容。
	RawText() string

	// SetRawText 设置原始内容。
	SetRawText(string)

	// AppendRawText 添加原始内容。
	AppendRawText(string)

	// Value 返回节点值。节点值即根据原始内容处理后的值。
	Value() items

	// SetValue 设置节点值。
	SetValue(items)

	// AppendValue 添加节点值。
	AppendValue(items)

	// Tokens 返回所有 tokens。
	Tokens() items

	// SetTokens 设置 tokens。
	SetTokens(items)

	// AppendTokens 添加 tokens。
	AppendTokens(items)

	// IsOpen 返回节点是否是打开的。
	IsOpen() bool

	// IsClosed 返回节点是否是关闭的。
	IsClosed() bool

	// Close 关闭节点。
	Close()

	// Finalize 节点最终化处理。比如围栏代码块提取 info 部分；HTML 代码块剔除结尾空格；段落需要解析链接引用定义等。
	Finalize(*Context)

	// Continue 用于判断节点是否可以继续处理，比如块引用需要 >，缩进代码块需要 4 空格，围栏代码块需要 ```。
	// 如果可以继续处理返回 0，如果不能接续处理返回 1，如果返回 2（仅在围栏代码块闭合时）则说明可以继续下一行处理了。
	Continue(*Context) int

	// AcceptLines 判断是否节点是否可以接受更多的文本行。比如 HTML 块、代码块和段落是可以接受更多的文本行的。
	AcceptLines() bool

	// CanContain 判断是否能够包含 NodeType 指定类型的节点。 比如列表节点（一种块级容器）只能包含列表项节点，
	// 块引用节点（另一种块级容器）可以包含任意节点；段落节点（一种叶子块节点）不能包含任何其他块级节点。
	CanContain(int) bool

	// LastLineBlank 判断节点最后一行是否是空行。
	LastLineBlank() bool

	// SetLastLineBlank 设置节点最后一行是否是空行。
	SetLastLineBlank(lastLineBlank bool)

	// LastLineChecked 返回最后一行是否检查过。在判断列表是紧凑或松散模式时作为标识用。
	LastLineChecked() bool

	// SetLastLineChecked 设置最后一行是否检查过。
	SetLastLineChecked(bool)
}

// BaseNode 描述了节点基础结构。
type BaseNode struct {
	typ             int    // 节点类型
	parent          Node   // 父节点
	previous        Node   // 前一个兄弟节点
	next            Node   // 后一个兄弟节点
	firstChild      Node   // 第一个子节点
	lastChild       Node   // 最后一个子节点
	rawText         string // 原始内容
	value           items // 原始内容处理后的值
	tokens          items  // 词法分析结果 tokens
	close           bool   // 标识是否关闭
	lastLineBlank   bool   // 标识最后一行是否是空行
	lastLineChecked bool   // 标识最后一行是否检查过
}

func (n *BaseNode) Type() int {
	return n.typ
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

func (n *BaseNode) Finalize(context *Context) {
}

func (n *BaseNode) Continue(context *Context) int {
	return 0
}

func (n *BaseNode) AcceptLines() bool {
	return false
}

func (n *BaseNode) CanContain(nodeType int) bool {
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

func (n *BaseNode) Value() items {
	return n.value
}

func (n *BaseNode) SetValue(value items) {
	n.value = value
}

func (n *BaseNode) AppendValue(value items) {
	n.value = append(n.value, value...)
}

func (n *BaseNode) Tokens() items {
	return n.tokens
}

func (n *BaseNode) SetTokens(tokens items) {
	n.tokens = tokens
}
func (n *BaseNode) AppendTokens(tokens items) {
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
	sibling.SetParent(this.Parent())
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

const (
	NodeRoot          = iota // 根节点类
	NodeParagraph            // 段落节点
	NodeHeading              // 标题节点
	NodeThematicBreak        // 分隔线节点
	NodeBlockquote           // 块引用节点
	NodeList                 // 列表节点
	NodeListItem             // 列表项节点
	NodeHTMLBlock            // HTML 块节点
	NodeInlineHTML           // 内联 HTML节点
	NodeCodeBlock            // 代码块节点
	NodeText                 // 文本节点
	NodeEmphasis             // 强调节点
	NodeStrong               // 加粗节点
	NodeCodeSpan             // 代码节点
	NodeHardBreak            // 硬换行节点
	NodeSoftBreak            // 软换行节点
	NodeLink                 // 链接节点
	NodeImage                // 图片节点
)
