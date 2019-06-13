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


func (t *Tree) parseIndentCode() Node {
	ret := &Code{NodeCode, 0, "", items{},t, "", "", ""}
	var code string
Loop:
	for {
		for i := 0; i < 4; {
			token := t.nextToken()
			switch token.typ {
			case itemSpace:
				i++
			case itemTab:
				i += 4
			default:
				break
			}
		}

		token := t.nextToken()
		for ; itemBacktick != token.typ && itemEOF != token.typ; token = t.nextToken() {
			code += token.val
			if itemNewline == token.typ {
				spaces, tabs, tokens, _ := t.nextNonWhitespace()
				if 1 > tabs && 4 > spaces {
					t.backup()
					break Loop
				} else {
					t.backups(tokens)
					continue Loop
				}
			}
			ret.items = append(ret.items, token)
		}
	}

	ret.Value = code
	ret.RawText = RawText(code)

	return ret
}