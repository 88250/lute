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

func (t *Tree) parseATXHeading() (ok bool, markers, content items, level int) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if itemCrosshatch != term(marker) {
		return
	}

	level = tokens.accept(itemCrosshatch)
	if 6 < level {
		return
	}

	if level < len(tokens) && !isWhitespace(term(tokens[level])) {
		return
	}

	markers = t.context.currentLine[t.context.nextNonspace : t.context.nextNonspace+level+1]

	content = make(items, 0, 256)
	_, tokens = trimLeft(tokens)
	_, tokens = trimLeft(tokens[level:])
	for _, token := range tokens {
		if itemNewline == term(token) {
			break
		}
		content = append(content, token)
	}

	_, content = trimRight(content)
	closingCrosshatchIndex := len(content) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if itemCrosshatch == term(content[closingCrosshatchIndex]) {
			continue
		}
		if itemSpace == term(content[closingCrosshatchIndex]) {
			break
		} else {
			closingCrosshatchIndex = len(content)
			break
		}
	}

	if 0 >= closingCrosshatchIndex {
		content = make(items, 0, 0)
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
	if itemEqual != term(marker) && itemHyphen != term(marker) {
		return
	}

	markers := 0
	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if itemEqual != term(token) && itemHyphen != term(token) {
			return
		}

		if 0 != term(marker) {
			if term(marker) != term(token) {
				return
			}
		} else {
			marker = token
		}
		markers++
	}

	level = 1
	if itemHyphen == term(marker) {
		level = 2
	}
	return
}
