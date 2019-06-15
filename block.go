// Lute - A structural markdown engine.
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

func (t *Tree) parseBlocks() {
	curNode := t.context.CurNode
	count := 0
	for line := t.nextLine(); ; {
		t.parseBlock(line)
		t.context.CurNode = curNode

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		count++
		if count > 2 {
			break
		}
	}
}

func (t *Tree) parseBlock(line items) (ret Node) {
	curNode := t.context.CurNode

	if t.isThematicBreak(line) {
		ret = t.parseThematicBreak(line)
	} else if t.isList(line) {
		ret = t.parseList(line)
	} else if t.isATXHeading(line) {
		ret = t.parseHeading(line)
	} else if t.isBlockquote(line) {
		ret = t.parseBlockquote(line)
	} else if t.isIndentCode(line) {
		ret = t.parseIndentCode(line)
	} else if t.isBlankLine(line) {
		return
	} else {
		ret = t.parseParagraph(line)
	}

	curNode.Append(ret)

	return
}
