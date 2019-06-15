// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
			case itemStr, itemNewline:
				n, tokens = t.parseText(tokens)
			case itemBacktick:
				n, tokens = t.parseInlineCode(tokens)
			case itemAsterisk:
				n, tokens = t.parseEmOrStrong(tokens)
			default:
				break
			}

			if nil != n {
				c.Append(n)
			}

			if 1 > len(tokens) || tokens.isEOF() {
				break Block
			}
		}
	}
}

func (t *Tree) parseEmphasis(tokens items) (ret Node, remains items) {
	var text string
	var textTokens = items{}
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if itemAsterisk == token.typ {
			remains = tokens[i+1:]
			break
		}
		text += token.val
		textTokens = append(textTokens, token)
	}
	c, _ := t.parseText(textTokens)
	ret = &Emphasis{NodeEmphasis, tokens[0].pos, RawText(text), textTokens, t, Children{}}
	ret.Append(c)

	return
}

func (t *Tree) parseStrong(tokens items) (ret Node, remains items) {
	var text string
	var textTokens = items{}
	for i := 2; i < len(tokens); i++ {
		token := tokens[i]
		if i < len(tokens) - 2 {
		if itemAsterisk == token.typ && itemAsterisk == tokens[i+1].typ {
			remains = tokens[i+2:]
			break
		}
		}
		text += token.val
		textTokens = append(textTokens, token)
	}
	c, _ := t.parseText(textTokens)
	ret = &Strong{NodeStrong, tokens[0].pos, RawText(text), textTokens, t, Children{}}
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
	var textTokens items
	for i := 0; i < len(tokens); i++ {
		token = tokens[i]
		if itemStr != token.typ && itemNewline != token.typ {
			remains = tokens[i:]
			break
		}
		text += token.val
		textTokens = append(textTokens, token)
	}
	ret = &Text{NodeText, token.pos, RawText(text), textTokens, t, text}

	return
}

func (t *Tree) parseInlineCode(tokens []*item) (ret Node, remains items) {
	marker := tokens[0]
	if !t.matchEnd(tokens[1:], marker) {
		marker.typ = itemStr

		return nil, tokens
	}

	var text string
	var textTokens = items{}

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if itemNewline == token.typ {
			text += " "
		} else {
			if itemBacktick == token.typ {
				remains = tokens[i+1:]
				break
			}
			text += token.val
		}
		textTokens = append(textTokens, token)
	}

	ret = &InlineCode{NodeInlineCode, RawText(text), textTokens, t, text}

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

func (t *Tree) matchEnd(tokens items, marker *item) bool {
	for _, token := range tokens {
		if token.typ == marker.typ && token.val == marker.val {
			return true
		}
	}

	return false
}