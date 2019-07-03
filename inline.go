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
	t.parseBlockInlines(t.Root.Children())
}

func (t *Tree) parseBlockInlines(blocks Children) {
	for _, block := range blocks {
		cType := block.Type()
		switch cType {
		case NodeCode, NodeInlineCode, NodeThematicBreak:
			continue
		}

		cs := block.Children()
		if 0 < len(cs) {
			t.parseBlockInlines(cs)

			continue
		}

		tokens := block.Tokens()
		for {
			token := tokens[0]
			var n Node
			switch token.typ {
			case itemBacktick:
				n, tokens = t.parseInlineCode(tokens)
			case itemAsterisk, itemUnderscore:
				n = &Text{NodeType: NodeText, Value: token.val}
				tokens = tokens[1:]
			}

			if nil != n {
				block.Append(n)
			}

			if 1 > len(tokens) || tokens.isEOF() {
				break
			}
		}

	}
}

func (t *Tree) parseEmOrStrong(stack *delimiterStack) (ret Node, remains items) {
	return
}

//func (t *Tree) parseText(token *item) (ret Node) {
//	var text string
//	var textTokens items
//	for i := 0; i < len(tokens); i++ {
//		token = tokens[i]
//		if itemHyphen != token.typ && itemEqual != token.typ && itemPlus != token.typ && itemStr != token.typ && itemNewline != token.typ {
//			break
//		}
//		text += token.val
//		textTokens = append(textTokens, token)
//	}
//	ret = &Text{NodeText, RawText(text), textTokens, t, text}
//
//	return
//}

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
