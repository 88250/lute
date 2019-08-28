// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import "bytes"

// parseBlocks 解析并生成块级节点。
func (t *Tree) parseBlocks() {
	t.context.tip = t.Root
	t.context.linkRefDef = map[string]*Link{}
	for line := t.lex.nextLine(); nil != line; line = t.lex.nextLine() {
		t.incorporateLine(line)
	}
	for nil != t.context.tip {
		t.context.finalize(t.context.tip)
	}
}

// incorporateLine 处理文本行 line 并把生成的块级节点挂到树上。
func (t *Tree) incorporateLine(line items) {
	t.context.oldtip = t.context.tip
	t.context.offset = 0
	t.context.column = 0
	t.context.blank = false
	t.context.partiallyConsumedTab = false
	t.context.currentLine = line
	t.context.currentLineLen = len(t.context.currentLine)

	allMatched := true
	var container Node
	container = t.Root
	lastChild := container.LastChild()
	for ; nil != lastChild && lastChild.IsOpen(); lastChild = container.LastChild() {
		container = lastChild
		t.context.findNextNonspace()

		switch container.Continue(t.context) {
		case 0: // we've matched, keep going
			break
		case 1: // we've failed to match a block
			allMatched = false
			break
		case 2: // we've hit end of line for fenced code close and can return
			return
		}

		if !allMatched {
			container = container.Parent() // back up to last matching block
			break
		}
	}

	t.context.allClosed = container == t.context.oldtip
	t.context.lastMatchedContainer = container

	matchedLeaf := container.Type() != NodeParagraph && container.AcceptLines()
	var startsLen = len(blockStarts)

	// 除非最后一个匹配到的是代码块，否则的话就起始一个新的块级节点
	for !matchedLeaf {
		t.context.findNextNonspace()

		// 如果不由潜在的节点标记开头 ^[#`~*+_=<>0-9-]，则说明不用继续迭代生成子节点
		// 这里仅做简单判断的话可以略微提升一些性能
		maybeMarker := t.context.currentLine[t.context.nextNonspace]
		if !t.context.indented &&
			itemHyphen != maybeMarker && itemAsterisk != maybeMarker && itemPlus != maybeMarker && // 无序列表
			!isDigit(maybeMarker) && // 有序列表
			itemBacktick != maybeMarker && itemTilde != maybeMarker && // 代码块
			itemCrosshatch != maybeMarker && // ATX 标题
			itemGreater != maybeMarker && // 块引用
			itemLess != maybeMarker && // HTML 块
			itemUnderscore != maybeMarker && itemEqual != maybeMarker { // Setext 标题
			t.context.advanceNextNonspace()
			break
		}

		// 逐个尝试是否可以起始一个块级节点
		var i = 0
		for i < startsLen {
			var res = blockStarts[i](t, container)
			if res == 1 { // 匹配到容器块，继续迭代下降过程
				container = t.context.tip
				break
			} else if res == 2 { // 匹配到叶子块，跳出迭代下降过程
				container = t.context.tip
				matchedLeaf = true
				break
			} else { // 没有匹配到，继续用下一个起始块模式进行匹配
				i++
			}
		}

		if i == startsLen { // nothing matched
			t.context.advanceNextNonspace()
			break
		}
	}

	// offset 后余下的内容算作是文本行，需要将其添加到相应的块节点上

	if !t.context.allClosed && !t.context.blank && t.context.tip.Type() == NodeParagraph {
		// 该行是段落延续文本，直接添加到当前末梢段落上
		t.addLine()
	} else {
		// 最终化未匹配的块
		t.context.closeUnmatchedBlocks()

		if t.context.blank && nil != container.LastChild() {
			container.LastChild().SetLastLineBlank(true)
		}

		typ := container.Type()
		isFenced := NodeCodeBlock == typ && container.(*CodeBlock).isFenced

		// 空行判断，主要是为了判断列表是紧凑模式还是松散模式
		var lastLineBlank = t.context.blank &&
			!(typ == NodeBlockquote || // 块引用行肯定不会是空行因为至少有一个 >
				(typ == NodeCodeBlock && isFenced) || // 围栏代码块不计入空行判断
				(typ == NodeListItem && nil == container.FirstChild())) // 内容为空的列表项也不计入空行判断
		// 因为列表是块级容器（可进行嵌套），所以需要在父节点方向上传播 lastLineBlank
		// lastLineBlank 目前仅在判断列表紧凑模式上使用
		for cont := container; nil != cont; cont = cont.Parent() {
			cont.SetLastLineBlank(lastLineBlank)
		}

		if container.AcceptLines() {
			t.addLine()
			if typ == NodeHTMLBlock {
				// HTML 块（类型 1-5）需要检查是否满足闭合条件
				html := container.(*HTMLBlock)
				if html.hType >= 1 && html.hType <= 5 {
					if t.isHTMLBlockClose(t.context.currentLine[t.context.offset:], html.hType) {
						t.context.finalize(container)
					}
				}
			}
		} else if t.context.offset < t.context.currentLineLen && !t.context.blank {
			// 普通段落开始
			t.context.addChild(&Paragraph{BaseNode: &BaseNode{typ: NodeParagraph, tokens: make(items, 0, 256)}})
			t.context.advanceNextNonspace()
			t.addLine()
		}
	}
}

// blockStartFunc 定义了用于判断块是否开始的函数签名。
type blockStartFunc func(t *Tree, container Node) int

// blockStarts 定义了一系列函数，每个函数用于判断某种块节点是否可以开始，返回值：
// 0：不匹配
// 1：匹配到块容器，需要继续迭代下降
// 2：匹配到叶子块
var blockStarts = []blockStartFunc{
	// 判断块引用（>）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented {
			token := t.context.currentLine.peek(t.context.nextNonspace)
			if itemEnd != token && itemGreater == token {
				t.context.advanceNextNonspace()
				t.context.advanceOffset(1, false)
				// > 后面的空格是可选的
				token = t.context.currentLine.peek(t.context.offset)
				if itemSpace == token || itemTab == token {
					t.context.advanceOffset(1, true)
				}

				t.context.closeUnmatchedBlocks()
				t.context.addChild(&Blockquote{BaseNode: &BaseNode{typ: NodeBlockquote}})
				return 1
			}
		}
		return 0
	},

	// 判断 ATX 标题（#）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented {
			if heading := t.parseATXHeading(); nil != heading {
				t.context.advanceNextNonspace()
				t.context.advanceOffset(len(heading.tokens), false)
				t.context.closeUnmatchedBlocks()
				t.context.addChild(heading)
				t.context.advanceOffset(t.context.currentLineLen-t.context.offset, false)
				return 2
			}
		}
		return 0
	},

	// 判断围栏代码块（```）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented {
			if codeBlock := t.parseFencedCode(); nil != codeBlock {
				t.context.closeUnmatchedBlocks()
				t.context.addChild(codeBlock)
				t.context.advanceNextNonspace()
				t.context.advanceOffset(codeBlock.fenceLength, false)
				return 2
			}
		}
		return 0
	},

	// 判断 HTML 块（<）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented && t.context.currentLine.peek(t.context.nextNonspace) == itemLess {
			tokens := t.context.currentLine[t.context.nextNonspace:]
			if html := t.parseHTML(tokens); nil != html {
				t.context.closeUnmatchedBlocks()
				t.context.addChild(html)
				return 2
			}
		}
		return 0
	},

	// 判断 Setext 标题（- =）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented && container.Type() == NodeParagraph {
			if heading := t.parseSetextHeading(); nil != heading {
				t.context.closeUnmatchedBlocks()
				// 解析链接引用定义
				for tokens := container.Tokens(); 0 < len(tokens) && itemOpenBracket == tokens[0]; tokens = container.Tokens() {
					if remains := t.context.parseLinkRefDef(tokens); nil != remains {
						container.SetTokens(remains)
					} else {
						break
					}
				}

				if value := container.Tokens(); 0 < len(value) {
					heading.tokens = bytes.TrimSpace(value)
					container.InsertAfter(container, heading)
					container.Unlink()
					t.context.tip = heading
					t.context.advanceOffset(t.context.currentLineLen-t.context.offset, false)
					return 2
				}
			}
		}
		return 0
	},

	// 判断分隔线（--- ***）是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented {
			if thematicBreak := t.parseThematicBreak(); nil != thematicBreak {
				t.context.closeUnmatchedBlocks()
				t.context.addChild(thematicBreak)
				t.context.advanceOffset(t.context.currentLineLen-t.context.offset, false)
				return 2
			}
		}
		return 0
	},

	// 判断列表、列表项（* - + 1.）或者任务列表项是否开始
	func(t *Tree, container Node) int {
		if !t.context.indented || container.Type() == NodeList {
			data := t.parseListMarker(container)
			if nil == data {
				return 0
			}

			t.context.closeUnmatchedBlocks()

			listsMatch := container.Type() == NodeList && t.context.listsMatch(container.(*List).listData, data)
			if t.context.tip.Type() != NodeList || !listsMatch {
				t.context.addChild(&List{&BaseNode{typ: NodeList}, data})
			}
			listItem := &ListItem{&BaseNode{typ: NodeListItem}, data}
			t.context.addChild(listItem)

			if 1 == listItem.listData.typ {
				// 修正有序列表项序号
				if prev, ok := listItem.previous.(*ListItem); ok {
					listItem.num = prev.num + 1
				} else {
					listItem.num = data.start
				}
			}

			return 1
		}
		return 0
	},

	// 判断缩进代码块（    code）是否开始
	func(t *Tree, container Node) int {
		if t.context.indented && t.context.tip.Type() != NodeParagraph && !t.context.blank {
			t.context.advanceOffset(4, true)
			t.context.closeUnmatchedBlocks()
			t.context.addChild(&CodeBlock{BaseNode: &BaseNode{typ: NodeCodeBlock}})
			return 2
		}
		return 0
	},
}

// addLine 用于在当前的末梢节点 context.tip 上添加迭代行剩余的所有 tokens。
// 调用该方法前必须确认末梢 tip 能够接受新行。
func (t *Tree) addLine() {
	if t.context.partiallyConsumedTab {
		t.context.offset++ // skip over tab
		// add space characters:
		var charsToTab = 4 - (t.context.column % 4)
		for i := 0; i < charsToTab; i++ {
			t.context.tip.AppendTokens(items{itemSpace})
		}
	}
	t.context.tip.AppendTokens(t.context.currentLine[t.context.offset:])
}
