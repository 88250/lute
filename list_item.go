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

func (listItem *BaseNode) ListItemContinue(context *Context) int {
	if context.blank {
		if nil == listItem.firstChild {
			// Blank line after empty list item
			return 1
		}

		context.advanceNextNonspace()
	} else if context.indent >= listItem.markerOffset+listItem.padding {
		context.advanceOffset(listItem.markerOffset+listItem.padding, true)
	} else {
		return 1
	}
	return 0
}
