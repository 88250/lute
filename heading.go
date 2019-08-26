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

// Heading 描述了标题节点结构。
type Heading struct {
	*BaseNode
	Level int // 1~6
}

func (heading *Heading) Continue(context *Context) int {
	return 1
}

func (heading *Heading) CanContain(nodeType int) bool {
	return false
}

func (t *Tree) parseATXHeading() (ret *Heading) {
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

	heading := &Heading{&BaseNode{typ: NodeHeading}, level}
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

func (t *Tree) parseSetextHeading() (ret *Heading) {
	start := t.context.nextNonspace
	marker := t.context.currentLine[start]
	if itemEqual != marker && itemHyphen != marker {
		return nil
	}

	end := t.context.currentLineLen - 2
	for ; 0 <= end; end-- {
		if token := t.context.currentLine[end]; itemSpace != token && itemTab != token {
			break
		}
	}

	markers := 0
	for ; start < end; start++ {
		token := t.context.currentLine[start]
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

	ret = &Heading{&BaseNode{typ: NodeHeading}, 1}
	if itemHyphen == marker {
		ret.Level = 2
	}

	return
}
