// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
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

func (t *Tree) parseList() Node {
	marker := t.next()

	token := t.peek()
	list := &List{
		NodeList, token.pos, t, Children{},
		false,
		1,
		false,
		marker.val,
	}

	loose := false
	for {
		c := t.parseListItem(len(marker.val))
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
		if marker.val != token.val {
			break
		}
	}

	list.Spread = loose

	return list
}

func (t *Tree) parseListItem(indentSpaces int) *ListItem {
	token := t.peek()
	if itemEOF == token.typ {
		return nil
	}

	ret := &ListItem{
		NodeListItem, token.pos, t, Children{},
		false,
		false,
		indentSpaces,
	}
	t.CurNode = ret

	paragraphs := 0
	for {
		c := t.parseBlockContent()
		if nil == c {
			break
		}
		ret.append(c)

		if NodeParagraph == c.Type() || NodeCode == c.Type() {
			paragraphs++
		}

		spaces, tabs, tokens := t.nextNonWhitespace()
		totalSpaces := spaces + tabs*4
		if totalSpaces < indentSpaces {
			t.backups(tokens)
			break
		}
		if totalSpaces == indentSpaces {
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
