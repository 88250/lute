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

// delimiter 描述了强调、链接和图片解析过程中用到的分隔符（[, ![, *, _, ~）相关信息。
type delimiter struct {
	node           *Node      // 分隔符对应的文本节点
	typ            byte       // 分隔符字节 [*_~
	num            int        // 分隔符字节数
	originalNum    int        // 原始分隔符字节数
	canOpen        bool       // 是否是开始分隔符
	canClose       bool       // 是否是结束分隔符
	previous, next *delimiter // 双向链表前后节点

	active            bool
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
	block.AppendChild(node)

	// 将这个分隔符入栈
	if delim.canOpen || delim.canClose {
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
		if nil != ctx.delimiters.previous {
			ctx.delimiters.previous.next = ctx.delimiters
		}
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
	for nil != closer && closer.previous != stackBottom {
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

			openerTokens := openerInl.tokens[len(openerInl.tokens)-useDelims:]
			text := openerInl.tokens[0 : len(openerInl.tokens)-useDelims]
			openerInl.tokens = text
			closerTokens := closerInl.tokens[len(closerInl.tokens)-useDelims:]
			text = closerInl.tokens[0 : len(closerInl.tokens)-useDelims]
			closerInl.tokens = text

			openMarker := &Node{tokens: openerTokens, close: true}
			emStrongDel := &Node{close: true}
			closeMarker := &Node{tokens: closerTokens, close: true}
			if 1 == useDelims {
				if itemAsterisk == closercc {
					emStrongDel.typ = NodeEmphasis
					openMarker.typ = NodeEmA6kOpenMarker
					closeMarker.typ = NodeEmA6kCloseMarker
				} else if itemUnderscore == closercc {
					emStrongDel.typ = NodeEmphasis
					openMarker.typ = NodeEmU8eOpenMarker
					closeMarker.typ = NodeEmU8eCloseMarker
				} else if itemTilde == closercc {
					if t.context.option.GFMStrikethrough {
						emStrongDel.typ = NodeStrikethrough
						openMarker.typ = NodeStrikethrough1OpenMarker
						closeMarker.typ = NodeStrikethrough1CloseMarker
					}
				}
			} else {
				if itemAsterisk == closercc {
					emStrongDel.typ = NodeStrong
					openMarker.typ = NodeStrongA6kOpenMarker
					closeMarker.typ = NodeStrongA6kCloseMarker
				} else if itemUnderscore == closercc {
					emStrongDel.typ = NodeStrong
					openMarker.typ = NodeStrongU8eOpenMarker
					closeMarker.typ = NodeStrongU8eCloseMarker
				} else if itemTilde == closercc {
					if t.context.option.GFMStrikethrough {
						emStrongDel.typ = NodeStrikethrough
						openMarker.typ = NodeStrikethrough2OpenMarker
						closeMarker.typ = NodeStrikethrough2CloseMarker
					}
				}
			}

			tmp := openerInl.next
			for nil != tmp && tmp != closerInl {
				next := tmp.next
				tmp.Unlink()
				emStrongDel.AppendChild(tmp)
				tmp = next
			}

			emStrongDel.PrependChild(openMarker) // 插入起始标记符
			emStrongDel.AppendChild(closeMarker) // 插入结束标记符
			openerInl.InsertAfter(emStrongDel)

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
			openersBottom[closercc] = oldCloser.previous
			if !oldCloser.canOpen {
				// We can remove a closer that can't be an opener,
				// once we've seen there's no matching opener:
				t.removeDelimiter(oldCloser, ctx)
			}
		}
	}

	// 移除所有分隔符
	for nil != ctx.delimiters && ctx.delimiters != stackBottom {
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
			tokenBefore, _ = utf8.DecodeLastRune(ctx.tokens[:startPos])
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

	afterIsWhitespace := isUnicodeWhitespace(tokenAfter)
	afterIsPunct := unicode.IsPunct(tokenAfter) || unicode.IsSymbol(tokenAfter)
	beforeIsWhitespace := isUnicodeWhitespace(tokenBefore)
	beforeIsPunct := unicode.IsPunct(tokenBefore) || unicode.IsSymbol(tokenBefore)

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
	if nil != delim.previous {
		delim.previous.next = delim.next
	}
	if nil == delim.next {
		ctx.delimiters = delim.previous // 栈顶
	} else {
		delim.next.previous = delim.previous
	}
	return
}
