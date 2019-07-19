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
	curNode := t.context.CurNodes.peek()
	for line := t.nextLine(); ; {
		n := t.parseBlock(line)
		if nil != n {
			curNode.AppendChild(curNode, n)
		}
		curNode = t.context.CurNodes.peek()

		line = t.nextLine()
		if line.isEOF() {
			break
		}
	}

	for child := t.Root.FirstChild(); nil != child; child = child.Next() {
		child.Close()
	}
}

func (t *Tree) parseBlock(line items) (ret Node) {
	if 1 > len(line) || line.isEOF() {
		return
	}

	atxHeadingLevel := 0
	htmlType := -1

	if line.isBlankLine() {
	} else if t.isIndentCode(line) {
		ret = t.parseIndentCode(line)
	} else if t.isFencedCode(line) {
		ret = t.parseFencedCode(line)
	} else if t.isThematicBreak(line) {
		ret = t.parseThematicBreak(line)
	} else if t.isATXHeading(line, &atxHeadingLevel) {
		ret = t.parseATXHeading(line, atxHeadingLevel)
	} else if t.isBlockquote(line) {
		ret = t.parseBlockquote(line)
	} else if isList, _ := t.isList(line); isList {
		if NodeList == t.context.CurNodes.peek().Type() {
			ret = t.parseListItem(line)
		} else {
			ret = t.parseList(line)
		}
	} else if t.isHTML(line, &htmlType) {
		ret = t.parseHTML(line, htmlType)
	} else if t.parseLinkRefDef(line) {
	} else {
		ret = t.parseParagraph(line)
	}

	return
}
