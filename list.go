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

type ListType int

const (
	ListTypeBullet  = 0
	ListTypeOrdered = 1
)

type ListData struct {
	typ          ListType
	tight        bool
	bulletChar   string
	start        int
	delimiter    string
	padding      int
	markerOffset int
}

type List struct {
	*BaseNode
	*ListData
}

func (list *List) CanContain(nodeType NodeType) bool {
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

func (t *Tree) parseOrderedListMarker(tokens items) (marker items) {
	var i int
	var token item
	for ; ; i++ {
		token = tokens[i]
		marker = append(marker, token)
		if !token.isDigit() || 8 < i {
			break
		}
	}

	if 1 > len(marker) || itemDot == token || itemCloseParen == token {
		return nil
	}

	return
}

// Parse a list marker and return data on the marker (type,
// start, delimiter, bullet character, padding) or null.
func (t *Tree) parseListMarker(container Node) *ListData {
	if t.context.indent >= 4 {
		return nil
	}
	tokens := t.context.currentLine[t.context.nextNonspace:]
	data := &ListData{
		typ:          ListTypeBullet,
		tight:        true, // lists are tight by default
		markerOffset: t.context.indent,
	}

	markerLength := 1
	marker := items{tokens[0]}
	if itemPlus == marker[0] || itemHyphen == marker[0] || itemAsterisk == marker[0] {
		data.typ = ListTypeBullet
		data.bulletChar = string(marker[0])
	} else if marker = t.parseOrderedListMarker(tokens); nil != marker {
		if container.Type() != NodeParagraph || '1' == marker[0] {
			data.typ = ListTypeOrdered
			data.start, _ = strconv.Atoi(marker.rawText())
			markerLength = 2
			if itemDot == tokens[1] {
				data.delimiter = "."
			} else if itemCloseParen == tokens[1] {
				data.delimiter = ")"
			} else {
				return nil
			}
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
		if t.context.column-spacesStartCol >= 5 || itemEOF == nextc || (itemSpace != nextc && itemTab != nextc) {
			break
		}
	}

	token := t.context.currentLine.peek(t.context.offset)
	var blank_item = itemEOF == token || itemNewline == token
	var spaces_after_marker = t.context.column - spacesStartCol
	if spaces_after_marker >= 5 || spaces_after_marker < 1 || blank_item {
		data.padding = len(marker) + 1
		t.context.column = spacesStartCol
		t.context.offset = spacesStartOffset
		if token = t.context.currentLine.peek(t.context.offset); itemSpace == token || itemTab == token {
			t.context.advanceOffset(1, true)
		}
	} else {
		data.padding = len(marker) + spaces_after_marker
	}
	if data.typ == ListTypeOrdered {
		data.padding++ // 加上分隔符 . 或者 ) 为 1 的长度
	}
	return data
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
