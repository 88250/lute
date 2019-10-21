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

import "strconv"

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
	sn := en
	baseOffset := 0
	if 0 < len(sn.tokens) {
		baseOffset = sn.tokens[0].Offset()
	}

	startOffset = startOffset - baseOffset
	endOffset = endOffset - baseOffset
	startOffset, endOffset = renderer.runeOffset(itemsToBytes(sn.tokens), startOffset, endOffset)
	sn.caretStartOffset = strconv.Itoa(startOffset)
	en.caretEndOffset = strconv.Itoa(endOffset)
	renderer.expand(sn)

	//var renderer *VditorRenderer
	//
	//switch blockType {
	//case NodeListItem:
	//	markerPart := param["marker"].(string)
	//	listType := 0
	//	num := 1
	//	marker := "*"
	//	delim := ""
	//	listItem := &Node{typ: NodeListItem}
	//	if isASCIILetterNum(markerPart[0]) {
	//		listType = 1 // 有序列表
	//		if strings.Contains(markerPart, ")") {
	//			delim = ")"
	//		} else {
	//			delim = "."
	//		}
	//		markerParts := strings.SplitN(markerPart, delim, 2)
	//		num, _ = strconv.Atoi(markerParts[0])
	//		num++
	//		marker = strconv.Itoa(num)
	//	} else {
	//		marker = string(markerPart[0])
	//	}
	//	listItem.listData = &listData{typ: listType, marker: strToItems(marker)}
	//	listItem.expand = true
	//	if 1 == listType {
	//		listItem.delimiter = newItem(delim[0], 0, 0, 0)
	//	}
	//
	//	text := &Node{typ: NodeText, tokens: strToItems("")}
	//	text.caretStartOffset = "0"
	//	text.caretEndOffset = "0"
	//	p := &Node{typ: NodeParagraph}
	//	p.AppendChild(text)
	//	listItem.AppendChild(p)
	//	listItem.tight = true
	//	renderer = lute.newVditorRenderer(listItem)
	//	var output []byte
	//	output, err = renderer.Render()
	//	if nil != err {
	//		return
	//	}
	//	html = string(output)
	//default:
	//	renderer = lute.newVditorRenderer(nil)
	//	renderer.writer.WriteString("<p data-ntype=\"" + NodeParagraph.String() + "\" data-mtype=\"" + renderer.mtype(NodeParagraph) + "\" data-cso=\"0\" data-ceo=\"0\">" +
	//		"<br><span class=\"newline\">\n\n</span></p>")
	//	html = renderer.writer.String()
	//}
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
