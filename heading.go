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

func (heading *Heading) Continue(context *Context) int {
	return 1
}

func (heading *Heading) CanContain(nodeType NodeType) bool {
	return false
}

func (t *Tree) parseATXHeading() (ret *Heading) {
	tokens := t.context.currentLine[t.context.nextNonspace:]
	marker := tokens[0]
	if itemCrosshatch != marker.typ {
		return
	}

	level := tokens.accept(itemCrosshatch)
	if 6 < level {
		return
	}

	if !tokens[level].isWhitespace() {
		return
	}

	heading := &Heading{&BaseNode{typ: NodeHeading}, level}
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

func (t *Tree) parseSetextHeading() (ret *Heading) {
	markers := 0
	var marker *item
	for i := t.context.nextNonspace; i < t.context.currentLineLen-1; i++ {
		token := t.context.currentLine[i]
		if itemSpace == token.typ || itemTab == token.typ {
			continue
		}

		if itemEqual != token.typ && itemHyphen != token.typ {
			return nil
		}

		if nil != marker {
			if marker.typ != token.typ {
				return nil
			}
		} else {
			marker = token
		}
		markers++
	}

	if nil == marker {
		return nil
	}

	ret = &Heading{&BaseNode{typ: NodeHeading}, 1}
	if itemHyphen == marker.typ {
		ret.Level = 2
	}

	return
}
