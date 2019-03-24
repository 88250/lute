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

func (t *Tree) parseBlockquote() Node {
	token := t.peek()
	ret := &Blockquote{NodeBlockquote, token.pos, "", items{}, Children{}}
	t.next() // consume >

	t.nextNonWhitespace()
	t.backup()
	for {
		token = t.peek()
		if itemEOF == token.typ {
			break
		}
		if itemStr == token.typ {
			c := t.parseParagraph()
			ret.Append(c)
			return ret
		}

		ret.RawText += RawText(token.val)
		ret.items = append(ret.items, token)
		if itemNewline == token.typ {
			t.next()
			if token = t.peek(); itemNewline == token.typ || itemEOF == token.typ {
				t.next()
				break
			}
			t.backup()
		}
	}

	return ret
}
