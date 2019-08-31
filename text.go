// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

// mergeText 合并 node 中所有（包括子节点）连续的文本节点。
// 合并后顺便进行中文排版优化以及 GFM 自动邮件链接识别。
func (t *Tree) mergeText(node *BaseNode) {
	if nil == node {
		return
	}

	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ {
			// 逐个合并后续兄弟节点
			for nil != next && NodeText == next.typ {
				child.AppendTokens(next.Tokens())
				next.Unlink()
				next = child.next
			}
		} else {
			t.mergeText(child) // 递归处理子节点
		}
		child = next
	}
}

//
//func (n *Text) Type() int {
//	return NodeText
//}
//
//func (n *Text) IsOpen() bool {
//	return !n.close
//}
//
//func (n *Text) IsClosed() bool {
//	return n.close
//}
//
//func (n *Text) Close() {
//	n.close = true
//}
//
//func (n *Text) Finalize(context *Context) {
//}
//
//func (n *Text) Continue(context *Context) int {
//	return 0
//}
//
//func (n *Text) AcceptLines() bool {
//	return false
//}
//
//func (n *Text) CanContain(nodeType int) bool {
//	return NodeListItem != nodeType
//}
//
//func (n *Text) LastLineBlank() bool {
//	return n.lastLineBlank
//}
//
//func (n *Text) SetLastLineBlank(lastLineBlank bool) {
//	n.lastLineBlank = lastLineBlank
//}
//
//func (n *Text) LastLineChecked() bool {
//	return n.lastLineChecked
//}
//
//func (n *Text) SetLastLineChecked(lastLineChecked bool) {
//	n.lastLineChecked = lastLineChecked
//}
//
//func (n *Text) Unlink() {
//	if nil != n.previous {
//		n.previous.SetNext(n.next)
//	} else if nil != n.parent {
//		n.parent.SetFirstChild(n.next)
//	}
//	if nil != n.next {
//		n.next.SetPrevious(n.previous)
//	} else if nil != n.parent {
//		n.parent.SetLastChild(n.previous)
//	}
//	n.parent = nil
//	n.next = nil
//	n.previous = nil
//}
//
//func (n *Text) Parent() *BaseNode {
//	return n.parent
//}
//
//func (n *Text) SetParent(parent *BaseNode) {
//	n.parent = parent
//}
//
//func (n *Text) Next() *BaseNode {
//	return n.next
//}
//
//func (n *Text) SetNext(next *BaseNode) {
//	n.next = next
//}
//
//func (n *Text) Previous() *BaseNode {
//	return n.previous
//}
//
//func (n *Text) SetPrevious(previous *BaseNode) {
//	n.previous = previous
//}
//
//func (n *Text) FirstChild() *BaseNode {
//	return n.firstChild
//}
//
//func (n *Text) SetFirstChild(firstChild *BaseNode) {
//	n.firstChild = firstChild
//}
//
//func (n *Text) LastChild() *BaseNode {
//	return n.lastChild
//}
//
//func (n *Text) SetLastChild(lastChild *BaseNode) {
//	n.lastChild = lastChild
//}
//
//func (n *Text) RawText() string {
//	return n.rawText
//}
//
//func (n *Text) SetRawText(rawText string) {
//	n.rawText = rawText
//}
//
//func (n *Text) AppendRawText(rawText string) {
//	n.rawText += rawText
//}
//
//func (n *Text) Tokens() items {
//	return n.tokens
//}
//
//func (n *Text) SetTokens(tokens items) {
//	n.tokens = tokens
//}
//func (n *Text) AppendTokens(tokens items) {
//	n.tokens = append(n.tokens, tokens...)
//}
//
//func (n *Text) InsertAfter(this *BaseNode, sibling *BaseNode) {
//	sibling.Unlink()
//	sibling.SetNext(n.next)
//	if nil != sibling.Next() {
//		sibling.Next().SetPrevious(sibling)
//	}
//	sibling.SetPrevious(this)
//	n.next = sibling
//	sibling.SetParent(n.parent)
//	if nil == sibling.Next() {
//		sibling.Parent().SetLastChild(sibling)
//	}
//}
//
//func (n *Text) InsertBefore(this *BaseNode, sibling *BaseNode) {
//	sibling.Unlink()
//	sibling.SetPrevious(n.previous)
//	if nil != sibling.Previous() {
//		sibling.Previous().SetNext(sibling)
//	}
//	sibling.SetNext(this)
//	n.previous = sibling
//	sibling.SetParent(n.parent)
//	if nil == sibling.Previous() {
//		sibling.Parent().SetFirstChild(sibling)
//	}
//}
//
//func (n *Text) AppendChild(this, child *BaseNode) {
//	child.Unlink()
//	child.SetParent(this)
//	if nil != n.lastChild {
//		n.lastChild.SetNext(child)
//		child.SetPrevious(n.lastChild)
//		n.lastChild = child
//	} else {
//		n.firstChild = child
//		n.lastChild = child
//	}
//}
