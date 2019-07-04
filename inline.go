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
	delimiters := &delimiter{}
	t.parseBlockInlines(t.Root.Children(), delimiters)

	t.parseEmphasis(nil, delimiters)
}

func (t *Tree) parseBlockInlines(blocks []*Node, delimiters *delimiter) {
	for _, block := range blocks {
		cType := block.NodeType
		switch cType {
		case NodeCode, NodeInlineCode, NodeThematicBreak:
			continue
		}

		cs := block.Children()
		if 0 < len(cs) {
			t.parseBlockInlines(cs, delimiters)

			continue
		}

		tokens := block.Tokens
		pos := 0

		for {
			token := tokens[pos]
			var n *Node
			switch token.typ {
			case itemBacktick:
				n = t.parseInlineCode(tokens, &pos)
			case itemAsterisk, itemUnderscore:
				n = t.parseDelimiter(tokens, &pos, delimiters)
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

func (t *Tree) parseEmphasis(stackBottom *delimiter, delimiters *delimiter) {
	var opener, closer, old_closer *delimiter
	var opener_inl, closer_inl *Node
	var tempstack *delimiter
	var use_delims int
	var tmp, next *delimiter
	var opener_found bool
	var openers_bottom = map[itemType]*delimiter{}
	var odd_match = false

	openers_bottom[itemUnderscore] = stackBottom
	openers_bottom[itemAsterisk] = stackBottom

	// find first closer above stack_bottom:
	closer = delimiters
	for closer != nil && closer.previous != stackBottom {
		closer = closer.previous
	}

	// move forward, looking for closers, and handling each
	for closer != nil {
		var closercc = closer.typ
		if !closer.canClose {
			continue
		}

		// found emphasis closer. now look back for first matching opener:
		opener = closer.previous
		opener_found = false
		for nil != opener && opener != stackBottom && opener != openers_bottom[closercc] {
			odd_match = (closer.canOpen || opener.canClose) && closer.originalNum%3 != 0 && (opener.originalNum+closer.originalNum)%3 == 0
			if opener.typ == closer.typ && opener.canOpen && !odd_match {
				opener_found = true
				break
			}
			opener = opener.previous
		}
		old_closer = closer

		if itemAsterisk == closercc || itemUnderscore == closercc {
			if !opener_found {
				closer = closer.next
			} else {
				// calculate actual number of delimiters used from closer
				if closer.num >= 2 && opener.num >= 2 {
					use_delims = 2
				} else {
					use_delims = 1
				}

				opener_inl = opener.node
				closer_inl = closer.node

				// remove used delimiters from stack elts and inlines
				opener.num -= use_delims
				closer.num -= use_delims

				text := opener_inl.RawText[0 : len(opener_inl.RawText)-use_delims]
				opener_inl.RawText = text

				text = closer_inl.RawText[0 : len(closer_inl.RawText)-use_delims]
				closer_inl.RawText = text

				// build contents for new emph element
				var emph *Node
				if 1 == use_delims {
					emph = &Node{NodeType: NodeEmphasis}
					_ = &Emphasis{Node: emph}
				} else {
					emph = &Node{NodeType: NodeStrong}
					_ = &Strong{Node: emph}
				}

				tmp.node = opener_inl.Next
				for nil != tmp && tmp.node != closer_inl {
					next = tmp.next
					tmp.node.Unlink()
					emph.Append(tmp.node)
					tmp = next
				}

				opener_inl.InsertAfter(emph)

				// remove elts between opener and closer in delimiters stack
				if opener.next != closer {
					opener.next = closer
					closer.previous = opener
				}

				// if opener has 0 delims, remove it and the inline
				if opener.num == 0 {
					opener_inl.Unlink()
					delimiters.remove(opener)
				}

				if closer.num == 0 {
					closer_inl.Unlink()
					tempstack = closer.next
					delimiters.remove(closer)
					closer = tempstack
				}
			}
		}
		if !opener_found && !odd_match {
			// Set lower bound for future searches for openers:
			// We don't do this with odd_match because a **
			// that doesn't match an earlier * might turn into
			// an opener, and the * might be matched by something
			// else.
			openers_bottom[closercc] = old_closer.previous
			if !old_closer.canOpen {
				// We can remove a closer that can't be an opener,
				// once we've seen there's no matching opener:
				delimiters.remove(old_closer)
			}
		}
	}

	// remove all delimiters
	for delimiters != nil && delimiters != stackBottom {
		delimiters.remove(delimiters)
	}
}

func (t *Tree) parseDelimiter(tokens items, pos *int, delimiters *delimiter) (ret *Node) {
	startPos := *pos
	delim := t.scanDelimiter(tokens, pos)

	subTokens, text := t.extractTokens(tokens, startPos, *pos)
	ret = &Node{NodeType: NodeText, RawText: text}
	_ = &Text{ret, subTokens, t, text}
	delim.node = ret

	// Add entry to stack for this opener
	delimiters = delim
	if delimiters.previous != nil {
		delimiters.previous.next = delimiters
	}

	return
}

func (t *Tree) extractTokens(tokens items, startPos, endPos int) (subTokens items, text string) {
	for i := startPos; i < endPos; i++ {
		text += tokens[i].val
		subTokens = append(subTokens, tokens[i])
	}

	return
}

func (t *Tree) scanDelimiter(tokens items, pos *int) *delimiter {
	startPos := *pos
	token := tokens[startPos]
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
	index := startPos - 1
	if 0 < index {
		tokenBefore = tokens[index]
	}
	index = startPos + 1
	if len(tokens) > index {
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

	return &delimiter{typ: token.typ, num: delimitersCount, active: true, canOpen: canOpen, canClose: canClose}
}

func (t *Tree) parseInlineCode(tokens items, pos *int) (ret *Node) {
	marker := tokens[0]
	if !t.matchEnd(tokens[1:], marker) {
		marker.typ = itemStr
		*pos++

		ret = &Node{NodeType: NodeText, RawText: marker.val}
		_ = &Text{ret, nil, t, marker.val}

		return
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

	ret = &Node{NodeType: NodeInlineCode, RawText: text}
	_ = &InlineCode{ret, textTokens, t, text}

	return
}

func (t *Tree) parseText(tokens items, pos *int) (ret *Node) {
	token := tokens[*pos]
	*pos++

	ret = &Node{NodeType: NodeText, RawText: token.val}
	_ = &Text{ret, nil, t, token.val}

	return
}

func (t *Tree) parseCode(tokens items) (ret *Node, remains items) {
	i := 1
	token := tokens[i]
	pos := token.pos
	var code string
	for ; itemBacktick != token.typ && itemEOF != token.typ; token = tokens[i] {
		code += token.val
		i++
	}

	ret = &Node{NodeType: NodeCode}
	_ = &Code{ret, pos, items{}, t, code, "", ""}
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
