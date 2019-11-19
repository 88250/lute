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

func (t *Tree) parseThematicBreak() (ok bool, markers []byte) {
	markerCnt := 0
	var marker byte
	for i := t.context.nextNonspace; i < t.context.currentLineLen-1; i++ {
		token := t.context.currentLine[i]
		term := token
		markers = append(markers, t.context.currentLine[i])
		if itemSpace == term || itemTab == term {
			continue
		}

		if itemHyphen != term && itemUnderscore != term && itemAsterisk != term {
			return
		}

		if 0 != marker {
			if marker != term {
				return
			}
		} else {
			marker = term
		}
		markerCnt++
	}

	return 3 <= markerCnt, markers
}
