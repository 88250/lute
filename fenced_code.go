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

func (t *Tree) parseFencedCode(line items) (ret Node) {
	_, line = line.trimLeftSpace()
	marker := line[0]
	n := line.accept(marker.typ)
	line = line[n:]
	infoStr := line.trim().rawText()
	line = t.nextLine()
	baseNode := &BaseNode{typ: NodeCode}
	code := &Code{baseNode, "", infoStr}
	var codeValue string

	for {
		var spaces, tabs int
		for i := 0; i < n && i < len(line); i++ {
			token := line[i]
			if itemSpace == token.typ {
				spaces++
			} else if itemTab == token.typ {
				tabs++
			}
			if 3 < spaces || 0 < tabs {
				line = line[i+1:]
				break
			}
		}

		for i := 0; i < len(line); i++ {
			token := line[i]
			codeValue += EscapeHTML(token.val)
			//if token.isNewline() {
			//	newlines, nonNewline := t.nonNewline()
			//	if nonNewline.isEOF() {
			//		break Loop
			//	}
			//
			//	codeValue += newlines.rawText()
			//	newlines = append(newlines, token)
			//	code.tokens = append(code.tokens, newlines...)
			//	line = nonNewline
			//	spaces, tabs, _, _ := t.nonWhitespace(line)
			//	if 1 > tabs && 4 > spaces {
			//		t.backupLine(line)
			//		break Loop
			//	} else {
			//		continue Loop
			//	}
			//}
			code.tokens = append(code.tokens, token)
		}

		line = t.nextLine()
		if t.isFencedCodeClose(line, marker, n) {
			break
		}
	}

	code.Value = codeValue
	ret = code

	return
}

func (t *Tree) isFencedCodeClose(line items, openMarker *item, num int) bool {
	spaces, line := line.trimLeftSpace()
	if t.context.IndentSpaces+3 < spaces {
		return false
	}

	closeMarker := line[0]
	if closeMarker.typ != openMarker.typ {
		return false
	}
	if num > line.accept(closeMarker.typ) {
		return false
	}

	return true
}

func (t *Tree) isFencedCode(line items) bool {
	spaces, line := line.trimLeftSpace()
	if t.context.IndentSpaces+3 < spaces {
		return false
	}

	marker := line[0]
	if itemBacktick != marker.typ && itemTilde != marker.typ {
		return false
	}

	if 3 > line.accept(marker.typ) {
		return false
	}

	return true
}
