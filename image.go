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

// parseBang 解析 !，可能是图片标记符开始 ![ 也可能是普通文本 !。
func (t *Tree) parseBang(ctx *InlineContext) (ret *Node) {
	var startPos = ctx.pos
	ctx.pos++
	if ctx.pos < ctx.tokensLen && itemOpenBracket == ctx.tokens[ctx.pos] {
		ctx.pos++
		ret = &Node{Type: NodeText, Tokens: ctx.tokens[startPos:ctx.pos]}
		// 将图片开始标记符入栈
		t.addBracket(ret, startPos+2, true, ctx)
		return
	}

	ret = &Node{Type: NodeText, Tokens: ctx.tokens[startPos:ctx.pos]}
	return
}
