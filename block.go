// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

func (t *Tree) parseBlocks() {
	curNode := t.context.CurNode
	for token := t.peek(); itemEOF != token.typ; token = t.peek() {
		t.parseBlock()
		t.context.CurNode = curNode
	}
}

func (t *Tree) parseBlock() (ret Node) {
	curNode := t.context.CurNode
	line := t.nextLineEnding()
	t.backups(line)

	if t.isThematicBreak(line) {
		ret = t.parseThematicBreak()
	} else if t.isList(line) {
		ret = t.parseList()
	} else if t.isATXHeading(line) {
		ret = t.parseHeading()
	} else if t.isBlockquote(line) {
		ret = t.parseBlockquote()
	} else if t.isIndentCode(line) {
		ret = t.parseIndentCode()
	} else if t.isBlankLine(line) {
		return
	} else {
		ret = t.parseParagraph()
	}

	curNode.Append(ret)

	return
}
