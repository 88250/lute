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

type ListItem struct {
	NodeType
	int
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children

	Checked bool
	Tight   bool

	Spaces int
}

func (n *ListItem) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *ListItem) HTML() string {
	var content string
	for _, c := range n.Subnodes {
		if n.Tight && NodeParagraph == c.Type() {
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

func newListItem(indentSpaces int, t *Tree, token *item) *ListItem {
	ret := &ListItem{
		NodeListItem, token.pos, "", items{}, t, t.context.CurNode, Children{},
		false,
		true,
		indentSpaces,
	}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseListItem(line items) *ListItem {
	indentSpaces := t.context.IndentSpaces
	ret := newListItem(indentSpaces, t, line[0])
	blankLineBetweenBlocks := false
	for {
		c := t.parseBlock(line)
		if nil == c {
			continue
		}

		blankLines := t.skipBlankLines()
		if 1 <= blankLines && !blankLineBetweenBlocks {
			blankLineBetweenBlocks = true
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		spaces, tabs, _, _ := t.nonWhitespace(line)
		totalSpaces := spaces + tabs*4
		if totalSpaces < indentSpaces {
			t.backupLine(line)

			break
		}

		line = indentOffset(line, indentSpaces, t)
	}

	if 1 < len(ret.Subnodes) && blankLineBetweenBlocks {
		ret.Tight = false
	}

	return ret
}
