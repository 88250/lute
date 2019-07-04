// Lute - A structured markdown engine.
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
		pos := 0
		stack := &delimiterStack{}
		for {
			token := tokens[pos]
			var n Node
			switch token.typ {
			case itemBacktick:
				n = t.parseInlineCode(tokens, &pos)
			case itemAsterisk, itemUnderscore:
				n = t.parseDelimiter(tokens, &pos, stack)
			case itemStr:
				n = t.parseText(tokens, &pos)
			}

			if nil != n {
				block.Append(n)
			}

			if 1 > len(tokens) || tokens[pos].isEOF() {
				break
			}
		}

	}
}

func (t *Tree) parseDelimiter(tokens items, pos *int, stack *delimiterStack) (ret Node) {
	startPos := *pos
	delim := t.scanDelimiter(tokens, pos)
	stack.push(delim)

	subTokens, text := t.extractTokens(tokens, startPos, *pos)
	ret = &Text{NodeText, RawText(text), subTokens, t, text}

	return
}

func (t *Tree) extractTokens(tokens items, startPos, endPos int) (subTokens items, text string) {
	for i:=startPos;i<endPos;i++ {
		text+=tokens[i].val
		subTokens = append(subTokens, tokens[i])
	}

	return
}

func (t *Tree) scanDelimiter(tokens items, pos *int) *delimiter {
	token := tokens[*pos]
	delimitersCount := 0
	for i := *pos; i < len(tokens); i++ {
		if token.val == tokens[i].val {
			delimitersCount++
			*pos++
		} else {
			break
		}
	}

	var tokenBefore, tokenAfter *item
	index := *pos - 1
	if 0 < index {
		tokenBefore = tokens[index]
	}
	index = *pos + 1
	if len(tokens) < index {
		tokenAfter = tokens[index]
	}

	var beforeIsPunct, beforeIsWhitespace, afterIsPunct, afterIsWhitespace, canOpen, canClose bool
	if nil != tokenBefore {
		beforeIsWhitespace = tokenBefore.isWhitespace()
		beforeIsPunct = tokenBefore.isPunct()
	}
	if nil != tokenAfter {
		afterIsWhitespace = tokenAfter.isWhitespace()
		afterIsPunct = tokenAfter.isPunct()
	}

	isLeftFlanking := !afterIsWhitespace && (!afterIsPunct || beforeIsWhitespace || beforeIsPunct)
	isRightFlanking := !beforeIsWhitespace && (!beforeIsPunct || afterIsWhitespace || afterIsPunct)
	if itemUnderscore == token.typ {
		canOpen = isLeftFlanking && (!isRightFlanking || beforeIsPunct)
		canClose = isRightFlanking && (!isLeftFlanking || afterIsPunct)
	} else {
		canOpen = isLeftFlanking
		canClose = isRightFlanking
	}

	return &delimiter{typ: token.val, num: delimitersCount, active: true, canOpen: canOpen, canClose: canClose}
}

func (t *Tree) parseInlineCode(tokens items, pos *int) (ret Node) {
	marker := tokens[0]
	if !t.matchEnd(tokens[1:], marker) {
		marker.typ = itemStr
		*pos++

		return &Text{NodeText, RawText(marker.val), nil, t, marker.val}
	}

	var text string
	var textTokens = items{}

	for i := *pos; i < len(tokens); i++ {
		token := tokens[i]
		if itemNewline == token.typ {
			text += " "
		} else {
			if itemBacktick == token.typ {
				*pos = i
				break
			}
			text += token.val
		}
		textTokens = append(textTokens, token)
	}

	ret = &InlineCode{NodeInlineCode, RawText(text), textTokens, t, text}

	return
}

func (t *Tree) parseText(tokens items, pos *int) (ret Node) {
	token := tokens[*pos]
	*pos++

	return &Text{NodeText, RawText(token.val), nil, t, token.val}
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
