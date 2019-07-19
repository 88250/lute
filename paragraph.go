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

func (t *Tree) parseParagraph(line items) {
	baseNode := &BaseNode{typ: NodeParagraph}
	p := &Paragraph{baseNode, "<p>", "</p>"}

	for {
		_, line = line.trimLeft()
		p.tokens = append(p.tokens, line...)

		line = t.nextLine()
		if line.isBlankLine() {
			t.backupLine(line)
			break
		}

		if level := t.isSetextHeading(line); 0 < level {
			t.parseSetextHeading(p, level)

			return
		}

		startIndentSpaces := line.spaceCountLeft()

		tokens := t.indentOffset(line, t.context.IndentSpaces)
		if isInterrup, tokens := t.interruptParagraph(startIndentSpaces, tokens); isInterrup {
			t.backupLine(line)

			break
		} else {
			line = tokens
		}
	}
	p.tokens = p.tokens.trimRight()
	t.context.AppendChild(p)

	return
}

func (t *Tree) interruptParagraph(startIndentSpaces int, line items) (ret bool, tokens items) {
	tokens = line
	if t.isIndentCode(line) {
		return
	}

	if t.isThematicBreak(line) {
		ret = true
		return
	}

	level := 0
	if t.isATXHeading(line, &level) {
		ret = true
		return
	}

	if isList, marker := t.isList(line); isList {
		if t.context.CurrentContainer().Is(NodeListItem) {
			if 2 < t.context.IndentSpaces && 3 < startIndentSpaces && t.context.IndentSpaces > startIndentSpaces {
				return
			}

			ret = true
			return
		}

		markerLen := len(marker)
		if line[markerLen:].isBlankLine() {
			return
		}

		_, marker, delim, _, _ := t.parseListItemMarker(line, nil)
		if " " != delim && "1" != marker[:markerLen-1] {
			return
		}

		ret = true
		return
	}

	if t.isFencedCode(line) {
		ret = true
		return
	}

	pos := line.index(itemGreater)
	if 0 < pos {
		maybeTag := line[:pos+1]
		htmlType := -1
		if t.isHTML(maybeTag, &htmlType) && 7 != htmlType {
			ret = true
			return
		}
	}

	container := t.context.CurrentContainer()
	if container.Is(NodeBlockquote) {
		blockquote := container.(*Blockquote)
		level := t.blockquoteMarkerCount(line)
		if 0 == level {
			return
		}

		if blockquote.level != level {
			ret = true
			return
		} else {
			tokens = t.decBlockquoteMarker(line)
			return
		}
	} else {
		if t.isBlockquote(line) {
			ret = true
			return
		}
	}

	return
}
