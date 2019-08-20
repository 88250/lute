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

// Text 描述了文本节点结构。
type Text struct {
	// *BaseNode
	// 这里通过减少嵌套结构体的构造来优化性能。
	// 整颗语法书上可能 95% 都是文本节点。此时如果文本节点嵌套了基础节点 BaseNode 的话（嵌套的目的是为了复用 BaseNode 的字段和接口实现），
	// 构造文本节点就需要创建两次对象，将会重复调用 `runtime.newobject` 降低性能。
	// 优化方案就是去掉嵌套的 BaseNode，将 BaseNode 的结构在 Text 中再做一次。
	// 总的来说，不构造对象可以换来巨大的性能提升，但代价就是降低代码可读性，并且看上去会显得有些僵硬。

	parent          Node   // 父节点
	previous        Node   // 前一个兄弟节点
	next            Node   // 后一个兄弟节点
	firstChild      Node   // 第一个子节点
	lastChild       Node   // 最后一个子节点
	rawText         string // 原始内容
	tokens          items  // 词法分析结果 tokens，语法分析阶段会继续操作这些 tokens
	close           bool   // 标识是否关闭
	lastLineBlank   bool   // 标识最后一行是否是空行
	lastLineChecked bool   // 标识最后一行是否检查过
}

func (n *Text) Type() int {
	return NodeText
}

func (n *Text) IsOpen() bool {
	return !n.close
}

func (n *Text) IsClosed() bool {
	return n.close
}

func (n *Text) Close() {
	n.close = true
}

func (n *Text) Finalize(context *Context) {
}

func (n *Text) Continue(context *Context) int {
	return 0
}

func (n *Text) AcceptLines() bool {
	return false
}

func (n *Text) CanContain(nodeType int) bool {
	return NodeListItem != nodeType
}

func (n *Text) LastLineBlank() bool {
	return n.lastLineBlank
}

func (n *Text) SetLastLineBlank(lastLineBlank bool) {
	n.lastLineBlank = lastLineBlank
}

func (n *Text) LastLineChecked() bool {
	return n.lastLineChecked
}

func (n *Text) SetLastLineChecked(lastLineChecked bool) {
	n.lastLineChecked = lastLineChecked
}

func (n *Text) Unlink() {
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

func (n *Text) Parent() Node {
	return n.parent
}

func (n *Text) SetParent(parent Node) {
	n.parent = parent
}

func (n *Text) Next() Node {
	return n.next
}

func (n *Text) SetNext(next Node) {
	n.next = next
}

func (n *Text) Previous() Node {
	return n.previous
}

func (n *Text) SetPrevious(previous Node) {
	n.previous = previous
}

func (n *Text) FirstChild() Node {
	return n.firstChild
}

func (n *Text) SetFirstChild(firstChild Node) {
	n.firstChild = firstChild
}

func (n *Text) LastChild() Node {
	return n.lastChild
}

func (n *Text) SetLastChild(lastChild Node) {
	n.lastChild = lastChild
}

func (n *Text) RawText() string {
	return n.rawText
}

func (n *Text) SetRawText(rawText string) {
	n.rawText = rawText
}

func (n *Text) AppendRawText(rawText string) {
	n.rawText += rawText
}

func (n *Text) Tokens() items {
	return n.tokens
}

func (n *Text) SetTokens(tokens items) {
	n.tokens = tokens
}
func (n *Text) AppendTokens(tokens items) {
	n.tokens = append(n.tokens, tokens...)
}

func (n *Text) InsertAfter(this Node, sibling Node) {
	sibling.Unlink()
	sibling.SetNext(n.next)
	if nil != sibling.Next() {
		sibling.Next().SetPrevious(sibling)
	}
	sibling.SetPrevious(this)
	n.next = sibling
	sibling.SetParent(n.parent)
	if nil == sibling.Next() {
		sibling.Parent().SetLastChild(sibling)
	}
}

func (n *Text) InsertBefore(this Node, sibling Node) {
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

func (n *Text) AppendChild(this, child Node) {
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
