// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/util"
	"strconv"
)

func listFinalize(list *ast.Node) {
	item := list.FirstChild

	// 检查子列表项之间是否包含空行，包含的话说明该列表是非紧凑的，即松散的
	for nil != item {
		if endsWithBlankLine(item) && nil != item.Next {
			list.Tight = false
			break
		}

		var subitem = item.FirstChild
		for nil != subitem {
			if endsWithBlankLine(subitem) &&
				(nil != item.Next || nil != subitem.Next) {
				list.Tight = false
				break
			}
			subitem = subitem.Next
		}
		item = item.Next
	}
}

var items1 = util.StrToBytes("1")

// parseListMarker 用于解析泛列表（列表、列表项或者任务列表）标记符。
func (t *Tree) parseListMarker(container *ast.Node) *ast.ListData {
	if t.context.indent >= 4 {
		return nil
	}

	ln := t.context.currentLine
	tokens := ln[t.context.nextNonspace:]
	data := &ast.ListData{
		Typ:          0,                // 默认无序列表
		Tight:        true,             // 默认紧凑模式
		MarkerOffset: t.context.indent, // 设置前置相对缩进
		Num:          -1,               // 假设有序列表起始为 -1，后面会进行计算赋值
	}

	markerLength := 1
	marker := []byte{tokens[0]}
	var delim byte
	if itemPlus == marker[0] || itemHyphen == marker[0] || itemAsterisk == marker[0] {
		data.BulletChar = marker[0]
	} else if marker, delim = t.parseOrderedListMarker(tokens); nil != marker {
		if container.Type != ast.NodeParagraph || bytes.Equal(items1, marker) {
			data.Typ = 1 // 有序列表
			data.Start, _ = strconv.Atoi(util.BytesToStr(marker))
			markerLength = len(marker) + 1
			data.Delimiter = delim
		} else {
			return nil
		}
	} else {
		return nil
	}

	data.Marker = marker

	var token = ln[t.context.nextNonspace+markerLength]
	// 列表项标记符后必须是空白字符
	if !isWhitespace(token) {
		return nil
	}

	// 如果要打断段落，则列表项内容部分不能为空
	if container.Type == ast.NodeParagraph && itemNewline == token {
		return nil
	}

	// 到这里说明满足列表规则，开始解析并计算内部缩进空格数
	t.context.advanceNextNonspace()             // 把起始下标移动到标记符起始位置
	t.context.advanceOffset(markerLength, true) // 把结束下标移动到标记符结束位置
	spacesStartCol := t.context.column
	spacesStartOffset := t.context.offset
	for {
		t.context.advanceOffset(1, true)
		token = peek(ln, t.context.offset)
		if t.context.column-spacesStartCol >= 5 || 0 == (token) || (itemSpace != token && itemTab != token) {
			break
		}
	}

	token = peek(ln, t.context.offset)
	var isBlankItem = 0 == token || itemNewline == token
	var spacesAfterMarker = t.context.column - spacesStartCol
	if spacesAfterMarker >= 5 || spacesAfterMarker < 1 || isBlankItem {
		data.Padding = markerLength + 1
		t.context.column = spacesStartCol
		t.context.offset = spacesStartOffset
		if token = peek(ln, t.context.offset); itemSpace == token || itemTab == token {
			t.context.advanceOffset(1, true)
		}
	} else {
		data.Padding = markerLength + spacesAfterMarker
	}

	if !isBlankItem {
		// 判断是否是任务列表项
		content := ln[t.context.offset:]
		if 3 <= len(content) { // 至少需要 [ ] 或者 [x] 3 个字符
			if itemOpenBracket == content[0] && ('x' == content[1] || 'X' == content[1] || itemSpace == content[1]) && itemCloseBracket == content[2] {
				data.Typ = 3
				data.Checked = 'x' == content[1] || 'X' == content[1]
			}
		}
	}
	return data
}

func (t *Tree) parseOrderedListMarker(tokens []byte) (marker []byte, delimiter byte) {
	length := len(tokens)
	var i int
	var token byte
	for ; i < length; i++ {
		token = tokens[i]
		if !isDigit(token) || 8 < i {
			delimiter = token
			break
		}
		marker = append(marker, token)
	}

	if 1 > len(marker) || (itemDot != delimiter && itemCloseParen != delimiter) {
		return nil, 0
	}

	return
}

// endsWithBlankLine 判断块节点 block 是否是空行结束。如果 block 是列表或者列表项则迭代下降进入判断。
func endsWithBlankLine(block *ast.Node) bool {
	for nil != block {
		if block.LastLineBlank {
			return true
		}
		t := block.Type
		if !block.LastLineChecked && (t == ast.NodeList || t == ast.NodeListItem) {
			block.LastLineChecked = true
			block = block.LastChild
		} else {
			block.LastLineChecked = true
			break
		}
	}

	return false
}
