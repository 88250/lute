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

type ThematicBreak struct {
	*BaseNode
}

func (t *Tree) parseThematicBreak(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeThematicBreak, tokens:line}
	ret = &ThematicBreak{baseNode}

	return
}

// A line consisting of 0-3 spaces of indentation, followed by a sequence of three or more matching -, _, or * characters, each followed optionally by any number of spaces or tabs, forms a thematic break.
// https://spec.commonmark.org/0.29/#thematic-break
func (t *Tree) isThematicBreak(line items) bool {
	if 3 > len(line) {
		return false
	}

	tokens := t.removeSpacesTabs(line)
	tokens = tokens[:len(tokens)-1] // remove tailing newline
	length := len(tokens)
	if 3 > length {
		return false
	}

	marker := tokens[0]
	if "-" != marker.val && "_" != marker.val && "*" != marker.val {
		return false
	}

	for i := 1; i < length; i++ {
		token := tokens[i]
		if marker.typ != token.typ {
			return false
		}
	}

	return true
}
