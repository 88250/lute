// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

var dollar = strToBytes("$")

func (t *Tree) parseInlineMath(ctx *InlineContext) (ret *Node) {
	if 3 > ctx.tokensLen {
		ctx.pos++
		return &Node{Type: NodeText, Tokens: dollar}
	}

	startPos := ctx.pos
	blockStartPos := startPos
	dollars := 0
	for ; blockStartPos < ctx.tokensLen && itemDollar == ctx.tokens[blockStartPos]; blockStartPos++ {
		dollars++
	}
	if 2 <= dollars {
		// 块节点
		matchBlock := false
		blockEndPos := blockStartPos + dollars
		var token byte
		for ; blockEndPos < ctx.tokensLen; blockEndPos++ {
			token = ctx.tokens[blockEndPos]
			if itemDollar == token && blockEndPos < ctx.tokensLen-1 && itemDollar == ctx.tokens[blockEndPos+1] {
				matchBlock = true
				break
			}
		}
		if matchBlock {
			ret = &Node{Type: NodeMathBlock}
			ret.AppendChild(&Node{Type: NodeMathBlockOpenMarker})
			ret.AppendChild(&Node{Type: NodeMathBlockContent, Tokens: ctx.tokens[blockStartPos:blockEndPos]})
			ret.AppendChild(&Node{Type: NodeMathBlockCloseMarker})
			ctx.pos = blockEndPos + 2
			return
		}
	}

	if !t.context.option.InlineMathAllowDigitAfterOpenMarker && ctx.tokensLen > startPos+1 && isDigit(ctx.tokens[startPos+1]) { // $ 后面不能紧跟数字
		ctx.pos += 3
		return &Node{Type: NodeText, Tokens: ctx.tokens[startPos : startPos+3]}
	}

	endPos := t.matchInlineMathEnd(ctx.tokens[startPos+1:])
	if 1 > endPos {
		ctx.pos++
		ret = &Node{Type: NodeText, Tokens: dollar}
		return
	}

	endPos = startPos + endPos + 2

	tokens := ctx.tokens[startPos+1 : endPos-1]
	if 1 > len(trimWhitespace(tokens)) {
		ctx.pos++
		return &Node{Type: NodeText, Tokens: dollar}
	}

	ret = &Node{Type: NodeInlineMath}
	ret.AppendChild(&Node{Type: NodeInlineMathOpenMarker})
	ret.AppendChild(&Node{Type: NodeInlineMathContent, Tokens: tokens})
	ret.AppendChild(&Node{Type: NodeInlineMathCloseMarker})

	ctx.pos = endPos
	return
}

func (t *Tree) matchInlineMathEnd(tokens []byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if itemDollar == tokens[pos] {
			if pos < length-1 {
				if !isDigit(tokens[pos+1]) {
					return pos
				}
			} else {
				return pos
			}
		}
	}
	return -1
}
