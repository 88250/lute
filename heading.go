// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

import (
	"bytes"
)

func (t *Tree) parseATXHeading() (ok bool, markers, content []byte, level int) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if itemCrosshatch != marker {
		return
	}

	level = accept(tokens, itemCrosshatch)
	if 6 < level {
		return
	}

	if level < len(tokens) && !isWhitespace(tokens[level]) {
		return
	}

	markers = t.context.currentLine[t.context.nextNonspace : t.context.nextNonspace+level+1]

	content = make([]byte, 0, 256)
	_, tokens = trimLeft(tokens)
	_, tokens = trimLeft(tokens[level:])
	for _, token := range tokens {
		if itemNewline == token {
			break
		}
		content = append(content, token)
	}

	_, content = trimRight(content)
	closingCrosshatchIndex := len(content) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if itemCrosshatch == content[closingCrosshatchIndex] {
			continue
		}
		if itemSpace == content[closingCrosshatchIndex] {
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
		_, content = trimRight(content)
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
	ln := trimWhitespace(t.context.currentLine)
	start := 0
	marker := ln[start]
	if itemEqual != marker && itemHyphen != marker {
		return
	}

	var caretInLn bool
	if t.context.option.VditorWYSIWYG {
		if bytes.Contains(ln, []byte(caret)) {
			caretInLn = true
			ln = bytes.ReplaceAll(ln, []byte(caret), []byte(""))
		}
	}

	markers := 0
	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if itemEqual != token && itemHyphen != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
		markers++
	}

	level = 1
	if itemHyphen == marker {
		level = 2
	}

	if t.context.option.VditorWYSIWYG && caretInLn {
		t.context.oldtip.tokens = trimWhitespace(t.context.oldtip.tokens)
		t.context.oldtip.AppendTokens([]byte(caret))
	}

	return
}
