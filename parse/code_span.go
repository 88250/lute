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
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

func (t *Tree) parseCodeSpan(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	startPos := ctx.pos
	n := 0
	for ; startPos+n < ctx.tokensLen; n++ {
		if lex.ItemBacktick != ctx.tokens[startPos+n] {
			break
		}
	}

	backticks := ctx.tokens[startPos : startPos+n]
	if ctx.tokensLen <= startPos+n {
		ctx.pos += n
		ret = &ast.Node{Type: ast.NodeText, Tokens: backticks}
		return
	}
	openMarker := &ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: backticks}

	endPos := t.matchCodeSpanEnd(ctx.tokens[startPos+n:], n)
	if 1 > endPos {
		ctx.pos += n
		ret = &ast.Node{Type: ast.NodeText, Tokens: backticks}
		return
	}
	endPos = startPos + endPos + n
	closeMarker := &ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: ctx.tokens[endPos : endPos+n]}

	textTokens := ctx.tokens[startPos+n : endPos]
	textTokens = lex.ReplaceAll(textTokens, lex.ItemNewline, lex.ItemSpace)
	if 2 < len(textTokens) && lex.ItemSpace == textTokens[0] && lex.ItemSpace == textTokens[len(textTokens)-1] && !lex.IsBlankLine(textTokens) {
		// 如果首尾是空格并且整行不是空行时剔除首尾的一个空格
		openMarker.Tokens = append(openMarker.Tokens, textTokens[0])
		closeMarker.Tokens = ctx.tokens[endPos-1 : endPos+n]
		textTokens = textTokens[1 : len(textTokens)-1]
	}

	if t.Context.ParseOption.GFMTable {
		if ast.NodeTableCell == block.Type {
			// 表格中的代码中带有管道符的处理 https://github.com/88250/lute/issues/63
			textTokens = bytes.ReplaceAll(textTokens, []byte("\\|"), []byte("|"))
		}
	}

	ret = &ast.Node{Type: ast.NodeCodeSpan, CodeMarkerLen: n}
	ret.AppendChild(openMarker)
	ret.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent, Tokens: textTokens})
	ret.AppendChild(closeMarker)
	ctx.pos = endPos + n
	return
}

func (t *Tree) matchCodeSpanEnd(tokens []byte, num int) (pos int) {
	length := len(tokens)
	for pos < length {
		l := lex.Accept(tokens[pos:], lex.ItemBacktick)
		if num == l {
			next := pos + l
			if length-1 > next && lex.ItemBacktick == tokens[next] {
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
