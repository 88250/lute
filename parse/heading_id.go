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
	"github.com/88250/lute/editor"
	"github.com/88250/lute/util"
)

var openCurlyBrace = util.StrToBytes("{")
var closeCurlyBrace = util.StrToBytes("}")

func (t *Tree) parseHeadingID(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	if !t.Context.ParseOption.HeadingID || ast.NodeHeading != block.Type || 3 > ctx.tokensLen {
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
		if !bytes.HasSuffix(content, []byte("}"+editor.Caret)) && bytes.HasSuffix(content, editor.CaretTokens) {
			// # foo {id}b‸
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
		}
	}

	if t.Context.ParseOption.VditorWYSIWYG {
		content = bytes.ReplaceAll(content, editor.CaretTokens, nil)
	}
	id := content[curlyBracesStart+1 : curlyBracesEnd]
	ctx.pos += curlyBracesEnd + 1
	if nil != block.LastChild {
		block.LastChild.Tokens = bytes.TrimRight(block.LastChild.Tokens, " ")
	}
	return &ast.Node{Type: ast.NodeHeadingID, Tokens: id}
}
