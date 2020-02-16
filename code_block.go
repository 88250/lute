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

import "bytes"

func (codeBlock *Node) CodeBlockContinue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent
	if codeBlock.IsFencedCodeBlock {
		if ok, closeFence := codeBlock.isFencedCodeClose(ln[context.nextNonspace:], codeBlock.CodeBlockFenceChar, codeBlock.CodeBlockFenceLen); indent <= 3 && ok {
			codeBlock.CodeBlockCloseFence = closeFence
			context.finalize(codeBlock, context.lineNum)
			return 2
		} else {
			// 跳过围栏标记符之前可能存在的空格
			var i = codeBlock.CodeBlockFenceOffset
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
	} else { // 缩进代码块
		if indent >= 4 {
			context.advanceOffset(4, true)
		} else if context.blank {
			context.advanceNextNonspace()
		} else {
			return 1
		}
	}
	return 0
}

func (codeBlock *Node) CodeBlockFinalize(context *Context) {
	if codeBlock.IsFencedCodeBlock {
		content := codeBlock.Tokens
		length := len(content)
		if 1 > length {
			return
		}

		var i int
		for ; i < length; i++ {
			if itemNewline == content[i] {
				break
			}
		}
		codeBlock.Tokens = content[i+1:]
	} else { // 缩进代码块
		codeBlock.Tokens = replaceNewlineSpace(codeBlock.Tokens)
	}
}

var codeBlockBacktick = strToBytes("`")

func (t *Tree) parseFencedCode() (ok bool, fenceChar byte, fenceLen int, fenceOffset int, openFence, codeBlockInfo []byte) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemBacktick != marker && itemTilde != marker {
		return
	}

	fenceChar = marker
	for i := t.context.nextNonspace; i < t.context.currentLineLen && fenceChar == t.context.currentLine[i]; i++ {
		fenceLen++
	}

	if 3 > fenceLen {
		return
	}

	openFence = t.context.currentLine[t.context.nextNonspace : t.context.nextNonspace+fenceLen]

	var info []byte
	infoTokens := t.context.currentLine[t.context.nextNonspace+fenceLen:]
	if itemBacktick == marker && bytes.Contains(infoTokens, codeBlockBacktick) {
		// info 部分不能包含 `
		return
	}
	info = trimWhitespace(infoTokens)
	info = unescapeString(info)
	return true, fenceChar, fenceLen, t.context.indent, openFence, info
}

func (codeBlock *Node) isFencedCodeClose(tokens []byte, openMarker byte, num int) (ok bool, closeFence []byte) {
	closeMarker := tokens[0]
	if closeMarker != openMarker {
		return false, nil
	}
	if num > accept(tokens, closeMarker) {
		return false, nil
	}
	tokens = trimWhitespace(tokens)
	for _, token := range tokens {
		if token != openMarker {
			return false, nil
		}
	}
	closeFence = tokens
	return true, closeFence
}
