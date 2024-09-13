// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
)

// delimiter 描述了强调、链接和图片解析过程中用到的分隔符（[, ![, *, _, ~）相关信息。
type delimiter struct {
	node           *ast.Node  // 分隔符对应的文本节点
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

// 嵌套强调和链接的解析算法的中文解读可参考这里 https://ld246.com/article/1566893557720

// handleDelim 将分隔符 *_~ 入栈。
func (t *Tree) handleDelim(block *ast.Node, ctx *InlineContext) {
	startPos := ctx.pos
	delim := t.scanDelims(ctx)

	text := ctx.tokens[startPos:ctx.pos]
	node := &ast.Node{Type: ast.NodeText, Tokens: text}
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
	if nil == ctx.delimiters {
		return
	}

	var opener, closer, oldCloser *delimiter
	var openerInl, closerInl *ast.Node
	var tempStack *delimiter
	var useDelims int
	var openerFound bool
	var openersBottom = map[byte]*delimiter{}
	var oddMatch = false

	openersBottom[lex.ItemUnderscore] = stackBottom
	openersBottom[lex.ItemAsterisk] = stackBottom
	openersBottom[lex.ItemTilde] = stackBottom
	openersBottom[lex.ItemEqual] = stackBottom
	openersBottom[lex.ItemCrosshatch] = stackBottom
	openersBottom[lex.ItemCaret] = stackBottom

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

			if t.Context.ParseOption.GFMStrikethrough || t.Context.ParseOption.Sub {
				if lex.ItemTilde == closercc && opener.num != closer.num {
					break
				}
			} else {
				if lex.ItemTilde == closercc {
					break
				}
			}

			if t.Context.ParseOption.Sup {
				if lex.ItemCaret == closercc && opener.num != closer.num {
					break
				}
			} else {
				if lex.ItemCaret == closercc {
					break
				}
			}

			if !t.Context.ParseOption.InlineAsterisk {
				if lex.ItemAsterisk == closercc {
					break
				}
			}

			if !t.Context.ParseOption.InlineUnderscore {
				if lex.ItemUnderscore == closercc {
					break
				}
			}

			if t.Context.ParseOption.Mark {
				if lex.ItemEqual == closercc && opener.num != closer.num {
					break
				}
			} else {
				if lex.ItemEqual == closercc {
					break
				}
			}

			if t.Context.ParseOption.Tag {
				if lex.ItemCrosshatch == closercc && opener.num != closer.num {
					break
				}
			} else {
				if lex.ItemCrosshatch == closercc {
					break
				}
			}

			// remove used delimiters from stack elts and inlines
			opener.num -= useDelims
			closer.num -= useDelims

			openerTokens := openerInl.Tokens[len(openerInl.Tokens)-useDelims:]
			text := openerInl.Tokens[0 : len(openerInl.Tokens)-useDelims]
			openerInl.Tokens = text
			closerTokens := closerInl.Tokens[len(closerInl.Tokens)-useDelims:]
			text = closerInl.Tokens[0 : len(closerInl.Tokens)-useDelims]
			closerInl.Tokens = text

			openMarker := &ast.Node{Tokens: openerTokens, Close: true}
			emStrongDelMark := &ast.Node{Close: true}
			closeMarker := &ast.Node{Tokens: closerTokens, Close: true}
			if 1 == useDelims {
				if lex.ItemAsterisk == closercc {
					emStrongDelMark.Type = ast.NodeEmphasis
					openMarker.Type = ast.NodeEmA6kOpenMarker
					closeMarker.Type = ast.NodeEmA6kCloseMarker
				} else if lex.ItemUnderscore == closercc {
					emStrongDelMark.Type = ast.NodeEmphasis
					openMarker.Type = ast.NodeEmU8eOpenMarker
					closeMarker.Type = ast.NodeEmU8eCloseMarker
				} else if lex.ItemTilde == closercc {
					if t.Context.ParseOption.Sub { // 优先下标
						emStrongDelMark.Type = ast.NodeSub
						openMarker.Type = ast.NodeSubOpenMarker
						closeMarker.Type = ast.NodeSubCloseMarker
					} else if t.Context.ParseOption.GFMStrikethrough && t.Context.ParseOption.GFMStrikethrough1 {
						emStrongDelMark.Type = ast.NodeStrikethrough
						openMarker.Type = ast.NodeStrikethrough1OpenMarker
						closeMarker.Type = ast.NodeStrikethrough1CloseMarker
					}
				} else if lex.ItemEqual == closercc {
					if t.Context.ParseOption.Mark {
						emStrongDelMark.Type = ast.NodeMark
						openMarker.Type = ast.NodeMark1OpenMarker
						closeMarker.Type = ast.NodeMark1CloseMarker
					}
				} else if lex.ItemCrosshatch == closercc {
					if t.Context.ParseOption.Tag {
						emStrongDelMark.Type = ast.NodeTag
						openMarker.Type = ast.NodeTagOpenMarker
						closeMarker.Type = ast.NodeTagCloseMarker
					}
				} else if lex.ItemCaret == closercc {
					if t.Context.ParseOption.Sup {
						emStrongDelMark.Type = ast.NodeSup
						openMarker.Type = ast.NodeSupOpenMarker
						closeMarker.Type = ast.NodeSupCloseMarker
					}
				}
			} else {
				if lex.ItemAsterisk == closercc {
					emStrongDelMark.Type = ast.NodeStrong
					openMarker.Type = ast.NodeStrongA6kOpenMarker
					closeMarker.Type = ast.NodeStrongA6kCloseMarker
				} else if lex.ItemUnderscore == closercc {
					emStrongDelMark.Type = ast.NodeStrong
					openMarker.Type = ast.NodeStrongU8eOpenMarker
					closeMarker.Type = ast.NodeStrongU8eCloseMarker
				} else if lex.ItemTilde == closercc {
					if t.Context.ParseOption.GFMStrikethrough {
						emStrongDelMark.Type = ast.NodeStrikethrough
						openMarker.Type = ast.NodeStrikethrough2OpenMarker
						closeMarker.Type = ast.NodeStrikethrough2CloseMarker
					}
				} else if lex.ItemEqual == closercc {
					if t.Context.ParseOption.Mark {
						emStrongDelMark.Type = ast.NodeMark
						openMarker.Type = ast.NodeMark2OpenMarker
						closeMarker.Type = ast.NodeMark2CloseMarker
					}
				}
			}

			tmp := openerInl.Next
			for nil != tmp && tmp != closerInl {
				next := tmp.Next
				tmp.Unlink()
				emStrongDelMark.AppendChild(tmp)
				tmp = next
			}

			emStrongDelMark.PrependChild(openMarker) // 插入起始标记符
			emStrongDelMark.AppendChild(closeMarker) // 插入结束标记符
			openerInl.InsertAfter(emStrongDelMark)

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

	tokenBefore, tokenAfter := rune(lex.ItemNewline), rune(lex.ItemNewline)
	if 0 < startPos {
		c := ctx.tokens[startPos-1]
		if c >= utf8.RuneSelf {
			tokenBefore, _ = utf8.DecodeLastRune(ctx.tokens[:startPos])
		} else {
			tokenBefore = rune(c)
		}

		if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG) && editor.Caret == string(tokenBefore) {
			// 跳过插入符位置向前看
			caretLen := len(editor.Caret)
			if 0 < startPos-caretLen {
				c = ctx.tokens[startPos-caretLen-1]
				if c >= utf8.RuneSelf {
					tokenBefore, _ = utf8.DecodeLastRune(ctx.tokens[:startPos-caretLen])
				} else {
					tokenBefore = rune(c)
				}
			}
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

	afterIsWhitespace := lex.IsUnicodeWhitespace(tokenAfter)
	afterIsPunct := unicode.IsPunct(tokenAfter) || unicode.IsSymbol(tokenAfter)
	if (lex.ItemAsterisk == token && '~' == tokenAfter) || (lex.ItemTilde == token && '*' == tokenAfter) ||
		(lex.ItemCaret == token && ('+' == tokenAfter || '-' == tokenAfter)) ||
		(lex.ItemTilde == token && ('+' == tokenAfter || '-' == tokenAfter)) {
		afterIsPunct = false
	}
	beforeIsWhitespace := lex.IsUnicodeWhitespace(tokenBefore)
	beforeIsPunct := unicode.IsPunct(tokenBefore) || unicode.IsSymbol(tokenBefore)
	if (lex.ItemAsterisk == token && '~' == tokenBefore) || (lex.ItemTilde == token && '*' == tokenBefore) ||
		(lex.ItemCaret == token && ('+' == tokenBefore || '-' == tokenBefore)) ||
		(lex.ItemTilde == token && ('+' == tokenBefore || '-' == tokenBefore)) {
		beforeIsPunct = false
	}

	if t.Context.ParseOption.ProtyleWYSIWYG {
		// Markdown 中 ** 加粗失效问题 https://ld246.com/article/1597581380183
		afterIsPunct, beforeIsPunct = false, false
	}

	isLeftFlanking := !afterIsWhitespace && (!afterIsPunct || beforeIsWhitespace || beforeIsPunct)
	isRightFlanking := !beforeIsWhitespace && (!beforeIsPunct || afterIsWhitespace || afterIsPunct)
	var canOpen, canClose bool
	if lex.ItemUnderscore == token {
		canOpen = isLeftFlanking && (!isRightFlanking || beforeIsPunct)
		canClose = isRightFlanking && (!isLeftFlanking || afterIsPunct)
	} else {
		if t.Context.ParseOption.Mark && lex.ItemEqual == token && 2 != delimitersCount { // ==Mark== 标记使用两个等号
			canOpen = false
			canClose = false
		} else if t.Context.ParseOption.Tag && lex.ItemCrosshatch == token && 1 != delimitersCount { // #Tag# 标记使用一个井号
			canOpen = false
			canClose = false
		} else if t.Context.ParseOption.Sup && lex.ItemCaret == token && 1 != delimitersCount { // ^Sup^ 标记使用一个 ^
			canOpen = false
			canClose = false
		} else if t.Context.ParseOption.Sub && lex.ItemTilde == token {
			if t.Context.ParseOption.GFMStrikethrough && 3 == delimitersCount { // 单独处理 ~~~foo~~~ 的情况，即下标嵌套删除线
				canOpen = isLeftFlanking
				canClose = isRightFlanking
			} else if 1 != delimitersCount { // ~Sub~ 标记使用一个 ~
				canOpen = false
				canClose = false
				if t.Context.ParseOption.GFMStrikethrough && 2 == delimitersCount {
					canOpen = isLeftFlanking
					canClose = isRightFlanking
				}
			} else {
				canOpen = isLeftFlanking
				canClose = isRightFlanking
			}
		} else if t.Context.ParseOption.GFMStrikethrough && lex.ItemTilde == token {
			if 1 == delimitersCount {
				if !t.Context.ParseOption.GFMStrikethrough1 {
					canOpen = false
					canClose = false
				} else {
					canOpen = isLeftFlanking
					canClose = isRightFlanking
				}
			} else {
				canOpen = isLeftFlanking
				canClose = isRightFlanking
			}
		} else {
			canOpen = isLeftFlanking
			canClose = isRightFlanking
		}
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
