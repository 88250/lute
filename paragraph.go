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

import (
	"fmt"
	"strings"
)

type Paragraph struct {
	*BaseNode
	int
	items []*item
	*Tree

	OpenTag, CloseTag string
}

func (n *Paragraph) HTML() string {
	content := html(n.Children())

	if "" != n.OpenTag {
		return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag+"\n", content)
	}

	return fmt.Sprintf(n.OpenTag+"%s"+n.CloseTag, content)
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
	n.rawText = strings.TrimSpace(n.rawText)
}

func (t *Tree) parseParagraph(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeParagraph}
	p := &Paragraph{baseNode, line[0].pos, nil, t, "<p>", "</p>"}
	defer p.trim()

	for {
		p.items = append(p.items, items(line)...)
		p.rawText += line.rawText()
		line = t.nextLine()
		if t.interruptParagrah(line) {
			t.backupLine(line)

			break
		}
	}

	return ret
}

func (t *Tree) interruptParagrah(line items) bool {
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
