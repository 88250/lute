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

func newBlockquote(t *Tree, token *item) (ret Node) {
	baseNode := &BaseNode{typ: NodeBlockquote, parent: t.context.CurNode}
	ret = &Blockquote{baseNode}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseBlockquote(line items) (ret Node) {
	token := line[0]
	indentSpaces := t.context.IndentSpaces + 2

	ret = newBlockquote(t, token)
	line = t.indentOffset(line[1:], indentSpaces)
	for {
		n := t.parseBlock(line)
		if nil == n {
			break
		}
		ret.AppendChild(ret, n)

		line = t.nextLine()
		if t.isThematicBreak(line) || t.isBlockquoteClose(line) {
			t.backupLine(line)
			break
		}
	}

	return
}

func (t *Tree) isBlockquote(line items) bool {
	if 2 > len(line) { // at least > and newline
		return false
	}

	_, marker := line.firstNonSpace()
	if ">" != marker.val {
		return false
	}

	return true
}

func (t *Tree) removeStartBlockquoteMarker(line items) (ret items) {
	if NodeBlockquote != t.context.CurNode.Type() {
		return line
	}

	_, ret = line[1:].trimLeft()

	return
}

func (t *Tree) isBlockquoteClose(line items) bool {
	if line.isEOF() || NodeBlockquote != t.context.CurNode.Type() {
		return false
	}

	if itemNewline == line[0].typ {
		return true
	}

	return true
}
