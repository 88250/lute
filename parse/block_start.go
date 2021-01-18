// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

// blockStartFunc 定义了用于判断块是否开始的函数签名。
type blockStartFunc func(t *Tree, container *ast.Node) int

// blockStarts 定义了一系列函数，每个函数用于判断某种块节点是否可以开始，返回值：
// 0：不匹配
// 1：匹配到容器块，需要继续迭代下降
// 2：匹配到叶子块
var blockStarts = []blockStartFunc{

	// 判断块引用（>）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		marker := lex.Peek(t.Context.currentLine, t.Context.nextNonspace)
		if lex.ItemGreater != marker {
			return 0
		}

		markers := []byte{marker}
		t.Context.advanceNextNonspace()
		t.Context.advanceOffset(1, false)
		// > 后面的空格是可选的
		whitespace := lex.Peek(t.Context.currentLine, t.Context.offset)
		withSpace := lex.ItemSpace == whitespace || lex.ItemTab == whitespace
		if withSpace {
			t.Context.advanceOffset(1, true)
			markers = append(markers, whitespace)
		}
		t.Context.closeUnmatchedBlocks()
		t.Context.addChild(ast.NodeBlockquote)
		t.Context.addChildMarker(ast.NodeBlockquoteMarker, markers)
		return 1
	},

	// 判断 ATX 标题（#）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		if ok, markers, content, level := t.parseATXHeading(); ok {
			t.Context.advanceNextNonspace()
			t.Context.advanceOffset(len(content), false)
			t.Context.closeUnmatchedBlocks()
			heading := t.Context.addChild(ast.NodeHeading)
			heading.HeadingLevel = level
			heading.Tokens = content
			crosshatchMarker := &ast.Node{Type: ast.NodeHeadingC8hMarker, Tokens: markers}
			heading.AppendChild(crosshatchMarker)
			t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
			return 2
		}
		return 0
	},

	// 判断围栏代码块（```）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		if ok, codeBlockFenceChar, codeBlockFenceLen, codeBlockFenceOffset, codeBlockOpenFence, codeBlockInfo := t.parseFencedCode(); ok {
			t.Context.closeUnmatchedBlocks()
			container := t.Context.addChild(ast.NodeCodeBlock)
			container.IsFencedCodeBlock = true
			container.CodeBlockFenceLen = codeBlockFenceLen
			container.CodeBlockFenceChar = codeBlockFenceChar
			container.CodeBlockFenceOffset = codeBlockFenceOffset
			container.CodeBlockOpenFence = codeBlockOpenFence
			container.CodeBlockInfo = codeBlockInfo
			t.Context.advanceNextNonspace()
			t.Context.advanceOffset(codeBlockFenceLen, false)
			return 2
		}
		return 0
	},

	// 判断 Setext 标题（- =）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented || ast.NodeParagraph != container.Type {
			return 0
		}
		level := t.parseSetextHeading()
		if 0 == level {
			return 0
		}

		if t.Context.ParseOption.GFMTable {
			// 尝试解析表，因为可能出现如下情况：
			//
			//   0
			//   -:
			//   -
			//
			// 前两行可以解析出一个只有一个单元格的表。
			// Empty list following GFM Table makes table broken https://github.com/b3log/lute/issues/9
			table := t.Context.parseTable0(container.Tokens)
			if nil != table {
				// 将该段落节点转成表节点
				container.Type = ast.NodeTable
				container.TableAligns = table.TableAligns
				for tr := table.FirstChild; nil != tr; {
					nextTr := tr.Next
					container.AppendChild(tr)
					tr = nextTr
				}
				container.Tokens = nil
				return 0
			}
		}

		t.Context.closeUnmatchedBlocks()
		// 解析链接引用定义
		for tokens := container.Tokens; 0 < len(tokens) && lex.ItemOpenBracket == tokens[0]; tokens = container.Tokens {
			if remains := t.Context.parseLinkRefDef(tokens); nil != remains {
				container.Tokens = remains
			} else {
				break
			}
		}

		if 0 < len(container.Tokens) {
			child := &ast.Node{Type: ast.NodeHeading, HeadingLevel: level, HeadingSetext: true}
			child.Tokens = lex.TrimWhitespace(container.Tokens)
			container.InsertAfter(child)
			container.Unlink()
			t.Context.Tip = child
			t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
			return 2
		}
		return 0
	},

	// 判断 HTML 块（<）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		if lex.ItemLess != lex.Peek(t.Context.currentLine, t.Context.nextNonspace) {
			return 0
		}

		if t.Context.ParseOption.VditorWYSIWYG {
			if bytes.Contains(t.Context.currentLine, []byte("vditor-comment")) {
				return 0
			}
		}

		tokens := t.Context.currentLine[t.Context.nextNonspace:]
		if htmlType := t.parseHTML(tokens); 0 != htmlType {
			t.Context.closeUnmatchedBlocks()
			block := t.Context.addChild(ast.NodeHTMLBlock)
			block.HtmlBlockType = htmlType
			return 2
		}
		return 0
	},

	// 判断 YAML Front Matter（---）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.YamlFrontMatter || t.Context.indented || nil != t.Root.FirstChild {
			return 0
		}

		if t.parseYamlFrontMatter() {
			node := &ast.Node{Type: ast.NodeYamlFrontMatter}
			t.Root.AppendChild(node)
			t.Context.Tip = node
			return 2
		}
		return 0
	},

	// 判断分隔线（--- ***）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		if ok, caretTokens := t.parseThematicBreak(); ok {
			t.Context.closeUnmatchedBlocks()
			thematicBreak := t.Context.addChild(ast.NodeThematicBreak)
			thematicBreak.Tokens = caretTokens
			t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
			return 2
		}
		return 0
	},

	// 判断列表、列表项（* - + 1.）或者任务列表项是否开始。
	func(t *Tree, container *ast.Node) int {
		if ast.NodeList != container.Type && t.Context.indented {
			return 0
		}

		data := t.parseListMarker(container)
		if nil == data {
			return 0
		}

		t.Context.closeUnmatchedBlocks()

		listsMatch := container.Type == ast.NodeList && t.Context.listsMatch(container.ListData, data)
		if t.Context.Tip.Type != ast.NodeList || !listsMatch {
			list := t.Context.addChild(ast.NodeList)
			list.ListData = data
		}
		listItem := t.Context.addChild(ast.NodeListItem)
		listItem.ListData = data
		listItem.Tokens = data.Marker
		if 1 == listItem.ListData.Typ || (3 == listItem.ListData.Typ && 0 == listItem.ListData.BulletChar) {
			// 修正有序列表项序号
			prev := listItem.Previous
			if nil != prev {
				listItem.Num = prev.Num + 1
			} else {
				listItem.Num = data.Start
			}
		}
		return 1
	},

	// 判断数学公式块（$$）是否开始。
	func(t *Tree, container *ast.Node) int {
		if t.Context.indented {
			return 0
		}

		if ok, mathBlockDollarOffset := t.parseMathBlock(); ok {
			t.Context.closeUnmatchedBlocks()
			block := t.Context.addChild(ast.NodeMathBlock)
			block.MathBlockDollarOffset = mathBlockDollarOffset
			t.Context.advanceNextNonspace()
			t.Context.advanceOffset(mathBlockDollarOffset, false)
			return 2
		}
		return 0
	},

	// 判断缩进代码块（    code）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.indented {
			return 0
		}

		if t.Context.Tip.Type != ast.NodeParagraph && !t.Context.blank {
			t.Context.advanceOffset(4, true)
			t.Context.closeUnmatchedBlocks()
			t.Context.addChild(ast.NodeCodeBlock)
			return 2
		}
		return 0
	},

	// 判断脚注定义（[^label]）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.Footnotes || t.Context.indented {
			return 0
		}

		marker := lex.Peek(t.Context.currentLine, t.Context.nextNonspace)
		if lex.ItemOpenBracket != marker {
			return 0
		}
		caret := lex.Peek(t.Context.currentLine, t.Context.nextNonspace+1)
		if lex.ItemCaret != caret {
			return 0
		}

		label := []byte{lex.ItemCaret}
		var token byte
		var i int
		for i = t.Context.nextNonspace + 2; i < t.Context.currentLineLen; i++ {
			token = t.Context.currentLine[i]
			if lex.ItemSpace == token || lex.ItemNewline == token || lex.ItemTab == token {
				return 0
			}
			if lex.ItemCloseBracket == token {
				break
			}
			label = append(label, token)
		}
		if i >= t.Context.currentLineLen {
			return 0
		}
		if lex.ItemColon != t.Context.currentLine[i+1] {
			return 0
		}
		t.Context.advanceOffset(1, false)

		t.Context.closeUnmatchedBlocks()
		t.Context.advanceOffset(len(label)+2, true)

		if ast.NodeFootnotesDefBlock != t.Context.Tip.Type {
			t.Context.addChild(ast.NodeFootnotesDefBlock)
		}

		def := t.Context.addChild(ast.NodeFootnotesDef)
		def.Tokens = label
		return 1
	},

	// 判断 kramdown 块级内联属性列表（{: attrs}）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.KramdownBlockIAL || t.Context.indented {
			return 0
		}

		if ast.NodeListItem == t.Context.Tip.Type && nil == t.Context.Tip.FirstChild { // 列表项 IAL 由后续第一个段落块进行解析
			return 0
		}

		if ial := t.parseKramdownBlockIAL(); nil != ial {
			t.Context.closeUnmatchedBlocks()
			t.Context.offset = t.Context.currentLineLen                    // 整行过
			if 1 < len(ial) && "type" == ial[1][0] && "doc" == ial[1][1] { // 文档块 IAL
				t.Context.rootIAL = &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: t.Context.currentLine[t.Context.nextNonspace:]}
				t.Root.KramdownIAL = ial
				t.Root.ID = ial[0][1]
				t.ID = t.Root.ID
				return 2
			}

			lastMatchedContainer := t.Context.lastMatchedContainer
			if t.Context.allClosed {
				if ast.NodeDocument == lastMatchedContainer.Type || ast.NodeListItem == lastMatchedContainer.Type {
					lastMatchedContainer = t.Context.Tip.LastChild // 挂到最后一个子块上
					if nil == lastMatchedContainer {
						lastMatchedContainer = t.Context.lastMatchedContainer
					}
				}
			}
			lastMatchedContainer.KramdownIAL = ial
			lastMatchedContainer.ID = ial[0][1]
			node := t.Context.addChild(ast.NodeKramdownBlockIAL)
			node.Tokens = t.Context.currentLine[t.Context.nextNonspace:]
			return 2
		}
		return 0
	},

	// 判断内容块嵌入（!((id "text"))）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.BlockRef || t.Context.indented {
			return 0
		}

		node := t.parseBlockEmbed()
		if nil == node {
			return 0
		}

		t.Context.closeUnmatchedBlocks()

		for !t.Context.Tip.CanContain(ast.NodeBlockEmbed) {
			t.Context.finalize(t.Context.Tip) // 注意调用 finalize 会向父节点方向进行迭代
		}
		t.Context.Tip.AppendChild(node)
		t.Context.Tip = node
		return 2
	},

	// 判断内容块查询嵌入（!{{ SELECT * FROM blocks WHERE content LIKE '%待办%' }}）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.BlockRef || t.Context.indented {
			return 0
		}

		node := t.parseBlockQueryEmbed()
		if nil == node {
			return 0
		}

		t.Context.closeUnmatchedBlocks()

		for !t.Context.Tip.CanContain(ast.NodeBlockQueryEmbed) {
			t.Context.finalize(t.Context.Tip) // 注意调用 finalize 会向父节点方向进行迭代
		}
		t.Context.Tip.AppendChild(node)
		t.Context.Tip = node
		return 2
	},

	// 判断超级块（{{{ blocks }}}）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.ParseOption.SuperBlock || t.Context.indented {
			return 0
		}

		if ok, layout := t.parseSuperBlock(); ok {
			t.Context.closeUnmatchedBlocks()
			t.Context.addChild(ast.NodeSuperBlock)
			t.Context.addChildMarker(ast.NodeSuperBlockOpenMarker, nil)
			t.Context.addChildMarker(ast.NodeSuperBlockLayoutMarker, layout)
			t.Context.offset = t.Context.currentLineLen - 1 // 整行过
			return 1
		}
		return 0
	},
}
