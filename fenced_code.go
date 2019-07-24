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

func (t *Tree) parseFencedCode(line items) {
	marker := line[0]
	n := line.accept(marker.typ)
	line = line[n:]
	infoStr := line.trim().rawText()
	if "" != infoStr {
		infoStr = strings.Split(infoStr, " ")[0]
	}
	baseNode := &BaseNode{typ: NodeCode}
	code := &CodeBlock{baseNode, "", infoStr}

	line = t.nextLine()
	if line.isEOF() {
		return
	}
	if t.isFencedCodeClose(line, marker, n) {
		return
	}

	var codeValue string
	for {
		codeValue += line.rawText()
		code.tokens = append(code.tokens, line...)

		line = t.nextLine()
		if t.isFencedCodeClose(line, marker, n) {
			break
		}
	}

	code.Value = codeValue
}

func (t *Tree) isFencedCodeClose(line items, openMarker *item, num int) bool {
	if line.isEOF() {
		return true
	}

	closeMarker := line[0]
	if closeMarker.typ != openMarker.typ {
		return false
	}
	if num > line.accept(closeMarker.typ) {
		return false
	}
	if !line.trim().allAre(openMarker.typ) {
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
