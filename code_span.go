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

func (t *Tree) parseCodeSpan(ctx *InlineContext) (ret *Node) {
	startPos := ctx.pos
	n := 0
	for ; startPos+n < ctx.tokensLen; n++ {
		if itemBacktick != ctx.tokens[startPos+n] {
			break
		}
	}

	backticks := ctx.tokens[startPos : startPos+n]
	if ctx.tokensLen <= startPos+n {
		ctx.pos += n
		ret = &Node{Typ: NodeText, Tokens: backticks}
		return
	}
	openMarker := &Node{Typ: NodeCodeSpanOpenMarker, Tokens: backticks}

	endPos := t.matchCodeSpanEnd(ctx.tokens[startPos+n:], n)
	if 1 > endPos {
		ctx.pos += n
		ret = &Node{Typ: NodeText, Tokens: backticks}
		return
	}
	endPos = startPos + endPos + n
	closeMarker := &Node{Typ: NodeCodeSpanCloseMarker, Tokens: ctx.tokens[endPos : endPos+n]}

	textTokens := ctx.tokens[startPos+n : endPos]
	if !t.context.option.VditorWYSIWYG {
		textTokens = replaceAll(textTokens, itemNewline, itemSpace)
		if 2 < len(textTokens) && itemSpace == textTokens[0] && itemSpace == textTokens[len(textTokens)-1] && !isBlankLine(textTokens) {
			// 如果首尾是空格并且整行不是空行时剔除首尾的一个空格
			openMarker.Tokens = append(openMarker.Tokens, textTokens[0])
			closeMarker.Tokens = ctx.tokens[endPos-1 : endPos+n]
			textTokens = textTokens[1:]
			textTokens = textTokens[:len(textTokens)-1]
		}
	}

	ret = &Node{Typ: NodeCodeSpan, codeMarkerLen: n}
	ret.AppendChild(openMarker)
	ret.AppendChild(&Node{Typ: NodeCodeSpanContent, Tokens: textTokens})
	ret.AppendChild(closeMarker)
	ctx.pos = endPos + n
	return
}

func (t *Tree) matchCodeSpanEnd(tokens []byte, num int) (pos int) {
	length := len(tokens)
	for pos < length {
		l := accept(tokens[pos:], itemBacktick)
		if num == l {
			next := pos + l
			if length-1 > next && itemBacktick == tokens[next] {
				continue
			}
			return pos
		}
		if 0 < l {
			pos += l
		} else {
			pos++
		}
	}
	return -1
}
