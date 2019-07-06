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

type Heading struct {
	*BaseNode

	Level int
}

func (t *Tree) parseSetextHeading(p *Paragraph, level int) (ret Node) {
	baseNode := &BaseNode{typ: NodeHeading}
	ret = &Heading{baseNode, level}

	p.tokens = t.trimRight(p.tokens)
	p.rawText = p.tokens.rawText()
	text := &Text{BaseNode: &BaseNode{typ: NodeText, tokens: p.tokens}}
	ret.AppendChild(ret, text)

	return
}

func (t *Tree) parseATXHeading(line items, level int) (ret Node) {
	baseNode := &BaseNode{typ: NodeHeading}
	heading := &Heading{baseNode, level}
	ret = heading

	tokens := t.skipWhitespaces(line[level:])
	for _, token := range tokens {
		if itemEOF == token.typ || itemNewline == token.typ {
			break
		}

		heading.rawText += token.val
		heading.tokens = append(heading.tokens, token)
	}

	return
}

func (t *Tree) isATXHeading(line items) (level int) {
	if 2 > len(line) { // at least # and newline
		return
	}

	index, marker := t.firstNonSpace(line)
	if itemCrosshatch != marker.typ {
		return
	}

	line = line[index:]
	level = t.accept(line, itemCrosshatch)
	if !line[level].isWhitespace() {
		return
	}

	return
}

func (t *Tree) isSetextHeading(line items) (level int) {
	tokens := t.removeSpacesTabs(line)
	tokens = tokens[:len(tokens)-1] // remove tailing newline
	length := len(tokens)
	marker := tokens[0]
	if itemHyphen != marker.typ && itemEqual != marker.typ {
		return
	}

	for i := 1; i < length; i++ {
		token := tokens[i]
		if marker.typ != token.typ {
			return
		}
	}

	if itemEqual == marker.typ {
		level = 1
	} else {
		level = 2
	}

	return
}
