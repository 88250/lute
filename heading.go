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

import (
	"bytes"
	"github.com/88250/lute/lex"
)

func (t *Tree) parseATXHeading() (ok bool, markers, content []byte, level int) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if lex.ItemCrosshatch != marker {
		return
	}

	level = lex.Accept(tokens, lex.ItemCrosshatch)
	if 6 < level {
		return
	}

	if level < len(tokens) && !lex.IsWhitespace(tokens[level]) {
		return
	}

	markers = t.context.currentLine[t.context.nextNonspace : t.context.nextNonspace+level+1]

	content = make([]byte, 0, 256)
	_, tokens = lex.TrimLeft(tokens)
	_, tokens = lex.TrimLeft(tokens[level:])
	for _, token := range tokens {
		if lex.ItemNewline == token {
			break
		}
		content = append(content, token)
	}

	_, content = lex.TrimRight(content)
	closingCrosshatchIndex := len(content) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if lex.ItemCrosshatch == content[closingCrosshatchIndex] {
			continue
		}
		if lex.ItemSpace == content[closingCrosshatchIndex] {
			break
		} else {
			closingCrosshatchIndex = len(content)
			break
		}
	}

	if 0 >= closingCrosshatchIndex {
		content = make([]byte, 0, 0)
	} else if 0 < closingCrosshatchIndex {
		content = content[:closingCrosshatchIndex]
		_, content = lex.TrimRight(content)
	}

	if t.context.option.VditorWYSIWYG {
		if caret == string(content) || "" == string(content) {
			return
		}
	}

	ok = true
	return
}

func (t *Tree) parseSetextHeading() (level int) {
	ln := lex.TrimWhitespace(t.context.currentLine)
	start := 0
	marker := ln[start]
	if lex.ItemEqual != marker && lex.ItemHyphen != marker {
		return
	}

	var caretInLn bool
	if t.context.option.VditorWYSIWYG {
		if bytes.Contains(ln, []byte(caret)) {
			caretInLn = true
			ln = bytes.ReplaceAll(ln, []byte(caret), []byte(""))
		}
	}

	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if lex.ItemEqual != token && lex.ItemHyphen != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
	}

	level = 1
	if lex.ItemHyphen == marker {
		level = 2
	}

	if t.context.option.VditorWYSIWYG && caretInLn {
		t.context.oldtip.Tokens = lex.TrimWhitespace(t.context.oldtip.Tokens)
		t.context.oldtip.AppendTokens([]byte(caret))
	}

	return
}
