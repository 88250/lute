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

var dollar = strToItems("$")

func (t *Tree) parseInlineMath(ctx *InlineContext) (ret *Node) {
	if 2 > ctx.tokensLen {
		ctx.pos++
		return &Node{typ: NodeText, tokens: dollar}
	}

	startPos := ctx.pos
	blockStartPos := startPos
	dollars := 0
	for ; blockStartPos < ctx.tokensLen && itemDollar == ctx.tokens[blockStartPos].term(); blockStartPos++ {
		dollars++
	}
	if 2 <= dollars {
		// 块节点
		matchBlock := false
		blockEndPos := blockStartPos + dollars
		var token byte
		for ; blockEndPos < ctx.tokensLen; blockEndPos++ {
			token = ctx.tokens[blockEndPos].term()
			if itemDollar == token && blockEndPos < ctx.tokensLen-1 && itemDollar == ctx.tokens[blockEndPos+1].term() {
				matchBlock = true
				break
			}
		}
		if matchBlock {
			ret = &Node{typ: NodeMathBlock, tokens: ctx.tokens[blockStartPos:blockEndPos]}
			ctx.pos = blockEndPos + 2
			return
		}
	}

	if isDigit(ctx.tokens[startPos+1].term()) { // $ 后面不能紧跟数字
		ctx.pos += 2
		return &Node{typ: NodeText, tokens: ctx.tokens[startPos : startPos+2]}
	}

	endPos := t.matchInlineMathEnd(ctx.tokens[startPos+1:])
	if 1 > endPos {
		ctx.pos++
		ret = &Node{typ: NodeText, tokens: dollar}
		return
	}

	endPos = startPos + endPos + 2
	content := ctx.tokens[startPos+1 : endPos-1]
	ret = &Node{typ: NodeInlineMath, tokens: content}
	ctx.pos = endPos
	return
}

func (t *Tree) matchInlineMathEnd(tokens items) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if itemDollar == tokens[pos].term() {
			if pos < length-1 {
				if !isDigit(tokens[pos+1].term()) {
					return pos
				}
			} else {
				return pos
			}
		}
	}
	return -1
}
