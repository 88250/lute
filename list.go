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
	"strings"
)

type ListType int

const (
	ListTypeBullet  = 0
	ListTypeOrdered = 1
)

type List struct {
	NodeType
	Pos
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children

	ListType ListType
	Start    int
	Tight    bool

	IndentSpaces int // #4 Indentation https://spec.commonmark.org/0.28/#list-items
	Marker       string
	WNSpaces     int // W + N https://spec.commonmark.org/0.28/#list-items
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

type ListItem struct {
	NodeType
	Pos
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children

	Checked bool
	Spread  bool // loose or tight

	Spaces int
}

func (n *ListItem) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *ListItem) HTML() string {
	var content string
	for _, c := range n.Subnodes {
		if !n.Spread && NodeParagraph == c.Type() {
			p := c.(*Paragraph)
			p.OpenTag, p.CloseTag = "", ""
		}

		content += c.HTML()
	}

	if strings.Contains(content, "<ul>") {
		return fmt.Sprintf("<li>%s</li>\n", content)
	}

	if 1 < len(n.Subnodes) || strings.Contains(content, "<pre><code") {
		return fmt.Sprintf("<li>\n%s</li>\n", content)
	}

	return fmt.Sprintf("<li>%s</li>\n", content)
}

func (n *ListItem) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *ListItem) Children() Children {
	return n.Subnodes
}

func newListItem(indentSpaces int, t *Tree, token item) *ListItem {
	ret := &ListItem{
		NodeListItem, token.pos, "", items{}, t, t.context.CurNode, Children{},
		false,
		false,
		indentSpaces,
	}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseList() Node {
	spaces, tabs, tokens, firstNonWhitespace := t.nextNonWhitespace()
	marker := firstNonWhitespace.val
	token := t.peek()
	if !token.isWhitespace() {
		t.backup()
		return t.parseParagraph()
	}

	// Thematic breaks
	var backupTokens []item
	chars := 0
	for {
		token = t.next()
		backupTokens = append(backupTokens, token)
		if itemNewline == token.typ || itemEOF == token.typ {
			chars++
			break
		}
		if token.val != marker && itemTab != token.typ && itemSpace != token.typ {
			break
		}
		chars++
	}
	if chars == len(backupTokens) {
		return t.parseThematicBreak()
	}
	t.backups(backupTokens)

	indentSpaces := spaces + tabs*4
	spaces, tabs, tokens, firstNonWhitespace = t.nextNonWhitespace()
	w := len(marker)
	n := spaces + tabs*4
	wnSpaces := w + n
	t.backups(tokens)
	if 4 <= n { // rule 2 in https://spec.commonmark.org/0.28/#list-items
		indentOffset(tokens, w + 1, t)
	} else {
		indentOffset(tokens, indentSpaces+wnSpaces, t)
	}
	list := newList(indentSpaces, marker, wnSpaces, t, token)
	loose := false
	for {
		t.context.IndentSpaces = indentSpaces + wnSpaces
		c := t.parseListItem()
		if nil == c {
			break
		}
		list.Append(c)

		if c.Spread {
			loose = true
		}

		token := t.peek()
		if itemNewline == token.typ {
			spaces, tabs, tokens, firstNonWhitespace := t.nextNonWhitespace()
			indentSpaces := spaces + tabs*4
			if indentSpaces < t.context.IndentSpaces {
				t.backups(tokens)
				break
			}

			t.next()
			continue
		}
		if marker != token.val {
			break
		}
	}

	list.Tight = loose

	return list
}

func (t *Tree) parseListItem() *ListItem {
	token := t.peek()
	if itemEOF == token.typ {
		return nil
	}

	indentSpaces := t.context.IndentSpaces
	ret := newListItem(indentSpaces, t, token)
	paragraphs := 0
	for {
		c := t.parseBlock()
		if NodeParagraph == c.Type() || NodeCode == c.Type() {
			paragraphs++
		}

		if itemEOF == t.peek().typ {
			break
		}

		spaces, tabs, tokens, firstNonWhitespace := t.nextNonWhitespace()
		if "-" != firstNonWhitespace.val {
			// break
		}
		t.backups(tokens)
		totalSpaces := spaces + tabs*4
		if totalSpaces > indentSpaces {
			if 4 == totalSpaces && 2 != indentSpaces{ // 对齐列表优先级高于缩进代码块
				break
			}
		}

		if totalSpaces < indentSpaces {
			break
		}

		indentOffset(tokens, indentSpaces, t)

		if itemHyphen == firstNonWhitespace.typ {
			t.backups(tokens)
		}

	}

	if 1 < paragraphs {
		ret.Spread = true
	}

	return ret
}
