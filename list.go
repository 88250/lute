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

import (
	"strconv"
)

// listData 用于记录列表或列表项节点的附加信息。
type listData struct {
	typ          int   // 0：无序列表，1：有序列表，3：任务列表
	tight        bool  // 是否是紧凑模式
	bulletChar   items // 无序列表标识，* - 或者 +
	start        int   // 有序列表起始序号
	delimiter    *item // 有序列表分隔符，. 或者 )
	padding      int   // 列表内部缩进空格数（包含标识符长度，即规范中的 W+N）
	markerOffset int   // 标识符（* - + 或者 1 2 3）相对缩进空格数
	checked      bool  // 任务列表项是否勾选
	marker       items // 列表标识符
	num          int   // 有序列表项修正过的序号
}

func (list *Node) listFinalize(context *Context) {
	item := list.firstChild

	// 检查子列表项之间是否包含空行，包含的话说明该列表是非紧凑的，即松散的
	for nil != item {
		if list.endsWithBlankLine(item) && nil != item.next {
			list.tight = false
			break
		}

		var subitem = item.firstChild
		for nil != subitem {
			if list.endsWithBlankLine(subitem) &&
				(nil != item.next || nil != subitem.next) {
				list.tight = false
				break
			}
			subitem = subitem.next
		}
		item = item.next
	}
}

var items1 = strToItems("1")

// parseListMarker 用于解析泛列表（列表、列表项或者任务列表）标记符。
func (t *Tree) parseListMarker(container *Node) *listData {
	if t.context.indent >= 4 {
		return nil
	}

	ln := t.context.currentLine // 弄短点
	tokens := ln[t.context.nextNonspace:]
	data := &listData{
		typ:          0,                // 默认无序列表
		tight:        true,             // 默认紧凑模式
		markerOffset: t.context.indent, // 设置前置相对缩进
		num:          -1,               // 假设有序列表起始为 -1，后面会进行计算赋值
	}

	markerLength := 1
	marker := items{tokens[0]}
	var delim *item
	if itemPlus == term(marker[0]) || itemHyphen == term(marker[0]) || itemAsterisk == term(marker[0]) {
		data.bulletChar = marker
	} else if marker, delim = t.parseOrderedListMarker(tokens); nil != marker {
		if container.typ != NodeParagraph || equal(items1, marker) {
			data.typ = 1 // 有序列表
			data.start, _ = strconv.Atoi(itemsToStr(marker))
			markerLength = len(marker) + 1
			data.delimiter = delim
		} else {
			return nil
		}
	} else {
		return nil
	}

	data.marker = marker

	var token = ln[t.context.nextNonspace+markerLength]
	// 列表项标记符后必须是空白字符
	if !isWhitespace(term(token)) {
		return nil
	}

	// 如果要打断段落，则列表项内容部分不能为空
	if container.typ == NodeParagraph && itemNewline == term(ln[t.context.nextNonspace+markerLength]) {
		return nil
	}

	// 到这里说明满足列表规则，开始解析并计算内部缩进空格数
	t.context.advanceNextNonspace()             // 把起始下标移动到标记符起始位置
	t.context.advanceOffset(markerLength, true) // 把结束下标移动到标记符结束位置
	spacesStartCol := t.context.column
	spacesStartOffset := t.context.offset
	for {
		t.context.advanceOffset(1, true)
		token = ln.peek(t.context.offset)
		if t.context.column-spacesStartCol >= 5 || nil == token || (itemSpace != term(token) && itemTab != term(token)) {
			break
		}
	}

	token = ln.peek(t.context.offset)
	var isBlankItem = nil == token || itemNewline == term(token)
	var spacesAfterMarker = t.context.column - spacesStartCol
	if spacesAfterMarker >= 5 || spacesAfterMarker < 1 || isBlankItem {
		data.padding = markerLength + 1
		t.context.column = spacesStartCol
		t.context.offset = spacesStartOffset
		if token = ln.peek(t.context.offset); itemSpace == term(token) || itemTab == term(token) {
			t.context.advanceOffset(1, true)
		}
	} else {
		data.padding = markerLength + spacesAfterMarker
	}

	if !isBlankItem {
		// 判断是否是任务列表项
		content := ln[t.context.offset:]
		if 3 <= len(content) { // 至少需要 [ ] 或者 [x] 3 个字符
			if itemOpenBracket == term(content[0]) && ('x' == term(content[1]) || 'X' == term(content[1]) || itemSpace == term(content[1])) && itemCloseBracket == term(content[2]) {
				data.typ = 3
				data.checked = 'x' == term(content[1]) || 'X' == term(content[1])
			}
		}
	}

	return data
}

func (t *Tree) parseOrderedListMarker(tokens items) (marker items, delimiter *item) {
	length := len(tokens)
	var i int
	var token *item
	for ; i < length; i++ {
		token = tokens[i]
		if !isDigit(term(token)) || 8 < i {
			delimiter = token
			break
		}
		marker = append(marker, token)
	}

	if 1 > len(marker) || (itemDot != term(delimiter) && itemCloseParen != term(delimiter)) {
		return nil, nil
	}

	return
}

// endsWithBlankLine 判断块节点 block 是否是空行结束。如果 block 是列表或者列表项则迭代下降进入判断。
func (list *Node) endsWithBlankLine(block *Node) bool {
	for nil != block {
		if block.lastLineBlank {
			return true
		}
		t := block.typ
		if !block.lastLineChecked && (t == NodeList || t == NodeListItem) {
			block.lastLineChecked = true
			block = block.lastChild
		} else {
			block.lastLineChecked = true
			break
		}
	}

	return false
}
