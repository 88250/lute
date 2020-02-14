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

import "bytes"

func (footnotesDef *Node) footnotesContinue(context *Context) int {
	if context.blank {
		return 0
	}

	if 4 > context.indent {
		return 1
	}

	context.advanceOffset(4, true)
	return 0
}

func (context *Context) findFootnotesDef(label []byte) (int, *Node) {
	for i, n := range context.footnotesDefs {
		if bytes.EqualFold(label, n.tokens) {
			return i + 1, n
		}
	}
	return -1, nil
}
