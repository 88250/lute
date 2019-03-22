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

import "fmt"

type Paragraph struct {
	NodeType
	Pos
	*Tree
	Children

	OpenTag, CloseTag string
}

func (n *Paragraph) String() string {
	return fmt.Sprintf("%s", n.Children)
}

func (n *Paragraph) HTML() string {
	content := html(n.Children)

	if "" != n.OpenTag {
		return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag+"\n", content)
	}

	return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag, content)
}

func (n *Paragraph) append(c Node) {
	n.Children = append(n.Children, c)
}

func (n *Paragraph) trim() {
	size := len(n.Children)
	if 1 > size {
		return
	}

	initialNoneWhitespace := 0
	notBreak := true
	for i := initialNoneWhitespace; i < size/2; i++ {
		if NodeBreak == n.Children[i].Type() {
			initialNoneWhitespace++
			notBreak = false
		}
		if notBreak {
			break
		}
	}

	finalNoneWhitespace := size
	notBreak = true
	for i := finalNoneWhitespace - 1; size/2 <= i; i-- {
		if NodeBreak == n.Children[i].Type() {
			finalNoneWhitespace--
			notBreak = false
		}
		if notBreak {
			break
		}
	}

	n.Children = n.Children[initialNoneWhitespace:finalNoneWhitespace]
}

func (t *Tree) parseParagraph() Node {
	token := t.peek()

	ret := &Paragraph{NodeParagraph, token.pos, t, Children{}, "<p>", "</p>"}

	for {
		c := t.parsePhrasingContent()
		if nil == c {
			ret.trim()

			break
		}
		ret.append(c)

		if token = t.peek(); itemNewline == token.typ {
			t.next()
			if token = t.peek();itemNewline == token.typ || itemEOF == token.typ {
				t.next()
				break
			} else{
				_, _, tokens := t.nextNonWhitespace()
				last := tokens[len(tokens) - 1]
				if itemHyphen == last.typ {
					t.backups(tokens)
					break
				}

				continue
			}

			t.backup()
		}
	}

	return ret
}