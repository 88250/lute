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

func (t *Tree) parseBlocks() {
	for line := t.nextLine(); !line.isEOF(); line = t.nextLine() {
		t.processLine(line)
	}
}

func (t *Tree) processLine(line items) {
	spaces, tabs, remains := t.nonSpaceTab(line)
	indentSpaces := spaces + tabs*4
	node := t.parseBlock(remains)

	_ = indentSpaces

	var lastOpenNode Node
	for lastOpenNode = t.Root.FirstChild(); nil != lastOpenNode && lastOpenNode.IsOpen(); lastOpenNode = lastOpenNode.Next() {
		switch lastOpenNode.Type() {
		case NodeListItem:
		case NodeBlockquote:
		case NodeRoot:
			lastOpenNode.AppendChild(lastOpenNode, node)
		default:
		}
	}

}

func (t *Tree) parseBlock(tokens items) (node Node) {
	if tokens.isEOF() {
		return
	}

	node = t.parseParagraph(tokens)

	return
}
