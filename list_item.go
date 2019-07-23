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

type ListItem struct {
	*BaseNode

	Bullet bool
	Start  int
	Delim  string
	Tight  bool

	Marker            string
	StartIndentSpaces int
	IndentSpaces      int
}

func (n *ListItem) Close() {
	if n.close {
		return
	}

	for child := n.FirstChild(); nil != child; child = child.Next() {
		child.Close()
	}
}

func (t *Tree) parseListItem(tokens items) (ret Node) {
	spaces, tabs, remains := t.nonSpaceTab(tokens)
	startIndentSpaces := spaces + tabs*4
	token := remains[0]
	start := 0
	var marker, delim string
	var bullet bool
	if itemAsterisk == token.typ {
		if !remains[1].isWhitespace() {
			return
		}
		marker = "*"
		delim = " "
		bullet = true
		remains = remains[2:]
	} else if itemHyphen == token.typ {
		if !remains[1].isWhitespace() {
			return
		}
		marker = "-"
		delim = " "
		bullet = true
		remains = remains[2:]
	} else if itemPlus == token.typ {
		if !tokens[1].isWhitespace() {
			return
		}
		marker = "+"
		delim = " "
		bullet = true
		remains = remains[2:]
	} else if token.isNumInt() && 9 >= len(token.val) {
		if !remains[2].isWhitespace() {
			return
		}
		start, _ = strconv.Atoi(token.val)
		if itemDot == remains[1].typ {
			delim = "."
			marker = token.val + delim
			remains = remains[2:]
		} else if itemCloseParen == remains[1].typ {
			delim = ")"
			marker = token.val + delim
			remains = remains[2:]
		} else {
			return
		}
	} else {
		return
	}

	spaces, tabs, remains = t.nonSpaceTab(remains)
	w := len(marker)
	n := spaces + tabs*4
	if 4 < n {
		n = 1
	} else if 1 > n {
		n = 1
	}
	wnSpaces := w + n
	indentSpaces := startIndentSpaces + wnSpaces

	li := &ListItem{
		&BaseNode{typ: NodeListItem, tokens: items{}},
		bullet,
		start,
		delim,
		true,
		marker,
		startIndentSpaces,
		indentSpaces,
	}
	ret = li

	child := t.parseBlock(remains)
	child.SetLeftSpaces(indentSpaces)
	li.AppendChild(li, child)

	return
}
