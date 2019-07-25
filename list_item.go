// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

type ListItem struct {
	*BaseNode
	*ListData
}

func (listItem *ListItem) Continue(context *Context) int {
	if context.blank {
		if nil == listItem.firstChild {
			// Blank line after empty list item
			return 1
		} else {
			context.advanceNextNonspace()
		}
	} else if context.indent >= listItem.markerOffset+listItem.padding {
		context.advanceOffset(listItem.markerOffset+
			listItem.padding, true)
	} else {
		return 1
	}
	return 0
}
