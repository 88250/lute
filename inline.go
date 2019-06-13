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

func (t *Tree) parseInlines() {
	t.parseChildren(t.Root.Children())
}

func (t *Tree) parseChildren(children Children) {
	for _, c := range children {
		cType := c.Type()
		switch cType {
		case NodeCode, NodeInlineCode, NodeThematicBreak:
			continue
		}

		cs := c.Children()
		if 0 < len(cs) {
			t.parseChildren(cs)

			continue
		}

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
			case itemAsterisk:
				n, tokens = t.parseEmOrStrong(tokens)
			case itemEOF:
				break Block
			}

			c.Append(n)

			if 1 > len(tokens) {
				break Block
			}
		}
	}
}

func (t *Tree) parseEmphasis(tokens items) (ret Node, remains items) {
	token := tokens[0]

	rawText := RawText(token.val)
	ret = &Emphasis{NodeEmphasis, token.pos, rawText, items{}, t, Children{}}
	c, remains := t.parseText(tokens)
	ret.Append(c)

	remains = remains[1:]

	return
}

func (t *Tree) parseStrong(tokens items) (ret Node, remains items) {
	token := tokens[0]

	rawText := RawText(token.val)
	ret = &Strong{NodeStrong, token.pos, rawText, items{}, t, Children{}}
	c, remains := t.parseText(tokens)
	ret.Append(c)

	remains = remains[2:]

	return
}

func (t *Tree) parseEmOrStrong(tokens items) (ret Node, remains items) {
	if 3 > len(tokens) {
		return t.parseText(tokens)
	}

	tokens = tokens[1:]
	token := tokens[0]
	if itemAsterisk == token.typ {
		tokens = tokens[1:]
		ret, remains = t.parseStrong(tokens)
	} else {
		ret, remains = t.parseEmphasis(tokens)
	}

	return
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
		if itemNewline == token.typ {
			code += " "
		} else {
			code += token.val
		}
		i++
	}

	ret = &InlineCode{NodeInlineCode, pos, "", items{}, t, code}
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
