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

// ThematicBreak 描述了分隔线节点结构。
type ThematicBreak struct {
	*BaseNode
}

func (thematicBreak *ThematicBreak) Continue(context *Context) int {
	return 1
}

func (thematicBreak *ThematicBreak) CanContain(nodeType int) bool {
	return false
}

func (t *Tree) parseThematicBreak() (ret *ThematicBreak) {
	markers := 0
	var marker item
	for i := t.context.nextNonspace; i < t.context.currentLineLen-1; i++ {
		token := t.context.currentLine[i]
		if itemSpace == token || itemTab == token {
			continue
		}

		if itemHyphen != token && itemUnderscore != token && itemAsterisk != token {
			return nil
		}

		if itemEnd != marker {
			if marker != token {
				return nil
			}
		} else {
			marker = token
		}
		markers++
	}

	if 3 > markers {
		return nil
	}

	return &ThematicBreak{&BaseNode{typ: NodeThematicBreak}}
}
