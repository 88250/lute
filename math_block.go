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

func (mathBlock *Node) mathBlockContinue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent

	if indent <= 3 && mathBlock.isMathBlockClose(ln[context.nextNonspace:]) {
		context.finalize(mathBlock, context.lineNum)
		return 2
	} else {
		// 跳过 $ 之前可能存在的空格
		var i = mathBlock.mathBlockDollarOffset
		var token byte
		for i > 0 {
			token = peek(ln, context.offset)
			if itemSpace != token && itemTab != token {
				break
			}
			context.advanceOffset(1, true)
			i--
		}
	}
	return 0
}

var mathBlockMarker = strToBytes("$$")

func (mathBlock *Node) mathBlockFinalize(context *Context) {
	tokens := mathBlock.tokens[2:] // 剔除开头的两个 $$
	tokens = trimWhitespace(tokens)
	if bytes.HasSuffix(tokens, mathBlockMarker) {
		tokens = tokens[:len(tokens)-2] // 剔除结尾的两个 $$
	}
	mathBlock.tokens = nil
	mathBlock.AppendChild(&Node{typ: NodeMathBlockOpenMarker})
	mathBlock.AppendChild(&Node{typ: NodeMathBlockContent, tokens: tokens})
	mathBlock.AppendChild(&Node{typ: NodeMathBlockCloseMarker})
}

func (t *Tree) parseMathBlock() (ok bool, mathBlockDollarOffset int) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemDollar != marker {
		return
	}

	fenceChar := marker
	fenceLength := 0
	for i := t.context.nextNonspace; i < t.context.currentLineLen && fenceChar == t.context.currentLine[i]; i++ {
		fenceLength++
	}

	if 2 > fenceLength {
		return
	}

	return true, t.context.indent
}

func (mathBlock *Node) isMathBlockClose(tokens []byte) bool {
	closeMarker := tokens[0]
	if closeMarker != itemDollar {
		return false
	}
	if 2 > accept(tokens, closeMarker) {
		return false
	}
	tokens = trimWhitespace(tokens)
	for _, token := range tokens {
		if token != itemDollar {
			return false
		}
	}
	return true
}
