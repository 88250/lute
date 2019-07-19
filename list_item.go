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

type ListItem struct {
	*BaseNode

	Checked bool
	Tight   bool
}

func (n *ListItem) Close() {
	if n.close {
		return
	}

	for child := n.FirstChild(); nil != child; child = child.Next() {
		child.Close()
	}
}

func (t *Tree) newListItem(marker string, bullet bool, start int, delim string, startIndentSpaces, indentSpaces int) (ret Node) {
	baseNode := &BaseNode{typ: NodeListItem, tokens: items{}}
	ret = &ListItem{
		baseNode,
		false,
		true,
	}

	return
}

func (t *Tree) parseListItem(line items) (ret Node) {
	if line.isBlankLine() {
		ret = &ListItem{BaseNode: &BaseNode{typ: NodeListItem}, Tight: true}
		return
	}

	indentSpaces := t.context.IndentSpaces

	line, marker, delim, bullet, start, startIndentSpaces, w, n := t.parseListMarker(line)
	ret = t.newListItem(marker, bullet, start, delim, startIndentSpaces, startIndentSpaces+w+n)

	var blankLineIndices []int
	i := 0
	for ; ; i++ {
		n := t.parseBlock(line)
		if nil == n {
			break
		}
		ret.AppendChild(ret, n)
		t.context.IndentSpaces = indentSpaces

		blankLines := t.skipBlankLines()
		if 0 < len(blankLines) {
			blankLineIndices = append(blankLineIndices, i)
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		if 0 < t.blockquoteMarkerCount(line) && 0 < t.context.BlockquoteLevel {
			line = t.removeStartBlockquoteMarker(line, t.context.BlockquoteLevel)
		}

		if t.context.IndentSpaces <= line.spaceCountLeft() {
			line = t.indentOffset(line, t.context.IndentSpaces)
			continue
		}

		t.backupLine(line)

		break
	}

	if 1 < len(blankLineIndices) {
		ret.(*ListItem).Tight = false
	} else if 1 == len(blankLineIndices) {
		ret.(*ListItem).Tight = 1 == len(blankLineIndices) && blankLineIndices[0] == i
	}

	t.context.CurNodes.push(ret)

	return
}

func (t *Tree) parseListItemMarker(line items, list Node) (remains items, marker, delim string, startIndentSpaces, indentSpaces int) {
	remains, marker, delim, startIndentSpaces, indentSpaces = t.parseListItemMarker0(line)

	if remains.isBlankLine() {
		remains = t.nextLine()
		if remains.isBlankLine() {
			list.AppendChild(list, &ListItem{BaseNode: &BaseNode{typ: NodeListItem}, Tight: false})
			t.skipBlankLines()
			remains = t.nextLine()
			remains, marker, delim, startIndentSpaces, indentSpaces = t.parseListItemMarker0(remains)

			return
		}

		if isList, marker := t.isList(remains); isList {
			list.AppendChild(list, &ListItem{BaseNode: &BaseNode{typ: NodeListItem}, Tight: true})
			remains = remains[len(marker):]
		}

		remains = t.indentOffset(remains, t.context.IndentSpaces)
	}

	return
}

func (t *Tree) parseListItemMarker0(line items) (remains items, marker, delim string, startIndentSpaces, indentSpaces int) {
	spaces, tabs, tokens, firstNonWhitespace := t.nonWhitespace(line)
	var markers items
	markers = append(markers, firstNonWhitespace)
	line = line[len(tokens):]
	if firstNonWhitespace.isNumInt() {
		markers = append(markers, line[0])
		line = line[1:]
	}
	switch markers[len(markers)-1].typ {
	case itemAsterisk:
		delim = " "
	case itemHyphen:
		delim = " "
	case itemPlus:
		delim = " "
	case itemCloseParen:
		delim = " "
	case itemDot:
		delim = "."
	}
	startIndentSpaces = spaces + tabs*4
	marker = markers.rawText()
	spaces, tabs, _, firstNonWhitespace = t.nonWhitespace(line)

	w := len(marker)
	n := spaces + tabs*4
	if 4 < n {
		n = 1
	} else if 1 > n {
		n = 1
	}
	wnSpaces := w + n
	indentSpaces = startIndentSpaces + wnSpaces
	if line[0].isTab() {
		line = t.indentOffset(line, 2)
	} else {
		line = line[1:]
	}

	remains = line

	return
}
