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

func FootnotesContinue(footnotesDef *ast.Node, context *Context) int {
	if context.blank {
		return 0
	}

	if 4 > context.indent {
		return 1
	}

	context.advanceOffset(4, true)
	return 0
}

func (t *Tree) FindFootnotesDef(label []byte) (pos int, def *ast.Node) {
	pos = 0
	if t.Context.Option.VditorIR || t.Context.Option.VditorSV || t.Context.Option.VditorWYSIWYG {
		label = bytes.ReplaceAll(label, util.CaretTokens, nil)
	}
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeFootnotesDef != n.Type {
			return ast.WalkContinue
		}
		if bytes.EqualFold(n.Tokens, label) {
			pos++
			def = n
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return
}
