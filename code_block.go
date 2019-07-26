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

type CodeBlock struct {
	*BaseNode
	InfoStr string

	IsFenced    bool
	FenceChar   *item
	FenceLength int
	FenceOffset int
}

func (codeBlock *CodeBlock) Continue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent
	if codeBlock.IsFenced {
		if indent <= 3 && codeBlock.isFencedCodeClose(ln[context.nextNonspace:], codeBlock.FenceChar, codeBlock.FenceLength) {
			// closing fence - we're at end of line, so we can return
			context.finalize(codeBlock)
			return 2
		} else {
			// skip optional spaces of fence offset
			var i = codeBlock.FenceOffset
			for i > 0 && ln.peek(context.offset).isSpaceOrTab() {
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

func (codeBlock *CodeBlock) Finalize() {
	if codeBlock.IsFenced {
		// first line becomes info string
		var content = codeBlock.value
		var newlinePos = strings.Index(content, "\n")
		var firstLine = content[:newlinePos]
		var rest = content[newlinePos+1:]
		codeBlock.InfoStr = unescapeString(strings.TrimSpace(firstLine))
		codeBlock.value = rest
	} else { // indented
		codeBlock.value = strings.TrimRight(codeBlock.value, "\n ") + "\n"
	}
	codeBlock.tokens = nil
}

func (codeBlock *CodeBlock) AcceptLines() bool {
	return true
}

func (codeBlock *CodeBlock) CanContain(nodeType NodeType) bool {
	return false
}

func (codeBlock *CodeBlock) isFencedCodeClose(tokens items, openMarker *item, num int) bool {
	closeMarker := tokens[0]
	if closeMarker.typ != openMarker.typ {
		return false
	}
	if num > tokens.accept(closeMarker.typ) {
		return false
	}
	if !tokens.trim().allAre(openMarker.typ) {
		return false
	}

	return true
}

func (t *Tree) isFencedCode(line items) bool {
	if 3 > len(line) {
		return false
	}

	marker := line[0]
	if itemBacktick != marker.typ && itemTilde != marker.typ {
		return false
	}

	pos := line.accept(marker.typ)
	if 3 > pos {
		return false
	}

	infoStr := line[pos:]
	if itemBacktick == marker.typ && infoStr.contain(itemBacktick) {
		return false
	}

	return true
}
