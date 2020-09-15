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
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// parseBlocks 解析并生成块级节点。
func (t *Tree) parseBlocks() {
	t.Context.Tip = t.Root
	t.Context.LinkRefDefs = map[string]*ast.Node{}
	t.Context.FootnotesDefs = []*ast.Node{}
	lines := 0
	for line := t.lexer.NextLine(); nil != line; line = t.lexer.NextLine() {
		if t.Context.Option.VditorWYSIWYG || t.Context.Option.VditorIR || t.Context.Option.VditorSV {
			if !bytes.Equal(line, util.CaretNewlineTokens) && t.Context.Tip.ParentIs(ast.NodeListItem) && bytes.HasPrefix(line, util.CaretTokens) {
				// 光标在开头的话移动到上一行结尾，处理 https://github.com/Vanessa219/vditor/issues/633 中的一些情况
				if ast.NodeListItem == t.Context.Tip.Type {
					t.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: line})
					break
				} else {
					t.Context.Tip.Tokens = bytes.TrimSuffix(t.Context.Tip.Tokens, []byte("\n"))
					t.Context.Tip.Tokens = append(t.Context.Tip.Tokens, util.CaretNewlineTokens...)
				}
				line = line[len(util.CaretTokens):]
			}
			ln := []rune(string(line))
			if 4 < len(ln) && lex.IsDigit(byte(ln[0])) && ('、' == ln[1] || '）' == ln[1]) {
				// 列表标记符自动优化 https://github.com/Vanessa219/vditor/issues/68
				line = []byte(string(ln[0]) + ". " + string(ln[2:]))
			}
		}

		t.incorporateLine(line)
		lines++
	}
	for nil != t.Context.Tip {
		t.Context.finalize(t.Context.Tip, lines)
	}
}

// incorporateLine 处理文本行 line 并把生成的块级节点挂到树上。
func (t *Tree) incorporateLine(line []byte) {
	t.Context.oldtip = t.Context.Tip
	t.Context.offset = 0
	t.Context.column = 0
	t.Context.blank = false
	t.Context.partiallyConsumedTab = false
	t.Context.lineNum++
	t.Context.currentLine = line
	t.Context.currentLineLen = len(t.Context.currentLine)

	allMatched := true
	var container *ast.Node
	container = t.Root
	lastChild := container.LastChild
	for ; nil != lastChild && !lastChild.Close; lastChild = container.LastChild {
		container = lastChild
		t.Context.findNextNonspace()

		switch _continue(container, t.Context) {
		case 0: // 说明匹配可继续处理
			break
		case 1: // 匹配失败，不能继续处理
			allMatched = false
			break
		case 2: // 匹配围栏代码块闭合，处理下一行
			return
		}

		if !allMatched {
			container = container.Parent // 回到上一个匹配的块
			break
		}
	}

	t.Context.allClosed = container == t.Context.oldtip
	t.Context.lastMatchedContainer = container

	matchedLeaf := container.Type != ast.NodeParagraph && container.AcceptLines()
	startsLen := len(blockStarts)

	// 除非最后一个匹配到的是代码块，否则的话就起始一个新的块级节点
	for !matchedLeaf {
		t.Context.findNextNonspace()

		// 如果不由潜在的节点标记符开头 ^[#`~*+_=<>0-9-${]，则说明不用继续迭代生成子节点
		// 这里仅做简单判断的话可以提升一些性能
		maybeMarker := t.Context.currentLine[t.Context.nextNonspace]
		if !t.Context.indented && // 缩进代码块
			lex.ItemHyphen != maybeMarker && lex.ItemAsterisk != maybeMarker && lex.ItemPlus != maybeMarker && // 无序列表
			!lex.IsDigit(maybeMarker) && // 有序列表
			lex.ItemBacktick != maybeMarker && lex.ItemTilde != maybeMarker && // 代码块
			lex.ItemCrosshatch != maybeMarker && // ATX 标题
			lex.ItemGreater != maybeMarker && // 块引用
			lex.ItemLess != maybeMarker && // HTML 块
			lex.ItemUnderscore != maybeMarker && lex.ItemEqual != maybeMarker && // Setext 标题
			lex.ItemDollar != maybeMarker && // 数学公式
			lex.ItemOpenBracket != maybeMarker && // 脚注
			lex.ItemOpenCurlyBrace != maybeMarker && // kramdown 内联属性列表
			util.Caret[0] != maybeMarker { // Vditor 编辑器支持
			t.Context.advanceNextNonspace()
			break
		}

		// 逐个尝试是否可以起始一个块级节点
		i := 0
		for i < startsLen {
			res := blockStarts[i](t, container)
			if res == 1 { // 匹配到容器块，继续迭代下降过程
				container = t.Context.Tip
				break
			} else if res == 2 { // 匹配到叶子块，跳出迭代下降过程
				container = t.Context.Tip
				matchedLeaf = true
				break
			} else { // 没有匹配到，继续用下一个起始块模式进行匹配
				i++
			}
		}

		if i == startsLen { // 没有匹配到任何块
			t.Context.advanceNextNonspace()
			break
		}
	}

	// offset 后余下的内容算作是文本行，需要将其添加到相应的块节点上

	if !t.Context.allClosed && !t.Context.blank && t.Context.Tip.Type == ast.NodeParagraph {
		// 该行是段落延续文本，直接添加到当前末梢段落上
		t.addLine()
	} else {
		// 最终化未匹配的块
		t.Context.closeUnmatchedBlocks()

		if t.Context.blank && nil != container.LastChild {
			container.LastChild.LastLineBlank = true
		}

		typ := container.Type
		isFenced := ast.NodeCodeBlock == typ && container.IsFencedCodeBlock

		// 空行判断，主要是为了判断列表是紧凑模式还是松散模式
		lastLineBlank := t.Context.blank &&
			!(typ == ast.NodeFootnotesDef ||
				typ == ast.NodeBlockquote || // 块引用行肯定不会是空行因为至少有一个 >
				(typ == ast.NodeCodeBlock && isFenced) || // 围栏代码块不计入空行判断
				(typ == ast.NodeMathBlock) || // 数学公式块不计入空行判断
				(typ == ast.NodeListItem && nil == container.FirstChild)) // 内容为空的列表项也不计入空行判断
		// 因为列表是块级容器（可进行嵌套），所以需要在父节点方向上传播 LastLineBlank
		// LastLineBlank 目前仅在判断列表紧凑模式上使用
		for cont := container; nil != cont; cont = cont.Parent {
			cont.LastLineBlank = lastLineBlank
		}

		if container.AcceptLines() {
			t.addLine()
			switch typ {
			case ast.NodeHTMLBlock:
				// HTML 块（类型 1-5）需要检查是否满足闭合条件
				html := container
				if html.HtmlBlockType >= 1 && html.HtmlBlockType <= 5 {
					tokens := t.Context.currentLine[t.Context.offset:]
					if t.isHTMLBlockClose(tokens, html.HtmlBlockType) {
						t.Context.finalize(container, t.Context.lineNum)
					}
				}
			case ast.NodeMathBlock:
				// 数学公式块标记符没有换行的形式（$$foo$$）需要判断右边结尾的闭合标记符
				if 3 < len(container.Tokens) &&
					(bytes.HasSuffix(container.Tokens, MathBlockMarkerNewline) ||
						bytes.HasSuffix(container.Tokens, MathBlockMarker) ||
						bytes.HasSuffix(container.Tokens, MathBlockMarkerCaretNewline)) {
					t.Context.finalize(container, t.Context.lineNum)
				}
			}
		} else if t.Context.offset < t.Context.currentLineLen && !t.Context.blank {
			// 普通段落开始
			t.Context.addChild(ast.NodeParagraph, t.Context.offset)
			t.Context.advanceNextNonspace()
			t.addLine()
		}
	}
}

// blockStartFunc 定义了用于判断块是否开始的函数签名。
type blockStartFunc func(t *Tree, container *ast.Node) int

// blockStarts 定义了一系列函数，每个函数用于判断某种块节点是否可以开始，返回值：
// 0：不匹配
// 1：匹配到块容器，需要继续迭代下降
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
		if t.Context.Option.VditorWYSIWYG || t.Context.Option.VditorIR || t.Context.Option.VditorSV {
			// Vditor 三个模式都不能存在空的块引用
			ln := util.BytesToStr(t.Context.currentLine[t.Context.offset:])
			ln = strings.ReplaceAll(ln, util.Caret, "")
			if ln = strings.TrimSpace(ln); "" == ln {
				return 0
			}
		}
		t.Context.closeUnmatchedBlocks()
		t.Context.addChild(ast.NodeBlockquote, t.Context.nextNonspace)
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
			heading := t.Context.addChild(ast.NodeHeading, t.Context.nextNonspace)
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
			container := t.Context.addChild(ast.NodeCodeBlock, t.Context.nextNonspace)
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

		if t.Context.Option.GFMTable {
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

		tokens := t.Context.currentLine[t.Context.nextNonspace:]
		if htmlType := t.parseHTML(tokens); 0 != htmlType {
			t.Context.closeUnmatchedBlocks()
			block := t.Context.addChild(ast.NodeHTMLBlock, t.Context.offset)
			block.HtmlBlockType = htmlType
			return 2
		}
		return 0
	},

	// 判断 YAML Front Matter（---）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.Option.YamlFrontMatter || t.Context.indented || nil != t.Root.FirstChild {
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
			thematicBreak := t.Context.addChild(ast.NodeThematicBreak, t.Context.nextNonspace)
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
			list := t.Context.addChild(ast.NodeList, t.Context.nextNonspace)
			list.ListData = data
		}
		listItem := t.Context.addChild(ast.NodeListItem, t.Context.nextNonspace)
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
			block := t.Context.addChild(ast.NodeMathBlock, t.Context.nextNonspace)
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
			t.Context.addChild(ast.NodeCodeBlock, t.Context.offset)
			return 2
		}
		return 0
	},

	// 判断脚注定义（[^label]）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.Option.Footnotes || t.Context.indented {
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
		footnotesDef := t.Context.addChild(ast.NodeFootnotesDef, t.Context.nextNonspace)
		footnotesDef.Tokens = label
		lowerCaseLabel := bytes.ToLower(label)
		if _, def := t.Context.FindFootnotesDef(lowerCaseLabel); nil == def {
			t.Context.FootnotesDefs = append(t.Context.FootnotesDefs, footnotesDef)
		}
		return 1
	},

	// 判断 kramdown 内联属性列表（{: attrs}）是否开始。
	func(t *Tree, container *ast.Node) int {
		if !t.Context.Option.KramdownIAL || t.Context.indented {
			return 0
		}

		if ial := t.parseKramdownIAL(); nil != ial {
			t.Context.closeUnmatchedBlocks()
			lastMatchedContainer := t.Context.lastMatchedContainer
			if t.Context.allClosed && (ast.NodeDocument == lastMatchedContainer.Type || ast.NodeListItem == lastMatchedContainer.Type) {
				lastMatchedContainer = t.Context.Tip.LastChild // 挂到最后一个子块上
				if nil == lastMatchedContainer {
					lastMatchedContainer = t.Context.lastMatchedContainer
				}
			}
			lastMatchedContainer.KramdownIAL = ial
			t.Context.offset = t.Context.currentLineLen // 整行过
			node := t.Context.addChild(ast.NodeKramdownBlockIAL, t.Context.nextNonspace)
			node.Tokens = t.Context.currentLine
			return 2
		}
		return 0
	},
}

// addLine 用于在当前的末梢节点 context.Tip 上添加迭代行剩余的所有 Tokens。
// 调用该方法前必须确认末梢 tip 能够接受新行。
func (t *Tree) addLine() {
	if t.Context.partiallyConsumedTab {
		t.Context.offset++ // skip over tab
		// add space characters:
		charsToTab := 4 - (t.Context.column % 4)
		t.Context.Tip.AppendTokens(bytes.Repeat(util.StrToBytes(" "), charsToTab))
	}
	t.Context.Tip.AppendTokens(t.Context.currentLine[t.Context.offset:])
}

// _continue 判断节点是否可以继续处理，比如块引用需要 >，缩进代码块需要 4 空格，围栏代码块需要 ```。
// 如果可以继续处理返回 0，如果不能接续处理返回 1，如果返回 2（仅在围栏代码块闭合时）则说明可以继续下一行处理了。
func _continue(n *ast.Node, context *Context) int {
	switch n.Type {
	case ast.NodeCodeBlock:
		return CodeBlockContinue(n, context)
	case ast.NodeHTMLBlock:
		return HtmlBlockContinue(n, context)
	case ast.NodeParagraph:
		return ParagraphContinue(n, context)
	case ast.NodeListItem:
		return ListItemContinue(n, context)
	case ast.NodeBlockquote:
		return BlockquoteContinue(n, context)
	case ast.NodeMathBlock:
		return MathBlockContinue(n, context)
	case ast.NodeYamlFrontMatter:
		return YamlFrontMatterContinue(n, context)
	case ast.NodeFootnotesDef:
		return FootnotesContinue(n, context)
	case ast.NodeHeading, ast.NodeThematicBreak, ast.NodeKramdownBlockIAL:
		return 1
	}
	return 0
}
