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
	t.context.Line = line
	t.walkOpenBlock(t.Root)
}

func (t *Tree) walkOpenBlock(openBlock Node) WalkStatus {
	if nil != openBlock.FirstChild() && openBlock.FirstChild().IsOpen() {
		for openBlock = openBlock.FirstChild(); nil != openBlock && openBlock.IsOpen(); openBlock = openBlock.Next() {
			if WalkStop == t.walkOpenBlock(openBlock) {
				return WalkStop
			}
		}
	}

	if nil == openBlock || openBlock.IsClosed() {
		return WalkContinue
	}

	t.appendBlock(openBlock)

	return WalkStop
}

func (t *Tree) appendBlock(openBlock Node) {
	// indent offset

	switch openBlock.Type() {
	case NodeListItem:
	case NodeBlockquote:
	case NodeParagraph:
		lineNode := t.parseBlock(t.context.Line)
		switch lineNode.Type() {
		case NodeParagraph:
			openBlock.AddTokens(items{tNewLine})
			openBlock.AddTokens(lineNode.Tokens())
		case NodeBlankLine:
			openBlock.Close()
		case NodeList:
			lineNode = lineNode.FirstChild()
			fallthrough
		default:
			openBlock.Close()
			prev := openBlock.Previous()
			if nil == prev {
				openBlock = openBlock.Parent()
				openBlock.AppendChild(openBlock, lineNode)
			} else {
				parent := prev.Parent()
				parent.AppendChild(parent, lineNode)
			}
		}
	case NodeRoot:
		lineNode := t.parseBlock(t.context.Line)
		openBlock.AppendChild(openBlock, lineNode)
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
