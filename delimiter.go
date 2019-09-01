// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"unicode"
	"unicode/utf8"
)

// delimiter 描述了强调、链接和图片解析过程中用到的分隔符（[, ![, *, _）相关信息。
type delimiter struct {
	node           *Node      // the text node point to
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

// 嵌套强调和链接的解析算法的中文解读可参考这里 https://hacpai.com/article/1566893557720

// handleDelim 将分隔符 *_~ 入栈。
func (t *Tree) handleDelim(block *Node, ctx *InlineContext) {
	startPos := ctx.pos
	delim := t.scanDelims(ctx)

	text := ctx.tokens[startPos:ctx.pos]
	node := &Node{typ: NodeText, tokens: text}
	block.AppendChild(block, node)

	// 将这个分隔符入栈
	ctx.delimiters = &delimiter{
		typ:         delim.typ,
		num:         delim.num,
		originalNum: delim.num,
		node:        node,
		previous:    ctx.delimiters,
		next:        nil,
		canOpen:     delim.canOpen,
		canClose:    delim.canClose,
	}
	if ctx.delimiters.previous != nil {
		ctx.delimiters.previous.next = ctx.delimiters
	}
}

// processEmphasis 处理强调、加粗以及删除线。
func (t *Tree) processEmphasis(stackBottom *delimiter, ctx *InlineContext) {
	var opener, closer, oldCloser *delimiter
	var openerInl, closerInl *Node
	var tempStack *delimiter
	var useDelims int
	var openerFound bool
	var openersBottom = map[byte]*delimiter{}
	var oddMatch = false

	openersBottom[itemUnderscore] = stackBottom
	openersBottom[itemAsterisk] = stackBottom
	openersBottom[itemTilde] = stackBottom

	// find first closer above stack_bottom:
	closer = ctx.delimiters
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

			var emStrongDel *Node
			if 1 == useDelims {
				if itemAsterisk == closercc {
					emStrongDel = &Node{typ: NodeEmphasis, strongEmDelMarker: itemAsterisk, strongEmDelMarkenLen: 1}
				} else if itemUnderscore == closercc {
					emStrongDel = &Node{typ: NodeEmphasis, strongEmDelMarker: itemUnderscore, strongEmDelMarkenLen: 1}
				} else if itemTilde == closercc {
					if t.context.option.GFMStrikethrough {
						emStrongDel = &Node{typ: NodeStrikethrough, strongEmDelMarker: itemTilde, strongEmDelMarkenLen: 1}
					}
				}
			} else {
				if itemTilde != closercc {
					emStrongDel = &Node{typ: NodeStrong, strongEmDelMarker: closercc, strongEmDelMarkenLen: 2}
				} else {
					if t.context.option.GFMStrikethrough {
						emStrongDel = &Node{typ: NodeStrikethrough, strongEmDelMarker: closercc, strongEmDelMarkenLen: 2}
					}
				}
			}

			tmp := openerInl.next
			for nil != tmp && tmp != closerInl {
				next := tmp.next
				tmp.Unlink()
				emStrongDel.AppendChild(emStrongDel, tmp)
				tmp = next
			}

			openerInl.InsertAfter(openerInl, emStrongDel)

			// remove elts between opener and closer in delimiters stack
			if opener.next != closer {
				opener.next = closer
				closer.previous = opener
			}

			// if opener has 0 delims, remove it and the inline
			if opener.num == 0 {
				openerInl.Unlink()
				t.removeDelimiter(opener, ctx)
			}

			if closer.num == 0 {
				closerInl.Unlink()
				tempStack = closer.next
				t.removeDelimiter(closer, ctx)
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
				t.removeDelimiter(oldCloser, ctx)
			}
		}
	}

	// 移除所有分隔符
	for ctx.delimiters != nil && ctx.delimiters != stackBottom {
		t.removeDelimiter(ctx.delimiters, ctx)
	}
}

func (t *Tree) scanDelims(ctx *InlineContext) *delimiter {
	startPos := ctx.pos
	token := ctx.tokens[startPos]
	delimitersCount := 0
	for i := ctx.pos; i < ctx.tokensLen && token == ctx.tokens[i]; i++ {
		delimitersCount++
		ctx.pos++
	}

	tokenBefore, tokenAfter := rune(itemNewline), rune(itemNewline)
	if 0 < startPos {
		t := ctx.tokens[startPos-1]
		if t >= utf8.RuneSelf {
			tokenBefore, _ = utf8.DecodeLastRune(ctx.tokens[:startPos-1])
		} else {
			tokenBefore = rune(t)
		}
	}
	if ctx.tokensLen > ctx.pos {
		t := ctx.tokens[ctx.pos]
		if t >= utf8.RuneSelf {
			tokenAfter, _ = utf8.DecodeRune(ctx.tokens[ctx.pos:])
		} else {
			tokenAfter = rune(t)
		}
	}

	beforeIsWhitespace := isUnicodeWhitespace(tokenBefore)
	beforeIsPunct := unicode.IsPunct(tokenBefore)
	afterIsWhitespace := isUnicodeWhitespace(tokenAfter)
	afterIsPunct := unicode.IsPunct(tokenAfter)
	isLeftFlanking := !afterIsWhitespace && (!afterIsPunct || beforeIsWhitespace || beforeIsPunct)
	isRightFlanking := !beforeIsWhitespace && (!beforeIsPunct || afterIsWhitespace || afterIsPunct)
	var canOpen, canClose bool
	if itemUnderscore == token {
		canOpen = isLeftFlanking && (!isRightFlanking || beforeIsPunct)
		canClose = isRightFlanking && (!isLeftFlanking || afterIsPunct)
	} else {
		canOpen = isLeftFlanking
		canClose = isRightFlanking
	}

	return &delimiter{typ: token, num: delimitersCount, active: true, canOpen: canOpen, canClose: canClose}
}

func (t *Tree) removeDelimiter(delim *delimiter, ctx *InlineContext) (ret *delimiter) {
	if delim.previous != nil {
		delim.previous.next = delim.next
	}
	if delim.next == nil {
		ctx.delimiters = delim.previous // 栈顶
	} else {
		delim.next.previous = delim.previous
	}
	return
}
