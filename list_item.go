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

func newListItem(t *Tree) (ret Node) {
	baseNode := &BaseNode{typ: NodeListItem, tokens: items{}}
	ret = &ListItem{
		baseNode,
		false,
		true,
	}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseListItem(line items) (ret Node, withBlankLine bool) {
	ret = newListItem(t)
	for {
		n := t.parseBlock(line)
		if nil == n {
			line = t.nextLine()
			continue
		}
		ret.AppendChild(ret, n)

		blankLines := t.skipBlankLines()
		if 1 <= blankLines && !withBlankLine {
			withBlankLine = true
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		if 0 < t.blockquoteMarkerCount(line) && 0 < t.context.BlockquoteLevel {
			line = t.removeStartBlockquoteMarker(line, t.context.BlockquoteLevel)
		}

		spaces, tabs, _, _ := t.nonWhitespace(line)
		totalSpaces := spaces + tabs*4
		if totalSpaces < t.context.IndentSpaces {
			t.backupLine(line)
			break
		}

		line = t.indentOffset(line, t.context.IndentSpaces)
	}

	if 1 < len(ret.Children()) && withBlankLine {
		ret.(*ListItem).Tight = false
	}

	return
}

func (t *Tree) parseListItemMarker(line items, list Node) (remains items, marker, delim string, startIndentSpaces, indentSpaces int) {
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

	if line.isBlankLine() {
		line = t.nextLine()
		if line.isBlankLine() {
			list.AppendChild(list, &ListItem{BaseNode: &BaseNode{typ: NodeListItem}, Tight: true})

			return
		}

		if isList, marker := t.isList(line); isList {
			list.AppendChild(list, &ListItem{BaseNode: &BaseNode{typ: NodeListItem}, Tight: true})
			line = line[len(marker):]
		}

		line = t.indentOffset(line, t.context.IndentSpaces)
	}

	remains = line

	return
}
