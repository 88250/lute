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
	return
}
