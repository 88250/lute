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

import "bytes"

func (footnotesDef *Node) footnotesContinue(context *Context) int {
	if 4 <= context.indent {
		context.advanceOffset(4, true)
	} else {
		return 1
	}
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
