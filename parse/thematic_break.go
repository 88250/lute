// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"github.com/88250/lute/lex"
)

func (t *Tree) parseThematicBreak() (ok bool, caretTokens []byte) {
	markerCnt := 0
	var marker byte
	ln := t.Context.currentLine
	var caretInLn bool
	length := len(ln)
	for i := t.Context.nextNonspace; i < length-1; i++ {
		token := ln[i]
		if lex.ItemSpace == token || lex.ItemTab == token {
			continue
		}

		if lex.ItemHyphen != token && lex.ItemUnderscore != token && lex.ItemAsterisk != token {
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

	if t.Context.Option.VditorWYSIWYG && caretInLn {
		caretTokens = []byte(Caret)
	}

	return 3 <= markerCnt, caretTokens
}
