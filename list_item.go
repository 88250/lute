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

func (listItem *Node) ListItemContinue(context *Context) int {
	if context.blank {
		if nil == listItem.FirstChild { // 列表项后面是空的
			return 1
		}

		context.advanceNextNonspace()
	} else if context.indent >= listItem.MarkerOffset+listItem.Padding {
		context.advanceOffset(listItem.MarkerOffset+listItem.Padding, true)
	} else {
		return 1
	}
	return 0
}
