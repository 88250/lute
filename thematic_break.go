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

func (t *Tree) parseThematicBreak() (ok bool, caretTokens []byte) {
	markerCnt := 0
	var marker byte
	ln := t.context.currentLine
	var caretInLn bool
	if t.context.option.VditorWYSIWYG {
		if bytes.Contains(ln, []byte(caret)) {
			caretInLn = true
			ln = bytes.ReplaceAll(ln, []byte(caret), []byte(""))
		}
	}

	length := len(ln)
	for i := t.context.nextNonspace; i < length-1; i++ {
		token := ln[i]
		if itemSpace == token || itemTab == token {
			continue
		}

		if itemHyphen != token && itemUnderscore != token && itemAsterisk != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
		markerCnt++
	}

	if t.context.option.VditorWYSIWYG && caretInLn {
		caretTokens = []byte(caret)
	}

	return 3 <= markerCnt, caretTokens
}
