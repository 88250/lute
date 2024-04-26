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
	"github.com/88250/lute/util"
)

var dollar = util.StrToBytes("$")

func (t *Tree) parseInlineMath(ctx *InlineContext) (ret *ast.Node) {
	if 3 > ctx.tokensLen || !t.Context.ParseOption.InlineMath {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: dollar}
	}

	startPos := ctx.pos
	blockStartPos := startPos
	dollars := 0
	for ; blockStartPos < ctx.tokensLen && lex.ItemDollar == ctx.tokens[blockStartPos]; blockStartPos++ {
		dollars++
	}
	if 2 <= dollars {
		if t.Context.ParseOption.ProtyleWYSIWYG {
			// Protyle 不允许从行级派生块级
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: dollar}
		}

		// 块节点
		matchBlock := false
		blockEndPos := blockStartPos + dollars
		var token byte
		for ; blockEndPos < ctx.tokensLen; blockEndPos++ {
			token = ctx.tokens[blockEndPos]
			if lex.ItemDollar == token && blockEndPos < ctx.tokensLen-1 && lex.ItemDollar == ctx.tokens[blockEndPos+1] {
				matchBlock = true
				break
			}
		}
		if matchBlock {
			ret = &ast.Node{Type: ast.NodeMathBlock}
			ret.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
			ret.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: ctx.tokens[blockStartPos:blockEndPos]})
			ret.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
			ctx.pos = blockEndPos + 2
			return
		}
	}

	if !t.Context.ParseOption.InlineMathAllowDigitAfterOpenMarker && ctx.tokensLen > startPos+1 && lex.IsDigit(ctx.tokens[startPos+1]) { // $ 后面不能紧跟数字
		ctx.pos += 3
		return &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[startPos : startPos+3]}
	}

	endPos := t.matchInlineMathEnd(ctx.tokens[startPos+1:])
	if 1 > endPos {
		ctx.pos++
		ret = &ast.Node{Type: ast.NodeText, Tokens: dollar}
		return
	}

	if t.Context.ParseOption.TextMark {
		if bytes.Contains(ctx.tokens[startPos+1:startPos+endPos+1], []byte("<span")) {
			// 中间包含 span 节点的话打断公式，以 span 优先
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: dollar}
		}
	}

	endPos = startPos + endPos + 2

	tokens := ctx.tokens[startPos+1 : endPos-1]
	if 1 > len(lex.TrimWhitespace(tokens)) {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: dollar}
	}

	ret = &ast.Node{Type: ast.NodeInlineMath}
	ret.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
	ret.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: tokens})
	ret.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})

	ctx.pos = endPos
	return
}

func (t *Tree) matchInlineMathEnd(tokens []byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if lex.ItemDollar == tokens[pos] && 0 < pos && lex.ItemBackslash != tokens[pos-1] {
			if pos < length-1 {
				if !lex.IsDigit(tokens[pos+1]) || t.Context.ParseOption.InlineMathAllowDigitAfterOpenMarker {
					return pos
				}
			} else {
				return pos
			}
		} else if lex.ItemNewline == tokens[pos] {
			return -1
		}
	}
	return -1
}
