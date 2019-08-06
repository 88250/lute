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

// delimiter 描述了强调、链接和图片解析过程中用到的分隔符（[, ![, *, _）相关信息。
type delimiter struct {
	node           Node       // the text node point to
	typ            item       // the type of delimiter ([, ![, *, _)
	num            int        // the number of delimiters
	originalNum    int        // the original number of delimiters
	canOpen        bool       // whether the delimiter is a potential opener
	canClose       bool       // whether the delimiter is a potential closer
	previous, next *delimiter // doubly linked list

	active            bool // whether the delimiter is "active" (all are active to start)
	image             bool
	bracketAfter      bool
	index             int
	previousDelimiter *delimiter
}

func (t *Tree) scanDelims(tokens items) *delimiter {
	startPos := t.context.pos
	token := tokens[startPos]
	delimitersCount := 0
	for i := t.context.pos; i < len(tokens) && token == tokens[i]; i++ {
		delimitersCount++
		t.context.pos++
	}

	tokenBefore, tokenAfter := itemNewline, itemNewline
	if 0 != startPos {
		tokenBefore = tokens[startPos-1]
	}
	if len(tokens) > t.context.pos {
		tokenAfter = tokens[t.context.pos]
	}

	var beforeIsPunct, beforeIsWhitespace, afterIsPunct, afterIsWhitespace, canOpen, canClose bool
	if itemEnd != tokenBefore {
		beforeIsWhitespace = tokenBefore.isUnicodeWhitespace()
		beforeIsPunct = tokenBefore.isPunct()
	}
	if itemEnd != tokenAfter {
		afterIsWhitespace = tokenAfter.isUnicodeWhitespace()
		afterIsPunct = tokenAfter.isPunct()
	}

	isLeftFlanking := !afterIsWhitespace && (!afterIsPunct || beforeIsWhitespace || beforeIsPunct)
	isRightFlanking := !beforeIsWhitespace && (!beforeIsPunct || afterIsWhitespace || afterIsPunct)
	if itemUnderscore == token {
		canOpen = isLeftFlanking && (!isRightFlanking || beforeIsPunct)
		canClose = isRightFlanking && (!isLeftFlanking || afterIsPunct)
	} else {
		canOpen = isLeftFlanking
		canClose = isRightFlanking
	}

	return &delimiter{typ: token, num: delimitersCount, active: true, canOpen: canOpen, canClose: canClose}
}

func (t *Tree) handleDelim(block Node, tokens items) {
	startPos := t.context.pos
	delim := t.scanDelims(tokens)

	_, text := t.extractTokens(tokens, startPos, t.context.pos)
	node := &Text{&BaseNode{typ: NodeText, value: text}}
	block.AppendChild(block, node)

	// Add entry to stack for this opener
	t.context.delimiters = &delimiter{
		typ:         delim.typ,
		num:         delim.num,
		originalNum: delim.num,
		node:        node,
		previous:    t.context.delimiters,
		next:        nil,
		canOpen:     delim.canOpen,
		canClose:    delim.canClose,
	}
	if t.context.delimiters.previous != nil {
		t.context.delimiters.previous.next = t.context.delimiters
	}
}

func (t *Tree) processEmphasis(stackBottom *delimiter) {
	var opener, closer, old_closer *delimiter
	var opener_inl, closer_inl Node
	var tempstack *delimiter
	var use_delims int
	var opener_found bool
	var openers_bottom = map[item]*delimiter{}
	var odd_match = false

	openers_bottom[itemUnderscore] = stackBottom
	openers_bottom[itemAsterisk] = stackBottom

	// find first closer above stack_bottom:
	closer = t.context.delimiters
	for closer != nil && closer.previous != stackBottom {
		closer = closer.previous
	}

	// move forward, looking for closers, and handling each
	for nil != closer {
		var closercc = closer.typ
		if !closer.canClose {
			closer = closer.next
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

			text := opener_inl.Value()[0 : len(opener_inl.Value())-use_delims]
			opener_inl.SetValue(text)
			text = closer_inl.Value()[0 : len(closer_inl.Value())-use_delims]
			closer_inl.SetValue(text)

			// build contents for new emph element
			var emph Node
			if 1 == use_delims {
				emph = &Emphasis{&BaseNode{typ: NodeEmphasis}}
			} else {
				emph = &Strong{&BaseNode{typ: NodeStrong}}
			}

			tmp := opener_inl.Next()
			for nil != tmp && tmp != closer_inl {
				next := tmp.Next()
				tmp.Unlink()
				emph.AppendChild(emph, tmp)
				tmp = next
			}

			opener_inl.InsertAfter(opener_inl, emph)

			// remove elts between opener and closer in delimiters stack
			if opener.next != closer {
				opener.next = closer
				closer.previous = opener
			}

			// if opener has 0 delims, remove it and the inline
			if opener.num == 0 {
				opener_inl.Unlink()
				t.removeDelimiter(opener)
			}

			if closer.num == 0 {
				closer_inl.Unlink()
				tempstack = closer.next
				t.removeDelimiter(closer)
				closer = tempstack
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
				t.removeDelimiter(old_closer)
			}
		}
	}

	// remove all delimiters
	for t.context.delimiters != nil && t.context.delimiters != stackBottom {
		t.removeDelimiter(t.context.delimiters)
	}
}

func (t *Tree) removeDelimiter(delim *delimiter) (ret *delimiter) {
	if delim.previous != nil {
		delim.previous.next = delim.next
	}
	if delim.next == nil {
		// top of stack
		t.context.delimiters = delim.previous
	} else {
		delim.next.previous = delim.previous
	}

	return
}
