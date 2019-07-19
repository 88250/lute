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

import "strconv"

type ListType int

type List struct {
	*BaseNode

	Bullet bool
	Start  int
	Delim  string
	Tight  bool

	Marker            string
	startIndentSpaces int
	indentSpaces      int
}

func (n *List) Close() {
	if n.close {
		return
	}

	tight := true
	for child := n.FirstChild(); nil != child; child = child.Next() {
		if !child.(*ListItem).Tight {
			tight = false
			break
		}
	}
	n.Tight = tight
	for child := n.FirstChild(); nil != child; child = child.Next() {
		child.(*ListItem).Tight = tight
		child.Close()
	}

	n.close = true
}

func (t *Tree) parseList(line items) (ret Node) {
	indentSpaces := t.context.IndentSpaces

	remains, marker, delim, bullet, start, startIndentSpaces, w, n := t.parseListMarker(line)
	ret = &List{
		&BaseNode{typ: NodeList},
		bullet,
		start,
		delim,
		false,
		marker,
		startIndentSpaces,
		startIndentSpaces + w + n,
	}
	t.context.IndentSpaces += startIndentSpaces + w + n

	if remains.isBlankLine() {
		t.context.IndentSpaces = startIndentSpaces + w + 1
	}

	node := t.parseListItem(line)
	if nil == node {
		return
	}
	ret.AppendChild(ret, node)
	t.context.IndentSpaces = indentSpaces
	t.context.CurNode = ret

	return
}

func (t *Tree) parseListMarker(line items) (remains items, marker, delim string, bullet bool, start, startIndentSpaces, w, n int) {
	spaces, tabs, tokens, firstNonWhitespace := t.nonWhitespace(line)
	var markers items
	markers = append(markers, firstNonWhitespace)
	line = line[len(tokens):]
	bullet = true
	start = 1
	if firstNonWhitespace.isNumInt() {
		bullet = false
		start, _ = strconv.Atoi(firstNonWhitespace.val)
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
		delim = ")"
	case itemDot:
		delim = "."
	}
	startIndentSpaces = spaces + tabs*4
	marker = markers.rawText()
	spaces, tabs, _, firstNonWhitespace = t.nonWhitespace(line)
	w = len(marker)
	n = spaces + tabs*4
	if 4 < n {
		n = 1
	} else if 1 > n {
		n = 1
	}
	if line[0].isTab() {
		line = t.indentOffset(line, 2)
	} else {
		line = line[1:]
	}

	remains = line

	return
}

func (t *Tree) isList(line items) (isList bool, marker string) {
	if 2 > len(line) { // at least marker and newline
		return
	}

	_, line = line.trimLeft()
	if 1 > len(line) {
		return
	}

	firstNonWhitespace := line[0]

	if itemAsterisk == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "*"
		return
	} else if itemHyphen == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "-"
		return
	} else if itemPlus == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "+"
		return
	} else if firstNonWhitespace.isNumInt() && 9 >= len(firstNonWhitespace.val) {
		isList = line[2].isWhitespace()
		if itemDot == line[1].typ {
			marker = firstNonWhitespace.val + "."
		} else if itemCloseParen == line[1].typ {
			marker = firstNonWhitespace.val + ")"
		}
		return
	}

	return
}
