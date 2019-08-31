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

func (t *Tree) parseThematicBreak() (ret *BaseNode) {
	markers := 0
	var marker byte
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

	return &BaseNode{typ: NodeThematicBreak}
}
