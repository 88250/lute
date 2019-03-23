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
	t        *Tree
	Parent   Node
	Children Children

	ListType ListType
	Start    int
	Tight    bool

	IdentSpaces int
	Marker      string
}

func (n *List) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *List) HTML() string {
	content := html(n.Children)

	if NodeListItem == n.Parent.Type() {
		return fmt.Sprintf("\n<ul>\n%s</ul>\n", content)
	}

	return fmt.Sprintf("<ul>\n%s</ul>\n", content)
}

func (n *List) append(c Node) {
	n.Children = append(n.Children, c)
}

func newList(marker string, indentSpaces int, t *Tree, token item) *List {
	ret := &List{
		NodeList, token.pos, t, t.context.CurNode, Children{},
		ListTypeBullet,
		1,
		false,
		indentSpaces,
		marker,
	}
	t.context.CurNode = ret

	return ret
}

type ListItem struct {
	NodeType
	Pos
	t        *Tree
	Parent   Node
	Children Children

	Checked bool
	Spread  bool // loose or tight

	Spaces int
}

func (n *ListItem) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *ListItem) HTML() string {
	var content string
	for _, c := range n.Children {
		if !n.Spread && NodeParagraph == c.Type() {
			p := c.(*Paragraph)
			p.OpenTag, p.CloseTag = "", ""
		}

		content += c.HTML()
	}

	if strings.Contains(content, "<ul>") {
		return fmt.Sprintf("<li>%s</li>\n", content)
	}

	if 1 < len(n.Children) || strings.Contains(content, "<pre><code") {
		return fmt.Sprintf("<li>\n%s</li>\n", content)
	}

	return fmt.Sprintf("<li>%s</li>\n", content)
}

func (n *ListItem) append(c Node) {
	n.Children = append(n.Children, c)
}

func newListItem(indentSpaces int, t *Tree, token item) *ListItem {
	ret := &ListItem{
		NodeListItem, token.pos, t, t.context.CurNode, Children{},
		false,
		false,
		indentSpaces,
	}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseList() Node {
	spaces, tabs, tokens := t.nextNonWhitespace()

	marker := tokens[len(tokens)-1].val
	token := t.peek()
	if !token.isWhitespace() {
		t.backup()
		return t.parseEmOrStrong()
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

	offsetSpaces := t.expandSpaces()
	for i := 0; i < offsetSpaces && i < 5; i++ {
		t.next()
	}

	indentSpaces := spaces + tabs*4
	list := newList(marker, indentSpaces, t, token)

	loose := false
	for {
		t.context.IndentSpaces = indentSpaces + len(marker) + 1
		c := t.parseListItem()
		if nil == c {
			break
		}
		list.append(c)

		if c.Spread {
			loose = true
		}

		token := t.peek()
		if itemNewline == token.typ {
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
		c := t.parseParagraph()
		if nil == c {
			break
		}
		ret.append(c)

		if NodeParagraph == c.Type() || NodeCode == c.Type() {
			paragraphs++
		}

		if itemEOF == t.peek().typ {
			break
		}

		spaces, tabs, tokens := t.nextNonWhitespace()

		totalSpaces := spaces + tabs*4
		if totalSpaces < indentSpaces {
			t.backups(tokens)
			break
		} else if totalSpaces == indentSpaces {
			t.backup()
			continue
		}

		indentOffset(tokens, indentSpaces, t)
	}

	if 1 < paragraphs {
		ret.Spread = true
	}

	return ret
}
