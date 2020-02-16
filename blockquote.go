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

import "github.com/88250/lute/ast"

func BlockquoteContinue(blockquote *ast.Node, context *Context) int {
	var ln = context.currentLine
	if !context.indented && peek(ln, context.nextNonspace) == itemGreater {
		context.advanceNextNonspace()
		context.advanceOffset(1, false)
		if token := peek(ln, context.offset); itemSpace == token || itemTab == token {
			context.advanceOffset(1, true)
		}
		return 0
	}
	return 1
}
