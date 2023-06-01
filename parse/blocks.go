// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// parseBlocks 解析并生成块级节点。
func (t *Tree) parseBlocks() {
	t.Context.Tip = t.Root
	lines := 0
	for line := t.lexer.NextLine(); nil != line; line = t.lexer.NextLine() {
		if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
			if !bytes.Equal(line, editor.CaretNewlineTokens) && t.Context.Tip.ParentIs(ast.NodeListItem) && bytes.HasPrefix(line, editor.CaretTokens) {
				// 插入符在开头的话移动到上一行结尾，处理 https://github.com/Vanessa219/vditor/issues/633 中的一些情况
				if ast.NodeListItem == t.Context.Tip.Type {
					t.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: line})
					break
				} else {
					t.Context.Tip.Tokens = bytes.TrimSuffix(t.Context.Tip.Tokens, []byte("\n"))
					t.Context.Tip.Tokens = append(t.Context.Tip.Tokens, editor.CaretNewlineTokens...)
				}
				line = line[len(editor.CaretTokens):]
			}
		}

		t.incorporateLine(line)
		lines++
	}
	for nil != t.Context.Tip {
		t.Context.finalize(t.Context.Tip)
	}
}

func (t *Tree) BlockCount() (ret int) {
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if "" == n.ID || !n.IsBlock() {
			return ast.WalkContinue
		}

		ret++
		return ast.WalkContinue
	})
	return
}

func (t *Tree) DocBlockCount() (ret int) {
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if !n.IsChildBlockOf(t.Root, 1) {
			return ast.WalkContinue
		}

		ret++
		return ast.WalkContinue
	})
	return
}

// incorporateLine 处理文本行 line 并把生成的块级节点挂到树上。
func (t *Tree) incorporateLine(line []byte) {
	t.Context.oldtip = t.Context.Tip
	t.Context.offset = 0
	t.Context.column = 0
	t.Context.blank = false
	t.Context.partiallyConsumedTab = false
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
		case 3: // 匹配超级块闭合，处理下一行
			t.Context.closeSuperBlockChildren() // 闭合超级块下的子节点
			if ast.NodeSuperBlock != t.Context.Tip.Type {
				sb := t.Context.Tip.Parent
				sb.Close = true
				sb.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
				t.Context.Tip = sb.Parent
				t.Context.lastMatchedContainer = sb
			} else {
				t.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
				t.Context.Tip.Close = true
				t.Context.Tip = t.Context.Tip.Parent
				t.Context.lastMatchedContainer = t.Context.Tip
			}
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
	blockParsers := blockStarts()
	startsLen := len(blockParsers)

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
			lex.ItemSemicolon != maybeMarker && // 定义块
			lex.ItemCrosshatch != maybeMarker && // ATX 标题
			lex.ItemGreater != maybeMarker && // 块引用
			lex.ItemLess != maybeMarker && // HTML 块
			lex.ItemUnderscore != maybeMarker && lex.ItemEqual != maybeMarker && // Setext 标题
			lex.ItemDollar != maybeMarker && // 数学公式
			lex.ItemOpenBracket != maybeMarker && // 脚注
			lex.ItemOpenBrace != maybeMarker && // kramdown 内联属性列表或超级块开始
			lex.ItemCloseBrace != maybeMarker && // 超级块闭合
			lex.ItemBang != maybeMarker && "！"[0] != maybeMarker && // 内容块嵌入
			editor.Caret[0] != maybeMarker { // Vditor 编辑器支持
			t.Context.advanceNextNonspace()
			break
		}

		// 逐个尝试是否可以起始一个块级节点
		i := 0
		for i < startsLen {
			res := blockParsers[i](t, container)
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
				(typ == ast.NodeCustomBlock) || // 自定义块不计入空行判断
				(typ == ast.NodeMathBlock) || // 数学公式块不计入空行判断
				(typ == ast.NodeGitConflict) || // Git 冲突标记不计入空行判断
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
						t.Context.finalize(container)
					}
				}
			case ast.NodeMathBlock:
				// 数学公式块标记符没有换行的形式（$$foo$$）需要判断右边结尾的闭合标记符
				if 3 < len(container.Tokens) &&
					(bytes.HasSuffix(container.Tokens, MathBlockMarkerNewline) ||
						bytes.HasSuffix(container.Tokens, MathBlockMarker) ||
						bytes.HasSuffix(container.Tokens, MathBlockMarkerCaretNewline)) {
					t.Context.finalize(container)
				}
			}
		} else if t.Context.offset < t.Context.currentLineLen && !t.Context.blank {
			// 普通段落开始
			t.Context.addChild(ast.NodeParagraph)
			t.Context.advanceNextNonspace()
			t.addLine()
		}
	}
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

	startWithSpace := 1 < t.Context.currentLineLen && (' ' == t.Context.currentLine[0] || '\t' == t.Context.currentLine[0])
	docChildPara := ast.NodeDocument == t.Context.Tip.Parent.Type
	if t.Context.ParseOption.ParagraphBeginningSpace && startWithSpace && docChildPara {
		t.Context.Tip.AppendTokens(t.Context.currentLine)
	} else {
		t.Context.Tip.AppendTokens(t.Context.currentLine[t.Context.offset:])
	}
}

// _continue 判断节点是否可以继续处理，比如块引用需要 >，缩进代码块需要 4 空格，围栏代码块需要 ```。
// 如果可以继续处理返回 0，如果不能接续处理返回 1，如果返回 2（仅在围栏代码块、超级块或自定义块闭合时）则说明可以继续下一行处理了。
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
	case ast.NodeSuperBlock:
		return SuperBlockContinue(n, context)
	case ast.NodeGitConflict:
		return GitConflictContinue(n, context)
	case ast.NodeCustomBlock:
		return CustomBlockContinue(n, context)
	case ast.NodeHeading, ast.NodeThematicBreak, ast.NodeKramdownBlockIAL, ast.NodeLinkRefDefBlock, ast.NodeBlockQueryEmbed,
		ast.NodeIFrame, ast.NodeVideo, ast.NodeAudio, ast.NodeWidget, ast.NodeAttributeView:
		return 1
	}
	return 0
}
