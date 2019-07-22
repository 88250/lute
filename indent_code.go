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

func (t *Tree) parseIndentCode(tokens items) (ret Node) {
	spaces, tabs, remains := t.nonSpaceTab(tokens)
	if 4 > spaces && 1 > tabs {
		return
	}

	baseNode := &BaseNode{typ: NodeCode}
	code := &Code{baseNode, "", ""}
	code.Value += remains.rawText()

	ret = code

	return
}

func (t *Tree) nonNewline() (newlines items, line items) {
	for line = t.nextLine(); line.isBlankLine() && !line.isEOF(); line = t.nextLine() {
		if 5 > len(line) {
			_, line = line.trimLeftSpace()
		} else {
			line = t.indentOffset(line, 4)
		}

		newlines = append(newlines, line...)
	}

	return
}
