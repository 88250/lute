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

// Node 描述了节点结构。
type Node struct {
	// 不用接口实现的原因：
	//   1. 转换节点类型非常方便，只需修改 typ 属性
	//   2. 为了极致的性能，牺牲一点扩展性

	// 节点基础结构

	typ             int    // 节点类型
	parent          *Node  // 父节点
	previous        *Node  // 前一个兄弟节点
	next            *Node  // 后一个兄弟节点
	firstChild      *Node  // 第一个子节点
	lastChild       *Node  // 最后一个子节点
	rawText         string // 原始内容
	tokens          items  // 词法分析结果 tokens，语法分析阶段会继续操作这些 tokens
	close           bool   // 标识是否关闭
	lastLineBlank   bool   // 标识最后一行是否是空行
	lastLineChecked bool   // 标识最后一行是否检查过

	// 代码块

	isFencedCodeBlock    bool
	codeBlockFenceChar   byte
	codeBlockFenceLength int
	codeBlockFenceOffset int
	codeBlockInfo        items

	// HTML 块

	htmlBlockType int // 规范中定义的 HTML 块类型（1-7）

	// 列表、列表项

	*listData

	// 链接、图片

	destination items // 链接地址
	title       items // 链接标题

	// 任务列表项 [ ]、[x] 或者 [X]

	taskListItemChecked bool // 是否勾选

	// 表

	tableAligns    []int // 从左到右每个表格节点的对齐方式，0：默认对齐，1：左对齐，2：居中对齐，3：右对齐
	tableCellAlign int   // 表的单元格对齐方式

	// 标题

	headingLevel int // 1~6

	// 强调、加粗和删除线

	strongEmDelMarker    byte
	strongEmDelMarkenLen int
}

// Finalize 节点最终化处理。比如围栏代码块提取 info 部分；HTML 代码块剔除结尾空格；段落需要解析链接引用定义等。
func (n *Node) Finalize(context *Context) {
	switch n.typ {
	case NodeCodeBlock:
		n.CodeBlockFinalize(context)
	case NodeHTMLBlock:
		n.HTMLBlockFinalize(context)
	case NodeParagraph:
		n.ParagraphFinalize(context)
	case NodeList:
		n.ListFinalize(context)
	}
}

// Continue 判断节点是否可以继续处理，比如块引用需要 >，缩进代码块需要 4 空格，围栏代码块需要 ```。
// 如果可以继续处理返回 0，如果不能接续处理返回 1，如果返回 2（仅在围栏代码块闭合时）则说明可以继续下一行处理了。
func (n *Node) Continue(context *Context) int {
	switch n.typ {
	case NodeCodeBlock:
		return n.CodeBlockContinue(context)
	case NodeHTMLBlock:
		return n.HTMLBlockContinue(context)
	case NodeParagraph:
		return n.ParagraphContinue(context)
	case NodeListItem:
		return n.ListItemContinue(context)
	case NodeBlockquote:
		return n.BlockquoteContinue(context)
	case NodeHeading, NodeThematicBreak:
		return 1
	}

	return 0
}

// AcceptLines 判断是否节点是否可以接受更多的文本行。比如 HTML 块、代码块和段落是可以接受更多的文本行的。
func (n *Node) AcceptLines() bool {
	switch n.typ {
	case NodeParagraph, NodeCodeBlock, NodeHTMLBlock, NodeTable:
		return true
	}
	return false
}

// CanContain 判断是否能够包含 NodeType 指定类型的节点。 比如列表节点（一种块级容器）只能包含列表项节点，
// 块引用节点（另一种块级容器）可以包含任意节点；段落节点（一种叶子块节点）不能包含任何其他块级节点。
func (n *Node) CanContain(nodeType int) bool {
	switch n.typ {
	case NodeCodeBlock, NodeHTMLBlock, NodeParagraph, NodeHeading, NodeThematicBreak, NodeTable:
		return false
	case NodeList:
		return NodeListItem == nodeType
	}

	return NodeListItem != nodeType
}

// Unlink 用于将节点从树上移除，后一个兄弟节点会接替该节点。
func (n *Node) Unlink() {
	if nil != n.previous {
		n.previous.next = n.next
	} else if nil != n.parent {
		n.parent.firstChild = n.next
	}
	if nil != n.next {
		n.next.previous = n.previous
	} else if nil != n.parent {
		n.parent.lastChild = n.previous
	}
	n.parent = nil
	n.next = nil
	n.previous = nil
}

// RawText 返回原始内容。
func (n *Node) RawText() string {
	return n.rawText
}

// SetRawText 设置原始内容。
func (n *Node) SetRawText(rawText string) {
	n.rawText = rawText
}

// AppendRawText 添加原始内容。
func (n *Node) AppendRawText(rawText string) {
	n.rawText += rawText
}

// Tokens 返回所有 tokens。
func (n *Node) Tokens() items {
	return n.tokens
}

// SetTokens 设置 tokens。
func (n *Node) SetTokens(tokens items) {
	n.tokens = tokens
}

// AppendTokens 添加 tokens。
func (n *Node) AppendTokens(tokens items) {
	n.tokens = append(n.tokens, tokens...)
}

// InsertAfter 在当前节点后插入一个兄弟节点。
func (n *Node) InsertAfter(this *Node, sibling *Node) {
	sibling.Unlink()
	sibling.next = n.next
	if nil != sibling.next {
		sibling.next.previous = sibling
	}
	sibling.previous = this
	n.next = sibling
	sibling.parent = n.parent
	if nil == sibling.next {
		sibling.parent.lastChild = sibling
	}
}

// InsertBefore 在当前节点前插入一个兄弟节点。
func (n *Node) InsertBefore(this *Node, sibling *Node) {
	sibling.Unlink()
	sibling.previous = n.previous
	if nil != sibling.previous {
		sibling.previous.next = sibling
	}
	sibling.next = this
	n.previous = sibling
	sibling.parent = n.parent
	if nil == sibling.previous {
		sibling.parent.firstChild = sibling
	}
}

// AppendChild 添加一个子节点。
func (n *Node) AppendChild(this, child *Node) {
	child.Unlink()
	child.parent = this
	if nil != n.lastChild {
		n.lastChild.next = child
		child.previous = n.lastChild
		n.lastChild = child
	} else {
		n.firstChild = child
		n.lastChild = child
	}
}

const (
	// CommonMark

	NodeDocument      = iota // 根节点类
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

	// GFM

	NodeTaskListItemMarker // 任务列表项标记节点
	NodeStrikethrough      // 删除线节点
	NodeTable              // 表节点
	NodeTableHead          // 表头节点
	NodeTableRow           // 表行节点
	NodeTableCell          // 表格节点
)
