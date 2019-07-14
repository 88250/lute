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

func (t *Tree) parseIndentCode(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeCode}
	code := &Code{baseNode, "", ""}

	var chunks []items
	for {
		var spaces, tabs int
		for i := 0; i < 4; i++ {
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

		chunk := items{}
		chunk = append(chunk, line...)
		newlines, nonNewline := t.nonNewline()
		line = nonNewline
		if !t.isIndentCode(line) {
			chunks = append(chunks, chunk)
			if 0 < len(newlines) {
				t.backupLine(newlines)
			}
			t.backupLine(line)
			break
		}

		if 0 < len(newlines) {
			chunk = append(chunk, newlines...)
		}
		chunks = append(chunks, chunk)

		if t.blockquoteMarkerCount(line) < t.context.BlockquoteLevel {
			t.backupLine(line)
			break
		}

	}

	if 1 > len(chunks) {
		return nil
	}

	for _, chunk := range chunks {
		code.Value += chunk.rawText()
	}

	ret = code

	return
}

func (t *Tree) isIndentCode(line items) bool {
	if line.isBlankLine() {
		return false
	}

	var spaces int
	for _, token := range line {
		if itemSpace == token.typ {
			spaces++
			continue
		}
		if itemTab == token.typ {
			spaces += 4
			continue
		}

		break
	}

	return t.context.IndentSpaces+3 < spaces
}

func (t *Tree) nonNewline() (newlines items, line items) {
	for line = t.nextLine(); line.isBlankLine() && !line.isEOF(); line = t.nextLine() {
		if 5 > len(line) {
			_, line = line.trimLeftSpace()
		} else {
			line = t.indentOffset(line, 4)
		}

		newlines = append(newlines, line...)
	}

	return
}
