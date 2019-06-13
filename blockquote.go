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

type Blockquote struct {
	NodeType
	int
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children
}

func (n *Blockquote) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Blockquote) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<blockquote>\n%s</blockquote>\n", content)
}

func (n *Blockquote) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Blockquote) Children() Children {
	return n.Subnodes
}

func newBlockquote(t *Tree, token item) *Blockquote {
	ret := &Blockquote{
		NodeBlockquote, token.pos, "", items{}, t, t.context.CurNode, Children{}}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseBlockquote(line line) Node {
	token := line[0]
	indentSpaces := t.context.IndentSpaces + 2

	ret := newBlockquote(t, token)
	_, _, tokens, _ := t.nonWhitespace(line[1:])
	line = indentOffset(tokens, indentSpaces, t)
	for {
		c := t.parseBlock(line)
		if nil == c {
			break
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		//spaces, tabs, tokens, _ := t.nonWhitespace(line)
		//
		//totalSpaces := spaces + tabs*4
		//if totalSpaces < indentSpaces {
		//	t.backups(tokens)
		//	break
		//} else if totalSpaces == indentSpaces {
		//	t.backup()
		//	continue
		//}
		//
		//indentOffset(tokens, indentSpaces, t)
	}

	return ret
}

// https://spec.commonmark.org/0.29/#block-quotes
func (t *Tree) isBlockquote(line []item) bool {
	if 2 > len(line) { // at least > and newline
		return false
	}

	_, marker := t.firstNonSpace(line)
	if ">" != marker.val {
		return false
	}

	return true
}
