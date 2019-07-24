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

func (heading *Heading) CanContain(node Node) bool {
	return false
}

func (t *Tree) parseSetextHeading(p *Paragraph, level int) {
	baseNode := &BaseNode{typ: NodeHeading}
	heading := &Heading{baseNode, level}

	p.tokens = p.tokens.trimRight()
	text := &Text{BaseNode: &BaseNode{typ: NodeText, tokens: p.tokens}}
	heading.AppendChild(heading, text)

	return
}

func (t *Tree) parseATXHeading(tokens items) (ret Node) {
	if 2 > len(tokens) {
		return
	}

	index, marker := tokens.firstNonSpace()
	if itemCrosshatch != marker.typ {
		return
	}

	tokens = tokens[index:]
	level := tokens.accept(itemCrosshatch)
	if 6 < level {
		return
	}

	if !tokens[level].isWhitespace() {
		return
	}

	baseNode := &BaseNode{typ: NodeHeading}
	heading := &Heading{baseNode, level}

	_, tokens = tokens.trimLeft()
	_, tokens = tokens[level:].trimLeft()
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

	if 0 >= closingCrosshatchIndex {
		heading.tokens = nil
	} else if 0 < closingCrosshatchIndex {
		heading.tokens = heading.tokens[:closingCrosshatchIndex]
		heading.tokens = heading.tokens.trimRight()
	}

	ret = heading

	return
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

	if 1 == length && itemHyphen == marker.typ {
		return
	}

	for i := 1; i < length; i++ {
		token := line[i]
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
