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

type Paragraph struct {
	*BaseNode

	OpenTag, CloseTag string
}

func (t *Tree) parseParagraph(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeParagraph}
	p := &Paragraph{baseNode, "<p>", "</p>"}
	ret = p

	for {
		_, line = line.trimLeft()
		p.tokens = append(p.tokens, line...)

		line = t.nextLine()
		if line.isBlankLine() {
			t.backupLine(line)
			break
		}

		if level := t.isSetextHeading(line); 0 < level {
			ret = t.parseSetextHeading(p, level)

			return
		}

		tokens := t.indentOffset(line, t.context.IndentSpaces)
		if t.interruptParagraph(tokens) {
			t.backupLine(line)

			break
		}
	}
	p.tokens = p.tokens.trimRight()

	return
}

func (t *Tree) interruptParagraph(line items) bool {
	if t.isIndentCode(line) {
		return false
	}

	if t.isThematicBreak(line) {
		return true
	}

	level := 0
	if t.isATXHeading(line, &level) {
		return true
	}

	if isList, _ := t.isList(line); isList {
		if line[1:].isBlankLine() {
			return false
		}

		return true
	}

	if t.isFencedCode(line) {
		return true
	}

	pos := line.index(itemGreater)
	if 0 < pos {
		maybeTag := line[:pos+1]
		htmlType := -1
		if t.isHTML(maybeTag, &htmlType) && 7 != htmlType {
			return true
		}
	}

	if 0 < t.context.BlockquoteLevel {
		tokens := t.removeStartBlockquoteMarker(line, t.context.BlockquoteLevel)
		if tokens.isBlankLine() {
			return true
		}
	}

	if t.isBlockquote(line) {
		return true
	}

	return false
}
