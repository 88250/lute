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

import (
	"bytes"
)

func (mathBlock *Node) mathBlockContinue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent

	if indent <= 3 && mathBlock.isMathBlockClose(ln[context.nextNonspace:]) {
		context.finalize(mathBlock)
		return 2
	} else {
		// 跳过 $ 之前可能存在的空格
		var i = mathBlock.mathBlockDollarOffset
		var token byte
		for i > 0 {
			token = ln.peek(context.offset)
			if itemSpace != token && itemTab != token {
				break
			}
			context.advanceOffset(1, true)
			i--
		}
	}
	return 0
}

func (mathBlock *Node) mathBlockFinalize(context *Context) {
	tokens := bytes.TrimSpace(mathBlock.tokens)
	closing := items{itemDollar, itemDollar}
	if !bytes.HasSuffix(tokens, closing) {
		tokens = append(tokens, itemNewline)
		tokens = append(tokens, closing...)
	}
	mathBlock.tokens = tokens
}

var mathBlockDollar = items{itemDollar}

func (t *Tree) parseMathBlock() (ret *Node) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemDollar != marker {
		return nil
	}

	fenceChar := marker
	fenceLength := 0
	for i := t.context.nextNonspace; i < t.context.currentLineLen && fenceChar == t.context.currentLine[i]; i++ {
		fenceLength++
	}

	if 2 > fenceLength {
		return nil
	}

	ret = &Node{typ: NodeMathBlock, tokens: make(items, 0, 256), mathBlockDollarOffset: t.context.indent}

	return
}

func (mathBlock *Node) isMathBlockClose(tokens items) bool {
	closeMarker := tokens[0]
	if closeMarker != itemDollar {
		return false
	}
	if 2 > tokens.accept(closeMarker) {
		return false
	}
	tokens = bytes.TrimSpace(tokens)
	for _, token := range tokens {
		if token != itemDollar {
			return false
		}
	}
	return true
}
