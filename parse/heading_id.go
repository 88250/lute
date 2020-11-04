// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/util"
)

var openCurlyBrace = util.StrToBytes("{")
var closeCurlyBrace = util.StrToBytes("}")

func (t *Tree) parseHeadingID(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	if !t.Context.Option.HeadingID || ast.NodeHeading != block.Type || 3 > ctx.tokensLen {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
	}

	startPos := ctx.pos
	content := ctx.tokens[startPos:]
	curlyBracesEnd := bytes.Index(content, closeCurlyBrace)
	if 2 > curlyBracesEnd {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
	}

	curlyBracesStart := bytes.Index(content, []byte("{"))
	if 0 > curlyBracesStart {
		return nil
	}

	length := len(content)
	if length-1 != curlyBracesEnd {
		if !bytes.HasSuffix(content, []byte("}"+util.Caret)) && bytes.HasSuffix(content, util.CaretTokens) {
			// # foo {id}b‸
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
		}
	}

	if t.Context.Option.VditorWYSIWYG {
		content = bytes.ReplaceAll(content, util.CaretTokens, nil)
	}
	id := content[curlyBracesStart+1 : curlyBracesEnd]
	ctx.pos += curlyBracesEnd + 1
	block.LastChild.Tokens = bytes.TrimRight(block.LastChild.Tokens, " ")
	return &ast.Node{Type: ast.NodeHeadingID, Tokens: id}
}
