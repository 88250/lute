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

		line := c.Tokens()
	Block:
		for {
			token := line[0]
			var n Node
			switch token.typ {
			case itemStr, itemNewline:
				n, line = t.parseText(line)
			case itemBacktick:
				n, line = t.parseInlineCode(line)
			case itemAsterisk:
				n, line = t.parseEmphasis(line)
			default:
				break
			}

			if nil != n {
				c.Append(n)
			}

			if 1 > len(line) || line.isEOF(){
				break Block
			}
		}
	}
}

func (t *Tree) parseEmphasis(tokens items) (ret Node, remains items) {
	token := tokens[0]

	rawText := RawText(token.val)
	ret = &Emphasis{NodeEmphasis, token.pos, rawText, tokens, t, Children{}}
	c, remains := t.parseText(tokens)
	ret.Append(c)

	return
}

func (t *Tree) parseStrong(tokens items) (ret Node, remains items) {
	token := tokens[0]

	rawText := RawText(token.val)
	ret = &Strong{NodeStrong, token.pos, rawText, tokens, t, Children{}}
	c, remains := t.parseText(tokens)
	ret.Append(c)

	return
}

func (t *Tree) parseEmOrStrong(tokens items) (ret Node, remains items) {
	if itemAsterisk == tokens[0].typ && itemAsterisk == tokens[1].typ {
		ret, remains = t.parseStrong(tokens)
	} else {
		ret, remains = t.parseEmphasis(tokens)
	}

	return
}

func (t *Tree) parseText(tokens items) (ret Node, remains items) {
	token := tokens[0]
	var text string
	for i := 0; i < len(tokens); i++ {
		token = tokens[i]
		if itemAsterisk == token.typ {
			remains = tokens[i:]
			break
		}
		text += token.val
	}
	ret = &Text{NodeText, token.pos, RawText(text), items{}, t, text}

	return
}

func (t *Tree) parseInlineCode(tokens []*item) (ret Node, remains items) {
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

	return
}
