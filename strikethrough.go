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

// Strikethrough 描述了删除线节点结构。
type Strikethrough struct {
	*BaseNode
}

func (t *Tree) parseStrikethrough(tokens items) (ret Node) {
	length := len(tokens)
	if 3 > length {
		return &Text{tokens: tokens[:1]}
	}

	token := tokens[1]
	if itemTilde != token {
		return &Text{tokens: tokens[:1]}
	}
	token = tokens[2]
	if itemTilde == token || isWhitespace(token) { // 最多只能两个 ~ 开头，第三个 token 不能是空白
		return &Text{tokens: tokens[:2]}
	}

	remains := tokens[2:]
	length = len(remains)
	matched := false
	i := 0
	for ; i < length; i++ {
		token = remains[i]
		if itemTilde == token {
			if i < length-1 && itemTilde == remains[i+1] {
				if !isWhitespace(remains[i-1]) {
					matched = true
					break
				}
			}
		}
	}
	if matched {
		t.context.pos += 2 + i + 2
		return &Strikethrough{&BaseNode{typ: NodeStrikethrough, tokens: tokens[:t.context.pos]}}
	}

	t.context.pos += length
	return &Text{tokens: tokens}
}
