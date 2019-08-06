// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import "strings"

// CodeBlock 描述了代码块节点结构。
type CodeBlock struct {
	*BaseNode
	isFenced    bool
	fenceChar   item
	fenceLength int
	fenceOffset int
	info        string
}

func (codeBlock *CodeBlock) Continue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent
	if codeBlock.isFenced {
		if indent <= 3 && codeBlock.isFencedCodeClose(ln[context.nextNonspace:], codeBlock.fenceChar, codeBlock.fenceLength) {
			// closing fence - we're at end of line, so we can return
			context.finalize(codeBlock)
			return 2
		} else {
			// skip optional spaces of fence offset
			var i = codeBlock.fenceOffset
			var token item
			for i > 0 {
				token = ln.peek(context.offset)
				if itemSpace != token && itemTab != token {
					break
				}
				context.advanceOffset(1, true)
				i--
			}
		}
	} else { // indented
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

func (codeBlock *CodeBlock) Finalize(context *Context) {
	if codeBlock.isFenced {
		// first line becomes info string
		content := codeBlock.tokens
		var i int
		var token item
		for ; ; i++ {
			token = codeBlock.tokens[i]
			if itemNewline == token {
				break
			}
		}
		firstLine := content[:i]
		rest := content[i+1:]

		codeBlock.info = unescapeString(strings.TrimSpace(firstLine.rawText()))
		codeBlock.tokens = rest
	} else { // indented
		codeBlock.tokens = codeBlock.tokens.replaceNewlineSpace()
	}
	codeBlock.value = codeBlock.tokens.rawText()
	codeBlock.tokens = nil
}

func (codeBlock *CodeBlock) AcceptLines() bool {
	return true
}

func (codeBlock *CodeBlock) CanContain(nodeType int) bool {
	return false
}

func (t *Tree) parseFencedCode() (ret *CodeBlock) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemBacktick != marker && itemTilde != marker {
		return nil
	}

	fenceChar := marker
	fenceLength := 0
	for i := t.context.nextNonspace; fenceChar == t.context.currentLine[i]; i++ {
		fenceLength++
	}

	if 3 > fenceLength {
		return nil
	}

	var info string
	infoTokens := t.context.currentLine[t.context.nextNonspace+fenceLength:]
	if itemBacktick == marker {
		if !infoTokens.contain(itemBacktick) {
			info = infoTokens.trim().rawText()
		} else {
			return nil // info 部分不能包含 `
		}
	} else {
		info = infoTokens.trim().rawText()
	}

	info = unescapeString(info)
	ret = &CodeBlock{&BaseNode{typ: NodeCodeBlock, tokens: make([]item, 0, 256)},
		true, fenceChar, fenceLength, t.context.indent, info}

	return
}

func (codeBlock *CodeBlock) isFencedCodeClose(tokens items, openMarker item, num int) bool {
	closeMarker := tokens[0]
	if closeMarker != openMarker {
		return false
	}
	if num > tokens.accept(closeMarker) {
		return false
	}
	for _, token := range tokens.trim() {
		if token != openMarker {
			return false
		}
	}

	return true
}
