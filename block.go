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
	node := t.parseBlock(line)

	var lastOpenContainer Node
	var lastOpenBlock Node
	for lastOpenContainer = t.Root; nil != lastOpenContainer && lastOpenContainer.IsOpen(); lastOpenContainer = lastOpenContainer.Next() {
		lastOpenBlock = lastOpenContainer.FirstChild()
		if nil == lastOpenBlock {
			t.appendBlock(lastOpenContainer, node)
			return
		}

		for ; nil != lastOpenBlock && lastOpenBlock.IsOpen(); lastOpenBlock = lastOpenBlock.Next() {
			t.appendBlock(lastOpenBlock, node)
			return
		}
	}

}

func (t *Tree) appendBlock(lastOpenBlock, node Node) {
	switch lastOpenBlock.Type() {
	case NodeListItem:
	case NodeBlockquote:
	case NodeParagraph:
		switch node.Type() {
		case NodeParagraph:
			lastOpenBlock.AddTokens(items{tNewLine})
			lastOpenBlock.AddTokens(node.Tokens())
		default:
			lastOpenBlock.Close()
			prev := lastOpenBlock.Previous()
			if nil == prev {
				lastOpenBlock = lastOpenBlock.Parent()
			}

			lastOpenBlock.AppendChild(lastOpenBlock, node)
		}

		return
	case NodeRoot:
		lastOpenBlock.AppendChild(lastOpenBlock, node)
		return
	}
}

func (t *Tree) parseBlock(tokens items) (node Node) {
	node = t.parseBlankLine(tokens)
	if nil == node {
		node = t.parseIndentCode(tokens)
	}
	if nil == node {
		node = t.parseATXHeading(tokens)
	}
	if nil == node {
		node = t.parseList(tokens)
	}
	if nil == node {
		node = t.parseBlockquote(tokens)
	}
	if nil == node {
		node = t.parseParagraph(tokens)
	}

	return
}
