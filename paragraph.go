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

import (
	"fmt"
	"strings"
)

type Paragraph struct {
	NodeType
	int
	RawText
	items items
	*Tree
	Subnodes Children

	OpenTag, CloseTag string
}

func (n *Paragraph) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Paragraph) HTML() string {
	content := html(n.Subnodes)

	if "" != n.OpenTag {
		return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag+"\n", content)
	}

	return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag, content)
}

func (n *Paragraph) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Paragraph) Children() Children {
	return n.Subnodes
}

func (n *Paragraph) Tokens() items {
	return n.items
}

func (n *Paragraph) trim() {
	size := len(n.items)
	if 1 > size {
		return
	}

	initialNoneWhitespace := 0
	notBreak := true
	for i := initialNoneWhitespace; i < size/2 && notBreak; i++ {
		if n.items[i].isWhitespace() {
			initialNoneWhitespace++
			notBreak = true
		} else {
			notBreak = false
		}
	}

	finalNoneWhitespace := size
	notBreak = true
	for i := finalNoneWhitespace - 1; size/2 <= i && notBreak; i-- {
		if n.items[i].isWhitespace() {
			finalNoneWhitespace--
			notBreak = true
		} else {
			notBreak = false
		}
	}

	n.items = n.items[initialNoneWhitespace:finalNoneWhitespace]
	n.RawText = RawText(strings.TrimSpace(string(n.RawText)))
}

func (t *Tree) parseParagraph(line line) Node {
	ret := &Paragraph{NodeParagraph, line[0].pos, "", nil, t, Children{}, "<p>", "</p>"}
	defer ret.trim()

	for {
		ret.items = append(ret.items, items(line)...)
		line = t.nextLine()
		if t.interruptParagrah(line) {
			t.backupLine(line)

			break
		}
	}

	return ret
}

func (t *Tree) interruptParagrah(line []item) bool {
	if t.isBlankLine(line) {
		return true
	}

	/*
	 * 专题分隔线 `***` 打断段落
	 * ATX 标题 `# h` 打断段落，Setext 标题不打断，需要用空行分隔之前的内容
	 * 围栏代码块 <code>```</code> 打断段落
	 * 大部分 HTML 标签可打断段落，除了带属性的，比如 `<a `、`<img `
	 * 块引用 `>` 打断段落
	 * 第一个非空列表项打断段落（即新列表打断段落）
	 */
	if t.isThematicBreak(line) {
		return true
	}

	if t.isList(line) {
		return true
	}

	return false
}
