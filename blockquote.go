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
	baseNode := &BaseNode{typ: NodeBlockquote}
	ret = &Blockquote{baseNode}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseBlockquote(line items) (ret Node) {
	_, line = line.trimLeft()
	token := line[0]
	ret = newBlockquote(t, token)
	indentSpaces := t.context.IndentSpaces + 2
	t.context.BlockquoteLevel++
	line = line[1:]
	if itemSpace == line[0].typ {
		line = line[1:]
	} else {
		line = t.indentOffset(line, indentSpaces)
	}
	curNode := t.context.CurNode
	for {
		n := t.parseBlock(line)
		t.context.CurNode = curNode
		if nil == n {
			line = t.nextLine()
			if !t.isBlockquote(line) {
				t.backupLine(line)
				break
			}

			t.backupLine(line)
			line = t.skipBlankBlockquote()
			if line.isEOF() {
				break
			}

			continue
		}
		ret.AppendChild(ret, n)

		line = t.nextLine()
		if line.isEOF() {
			break
		}
		if t.isThematicBreak(line) {
			t.backupLine(line)
			break
		}

		for {
			if !t.isParagraphContinuation(line) {
				break
			}
			lastc := t.context.CurNode.LastChild()
			p := lastc.(*Paragraph)
			line = t.removeStartBlockquoteMarker(line)
			continuation := t.parseParagraph(line)
			p.tokens = append(p.tokens, tNewLine)
			p.tokens = append(p.tokens, continuation.Tokens()...)
			line = t.nextLine()
		}

		if !t.isBlockquote(line) {
			t.backupLine(line)
			break
		}

		_, line = line.trimLeft()
		line = t.removeStartBlockquoteMarker(line)
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

	for _, ret = line[1:].trimLeft(); 0 < len(ret) && (itemGreater == ret[0].typ || ret[0].isSpaceOrTab()); ret = ret[1:] {
	}

	return
}

func (t *Tree) isParagraphContinuation(line items) bool {
	lastc := t.context.CurNode.LastChild()
	if nil == lastc {
		return false
	}

	if NodeParagraph != lastc.Type() {
		return false
	}

	line = t.removeStartBlockquoteMarker(line)
	if t.interruptParagraph(line) {
		return false
	}

	if line.isBlankLine() {
		return false
	}

	return true
}

func (t *Tree) blockquoteMarkerCount(line items) (ret int) {
	_, line = line.trimLeft()
	for _, token := range line {
		if itemGreater == token.typ {
			ret++
		} else if itemSpace != token.typ && itemTab != token.typ {
			break
		}
	}

	return
}

func (t *Tree) skipBlankBlockquote() (line items) {
	for {
		line = t.nextLine()
		if line.isEOF() {
			return
		}

		if !t.isBlockquote(line) {
			return line
		}

		remains := t.removeStartBlockquoteMarker(line)
		if remains.isBlankLine() {
			continue
		}

		return remains
	}
}
