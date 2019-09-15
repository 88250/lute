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

import "bytes"

func (t *Tree) parseATXHeading() (content items, level int) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if itemCrosshatch != marker {
		return
	}

	level = tokens.accept(itemCrosshatch)
	if 6 < level {
		return
	}

	if level < len(tokens) && !isWhitespace(tokens[level]) {
		return
	}

	content = make(items, 0, 256)

	tokens = bytes.TrimLeft(tokens, " \t\n")
	tokens = bytes.TrimLeft(tokens[level:], " \t\n")
	for _, token := range tokens {
		if itemEnd == token || itemNewline == token {
			break
		}

		content = append(content, token)
	}

	content = bytes.TrimRight(content, " \t\n")
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
		content = make(items, 0, 0)
	} else if 0 < closingCrosshatchIndex {
		content = content[:closingCrosshatchIndex]
		content = bytes.TrimRight(content, " \t\n")
	}

	return
}

func (t *Tree) parseSetextHeading() (level int) {
	ln := bytes.TrimSpace(t.context.currentLine)
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

		if itemEnd != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
		markers++
	}

	if itemEnd == marker {
		return
	}

	level = 1
	if itemHyphen == marker {
		level = 2
	}
	return
}
