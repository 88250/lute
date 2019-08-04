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
	isFenced    bool
	fenceChar   string
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

func (codeBlock *CodeBlock) Finalize(context *Context) {
	if codeBlock.isFenced {
		// first line becomes info string
		content := codeBlock.value
		newlinePos := strings.Index(content, "\n")
		firstLine := content[:newlinePos]
		rest := content[newlinePos+1:]
		codeBlock.info = unescapeString(strings.TrimSpace(firstLine))
		codeBlock.value = rest
	} else { // indented
		i := len(codeBlock.value) - 1
		for ; 0 <= i && ('\n' == codeBlock.value[i] || ' ' == codeBlock.value[i]); i-- {
		}
		i++
		rest := codeBlock.value[i:]
		for 0 <= strings.Index(rest, "\n ") {
			rest = strings.ReplaceAll(rest, "\n ", "\n")
		}
		for 0 <= strings.Index(rest, "\n\n") {
			rest = strings.ReplaceAll(rest, "\n\n", "\n")
		}
		codeBlock.value = codeBlock.value[:i] + rest
	}
	codeBlock.tokens = nil
}

func (codeBlock *CodeBlock) AcceptLines() bool {
	return true
}

func (codeBlock *CodeBlock) CanContain(nodeType NodeType) bool {
	return false
}

func (t *Tree) parseFencedCode() (ret *CodeBlock) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemBacktick != marker.typ && itemTilde != marker.typ {
		return nil
	}

	fenceChar := marker.Value()
	fenceLength := 0
	for i := t.context.nextNonspace; fenceChar == t.context.currentLine[i].Value(); i++ {
		fenceLength++
	}

	if 3 > fenceLength {
		return nil
	}

	var info string
	infoTokens := t.context.currentLine[t.context.nextNonspace+fenceLength:]
	if itemBacktick == marker.typ {
		if !infoTokens.contain(itemBacktick) {
			info = infoTokens.trim().rawText()
		} else {
			return nil // info 部分不能包含 `
		}
	} else {
		info = infoTokens.trim().rawText()
	}

	info = unescapeString(info)
	ret = &CodeBlock{&BaseNode{typ: NodeCodeBlock},
		true, fenceChar, fenceLength, t.context.indent, info}

	return
}

func (codeBlock *CodeBlock) isFencedCodeClose(tokens items, openMarker string, num int) bool {
	closeMarker := tokens[0]
	if closeMarker.Value() != openMarker {
		return false
	}
	if num > tokens.accept(closeMarker.typ) {
		return false
	}
	for _, token := range tokens.trim() {
		if token.Value() != openMarker {
			return false
		}
	}

	return true
}
