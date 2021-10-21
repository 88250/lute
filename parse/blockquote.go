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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

// BlockquoteStart 判断块引用（>）是否开始。
func BlockquoteStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	marker := lex.Peek(t.Context.currentLine, t.Context.nextNonspace)
	if lex.ItemGreater != marker {
		return 0
	}

	markers := []byte{marker}
	t.Context.advanceNextNonspace()
	t.Context.advanceOffset(1, false)
	// > 后面的空格是可选的
	whitespace := lex.Peek(t.Context.currentLine, t.Context.offset)
	withSpace := lex.ItemSpace == whitespace || lex.ItemTab == whitespace
	if withSpace {
		t.Context.advanceOffset(1, true)
		markers = append(markers, whitespace)
	}
	t.Context.closeUnmatchedBlocks()
	t.Context.addChild(ast.NodeBlockquote)
	t.Context.addChildMarker(ast.NodeBlockquoteMarker, markers)
	return 1
}

func BlockquoteContinue(blockquote *ast.Node, context *Context) int {
	ln := context.currentLine
	if !context.indented && lex.Peek(ln, context.nextNonspace) == lex.ItemGreater {
		context.advanceNextNonspace()
		context.advanceOffset(1, false)
		if token := lex.Peek(ln, context.offset); lex.ItemSpace == token || lex.ItemTab == token {
			context.advanceOffset(1, true)
		}
		return 0
	}
	return 1
}
