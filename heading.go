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

	p.tokens = p.tokens.trimRight()
	text := &Text{BaseNode: &BaseNode{typ: NodeText, tokens: p.tokens}}
	ret.AppendChild(ret, text)

	return
}

func (t *Tree) parseATXHeading(line items, level int) (ret Node) {
	baseNode := &BaseNode{typ: NodeHeading}
	heading := &Heading{baseNode, level}
	ret = heading

	tokens := line.trimLeft()
	tokens = tokens[level:].trimLeft()
	for _, token := range tokens {
		if itemEOF == token.typ || itemNewline == token.typ {
			break
		}

		heading.tokens = append(heading.tokens, token)
	}

	heading.tokens = heading.tokens.trimRight()
	closingCrosshatchIndex := len(heading.tokens) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if itemCrosshatch == heading.tokens[closingCrosshatchIndex].typ {
			continue
		}

		if itemSpace == heading.tokens[closingCrosshatchIndex].typ {
			break
		} else {
			closingCrosshatchIndex = len(heading.tokens)
			break
		}
	}

	if 0 < closingCrosshatchIndex {
		heading.tokens = heading.tokens[:closingCrosshatchIndex]
		heading.tokens = heading.tokens.trimRight()
	}

	return
}

func (t *Tree) isATXHeading(line items, level *int) bool {
	len := len(line)
	if 2 > len { // at least # and newline
		return false
	}

	index, marker := line.firstNonSpace()
	if itemCrosshatch != marker.typ {
		return false
	}

	line = line[index:]
	*level = line.accept(itemCrosshatch)
	if 6 < *level {
		return false
	}

	if !line[*level].isWhitespace() {
		return false
	}

	return true
}

func (t *Tree) isSetextHeading(line items) (level int) {
	spaces, line := line.trimLeftSpace()
	if 3 < spaces {
		return
	}

	line = line.trimRight()
	length := len(line)
	marker := line[0]
	if itemHyphen != marker.typ && itemEqual != marker.typ {
		return
	}

	for i := 1; i < length; i++ {
		token := line[i]
		if marker.typ != token.typ {
			return
		}
	}

	parentType := t.context.CurNode.Type()
	if NodeBlockquote == parentType || NodeListItem == parentType {
		return
	}

	if itemEqual == marker.typ {
		level = 1
	} else {
		level = 2
	}

	return
}
