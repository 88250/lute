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

import (
	"strconv"
)

// List 描述了列表节点结构。
type List struct {
	*BaseNode
	*listData
}

type listData struct {
	typ          int   // 0：无序列表，1：有序列表
	tight        bool  // 是否是紧凑模式
	bulletChar   items // 无序列表标识，* - 或者 +
	start        int   // 有序列表开始序号
	delimiter    byte  // 有序列表分隔符，. 或者 )
	padding      int   // 列表内部缩进空格数，即无序列表标识或分隔符和后续第一个非空字符之间的空格数，规范里的 N
	markerOffset int   // 标识符（* - + 或者 1 2 3）缩进
}

func (list *List) CanContain(nodeType int) bool {
	return NodeListItem == nodeType
}

func (list *List) Finalize(context *Context) {
	item := list.firstChild
	for nil != item {
		// check for non-final list item ending with blank line:
		if list.endsWithBlankLine(item) && nil != item.Next() {
			list.tight = false
			break
		}

		// recurse into children of list item, to see if there are
		// spaces between any of them:
		var subitem = item.FirstChild()
		for nil != subitem {
			if list.endsWithBlankLine(subitem) &&
				(nil != item.Next() || nil != subitem.Next()) {
				list.tight = false
				break
			}
			subitem = subitem.Next()
		}
		item = item.Next()
	}
}

// Parse a list marker and return data on the marker (type,
// start, delimiter, bullet character, padding) or null.
func (t *Tree) parseListMarker(container Node) *listData {
	if t.context.indent >= 4 {
		return nil
	}
	tokens := t.context.currentLine[t.context.nextNonspace:]
	data := &listData{
		typ:          0,    // 无序列表
		tight:        true, // lists are tight by default
		markerOffset: t.context.indent,
	}

	markerLength := 1
	marker := items{tokens[0]}
	if itemPlus == marker[0] || itemHyphen == marker[0] || itemAsterisk == marker[0] {
		data.bulletChar = marker
	} else if marker, delim := t.parseOrderedListMarker(tokens); nil != marker {
		if container.Type() != NodeParagraph || "1" == marker.string() {
			data.typ = 1 // 有序列表
			data.start, _ = strconv.Atoi(marker.string())
			markerLength = len(marker) + 1
			data.delimiter = delim
		} else {
			return nil
		}
	} else {
		return nil
	}

	// make sure we have spaces after
	nextc := t.context.currentLine[t.context.nextNonspace+markerLength]
	if itemNewline != nextc && (itemSpace != nextc && itemTab != nextc) {
		return nil
	}

	// if it interrupts paragraph, make sure first line isn't blank
	if container.Type() == NodeParagraph && itemNewline == t.context.currentLine[t.context.nextNonspace+markerLength] {
		return nil
	}

	// we've got a match! advance offset and calculate padding
	t.context.advanceNextNonspace()             // to start of marker
	t.context.advanceOffset(markerLength, true) // to end of marker
	spacesStartCol := t.context.column
	spacesStartOffset := t.context.offset
	for {
		t.context.advanceOffset(1, true)
		nextc = t.context.currentLine.peek(t.context.offset)
		if t.context.column-spacesStartCol >= 5 || itemEnd == nextc || (itemSpace != nextc && itemTab != nextc) {
			break
		}
	}

	token := t.context.currentLine.peek(t.context.offset)
	var blank_item = itemEnd == token || itemNewline == token
	var spaces_after_marker = t.context.column - spacesStartCol
	if spaces_after_marker >= 5 || spaces_after_marker < 1 || blank_item {
		data.padding = markerLength + 1
		t.context.column = spacesStartCol
		t.context.offset = spacesStartOffset
		if token = t.context.currentLine.peek(t.context.offset); itemSpace == token || itemTab == token {
			t.context.advanceOffset(1, true)
		}
	} else {
		data.padding = markerLength + spaces_after_marker
	}

	return data
}

func (t *Tree) parseOrderedListMarker(tokens items) (marker items, delimiter byte) {
	var i int
	var token byte
	for ; ; i++ {
		token = tokens[i]
		if !isDigit(token) || 8 < i {
			delimiter = token
			break
		}
		marker = append(marker, token)
	}

	if 1 > len(marker) || (itemDot != delimiter && itemCloseParen != delimiter) {
		return nil, itemEnd
	}

	return
}

// Returns true if block ends with a blank line, descending if needed into lists and sublists.
func (list *List) endsWithBlankLine(block Node) bool {
	for nil != block {
		if block.LastLineBlank() {
			return true
		}
		t := block.Type()
		if !block.LastLineChecked() && (t == NodeList || t == NodeListItem) {
			block.SetLastLineChecked(true)
			block = block.LastChild()
		} else {
			block.SetLastLineChecked(true)
			break
		}
	}

	return false
}
