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

import "strconv"

// Node 描述了节点结构。
type Node struct {
	// 不用接口实现的原因：
	//   1. 转换节点类型非常方便，只需修改 typ 属性
	//   2. 为了极致的性能而牺牲扩展性

	// 节点基础结构

	typ             nodeType // 节点类型
	parent          *Node    // 父节点
	previous        *Node    // 前一个兄弟节点
	next            *Node    // 后一个兄弟节点
	firstChild      *Node    // 第一个子节点
	lastChild       *Node    // 最后一个子节点
	rawText         string   // 原始内容
	tokens          []byte   // 词法分析结果 tokens，语法分析阶段会继续操作这些 tokens
	close           bool     // 标识是否关闭
	lastLineBlank   bool     // 标识最后一行是否是空行
	lastLineChecked bool     // 标识最后一行是否检查过

	// 代码

	codeMarkerLen int // ` 个数，1 或 2

	// 代码块

	isFencedCodeBlock    bool
	codeBlockFenceChar   byte
	codeBlockFenceLen    int
	codeBlockFenceOffset int
	codeBlockOpenFence   []byte
	codeBlockInfo        []byte
	codeBlockCloseFence  []byte

	// HTML 块

	htmlBlockType int // 规范中定义的 HTML 块类型（1-7）

	// 列表、列表项

	*listData

	// 任务列表项 [ ]、[x] 或者 [X]

	taskListItemChecked bool // 是否勾选

	// 表

	tableAligns    []int // 从左到右每个表格节点的对齐方式，0：默认对齐，1：左对齐，2：居中对齐，3：右对齐
	tableCellAlign int   // 表的单元格对齐方式

	// 链接

	linkType int // 链接类型，0：内联链接 [foo](/bar)，1：链接引用 [label]，2：自动链接

	// 标题

	headingLevel int // 1~6

	// 数学公式块

	mathBlockDollarOffset int

	// Vditor 所见即所得支持
	expand           bool   // 是否需要展开节点
	caretStartOffset string // 光标插入起始偏移位置
	caretEndOffset   string // 光标插入结束偏移位置
}

// lastDeepestChild 返回 n 的最后一个最深子节点。
func (n *Node) lastDeepestChild() (ret *Node) {
	if nil == n.lastChild {
		return n
	}
	return n.lastChild.lastDeepestChild()
}

// firstDeepestChild 返回 n 的第一个最深的子节点。
func (n *Node) firstDeepestChild() (ret *Node) {
	if nil == n.firstChild {
		return n
	}
	return n.firstChild.firstDeepestChild()
}

// LinkDest 在 n 的子节点中查找 childType 指定类型的第一个子节点。
func (n *Node) ChildByType(childType nodeType) *Node {
	for c := n.firstChild; nil != c; c = c.next {
		if c.typ == childType {
			return c
		}
	}
	return nil
}

// Text 返回 n 及其文本子节点的文本值。
func (n *Node) Text() (ret string) {
	Walk(n, func(n *Node, entering bool) (status WalkStatus, e error) {
		if (NodeText == n.typ || NodeLinkText == n.typ) && entering {
			ret += bytesToStr(n.tokens)
		}
		return WalkContinue, nil
	})
	return
}

func (n *Node) NextNodeText() string {
	if nil == n.next {
		return ""
	}
	return n.next.Text()
}

func (n *Node) PreviousNodeText() string {
	if nil == n.previous {
		return ""
	}
	return n.previous.Text()
}

// Finalize 节点最终化处理。比如围栏代码块提取 info 部分；HTML 代码块剔除结尾空格；段落需要解析链接引用定义等。
func (n *Node) Finalize(context *Context) {
	switch n.typ {
	case NodeCodeBlock:
		n.codeBlockFinalize(context)
	case NodeHTMLBlock:
		n.htmlBlockFinalize(context)
	case NodeParagraph:
		n.paragraphFinalize(context)
	case NodeMathBlock:
		n.mathBlockFinalize(context)
	case NodeList:
		n.listFinalize(context)
	}
}

// Continue 判断节点是否可以继续处理，比如块引用需要 >，缩进代码块需要 4 空格，围栏代码块需要 ```。
// 如果可以继续处理返回 0，如果不能接续处理返回 1，如果返回 2（仅在围栏代码块闭合时）则说明可以继续下一行处理了。
func (n *Node) Continue(context *Context) int {
	switch n.typ {
	case NodeCodeBlock:
		return n.codeBlockContinue(context)
	case NodeHTMLBlock:
		return n.htmlBlockContinue(context)
	case NodeParagraph:
		return n.paragraphContinue(context)
	case NodeListItem:
		return n.listItemContinue(context)
	case NodeBlockquote:
		return n.blockquoteContinue(context)
	case NodeMathBlock:
		return n.mathBlockContinue(context)
	case NodeHeading, NodeThematicBreak:
		return 1
	}

	return 0
}

// AcceptLines 判断是否节点是否可以接受更多的文本行。比如 HTML 块、代码块和段落是可以接受更多的文本行的。
func (n *Node) AcceptLines() bool {
	switch n.typ {
	case NodeParagraph, NodeCodeBlock, NodeHTMLBlock, NodeTable, NodeMathBlock:
		return true
	}
	return false
}

// CanContain 判断是否能够包含 NodeType 指定类型的节点。 比如列表节点（一种块级容器）只能包含列表项节点，
// 块引用节点（另一种块级容器）可以包含任意节点；段落节点（一种叶子块节点）不能包含任何其他块级节点。
func (n *Node) CanContain(nodeType nodeType) bool {
	switch n.typ {
	case NodeCodeBlock, NodeHTMLBlock, NodeParagraph, NodeThematicBreak, NodeTable, NodeMathBlock:
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

// AppendTokens 添加 tokens。
func (n *Node) AppendTokens(tokens []byte) {
	n.tokens = append(n.tokens, tokens...)
}

// InsertAfter 在当前节点后插入一个兄弟节点。
func (n *Node) InsertAfter(sibling *Node) {
	sibling.Unlink()
	sibling.next = n.next
	if nil != sibling.next {
		sibling.next.previous = sibling
	}
	sibling.previous = n
	n.next = sibling
	sibling.parent = n.parent
	if nil == sibling.next {
		sibling.parent.lastChild = sibling
	}
}

// InsertBefore 在当前节点前插入一个兄弟节点。
func (n *Node) InsertBefore(sibling *Node) {
	sibling.Unlink()
	sibling.previous = n.previous
	if nil != sibling.previous {
		sibling.previous.next = sibling
	}
	sibling.next = n
	n.previous = sibling
	sibling.parent = n.parent
	if nil == sibling.previous {
		sibling.parent.firstChild = sibling
	}
}

// AppendChild 在 n 的子节点最后再添加一个子节点。
func (n *Node) AppendChild(child *Node) {
	child.Unlink()
	child.parent = n
	if nil != n.lastChild {
		n.lastChild.next = child
		child.previous = n.lastChild
		n.lastChild = child
	} else {
		n.firstChild = child
		n.lastChild = child
	}
}

// PrependChild 在 n 的子节点最前添加一个子节点。
func (n *Node) PrependChild(child *Node) {
	child.Unlink()
	child.parent = n
	if nil != n.firstChild {
		n.firstChild.previous = child
		child.next = n.firstChild
		n.firstChild = child
	} else {
		n.firstChild = child
		n.lastChild = child
	}
}

// List 将 n 及其所有子节点按深度优先遍历添加到结果列表 ret 中。
func (n *Node) List() (ret []*Node) {
	ret = make([]*Node, 0, 512)
	Walk(n, func(n *Node, entering bool) (status WalkStatus, e error) {
		if entering {
			ret = append(ret, n)
		}
		return WalkContinue, nil
	})
	return
}

// isMarker 判断 n  是否是排版类（比如强调加粗）标记节点。
func (n *Node) isMarker() bool {
	switch n.typ {
	case NodeEmA6kOpenMarker, NodeEmA6kCloseMarker, NodeEmU8eOpenMarker, NodeEmU8eCloseMarker,
		NodeStrongA6kOpenMarker, NodeStrongA6kCloseMarker, NodeStrongU8eOpenMarker, NodeStrongU8eCloseMarker,
		NodeStrikethrough1OpenMarker, NodeStrikethrough1CloseMarker, NodeStrikethrough2OpenMarker, NodeStrikethrough2CloseMarker:
		// TODO: 判断节点是否是标记节点
		return true
	}
	return false
}

type nodeType int

func (typ nodeType) String() string {
	return strconv.Itoa(int(typ))
}

const (
	// CommonMark

	NodeDocument                  nodeType = 0  // 根
	NodeParagraph                 nodeType = 1  // 段落
	NodeHeading                   nodeType = 2  // 标题
	NodeHeadingC8hMarker          nodeType = 3  // ATX 标题标记符 #
	NodeThematicBreak             nodeType = 4  // 分隔线
	NodeBlockquote                nodeType = 5  // 块引用
	NodeBlockquoteMarker          nodeType = 6  // 块引用标记符 >
	NodeList                      nodeType = 7  // 列表
	NodeListItem                  nodeType = 8  // 列表项
	NodeHTMLBlock                 nodeType = 9  // HTML 块
	NodeInlineHTML                nodeType = 10 // 内联 HTML
	NodeCodeBlock                 nodeType = 11 // 代码块
	NodeCodeBlockFenceOpenMarker  nodeType = 12 // 开始围栏代码块标记符 ```
	NodeCodeBlockFenceCloseMarker nodeType = 13 // 结束围栏代码块标记符 ```
	NodeCodeBlockFenceInfoMarker  nodeType = 14 // 围栏代码块信息标记符 info string
	NodeCodeBlockCode             nodeType = 15 // 围栏代码块代码
	NodeText                      nodeType = 16 // 文本
	NodeEmphasis                  nodeType = 17 // 强调
	NodeEmA6kOpenMarker           nodeType = 18 // 开始强调标记符 *
	NodeEmA6kCloseMarker          nodeType = 19 // 结束强调标记符 *
	NodeEmU8eOpenMarker           nodeType = 20 // 开始强调标记符 _
	NodeEmU8eCloseMarker          nodeType = 21 // 结束强调标记符 _
	NodeStrong                    nodeType = 22 // 加粗
	NodeStrongA6kOpenMarker       nodeType = 23 // 开始加粗标记符 **
	NodeStrongA6kCloseMarker      nodeType = 24 // 结束加粗标记符 **
	NodeStrongU8eOpenMarker       nodeType = 25 // 开始加粗标记符 __
	NodeStrongU8eCloseMarker      nodeType = 26 // 结束加粗标记符 __
	NodeCodeSpan                  nodeType = 27 // 代码
	NodeCodeSpanOpenMarker        nodeType = 28 // 开始代码标记符 `
	NodeCodeSpanContent           nodeType = 29 // 代码内容
	NodeCodeSpanCloseMarker       nodeType = 30 // 结束代码标记符 `
	NodeHardBreak                 nodeType = 31 // 硬换行
	NodeSoftBreak                 nodeType = 32 // 软换行
	NodeLink                      nodeType = 33 // 链接
	NodeImage                     nodeType = 34 // 图片
	NodeBang                      nodeType = 35 // !
	NodeOpenBracket               nodeType = 36 // [
	NodeCloseBracket              nodeType = 37 // ]
	NodeOpenParen                 nodeType = 38 // (
	NodeCloseParen                nodeType = 39 // )
	NodeLinkText                  nodeType = 40 // 链接文本
	NodeLinkDest                  nodeType = 41 // 链接地址
	NodeLinkTitle                 nodeType = 42 // 链接标题
	NodeLinkSpace                 nodeType = 43 // 链接地址和链接标题之间的空格

	// GFM

	NodeTaskListItemMarker        nodeType = 100 // 任务列表项标记符
	NodeStrikethrough             nodeType = 101 // 删除线
	NodeStrikethrough1OpenMarker  nodeType = 102 // 开始删除线标记符 ~
	NodeStrikethrough1CloseMarker nodeType = 103 // 结束删除线标记符 ~
	NodeStrikethrough2OpenMarker  nodeType = 104 // 开始删除线标记符 ~~
	NodeStrikethrough2CloseMarker nodeType = 105 // 结束删除线标记符 ~~
	NodeTable                     nodeType = 106 // 表
	NodeTableHead                 nodeType = 107 // 表头
	NodeTableRow                  nodeType = 108 // 表行
	NodeTableCell                 nodeType = 109 // 表格

	// Emoji

	NodeEmoji        nodeType = 200 // Emoji
	NodeEmojiUnicode nodeType = 201 // Emoji Unicode
	NodeEmojiImg     nodeType = 202 // Emoji 图片
	NodeEmojiAlias   nodeType = 203 // Emoji ASCII

	// 数学公式

	NodeMathBlock             nodeType = 300 // 数学公式块
	NodeMathBlockOpenMarker   nodeType = 301 // 开始数学公式块标记符 $$
	NodeMathBlockContent      nodeType = 302 // 数学公式块内容
	NodeMathBlockCloseMarker  nodeType = 303 // 结束数学公式块标记符 $$
	NodeInlineMath            nodeType = 304 // 内联数学公式
	NodeInlineMathOpenMarker  nodeType = 305 // 开始内联数学公式标记符 $
	NodeInlineMathContent     nodeType = 306 // 内联数学公式内容
	NodeInlineMathCloseMarker nodeType = 307 // 结束内联数学公式标记符 $
)
