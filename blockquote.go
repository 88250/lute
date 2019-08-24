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

// Blockquote 描述了块引用节点结构。
type Blockquote struct {
	*BaseNode
}

func (blockquote *Blockquote) Continue(context *Context) int {
	var ln = context.currentLine
	if !context.indented && ln.peek(context.nextNonspace) == itemGreater {
		context.advanceNextNonspace()
		context.advanceOffset(1, false)
		if token := ln.peek(context.offset); itemSpace == token || itemTab == token {
			context.advanceOffset(1, true)
		}

		return 0
	}
	return 1
}
