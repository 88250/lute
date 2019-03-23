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

func (t *Tree) parseInlines() {

	for _, c := range t.Root.Children {
		raw := c.Raw()

		tokens := c.Tokens()
	Block:
		for {
			token := tokens[0]
			var n Node
			switch token.typ {
			case itemStr:
				n, tokens = t.parseText(tokens)
			case itemBacktick:
				n, tokens = t.parseInlineCode(tokens)
			case itemEOF:
				break Block
			}

			c.Append(n)
		}

		fmt.Printf("%s", raw)
	}
}

func (t *Tree) parseText(tokens items) (n Node, remains items) {
	token := tokens[0]
	ret := &Text{NodeText, token.pos, RawText(token.val), items{}, t, token.val}

	return ret, tokens[1:]
}

func (t *Tree) parseInlineCode(tokens items) (ret Node, remains items) {
	i := 1
	token := tokens[i]
	pos := token.pos
	var code string
	for ; itemBacktick != token.typ && itemEOF != token.typ; token = tokens[i] {
		code += token.val
		i++
	}

	ret = &InlineCode{NodeInlineCode, pos, "", items{},t, code}
	remains = tokens[i+1:]

	if itemEOF == t.peek().typ {
		return
	}

	return
}

func (t *Tree) parseCode(tokens items) (ret Node, remains items) {
	i := 1
	token := tokens[i]
	pos := token.pos
	var code string
	for ; itemBacktick != token.typ && itemEOF != token.typ; token = tokens[i] {
		code += token.val
		i++
	}

	ret = &Code{NodeCode, pos, "", items{}, t, code, "", ""}
	remains = tokens[i+1:]

	if itemEOF == t.peek().typ {
		return
	}

	return
}
