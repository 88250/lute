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

type Blockquote struct {
	*BaseNode
}

func (t *Tree) parseBlockquote(tokens items) (ret Node) {
	if 2 > len(tokens) {
		return
	}

	_, marker := tokens.firstNonSpace()
	if itemGreater != marker.typ {
		return
	}

	ret = &Blockquote{&BaseNode{typ: NodeBlockquote}}
	tokens = tokens[1:]
	if tokens[0].isSpace() {
		tokens = tokens[1:]
	} else if tokens[0].isTab() {
		tokens = t.indentOffset(tokens, 2)
	}

	child := t.parseBlock(tokens)
	ret.AppendChild(ret, child)

	return
}
