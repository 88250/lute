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
	int
	*Tree

	Depth int
}

func (t *Tree) parseHeading(line items) (ret Node) {
	marker := line[0]

	baseNode := &BaseNode{typ: NodeHeading, tokens: items{}}
	heading := &Heading{baseNode, marker.pos, t, len(marker.val)}
	ret = heading

	tokens := t.skipWhitespaces(line[1:])
	for _, token := range tokens {
		if itemEOF == token.typ {
			break
		}
		if itemNewline == token.typ {
			break
		}

		heading.rawText += token.val
		heading.tokens = append(heading.tokens, token)
	}

	return
}

// https://spec.commonmark.org/0.29/#atx-headings
func (t *Tree) isATXHeading(line items) bool {
	if 2 > len(line) { // at least # and newline
		return false
	}

	_, marker := t.firstNonSpace(line)
	// TODO: # 后面还需要空格才能确认是否是列表
	if "#" != marker.val {
		return false
	}

	return true
}
