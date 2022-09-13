// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
)

// 判断分隔线（--- ***）是否开始。
func ThematicBreakStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if ok, caretTokens := t.parseThematicBreak(); ok {
		t.Context.closeUnmatchedBlocks()
		thematicBreak := t.Context.addChild(ast.NodeThematicBreak)
		thematicBreak.Tokens = caretTokens
		t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
		return 2
	}
	return 0
}

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

	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG) && caretInLn {
		caretTokens = editor.CaretTokens
	}
	return 3 <= markerCnt, caretTokens
}
