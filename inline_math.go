// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

var dollar = strToBytes("$")

func (t *Tree) parseInlineMath(ctx *InlineContext) (ret *Node) {
	if 3 > ctx.tokensLen {
		ctx.pos++
		return &Node{Typ: NodeText, Tokens: dollar}
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
			ret = &Node{Typ: NodeMathBlock}
			ret.AppendChild(&Node{Typ: NodeMathBlockOpenMarker})
			ret.AppendChild(&Node{Typ: NodeMathBlockContent, Tokens: ctx.tokens[blockStartPos:blockEndPos]})
			ret.AppendChild(&Node{Typ: NodeMathBlockCloseMarker})
			ctx.pos = blockEndPos + 2
			return
		}
	}

	if !t.context.option.InlineMathAllowDigitAfterOpenMarker && ctx.tokensLen > startPos+1 && isDigit(ctx.tokens[startPos+1]) { // $ 后面不能紧跟数字
		ctx.pos += 3
		return &Node{Typ: NodeText, Tokens: ctx.tokens[startPos : startPos+3]}
	}

	endPos := t.matchInlineMathEnd(ctx.tokens[startPos+1:])
	if 1 > endPos {
		ctx.pos++
		ret = &Node{Typ: NodeText, Tokens: dollar}
		return
	}

	endPos = startPos + endPos + 2

	tokens := ctx.tokens[startPos+1 : endPos-1]
	if 1 > len(trimWhitespace(tokens)) {
		ctx.pos++
		return &Node{Typ: NodeText, Tokens: dollar}
	}

	ret = &Node{Typ: NodeInlineMath}
	ret.AppendChild(&Node{Typ: NodeInlineMathOpenMarker})
	ret.AppendChild(&Node{Typ: NodeInlineMathContent, Tokens: tokens})
	ret.AppendChild(&Node{Typ: NodeInlineMathCloseMarker})

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
