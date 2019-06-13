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

func (t *Tree) parseIndentCode(line []item) Node {
	ret := &Code{NodeCode, 0, "", nil, t, "", "", ""}
	var code string
Loop:
	for {
		var spaces, tabs int
		for i := 0; i < 4; i++ {
			token := line[i]
			if itemSpace == token.typ {
				spaces++
			} else if itemTab == token.typ {
				tabs++
			}
			if 3 < spaces || 0 < tabs {
				line = line[i+1:]
				break
			}
		}

		token := line[0]
		for i := 0; itemBacktick != token.typ && itemEOF != token.typ; i++ {
			token := line[i]
			code += token.val
			if itemNewline == token.typ {
				line = t.nextLineEnding()
				spaces, tabs, tokens, _ := t.nonWhitespace(line)
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

// https://spec.commonmark.org/0.29/#indented-code-blocks
func (t *Tree) isIndentCode(line []item) bool {
	var tabs, spaces int
	for _, token := range line {
		if itemSpace == token.typ {
			spaces++
			continue
		}
		if itemTab == token.typ {
			tabs++
			continue
		}

		break
	}

	return 0 < tabs || 3 < spaces
}
