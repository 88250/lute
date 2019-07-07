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

type Paragraph struct {
	*BaseNode

	OpenTag, CloseTag string
}

func (t *Tree) parseParagraph(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeParagraph}
	p := &Paragraph{baseNode, "<p>", "</p>"}
	ret = p

	for {
		line = line.trimLeft()
		len := len(line)
		for i, token := range line {
			if itemBackslash != token.typ {
				p.tokens = append(p.tokens, token)
			} else if i < len-2 && !line[i+1].isASCIIPunct() {
				p.tokens = append(p.tokens, token)
			}
		}

		p.rawText += line.rawText()
		line = t.nextLine()
		if line.isBlankLine() {
			t.backupLine(line)
			break
		}

		if level := t.isSetextHeading(line); 0 < level {
			ret = t.parseSetextHeading(p, level)

			return
		}

		if t.interruptParagrah(line) {
			t.backupLine(line)

			break
		}
	}
	p.tokens = p.tokens.trimRight()
	p.rawText = p.tokens.rawText()

	return
}

func (t *Tree) interruptParagrah(line items) bool {
	if t.isIndentCode(line) {
		return false
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

	if 0 < t.isATXHeading(line) {
		return true
	}

	if t.isList(line) {
		return true
	}

	return false
}
