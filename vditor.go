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

// +build javascript

package lute

// RenderVditorDOM 用于渲染 Vditor DOM，start 和 end 是光标位置，从 0 开始。
func (lute *Lute) RenderVditorDOM(markdownText string, startOffset, endOffset int) (html string, err error) {
	var tree *Tree
	lute.VditorWYSIWYG = true
	markdownText = lute.endNewline(markdownText)
	tree, err = lute.parse("", []byte(markdownText))
	if nil != err {
		return
	}

	renderer := lute.newVditorRenderer(tree.Root)
	startOffset, endOffset = renderer.byteOffset(markdownText, startOffset, endOffset)
	renderer.mapSelection(tree.Root, startOffset, endOffset)
	var output []byte
	output, err = renderer.Render()
	html = string(output)
	return
}

// VditorOperation 用于在 markdownText 中 startOffset、endOffset 指定的选段上做出 operation 操作。
// 支持的 operation 如下：
//   * newline 换行
func (lute *Lute) VditorOperation(markdownText string, startOffset, endOffset int, operation string) (html string, err error) {
	var tree *Tree
	lute.VditorWYSIWYG = true
	markdownText = lute.endNewline(markdownText)
	tree, err = lute.parse("", []byte(markdownText))
	if nil != err {
		return
	}

	renderer := lute.newVditorRenderer(tree.Root)
	startOffset, endOffset = renderer.byteOffset(markdownText, startOffset, endOffset)

	var nodes []*Node
	for c := tree.Root.firstChild; nil != c; c = c.next {
		renderer.findSelection(c, startOffset, endOffset, &nodes)
	}

	if 1 > len(nodes) {
		// 当且仅当渲染空 Markdown 时
		nodes = append(nodes, tree.Root)
	}

	// TODO: 仅实现了光标插入位置节点获取，选段映射待实现

	en := renderer.nearest(nodes, endOffset)

	baseOffset := 0
	if 0 < len(en.tokens) {
		baseOffset = en.tokens[0].Offset()
	}
	endOffset = endOffset - baseOffset

	newNode := &Node{typ: en.typ, tokens: en.tokens[endOffset:]}
	en.tokens = en.tokens[:endOffset]

	// 在父节点方向上获取节点路径
	var pathNodes []*Node
	var parent *Node
	for parent = en.parent; ; parent = parent.parent {
		pathNodes = append(pathNodes, parent)
		if NodeDocument == parent.typ || NodeListItem == parent.typ || NodeBlockquote == parent.typ {
			// 遇到块容器则停止
			break
		}
	}

	// 克隆新的节点路径
	length := len(pathNodes)
	top := pathNodes[length-1]
	newPath := &Node{typ: top.typ}
	if NodeListItem == top.typ {
		newPath.listData = top.listData
	}
	for i := length - 2; 0 <= i; i-- {
		n := pathNodes[i]
		newNode := &Node{typ: n.typ}
		newPath.AppendChild(newNode)
		newPath = newNode
	}

	// 把选段及其后续节点挂到新路径下
	newPath.AppendChild(newNode)
	newNode.caretStartOffset = "0"
	newNode.caretEndOffset = "0"
	newNode.expand = true
	for next := en.next; nil != next; next = next.next {
		newPath.AppendChild(next)
	}

	// 把新路径作为路径同级兄弟挂到树上
	np := newPath
	for ; nil != np.parent && NodeDocument != np.parent.typ; np = np.parent {
	}
	if NodeDocument == parent.typ {
		parent.AppendChild(np)
	} else {
		parent.InsertAfter(np)
	}

	// 格式化还原为 Markdown 原文
	formatRenderer := lute.newFormatRenderer(tree.Root)
	var output []byte
	output, err = formatRenderer.Render()
	if nil != err {
		return
	}

	// 再次解析生成语法树
	tree, err = lute.parse("", output)
	if nil != err {
		return
	}

	// 进行最终渲染
	renderer = lute.newVditorRenderer(tree.Root)
	output, err = renderer.Render()
	html = string(output)
	return
}

// VditorDOMMarkdown 用于将 Vditor DOM 转换为 Markdown 文本。
// TODO：改为解析标准 DOM
func (lute *Lute) VditorDOMMarkdown(html string) (markdown string, err error) {
	tree, err := lute.parseVditorDOM(html)
	if nil != err {
		return
	}

	var formatted []byte
	renderer := lute.newFormatRenderer(tree.Root)
	formatted, err = renderer.Render()
	if nil != err {
		return
	}
	markdown = bytesToStr(formatted)
	return
}
