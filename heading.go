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

func (t *Tree) parseATXHeading() (ret *BaseNode) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if itemCrosshatch != marker {
		return
	}

	level := tokens.accept(itemCrosshatch)
	if 6 < level {
		return
	}

	if level < len(tokens) && !isWhitespace(tokens[level]) {
		return
	}

	heading := &BaseNode{typ: NodeHeading, headingLevel: level}
	tokens = bytes.TrimLeft(tokens, " \t\n")
	tokens = bytes.TrimLeft(tokens[level:], " \t\n")
	for _, token := range tokens {
		if itemEnd == token || itemNewline == token {
			break
		}

		heading.tokens = append(heading.tokens, token)
	}

	heading.tokens = bytes.TrimRight(heading.tokens, " \t\n")
	closingCrosshatchIndex := len(heading.tokens) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if itemCrosshatch == heading.tokens[closingCrosshatchIndex] {
			continue
		}

		if itemSpace == heading.tokens[closingCrosshatchIndex] {
			break
		} else {
			closingCrosshatchIndex = len(heading.tokens)
			break
		}
	}

	if 0 >= closingCrosshatchIndex {
		heading.tokens = nil
	} else if 0 < closingCrosshatchIndex {
		heading.tokens = heading.tokens[:closingCrosshatchIndex]
		heading.tokens = bytes.TrimRight(heading.tokens, " \t\n")
	}

	ret = heading

	return
}

func (t *Tree) parseSetextHeading() (ret *BaseNode) {
	ln := bytes.TrimSpace(t.context.currentLine)
	start := 0
	marker := ln[start]
	if itemEqual != marker && itemHyphen != marker {
		return nil
	}

	markers := 0
	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if itemEqual != token && itemHyphen != token {
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

	if itemEnd == marker {
		return nil
	}

	ret = &BaseNode{typ: NodeHeading, headingLevel: 1}
	if itemHyphen == marker {
		ret.headingLevel = 2
	}

	return
}
