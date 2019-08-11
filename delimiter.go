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
	typ            byte       // the type of delimiter ([, ![, *, _)
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

// handleDelim 将分隔符 *_~ 入栈。
func (t *Tree) handleDelim(block Node, tokens items) {
	startPos := t.context.pos
	delim := t.scanDelims(tokens)

	text := tokens[startPos:t.context.pos]
	node := &Text{tokens: text}
	block.AppendChild(block, node)

	// 将这个分隔符入栈
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

// processEmphasis 处理强调、加粗以及删除线。
func (t *Tree) processEmphasis(stackBottom *delimiter) {
	var opener, closer, oldCloser *delimiter
	var openerInl, closerInl Node
	var tempStack *delimiter
	var useDelims int
	var openerFound bool
	var openersBottom = map[byte]*delimiter{}
	var oddMatch = false

	openersBottom[itemUnderscore] = stackBottom
	openersBottom[itemAsterisk] = stackBottom
	openersBottom[itemTilde] = stackBottom

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
		openerFound = false
		for nil != opener && opener != stackBottom && opener != openersBottom[closercc] {
			oddMatch = (closer.canOpen || opener.canClose) && closer.originalNum%3 != 0 && (opener.originalNum+closer.originalNum)%3 == 0
			if opener.typ == closer.typ && opener.canOpen && !oddMatch {
				openerFound = true
				break
			}
			opener = opener.previous
		}
		oldCloser = closer

		if !openerFound {
			closer = closer.next
		} else {
			// calculate actual number of delimiters used from closer
			if closer.num >= 2 && opener.num >= 2 {
				useDelims = 2
			} else {
				useDelims = 1
			}

			openerInl = opener.node
			closerInl = closer.node

			// remove used delimiters from stack elts and inlines
			opener.num -= useDelims
			closer.num -= useDelims

			text := openerInl.Tokens()[0 : len(openerInl.Tokens())-useDelims]
			openerInl.SetTokens(text)
			text = closerInl.Tokens()[0 : len(closerInl.Tokens())-useDelims]
			closerInl.SetTokens(text)

			var emphStrongDel Node
			if 1 == useDelims {
				emphStrongDel = &Emphasis{&BaseNode{typ: NodeEmphasis}}
			} else {
				if itemTilde != closercc {
					emphStrongDel = &Strong{&BaseNode{typ: NodeStrong}}
				} else {
					emphStrongDel = &Strikethrough{&BaseNode{typ: NodeStrikethrough}}
				}
			}

			tmp := openerInl.Next()
			for nil != tmp && tmp != closerInl {
				next := tmp.Next()
				tmp.Unlink()
				emphStrongDel.AppendChild(emphStrongDel, tmp)
				tmp = next
			}

			openerInl.InsertAfter(openerInl, emphStrongDel)

			// remove elts between opener and closer in delimiters stack
			if opener.next != closer {
				opener.next = closer
				closer.previous = opener
			}

			// if opener has 0 delims, remove it and the inline
			if opener.num == 0 {
				openerInl.Unlink()
				t.removeDelimiter(opener)
			}

			if closer.num == 0 {
				closerInl.Unlink()
				tempStack = closer.next
				t.removeDelimiter(closer)
				closer = tempStack
			}
		}

		if !openerFound && !oddMatch {
			// Set lower bound for future searches for openers:
			// We don't do this with oddMatch because a **
			// that doesn't match an earlier * might turn into
			// an opener, and the * might be matched by something
			// else.
			openersBottom[closercc] = oldCloser.previous
			if !oldCloser.canOpen {
				// We can remove a closer that can't be an opener,
				// once we've seen there's no matching opener:
				t.removeDelimiter(oldCloser)
			}
		}
	}

	// 移除所有分隔符
	for t.context.delimiters != nil && t.context.delimiters != stackBottom {
		t.removeDelimiter(t.context.delimiters)
	}
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
		beforeIsWhitespace = isUnicodeWhitespace(tokenBefore)
		beforeIsPunct = isPunct(tokenBefore)
	}
	if itemEnd != tokenAfter {
		afterIsWhitespace = isUnicodeWhitespace(tokenAfter)
		afterIsPunct = isPunct(tokenAfter)
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

func (t *Tree) removeDelimiter(delim *delimiter) (ret *delimiter) {
	if delim.previous != nil {
		delim.previous.next = delim.next
	}
	if delim.next == nil {
		t.context.delimiters = delim.previous // 栈顶
	} else {
		delim.next.previous = delim.previous
	}
	return
}
