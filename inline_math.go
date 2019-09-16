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

var dollar = toItems("$")

func (t *Tree) parseInlineMath(ctx *InlineContext) (ret *Node) {
	startPos := ctx.pos
	endPos := t.matchInlineMathEnd(ctx.tokens[startPos+1:])
	if 1 > endPos {
		ctx.pos++
		ret = &Node{typ: NodeText, tokens: dollar}
		return
	}

	endPos = startPos + endPos + 2
	textTokens := ctx.tokens[startPos+1 : endPos-1]
	ret = &Node{typ: NodeInlineMath, tokens: textTokens}
	ctx.pos = endPos
	return
}

func (t *Tree) matchInlineMathEnd(tokens items) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if itemDollar == tokens[pos] {
			return pos
		}
	}
	return -1
}
