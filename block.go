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

func (t *Tree) parseBlocks() {
	curNode := t.context.CurNode
	for line := t.nextLine(); ; {
		n := t.parseBlock(line)
		if nil != n {
			curNode.AppendChild(curNode, n)
		}
		t.context.CurNode = curNode

		line = t.nextLine()
		if line.isEOF() {
			break
		}
	}
}

func (t *Tree) parseBlock(line items) (ret Node) {
	if t.isIndentCode(line) {
		ret = t.parseIndentCode(line)
	} else if t.isThematicBreak(line) {
		ret = t.parseThematicBreak(line)
	} else if level := t.isATXHeading(line); 0 < level {
		ret = t.parseATXHeading(line, level)
	} else if t.isBlockquote(line) {
		ret = t.parseBlockquote(line)
	} else if t.isList(line) {
		ret = t.parseList(line)
	} else if t.isBlankLine(line) {
		return
	} else {
		ret = t.parseParagraph(line)
	}

	return
}
