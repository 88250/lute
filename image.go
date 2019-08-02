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

type Image struct {
	*BaseNode
	Destination string
	Title       string
}

func (t *Tree) parseBang(tokens items) (ret Node) {
	var startPos = t.context.pos
	t.context.pos++
	if itemOpenBracket == tokens[t.context.pos].typ {
		t.context.pos++
		ret = &Text{&BaseNode{typ: NodeText, value: "!["}}
		// Add entry to stack for this opener
		t.addBracket(ret, startPos+2, true)
	} else {
		ret = &Text{&BaseNode{typ: NodeText, value: "!"}}
	}
	return
}
