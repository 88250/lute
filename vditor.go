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

import (
	"strconv"
	"strings"
)

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

// VditorNewline 用于在类型为 blockType 的块中进行换行生成新的 Vditor 节点。
// param 用于传递生成某些块换行所需的参数，比如在列表项中换行需要传列表项标记符。
func (lute *Lute) VditorNewline(blockType nodeType, param map[string]interface{}) (html string, err error) {
	var renderer *VditorRenderer

	switch blockType {
	case NodeListItem:
		markerPart := param["marker"].(string)
		listType := 0
		num := 1
		marker := "*"
		delim := ""
		listItem := &Node{typ: NodeListItem}
		if isASCIILetterNum(markerPart[0]) {
			listType = 1 // 有序列表
			if strings.Contains(markerPart, ")") {
				delim = ")"
			} else {
				delim = "."
			}
			markerParts := strings.SplitN(markerPart, delim, 2)
			num, _ = strconv.Atoi(markerParts[0])
			num++
			marker = strconv.Itoa(num)
		} else {
			marker = string(markerPart[0])
		}
		listItem.listData = &listData{typ: listType, marker: strToItems(marker)}
		listItem.expand = true
		if 1 == listType {
			listItem.delimiter = newItem(delim[0], 0, 0, 0)
		}

		text := &Node{typ: NodeText, tokens: strToItems("")}
		text.caretStartOffset = "0"
		text.caretEndOffset = "0"
		p := &Node{typ: NodeParagraph}
		p.AppendChild(text)
		listItem.AppendChild(p)
		listItem.tight = true
		renderer = lute.newVditorRenderer(listItem)
		var output []byte
		output, err = renderer.Render()
		if nil != err {
			return
		}
		html = string(output)
	default:
		renderer = lute.newVditorRenderer(nil)
		renderer.writer.WriteString("<p data-ntype=\"" + NodeParagraph.String() + "\" data-mtype=\"" + renderer.mtype(NodeParagraph) + "\" data-cso=\"0\" data-ceo=\"0\">" +
			"<br><span class=\"newline\">\n\n</span></p>")
		html = renderer.writer.String()
	}
	return
}
