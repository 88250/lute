// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

import (
	"fmt"
)

type ListType int

const (
	ListTypeBullet  = 0
	ListTypeOrdered = 1
)

type List struct {
	NodeType
	int
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children

	ListType ListType
	Start    int
	Tight    bool

	IndentSpaces int // #4 Indentation https://spec.commonmark.org/0.29/#list-items
	Marker       string
	WNSpaces     int // W + N https://spec.commonmark.org/0.29/#list-items
}

func (n *List) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *List) HTML() string {
	content := html(n.Subnodes)

	if NodeListItem == n.Parent.Type() {
		return fmt.Sprintf("\n<ul>\n%s</ul>\n", content)
	}

	return fmt.Sprintf("<ul>\n%s</ul>\n", content)
}

func (n *List) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *List) Children() Children {
	return n.Subnodes
}

func newList(indentSpaces int, marker string, wnSpaces int, t *Tree, token item) *List {
	ret := &List{
		NodeList, token.pos, "", items{}, t, t.context.CurNode, Children{},
		ListTypeBullet,
		1,
		false,
		indentSpaces,
		marker,
		wnSpaces,
	}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseList(line line) Node {
	spaces, tabs, tokens, firstNonWhitespace := t.nonWhitespace(line)
	marker := firstNonWhitespace
	indentSpaces := spaces + tabs*4
	spaces, tabs, tokens, firstNonWhitespace = t.nonWhitespace(line[len(tokens):])
	w := len(marker.val)
	n := spaces + tabs*4
	wnSpaces := w + n
	if 4 <= n { // rule 2 in https://spec.commonmark.org/0.29/#list-items
		line = indentOffset(tokens, w+1, t)
	} else {
		line = indentOffset(tokens, indentSpaces+wnSpaces, t)
	}
	list := newList(indentSpaces, marker.val, wnSpaces, t, marker)
	tight := false
	for {
		t.context.IndentSpaces = indentSpaces + wnSpaces
		c := t.parseListItem(line)
		if nil == c {
			break
		}
		list.Append(c)

		if c.Tight {
			tight = true
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		spaces, tabs, _, _ := t.nonWhitespace(line)
		indentSpaces := spaces + tabs*4
		if indentSpaces < t.context.IndentSpaces {
			t.backupLine(line)
			break
		}

		line = t.skipWhitespaces(line)
		if marker != line[0]{
			// TODO: 考虑有序列表序号递增
			break
		}
	}

	list.Tight = tight

	return list
}

// https://spec.commonmark.org/0.29/#lists
func (t *Tree) isList(line []item) bool {
	if 2 > len(line) { // at least marker and newline
		return false
	}

	_, marker := t.firstNonSpace(line)
	// TODO: marker 后面还需要空格才能确认是否是列表
	if "*" != marker.val && "-" != marker.val && "+" != marker.val {
		return false
	}

	return true
}
