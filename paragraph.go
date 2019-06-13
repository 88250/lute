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

type Paragraph struct {
	NodeType
	int
	RawText
	items items
	*Tree
	Subnodes Children

	OpenTag, CloseTag string
}

func (n *Paragraph) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Paragraph) HTML() string {
	content := html(n.Subnodes)

	if "" != n.OpenTag {
		return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag+"\n", content)
	}

	return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag, content)
}

func (n *Paragraph) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Paragraph) Children() Children {
	return n.Subnodes
}

func (n *Paragraph) Tokens() items {
	return n.items
}

func (n *Paragraph) trim() {
	size := len(n.items)
	if 1 > size {
		return
	}

	initialNoneWhitespace := 0
	notBreak := true
	for i := initialNoneWhitespace; i < size/2; i++ {
		if itemNewline == n.items[i].typ {
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
		if itemNewline == n.items[i].typ {
			finalNoneWhitespace--
			notBreak = false
		}
		if notBreak {
			break
		}
	}

	n.items = n.items[initialNoneWhitespace:finalNoneWhitespace]
	n.RawText = RawText(strings.TrimSpace(string(n.RawText)))
}

func (t *Tree) parseParagraph() Node {
	token := t.peek()
	ret := &Paragraph{NodeParagraph, token.pos, "", items{}, t, Children{}, "<p>", "</p>"}
Loop:
	for {
		token = t.nextToken()
		ret.items = append(ret.items, token)
		if itemEOF == token.typ {
			break
		}
		ret.RawText += RawText(token.val)

		if itemNewline == token.typ {
			// 判断是否为段落延续文本

			line := t.nextLineEnding()
			if t.interruptParagrah(line) {
				t.backups(line)

				break
			}

			//indentSpaces := spaces + tabs*4
			//if indentSpaces < t.context.IndentSpaces {
			//	t.backups(tokens)
			//	// TODO(D): 待确定是否还有用
			//	break
			//}
			switch firstNonWhitespace.typ {
			case itemEOF:
				break Loop
			case itemHyphen, itemAsterisk:
				t.backups(tokens[1:])
				break Loop
			default:
				t.backups(tokens)
				continue
			}
		} else {
			t.backups(tokens)
		}
	}

	ret.trim()

	return ret
}

func (t *Tree) interruptParagrah(line []item) bool {
	typ := line[0].typ
	if itemNewline == typ {
		return true
	}

	if itemHyphen == typ || itemAsterisk == typ {
		// TODO: marker 后面还需要空格才能确认是否是列表
		return true
	}

	return false
}
